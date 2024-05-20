package interfaces

import "gorm.io/gorm"

type DatabaseInterface interface {
	GetEntityManager() EntityManagerInterface
	GetMysqlEntityManager() EntityManagerInterface
}

type DatabaseProviderInterface interface {
	GetDatabase(configProvider ConfigProviderInterface) DatabaseInterface
}

type EntityManagerInterface interface {
	GetGormORM() *gorm.DB
}
