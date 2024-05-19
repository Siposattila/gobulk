package gorm

import (
	"regexp"

	"github.com/Siposattila/gobulk/internal/config"
	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EntityManager struct {
	GormORM *gorm.DB
}

func Gorm() *EntityManager {
	database, error := gorm.Open(sqlite.Open("gobulk.db"), &gorm.Config{})
	if error != nil {
		console.Fatal("Fatal error during connecting to database: " + error.Error())
	}
	console.Success("Connection to database was successful.")

	em := &EntityManager{
		GormORM: database,
	}

	err := database.AutoMigrate(
		&config.Config{},
		&email.Cache{},
		&email.Email{},
		&email.Last{},
	)
	if err != nil {
		console.Fatal("Fatal error during migration: " + err.Error())
	}

	createBasicConfiguration(em.GormORM)

	return em
}

func GormExternal(dsn *string) *EntityManager {
	match, _ := regexp.MatchString(`^[^:]+:[^@]+@tcp\([^:]+\:\d+\)\/[^?]+\?.*$`, *dsn)
	if !match {
		console.Fatal("Bad mysql DSN!")
	}

	database, error := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if error != nil {
		console.Fatal("Fatal error during connecting to database: " + error.Error())
	}
	console.Success("Connection to database was successful.")

	em := &EntityManager{
		GormORM: database,
	}

	return em
}

func GetConfig(g *gorm.DB) *config.Config {
	if g == nil {
		console.Fatal("You need to connect to the database first!")
	}

	var conf config.Config
	result := g.First(&conf)
	if result.RowsAffected != 1 {
		console.Fatal("Something went wrong when getting config!")
	}

	return &conf
}

func createBasicConfiguration(g *gorm.DB) {
	var conf config.Config
	var result = g.First(&conf)

	if result.RowsAffected != 1 {
		conf = config.Config{
			MysqlDSN:            "root:123456@tcp(localhost:3306)/xy?charset=utf8mb4&parseTime=True&loc=Local",
			SyncCron:            "0 0 * * *",
			EmailDSN:            "smtp://user:pass@localhost:1025",
			SendDelay:           4615,
			MysqlQuery:          "SELECT DISTINCT email, name FROM users;",
			CompanyName:         "GoBulk",
			HttpServerPort:      "2000",
			UnsubscribeEndpoint: "http://localhost:2000/unsub",
			ResubscribeEndpoint: "http://localhost:2000/resub",
		}
		g.Create(conf)
		console.Fatal("Configuration was not found! Basic configuration was created in the local db.")
	}
}
