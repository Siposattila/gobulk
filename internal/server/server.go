package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/Siposattila/gobulk/internal/logger"
)

var instance *server = &server{}

type server struct {
	database interfaces.DatabaseInterface
	config   interfaces.ConfigInterface
	init     sync.Once
}

func ctor(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) {
	instance.database = database
	instance.config = config
}

func GetServer(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) interfaces.ServerInterface {
	instance.init.Do(func() { ctor(database, config) })

	return instance
}

func (s *server) Run() {
	router := http.NewServeMux()
	router.HandleFunc("GET /bulk", s.bulk)
	router.HandleFunc("GET /unsub/{email}", s.unsubscribe)
	router.HandleFunc("GET /resub/{email}", s.resubscribe)

	server := http.Server{
		Addr:    ":" + s.config.GetHttpServerPort(),
		Handler: router,
	}

	logger.Success("Http server is listening on port :" + s.config.GetHttpServerPort())

	go func() {
		<-kill.KillCtx.Done()
		logger.Warning("Shutdown signal received shutting down http server.")
		server.Shutdown(context.Background())
	}()
	server.ListenAndServe()
}
