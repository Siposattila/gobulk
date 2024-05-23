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

func (a *App) GetSync() interfaces.SyncInterface             { return a.Sync }
func (a *App) GetBulk() interfaces.BulkInterface             { return a.Bulk }
func (a *App) GetValidation() interfaces.ValidationInterface { return a.Validation }
func (a *App) GetServer() interfaces.ServerInterface         { return a.Server }

func Init(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) *App {
	app := &App{}
	app.Sync = sync.Init(app, database, config)
	app.Bulk = bulk.Init(app, database, config)
	app.Validation = validation.Init(app, database)
	app.Server = server.GetServer(app, database, config)

	return app
}
