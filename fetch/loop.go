package fetch

import (
	"log"
	"time"

	"github.com/ArchitBhonsle/go-fpw/pipes"
)

func Loop(
	symbols []string,
	nRetries int,
	sleepInterval time.Duration,
	fetchInterval time.Duration,
	refetchInterval time.Duration,
	cleanup *pipes.Cleanup,
) (<-chan Fetched, <-chan error) {
	resc := make(chan Fetched)
	errc := make(chan error)

	fetchers := make([]Fetcher, 0, len(symbols))
	for _, symbol := range symbols {
		fetcher, err := NewFetcher(symbol, nRetries, refetchInterval)
		if err != nil {
			panic(err)
		}

		fetchers = append(fetchers, fetcher)
	}

	statusChecker, err := NewStatusChecker()
	if err != nil {
		panic(err)
	}

	cleanup.Add()
	go func() {
		defer func() {
			close(resc)
			close(errc)
			cleanup.Done()
		}()

		// wait for a multiple of delta
		now := time.Now()
		next := now.Round(sleepInterval)
		if next.Before(now) {
			next = next.Add(sleepInterval)
		}
		sleepTimer := time.NewTimer(next.Sub(now))

		fetcherIndex := 0
		fetchTimer := time.NewTimer(0)
		<-fetchTimer.C

		for {
			select {
			case <-cleanup.E:
				return
			case <-sleepTimer.C:
				fetcherIndex = 0
				sleepTimer.Reset(sleepInterval)

				if true || statusChecker.Check() {
					fetchTimer.Reset(fetchInterval)
				}
			case <-fetchTimer.C:
				if fetcherIndex == len(fetchers) {
					continue
				}

				log.Println("fetching")

				fetched, err := fetchers[fetcherIndex].Fetch()
				if err != nil {
					errc <- err
				}

				resc <- fetched

				fetcherIndex++
				fetchTimer.Reset(fetchInterval)
			}
		}
	}()

	return resc, errc
}
