package validate

import (
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
	var results []email.Email
	v.EM.GormORM.Where("valid = ? AND status = ?", email.EMAIL_INVALID, email.EMAIL_STATUS_ACTIVE).FindInBatches(&results, 100, func(tx *g.DB, batch int) error {
		for _, result := range results {
			result.ValidateEmail()
			v.EM.GormORM.Save(&result)
		}

		// Returning an error will stop further batch processing
		return nil
	})
	console.Success("Validation finished successfully!")
}
