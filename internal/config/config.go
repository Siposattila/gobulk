package config

import (
	"sync"

	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/logger"
	"gorm.io/gorm"
)

var instance *Config = &Config{}

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
	tx := database.GetEntityManager().GetGormORM().First(instance)
	if tx.Error == gorm.ErrRecordNotFound {
		instance.MysqlDSN = "root:123456@tcp(localhost:3306)/xy?charset=utf8mb4&parseTime=True&loc=Local"
		instance.SyncCron = "0 0 * * *"
		instance.EmailDSN = "smtp://user:pass@localhost:1025"
		instance.SendDelay = 1000
		instance.MysqlQuery = "SELECT DISTINCT email, name FROM users;"
		instance.CompanyName = "GoBulk"
		instance.HttpServerPort = "2000"
		instance.UnsubscribeEndpoint = "http://localhost:2000/unsub"
		instance.ResubscribeEndpoint = "http://localhost:2000/resub"

		database.GetEntityManager().GetGormORM().Create(instance)
		logger.Fatal("Configuration was not found! Basic configuration was created in the local database.")
	}
}

func GetConfig(database interfaces.DatabaseInterface) interfaces.ConfigInterface {
	instance.init.Do(func() { ctor(database) })

	return instance
}
