package main

import (
	"flag"

	"github.com/Siposattila/gobulk/internal/bulk"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/sync"
	"github.com/Siposattila/gobulk/internal/validate"
)

func main() {
	flag.Bool(
		"sync",
		false,
		`This flag will start gobulk's sync process which will sync the local db with the given mysql one (email, name).
        Can't run this with validate or bulk!`,
	)
	flag.Bool(
		"server",
		false,
		`This flag will start gobulk's unsubscribe server process which will let email owners to unsubscribe.
        You can only run this with or without sync!`,
	)
	flag.Bool("validate", false, "This flag will start gobulk's validate process which will validate the email addresses in local db.")
	flag.Bool("bulk", false, "This flag will start gobulk's bulk email sending process.")

	flag.Parse()

	email.ListenForKill()

	if isFlagPassed("sync") {
		if isFlagPassed("server") {
			server := bulk.InitForServer()
			go server.HttpServer()
		}

		sync := sync.Init()
		sync.Start()
	}

	if isFlagPassed("server") && !isFlagPassed("sync") {
		server := bulk.InitForServer()
		server.HttpServer()
	}

	if isFlagPassed("validate") && (!isFlagPassed("sync") && !isFlagPassed("server")) {
		validate := validate.Init()
		validate.Start()
	}

	if isFlagPassed("bulk") && (!isFlagPassed("sync") && !isFlagPassed("server")) {
		bulk := bulk.Init()
		bulk.Start()
	}

	return
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
