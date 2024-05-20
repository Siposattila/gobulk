package email

import (
	"github.com/Siposattila/gobulk/internal/interfaces"
	"gorm.io/gorm"
)

const (
	LAST_PROCESS_SEND     = 1
	LAST_PROCESS_VALIDATE = 2
)

type Last struct {
	Offset    int64
	ProcessID uint8
}

func NewLast(offset int64, processId uint8) *Last {
	return &Last{
		Offset:    offset,
		ProcessID: processId,
	}
}

func GetLast(database interfaces.DatabaseInterface, processId int) *Last {
	var last Last
	tx := database.GetEntityManager().GetGormORM().First(&last, "process_id = ?", processId)

	if tx.Error != gorm.ErrRecordNotFound {
		database.GetEntityManager().GetGormORM().Delete(last, "process_id = ?", processId)
	}

	return &last
}
