package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ArchitBhonsle/go-fpw/fetch"
	"github.com/ArchitBhonsle/go-fpw/options"
	"github.com/ArchitBhonsle/go-fpw/process"
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

	cleanExit := func() {
		fmt.Println("exitting")
		close(exit)
		cleanup.Wait()
	}
	defer cleanExit()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanExit()
		os.Exit(1)
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

	errors := merge(exit, &cleanup, fetchErrors, processErrors)

	for {
		select {
		case <-exit:
			return
		case pr := <-processResults:
			fmt.Printf("%v %v processed\n", pr.Timestamp, pr.Underlying)
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
