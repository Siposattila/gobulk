package interfaces

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

type EmailInterface interface {
	GetEmail() string
	GetName() string
	GetStatus() uint8
	SetStatus(status uint8)
	GetValid() uint8
	SetValid(valid uint8)
	GetSendStatus() uint8
	SetSendStatus(sendStatus uint8)
}

type EmailBodyInterface interface {
	GetSubject() string
	GetGreeting() string
	GetMessage() string
	GetFarewell() string
	GetCompany() string
	GetUnsubscribe() string
}

type EmailClientInterface interface {
	SetEmailBody(body EmailBodyInterface)
	Send(e EmailInterface)
}
