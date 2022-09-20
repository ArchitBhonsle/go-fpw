package main

import (
	"log"
	"time"

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

	// generator
	ticker := time.NewTicker(options.SleepInterval)

	// Does not work as expected and sacrifices type safety
	// pipes.Pipeline(
	// 	ticker.C,
	// 	fetch.NewFetcher(options.Symbol, options.NRetries, options.RefetchInterval),
	// 	process.Transform,
	// 	write.NewWriter("out/test.db"),
	// )

	// fetch
	fetchResults, fetchErrors := pipes.Pipe(
		ticker.C,
		fetch.NewFetcher(options.Symbol, options.NRetries, options.RefetchInterval),
		cleanup,
	)

	// process
	processResults, processErrors := pipes.PipeWithFanout(
		fetchResults,
		process.Transform,
		options.NProcessFanout,
		cleanup,
	)

	// write
	writeResults, writeErrors := pipes.PipeWithFanout(
		processResults,
		write.NewWriter("out/test.db"),
		options.NWriteFanout,
		cleanup,
	)

	errors := pipes.Merge(cleanup, fetchErrors, processErrors, writeErrors)

	defer func() {
		log.Println("exiting")
		cleanup.Cleanup()
	}()

	// consumer
	for {
		select {
		case <-cleanup.E:
			return
		case <-writeResults:
		case e := <-errors:
			panic(e)
		}
	}
}
