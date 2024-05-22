package email

import (
	"time"
)

const (
	EMAIL_STATUS_ACTIVE       = 1
	EMAIL_STATUS_INACTIVE     = 2
	EMAIL_STATUS_UNSUBSCRIBED = 3

	EMAIL_SEND_STATUS_NOT_SENT = 1
	EMAIL_SEND_STATUS_SENT     = 2

	EMAIL_UNKNOWN = 1
	EMAIL_VALID   = 2
	EMAIL_INVALID = 3
)

type Email struct {
	ID         uint `gorm:"primarykey"`
	Email      string
	Name       string
	Status     uint8 `gorm:"default:1"`
	Valid      uint8 `gorm:"default:1"`
	SendStatus uint8 `gorm:"default:1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewEmail(email *string, name *string) *Email {
	return &Email{
		Email: *email,
	}
}

func (e *Email) ValidateEmail() {
	time.Sleep(300 * time.Millisecond)

	validity := e.verifyEmail()
	e.Valid = validity
}
