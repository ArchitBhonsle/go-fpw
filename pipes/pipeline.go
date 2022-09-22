package pipes

func Pipeline(
	generator <-chan any,
	pipeFuncs ...any,
) (<-chan any, <-chan error) {
	cleanup := NewCleanup()

	errors := make([]<-chan error, 0)
	lastOutput := any(generator)
	for _, pipeFunc := range pipeFuncs {
		out, err := Pipe(
			lastOutput.(<-chan any),
			pipeFunc.(func(any) (any, error)),
			cleanup,
		)
		lastOutput = out
		errors = append(errors, err)
	}
	err := Merge(cleanup, errors...)

	return lastOutput.(<-chan any), err
}
