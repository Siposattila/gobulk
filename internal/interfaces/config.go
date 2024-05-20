package interfaces

type ConfigInterface interface {
	GetMysqlDSN() string
	GetMysqlQuery() string
	GetEmailDSN() string
	GetSyncCron() string
	GetSendDelay() uint16
	GetCompanyName() string
	GetHttpServerPort() string
	GetUnsubscribeEndpoint() string
	GetResubscribeEndpoint() string
}

type ConfigProviderInterface interface {
	GetConfig(database DatabaseInterface) ConfigInterface
}
