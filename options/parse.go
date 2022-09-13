package options

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Options struct {
	Symbols         []string
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

	symbolsS := os.Getenv("SYMBOLS")
	symbols := strings.Split(symbolsS, " ")
	if len(symbols) == 0 {
		panic("The number of arguments (symbols) should be greater than 0.")
	}

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

	fetchIntervalN, err := strconv.Atoi(os.Getenv("FETCH_INTERVAL"))
	if err != nil {
		panic(err)
	}
	fetchInterval := time.Second * time.Duration(fetchIntervalN)

	refetchIntervalN, err := strconv.Atoi(os.Getenv("REFETCH_INTERVAL"))
	if err != nil {
		panic(err)
	}
	refetchInterval := time.Second * time.Duration(refetchIntervalN)

	return Options{
		Symbols:         symbols,
		NProcessFanout:  nProcessFanout,
		NWriteFanout:    nWriteFanout,
		NRetries:        nRetries,
		SleepInterval:   sleepInterval,
		FetchInterval:   fetchInterval,
		RefetchInterval: refetchInterval,
	}
}
