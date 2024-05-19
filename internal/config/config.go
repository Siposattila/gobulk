package config

type Config struct {
	MysqlDSN               string
	MysqlQueryNameAndEmail string
	MysqlDateFieldName     string
	EmailDSN               string
	SyncCron               string
	SendDelay              uint16
	CompanyName            string
	UnsubscribeEndpoint    string
}
