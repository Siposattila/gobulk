package interfaces

type BulkInterface interface {
	StartConsole()
	Start(subject string, greeting string, message string, farewell string, shouldContinue bool)
}
