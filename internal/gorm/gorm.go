package gorm

import (
	"regexp"
	"strings"
	"sync"

	"github.com/Siposattila/gobulk/internal/config"
	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *Database = &Database{}

type Database struct {
	em                interfaces.EntityManagerInterface
	mem               interfaces.EntityManagerInterface
	mysqlDatabaseName string
	init              sync.Once
	configProvider    interfaces.ConfigProviderInterface
}

type DatabaseProvider struct{}

type EntityManager struct {
	GormORM *gorm.DB
}

func (d *Database) GetEntityManager() interfaces.EntityManagerInterface      { return d.em }
func (d *Database) GetMysqlEntityManager() interfaces.EntityManagerInterface { return d.mem }
func (d *Database) GetMysqlDatabaseName() string                             { return d.mysqlDatabaseName }

func (dp *DatabaseProvider) GetDatabase(configProvider interfaces.ConfigProviderInterface) interfaces.DatabaseInterface {
	db.init.Do(func() { ctor(configProvider) })

	return db
}

func (em *EntityManager) GetGormORM() *gorm.DB { return em.GormORM }

func ctor(configProvider interfaces.ConfigProviderInterface) {
	db.configProvider = configProvider
	db.gorm()
}

func (d *Database) gorm() {
	database, error := gorm.Open(sqlite.Open("gobulk.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if error != nil {
		console.Fatal("Fatal error during connecting to database: " + error.Error())
	}
	console.Success("Connection to the local database was successful.")

	err := database.AutoMigrate(
		&config.Config{},
		&email.Cache{},
		&email.Email{},
	)
	if err != nil {
		console.Fatal("Fatal error during migration: " + err.Error())
	}

	d.em = &EntityManager{
		GormORM: database,
	}

	d.gormExternal(d.configProvider.GetConfig(db).GetMysqlDSN())
}

func (d *Database) gormExternal(dsn string) {
	match, _ := regexp.MatchString(`^[^:]+:[^@]+@tcp\([^:]+\:\d+\)\/[^?]+\?.*$`, dsn)
	if !match {
		console.Fatal("Bad mysql DSN!")
	}

	d.mysqlDatabaseName = dsn[strings.Index(dsn, "/")+1 : strings.Index(dsn, "?")]

	database, error := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if error != nil {
		console.Fatal("Fatal error during connecting to database: " + error.Error())
	}
	console.Success("Connection to the mysql database was successful.")

	d.mem = &EntityManager{
		GormORM: database,
	}
}
