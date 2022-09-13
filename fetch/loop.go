package fetch

import (
	"time"
)

func Loop(
	symbols []string,
	nRetries int,
	sleepInterval time.Duration,
	fetchInterval time.Duration,
	refetchInterval time.Duration,
	exit <-chan struct{},
) (<-chan Fetched, <-chan error, error) {
	resc := make(chan Fetched)
	errc := make(chan error)

	fetchers := make([]Fetcher, 0, len(symbols))
	for _, symbol := range symbols {
		fetcher, err := NewFetcher(symbol, nRetries, refetchInterval)
		if err != nil {
			close(resc)
			close(errc)
			return resc, errc, err
		}

		fetchers = append(fetchers, fetcher)
	}

	statusChecker, err := NewStatusChecker()
	if err != nil {
		close(resc)
		close(errc)
		return resc, errc, err
	}

	go func() {
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
			case <-exit:
				close(resc)
				close(errc)
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

	return resc, errc, nil
}
