package write

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/ArchitBhonsle/go-fpw/process"
)

func Loop(
	processResults <-chan process.Data,
	db *sql.DB,
	exit <-chan struct{},
	cleanup *sync.WaitGroup,
) <-chan error {
	errc := make(chan error)

	cleanup.Add(1)
	go func() {
		defer func() {
			fmt.Println("cleanup: write.Loop")
			close(errc)
			cleanup.Done()
		}()

		for {
			select {
			case <-exit:
				return
			case processed := <-processResults:
				fmt.Println(time.Now().Format("15:04:05.000"), processed.Underlying, "writing")

				records := transform(processed)
				for _, record := range records {
					err := writeRecord(db, record)
					if err != nil {
						errc <- err
					}
				}
			}
		}
	}()

	return errc
}
