package validate

import (
	"os"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/schollz/progressbar/v3"
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
	console.Normal("Validation is started. This may take a long time!!!")

	last := email.GetLast(v.database, email.LAST_PROCESS_VALIDATE)
	offset := 0
	if last != nil {
		offset = int(last.Offset)
	}

	var (
		results []email.Email
		total   int64
	)
	v.database.GetEntityManager().GetGormORM().Find(&email.Email{}, "status = ?", email.EMAIL_STATUS_ACTIVE).Count(&total)

	bar := progressbar.Default(total)
	v.database.GetEntityManager().GetGormORM().Where(
		"status = ?",
		email.EMAIL_STATUS_ACTIVE,
	).Offset(offset).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			bar.Add(1)
			select {
			case <-kill.KillCtx.Done():
				last := email.NewLast(int64(offset), email.LAST_PROCESS_VALIDATE)
				v.database.GetEntityManager().GetGormORM().Create(last)
				console.Warning("Unexpected shutdown while validating emails. Saving last progress...")

				os.Exit(1)
			default:
				offset += 1
				result.ValidateEmail()
				v.database.GetEntityManager().GetGormORM().Save(result)
			}
		}

		// Returning an error will stop further batch processing
		return nil
	})

	console.Success("Validation finished successfully!")
}
