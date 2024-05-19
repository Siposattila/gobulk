package config

type Config struct {
	MysqlDSN            string
	MysqlQuery          string
	EmailDSN            string
	SyncCron            string
	SendDelay           uint16
	CompanyName         string
	HttpServerPort      string
	UnsubscribeEndpoint string
	ResubscribeEndpoint string
}
