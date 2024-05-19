package email

import (
	"time"
)

const (
	EMAIL_STATUS_ACTIVE       = 1
	EMAIL_STATUS_INACTIVE     = 2
	EMAIL_STATUS_UNSUBSCRIBED = 3

	EMAIL_VALID   = true
	EMAIL_INVALID = false
)

type Email struct {
	ID        uint `gorm:"primarykey"`
	Email     string
	Name      string
	Status    uint8 `gorm:"default:1"`
	Valid     bool  `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEmail(email *string, name *string) *Email {
	return &Email{
		Email: *email,
	}
}

func (e *Email) ValidateEmail() {
	time.Sleep(1 * time.Second)

	validity := e.verifyEmail()
	e.Valid = validity
}
