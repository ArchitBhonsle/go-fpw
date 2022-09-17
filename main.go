package main

import (
	"fmt"

	"github.com/ArchitBhonsle/go-fpw/fetch"
	"github.com/ArchitBhonsle/go-fpw/options"
	"github.com/ArchitBhonsle/go-fpw/pipes"
	"github.com/ArchitBhonsle/go-fpw/process"
	"github.com/ArchitBhonsle/go-fpw/write"
)

// three stages
// 1. fetch
// 2. process
// 3. write
// make this generic

func main() {
	options := options.ParseOptions()

	// setup the cleanup mechanism
	cleanup := pipes.NewCleanup()

	defer func() {
		fmt.Println("exiting")
		cleanup.Cleanup()
	}()

	// fetch
	fetchResults, fetchErrors := fetch.Loop(
		options.Symbols,
		options.NRetries,
		options.SleepInterval,
		options.FetchInterval,
		options.RefetchInterval,
		cleanup,
	)

	// process
	processResults, processErrors := pipes.PipeWithFanout(
		fetchResults,
		struct{}{},
		process.Transform,
		options.NProcessFanout,
		cleanup,
	)

	// write
	db := write.NewDB("out/test.db")
	writeErrors := write.Loop(processResults, db, cleanup)

	errors := pipes.Merge(cleanup, fetchErrors, processErrors, writeErrors)

	for {
		select {
		case <-cleanup.E:
			return
		case e := <-errors:
			panic(e)
		}
	}
}
