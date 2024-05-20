package validate

import (
	"os"

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
	console.Normal("Validation is started. This may take a long time!!!")

	last := v.getLast()
	offset := 0
	if last != nil {
		offset = int(last.Offset) - 1
	}

	var results []email.Email
	v.database.GetEntityManager().GetGormORM().Where(
		"valid = ? AND status = ?",
		email.EMAIL_INVALID,
		email.EMAIL_STATUS_ACTIVE,
	).Offset(offset).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			select {
			case <-kill.KillCtx.Done():
				last := email.NewLast(&tx.RowsAffected, email.LAST_PROCESS_VALIDATE)
				v.database.GetEntityManager().GetGormORM().Create(last)
				console.Warning("Unexpected shutdown while validating emails. Saving last progress...")

				os.Exit(1)
			default:
				result.ValidateEmail()
				v.database.GetEntityManager().GetGormORM().Save(&result)

				if result.Valid == email.EMAIL_INVALID {
					console.Warning("Email " + result.Email + " is not valid")
				} else {
					console.Normal("Email " + result.Email + " is validated")
				}
			}
		}

		// Returning an error will stop further batch processing
		return nil
	})
	console.Success("Validation finished successfully!")
}

// FIXME: code dup (bulk.go)
func (v *Validate) getLast() *email.Last {
	var last email.Last
	tx := v.database.GetEntityManager().GetGormORM().First(&last, "process_id = ?", email.LAST_PROCESS_VALIDATE)

	if tx.Error == nil {
		v.database.GetEntityManager().GetGormORM().Delete(last)
	}

	return &last
}
