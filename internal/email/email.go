package email

import (
	"time"

	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/logger"
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

func (e *Email) GetEmail() string               { return e.Email }
func (e *Email) GetName() string                { return e.Name }
func (e *Email) GetStatus() uint8               { return e.Status }
func (e *Email) SetStatus(status uint8)         { e.Status = status }
func (e *Email) GetValid() uint8                { return e.Valid }
func (e *Email) SetValid(valid uint8)           { e.Valid = valid }
func (e *Email) GetSendStatus() uint8           { return e.SendStatus }
func (e *Email) SetSendStatus(sendStatus uint8) { e.SendStatus = sendStatus }

func NewEmail(email *string, name *string) *Email {
	return &Email{
		Email: *email,
	}
}

func (e *Email) ValidateEmail() {
	time.Sleep(300 * time.Millisecond)

	validity := e.verifyEmail()
	e.Valid = validity
	if validity == interfaces.EMAIL_INVALID {
		logger.LogWarning(e.GetEmail() + " is not a valid email.")
	} else {
		logger.LogNormal(e.GetEmail() + " is a valid email.")
	}
}
