package write

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ArchitBhonsle/go-fpw/pipes"
	"github.com/ArchitBhonsle/go-fpw/process"
)

func Loop(
	processResults <-chan process.Data,
	db *sql.DB,
	cleanup *pipes.Cleanup,
) <-chan error {
	errc := make(chan error)

	cleanup.Add()
	go func() {
		defer func() {
			fmt.Println("cleanup: write.Loop")
			close(errc)
			cleanup.Done()
		}()

		for {
			select {
			case <-cleanup.E:
				return
			case processed := <-processResults:
				fmt.Println(time.Now().Format("15:04:05.000"), processed.Underlying, "writing")

				records := transform(processed)
				for _, record := range records {
					err := WriteRecord(record, db)
					if err != nil {
						errc <- err
					}
				}
			}
		}
	}()

	return errc
}
