package email

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ShutdownChan = make(chan os.Signal, 1)

const (
	LAST_PROCESS_SEND     = 1
	LAST_PROCESS_VALIDATE = 2
)

type Last struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	EmailID   uint
	ProcessID uint8
}

func NewLast(emailId *uint, processId uint8) *Last {
	return &Last{
		EmailID:   *emailId,
		ProcessID: processId,
	}
}

func ListenForKill() {
	signal.Notify(ShutdownChan, os.Interrupt, syscall.SIGTERM)
}
