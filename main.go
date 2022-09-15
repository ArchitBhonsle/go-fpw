package main

import (
	"fmt"
	"sync"

	"github.com/ArchitBhonsle/go-fpw/fetch"
	"github.com/ArchitBhonsle/go-fpw/options"
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
	exit := make(chan struct{})
	cleanup := sync.WaitGroup{}

	defer func() {
		fmt.Println("exiting")
		close(exit)
		cleanup.Wait()
	}()

	// fetch
	fetchResults, fetchErrors, err := fetch.Loop(
		options.Symbols,
		options.NRetries,
		options.SleepInterval,
		options.FetchInterval,
		options.RefetchInterval,
		exit,
		&cleanup,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// process
	processResults, processErrors, err := process.Loop(fetchResults, exit, &cleanup)
	if err != nil {
		fmt.Println(err)
		return
	}

	// write
	db := write.NewDB("out/test.db")
	writeErrors := write.Loop(processResults, db, exit, &cleanup)

	errors := merge(exit, &cleanup, fetchErrors, processErrors, writeErrors)

	for {
		select {
		case <-exit:
			return
		case e := <-errors:
			fmt.Println(e)
			return
		}
	}
}

func merge[T any](exit <-chan struct{}, cleanup *sync.WaitGroup, channels ...<-chan T) <-chan T {
	out := make(chan T)

	for i := range channels {
		c := channels[i]
		go func() {
			select {
			case <-exit:
				return
			case v := <-c:
				out <- v
			}
		}()
	}

	return out
}
