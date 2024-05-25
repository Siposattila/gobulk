package main

import (
	"flag"

	"github.com/Siposattila/gobulk/internal/app"
	"github.com/Siposattila/gobulk/internal/config"
	"github.com/Siposattila/gobulk/internal/database"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/Siposattila/gobulk/internal/logger"
)

var (
	version   = "v0.0.0"
	buildTime = "0000.00.00."
	buildHash = "0000000000"
)

func main() {
	flag.Bool("up", false, "This flag will start gobulk as a process.")
	flag.Bool(
		"sync",
		false,
		"This flag will start gobulk's sync process which will sync the local db with the given mysql one (email, name).",
	)
	flag.Bool("validate", false, "This flag will start gobulk's validate process which will validate the email addresses in local db.")
	flag.Bool("bulk", false, "This flag will start gobulk's bulk email sending process.")
	flag.Bool("version", false, "Version information of gobulk.")
	flag.Parse()

	kill.ListenForKill()

	var (
		conf   interfaces.ConfigInterface
		gobulk *app.App
	)
	if !isFlagPassed("version") {
		gobulk = app.Init(database.GetDatabase(conf), config.GetConfig(nil))
	}

	if !isFlagPassed("up") {
		if isFlagPassed("sync") {
			gobulk.Sync.Start()
		} else if isFlagPassed("validate") {
			gobulk.Validation.Start()
		} else if isFlagPassed("bulk") {
			gobulk.Bulk.StartConsole()
		} else if isFlagPassed("version") {
			logger.Normal("\nVersion: ", version, "\nBuild: ", buildTime, "\nHash: ", buildHash, "\n")
		}
	} else {
		// go app.Sync.Start()
		go gobulk.Server.Run()
		// go app.Validation.Start()

		<-kill.KillCtx.Done()
		logger.Warning("Shutdown completed.")
	}
}

func isFlagPassed(name string) bool {
	var found = false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}
