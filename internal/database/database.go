package database

import (
	"regexp"
	"strings"
	"sync"

	"github.com/Siposattila/gobulk/internal/config"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/Siposattila/gobulk/internal/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

var instance *database = &database{}

type database struct {
	em                interfaces.EntityManagerInterface
	mem               interfaces.EntityManagerInterface
	mysqlDatabaseName string
	init              sync.Once
	config            interfaces.ConfigInterface
}

type entityManager struct {
	GormORM *gorm.DB
}

func (d *database) GetEntityManager() interfaces.EntityManagerInterface      { return d.em }
func (d *database) GetMysqlEntityManager() interfaces.EntityManagerInterface { return d.mem }
func (d *database) GetMysqlDatabaseName() string                             { return d.mysqlDatabaseName }

func GetDatabase(config interfaces.ConfigInterface) interfaces.DatabaseInterface {
	instance.init.Do(func() { ctor(config) })

	return instance
}

func (em *entityManager) GetGormORM() *gorm.DB { return em.GormORM }

func ctor(config interfaces.ConfigInterface) {
	instance.config = config
	instance.database()
}

func (d *database) database() {
	database, error := gorm.Open(sqlite.Open("gobulk.db"), &gorm.Config{
		Logger: gl.Default.LogMode(gl.Silent),
	})
	if error != nil {
		logger.Fatal("Fatal error during connecting to database: " + error.Error())
	}
	logger.Success("Connection to the local database was successful.")

	err := database.AutoMigrate(
		&config.Config{},
		&email.Cache{},
		&email.Email{},
	)
	if err != nil {
		logger.Fatal("Fatal error during migration: " + err.Error())
	}

	d.em = &entityManager{
		GormORM: database,
	}

	go func() {
		<-kill.KillCtx.Done()
		logger.LogWarning("Shutdown signal received closing local database connection.")
		db, _ := database.DB()
		db.Close()
	}()

	d.databaseExternal(config.GetConfig(instance).GetMysqlDSN())
}

func (d *database) databaseExternal(dsn string) {
	match, _ := regexp.MatchString(`^[^:]+:[^@]+@tcp\([^:]+\:\d+\)\/[^?]+\?.*$`, dsn)
	if !match {
		logger.Fatal("Bad mysql DSN!")
	}

	d.mysqlDatabaseName = dsn[strings.Index(dsn, "/")+1 : strings.Index(dsn, "?")]

	database, error := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gl.Default.LogMode(gl.Silent),
	})
	if error != nil {
		logger.Fatal("Fatal error during connecting to database: " + error.Error())
	}
	logger.Success("Connection to the mysql database was successful.")

	d.mem = &entityManager{
		GormORM: database,
	}

	go func() {
		<-kill.KillCtx.Done()
		logger.LogWarning("Shutdown signal received closing external database connection.")
		db, _ := database.DB()
		db.Close()
	}()
}
