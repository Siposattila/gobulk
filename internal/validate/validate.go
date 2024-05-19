package validate

import (
	"os"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/gorm"
	g "gorm.io/gorm"
)

type Validate struct {
	EM *gorm.EntityManager
}

func Init() *Validate {
	validate := Validate{
		EM: gorm.Gorm(),
	}

	return &validate
}

func (v *Validate) Start() {
	console.Normal("Validation is started. This may take a long time!!!")

	last := v.getLast()
	offset := 0
	if last != nil {
		offset = int(last.EmailID) - 1
	}

	var results []email.Email
	v.EM.GormORM.Where(
		"valid = ? AND status = ?",
		email.EMAIL_INVALID,
		email.EMAIL_STATUS_ACTIVE,
	).Offset(offset).FindInBatches(&results, 100, func(tx *g.DB, batch int) error {
		for _, result := range results {
			select {
			case <-email.ShutdownChan:
				last := email.NewLast(&result.ID, email.LAST_PROCESS_VALIDATE)
				v.EM.GormORM.Create(last)
				console.Warning("Unexpected shutdown while validating emails. Saving last progress...")

				os.Exit(1)
			default:
				result.ValidateEmail()
				v.EM.GormORM.Save(&result)

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
	tx := v.EM.GormORM.First(&last, "process_id = ?", email.LAST_PROCESS_VALIDATE)

	if tx.Error == nil {
		v.EM.GormORM.Delete(last)
	}

	return &last
}
