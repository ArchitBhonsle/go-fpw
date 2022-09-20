package options

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Options struct {
	Symbol          string
	NProcessFanout  int
	NWriteFanout    int
	NRetries        int
	SleepInterval   time.Duration
	FetchInterval   time.Duration
	RefetchInterval time.Duration
}

func ParseOptions() Options {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	symbol := os.Getenv("SYMBOL")

	nProcessFanout, err := strconv.Atoi(os.Getenv("N_PROCESS_FANOUT"))
	if err != nil {
		panic(err)
	}

	nWriteFanout, err := strconv.Atoi(os.Getenv("N_WRITE_FANOUT"))
	if err != nil {
		panic(err)
	}

	nRetries, err := strconv.Atoi(os.Getenv("N_RETRIES"))
	if err != nil {
		panic(err)
	}

	sleepIntervalN, err := strconv.Atoi(os.Getenv("SLEEP_INTERVAL"))
	if err != nil {
		panic(err)
	}
	sleepInterval := time.Second * time.Duration(sleepIntervalN)

	refetchIntervalN, err := strconv.Atoi(os.Getenv("REFETCH_INTERVAL"))
	if err != nil {
		panic(err)
	}
	refetchInterval := time.Second * time.Duration(refetchIntervalN)

	return Options{
		Symbol:          symbol,
		NProcessFanout:  nProcessFanout,
		NWriteFanout:    nWriteFanout,
		NRetries:        nRetries,
		SleepInterval:   sleepInterval,
		RefetchInterval: refetchInterval,
	}
}
