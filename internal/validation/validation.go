package validation

import (
	"errors"
	"strconv"

	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/Siposattila/gobulk/internal/logger"
	"gorm.io/gorm"
)

type validation struct {
	app      interfaces.AppInterface
	database interfaces.DatabaseInterface
}

func Init(app interfaces.AppInterface, database interfaces.DatabaseInterface) interfaces.ValidationInterface {
	return &validation{
		app:      app,
		database: database,
	}
}

func (v *validation) Start() {
	logger.Normal("Validation is started. This may take a long time!")

	var (
		results []email.Email
		total   int64
	)
	v.database.GetEntityManager().GetGormORM().Find(
		&email.Email{},
		"status = ? AND valid = ?",
		interfaces.EMAIL_STATUS_ACTIVE, interfaces.EMAIL_UNKNOWN,
	).Count(&total)

	logger.Normal("Validating " + strconv.Itoa(int(total)) + " emails.")
	master := newMaster(total, 5)
	master.Start()

	v.database.GetEntityManager().GetGormORM().Where(
		"status = ? AND valid = ?",
		interfaces.EMAIL_STATUS_ACTIVE, interfaces.EMAIL_UNKNOWN,
	).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			select {
			case <-kill.KillCtx.Done():
				logger.Warning("Shutdown signal received shutting down validation process.")
				master.Stop()

				return errors.New("Shutdown")
			default:
				master.NewWork(func() {
					result.ValidateEmail()
					v.database.GetEntityManager().GetGormORM().Save(result)
				})
			}
		}

		return nil
	})
	master.Wait()

	logger.Success("Validation finished successfully!")
}
