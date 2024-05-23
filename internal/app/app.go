package app

import (
	"github.com/Siposattila/gobulk/internal/bulk"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/server"
	"github.com/Siposattila/gobulk/internal/sync"
	"github.com/Siposattila/gobulk/internal/validation"
)

type App struct {
	Sync       interfaces.SyncInterface
	Bulk       interfaces.BulkInterface
	Validation interfaces.ValidationInterface
	Server     interfaces.ServerInterface
}

func Init(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) *App {
	return &App{
		Sync:       sync.Init(database, config),
		Bulk:       bulk.Init(database, config),
		Validation: validation.Init(database),
		Server:     server.GetServer(database, config),
	}
}
