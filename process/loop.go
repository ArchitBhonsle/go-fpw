package process

import (
	"fmt"
	"sync"
	"time"

	"github.com/ArchitBhonsle/go-fpw/fetch"
)

func Loop(
	fetchResults <-chan fetch.Fetched,
	exit <-chan struct{},
	cleanup *sync.WaitGroup,
) (<-chan Data, <-chan error, error) {
	resc := make(chan Data)
	errc := make(chan error)

	cleanup.Add(1)
	go func() {
		defer func() {
			fmt.Println("cleanup: process.Loop")
			close(resc)
			close(errc)
			cleanup.Done()
		}()

		for {
			select {
			case <-exit:
				return
			case fetched := <-fetchResults:
				fmt.Println(time.Now().Format("15:04:05.000"), fetched.Records.Data[0].PE.Underlying, "processing")

				res, err := transform(fetched)
				if err != nil {
					errc <- err
				} else {
					resc <- res
				}
			}
		}
	}()

	return resc, errc, nil
}
