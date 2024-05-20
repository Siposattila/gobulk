package email

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
