package main

import (
	"fmt"

	"github.com/ArchitBhonsle/go-pipes/fetch"
	"github.com/ArchitBhonsle/go-pipes/options"
)

// three stages
// 1. fetch
// 2. process
// 3. write

func main() {
	options := options.ParseOptions()

	exit := make(chan struct{})

	fetchResults, fetchError, err := fetch.Loop(
		options.Symbols,
		options.NRetries,
		options.SleepInterval,
		options.FetchInterval,
		options.RefetchInterval,
		exit,
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println("exitted")
		close(exit)
	}()

	for {
		select {
		case fr := <-fetchResults:
			fmt.Println(fr.Records.Timestamp, fr.Records.UnderlyingValue)
		case fe := <-fetchError:
			fmt.Println(fe)
			exit <- struct{}{}
		}
	}
}
