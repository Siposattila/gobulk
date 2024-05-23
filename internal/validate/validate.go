package validate

import (
	"errors"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"gorm.io/gorm"
)

type Validate struct {
	database interfaces.DatabaseInterface
}

func Init(database interfaces.DatabaseInterface) *Validate {
	return &Validate{
		database: database,
	}
}

func (v *Validate) Start() {
	console.Normal("Validation is started. This may take a long time!")

	var (
		results []email.Email
		total   int64
	)
	v.database.GetEntityManager().GetGormORM().Find(
		&email.Email{},
		"status = ? AND valid = ?",
		email.EMAIL_STATUS_ACTIVE, email.EMAIL_UNKNOWN,
	).Count(&total)

	master := NewMaster(total, 5)
	master.Start()

	v.database.GetEntityManager().GetGormORM().Where(
		"status = ? AND valid = ?",
		email.EMAIL_STATUS_ACTIVE, email.EMAIL_UNKNOWN,
	).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			select {
			case <-kill.KillCtx.Done():
				master.Stop()
				console.Warning("Shutdown signal received shutting down validation process.")

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

	console.Success("Validation finished successfully!")
}
