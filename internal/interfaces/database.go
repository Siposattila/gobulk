package interfaces

import "gorm.io/gorm"

type DatabaseInterface interface {
	GetEntityManager() EntityManagerInterface
	GetMysqlEntityManager() EntityManagerInterface
	GetMysqlDatabaseName() string
}

type EntityManagerInterface interface {
	GetGormORM() *gorm.DB
}
