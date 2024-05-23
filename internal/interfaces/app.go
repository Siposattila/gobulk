package interfaces

type AppInterface interface {
	GetSync() SyncInterface
	GetBulk() BulkInterface
	GetValidation() ValidationInterface
	GetServer() ServerInterface
}
