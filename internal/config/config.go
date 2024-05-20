package config

import (
	"sync"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"gorm.io/gorm"
)

var conf *Config = &Config{}

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
	init                sync.Once
}

type ConfigProvider struct{}

func (c *Config) GetMysqlDSN() string            { return c.MysqlDSN }
func (c *Config) GetMysqlQuery() string          { return c.MysqlQuery }
func (c *Config) GetEmailDSN() string            { return c.EmailDSN }
func (c *Config) GetSyncCron() string            { return c.SyncCron }
func (c *Config) GetSendDelay() uint16           { return c.SendDelay }
func (c *Config) GetCompanyName() string         { return c.CompanyName }
func (c *Config) GetHttpServerPort() string      { return c.HttpServerPort }
func (c *Config) GetUnsubscribeEndpoint() string { return c.UnsubscribeEndpoint }
func (c *Config) GetResubscribeEndpoint() string { return c.ResubscribeEndpoint }

func ctor(database interfaces.DatabaseInterface) {
	tx := database.GetEntityManager().GetGormORM().First(conf)
	if tx.Error == gorm.ErrRecordNotFound {
		conf.MysqlDSN = "root:123456@tcp(localhost:3306)/xy?charset=utf8mb4&parseTime=True&loc=Local"
		conf.SyncCron = "0 0 * * *"
		conf.EmailDSN = "smtp://user:pass@localhost:1025"
		conf.SendDelay = 1000
		conf.MysqlQuery = "SELECT DISTINCT email, name FROM users;"
		conf.CompanyName = "GoBulk"
		conf.HttpServerPort = "2000"
		conf.UnsubscribeEndpoint = "http://localhost:2000/unsub"
		conf.ResubscribeEndpoint = "http://localhost:2000/resub"

		database.GetEntityManager().GetGormORM().Create(conf)
		console.Fatal("Configuration was not found! Basic configuration was created in the local database.")
	}
}

func (cp *ConfigProvider) GetConfig(database interfaces.DatabaseInterface) interfaces.ConfigInterface {
	conf.init.Do(func() { ctor(database) })

	return conf
}
