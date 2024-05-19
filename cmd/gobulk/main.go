package main

import (
	"flag"

	"github.com/Siposattila/gobulk/internal/bulk"
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
	flag.Bool("validate", false, "This flag will start gobulk's validate process which will validate the email addresses in local db.")
	flag.Bool("bulk", false, "This flag will start gobulk's bulk email sending process.")

	flag.Parse()

	if isFlagPassed("sync") {
		sync := sync.Init()
		sync.Start()
	}

	if isFlagPassed("validate") && !isFlagPassed("sync") {
		validate := validate.Init()
		validate.Start()
	}

	if isFlagPassed("bulk") && !isFlagPassed("sync") {
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
