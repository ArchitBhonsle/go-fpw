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
	defer func() {
		log.Println("exiting")
		cleanup.Cleanup()
	}()

	// generator
	ticker := time.NewTicker(options.SleepInterval)
	generator := make(chan any)
	go func() {
		for t := range ticker.C {
			generator <- t
		}
	}()

	// Does not work as expected and sacrifices type safety
	results, errors := pipes.Pipeline(
		generator,
		fetch.NewFetcherAny(options.Symbol, options.NRetries, options.RefetchInterval),
		process.TransformAny,
		write.NewWriterAny("out/test.db"),
	)

	// consumer
	for {
		select {
		case <-cleanup.E:
			return
		case r := <-results:
			log.Println(r)
		case e := <-errors:
			panic(e)
		}
	}
}
