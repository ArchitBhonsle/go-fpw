package pipes

import (
	"fmt"
	"reflect"
)

func Pipeline(
	generator any,
	pipeFuncs ...any,
) (<-chan struct{}, <-chan error) {
	// Validation complete

	// generator is a channel
	generatorType := reflect.TypeOf(generator)
	if generatorType.Kind() != reflect.Chan && generatorType.ChanDir() != reflect.RecvDir {
		panic("generator should be a recieve only channel")
	}

	// at least one pipeFunc
	if len(pipeFuncs) == 0 {
		panic("there should be at least one pipeFunc")
	}

	// generator matches the first pipeFunc
	generatorOutputType := generatorType.Elem()
	if generatorOutputType != reflect.TypeOf(pipeFuncs[0]).In(0) {
		panic("output of generator does not match the input of pipeFunc[0]")
	}

	// last pipeFunc outputs a struct{}
	if reflect.TypeOf(pipeFuncs[len(pipeFuncs)-1]).Out(0) != reflect.TypeOf(struct{}{}) {
		panic(fmt.Sprintf("output of pipeFunc[%v] should be of the type struct{}{}", len(pipeFuncs)-1))
	}

	// every pipeFunc's input matches the next pipeFunc's input
	for i := 0; i < len(pipeFuncs)-1; i++ {
		curr, next := pipeFuncs[i], pipeFuncs[i+1]
		currType, nextType := reflect.TypeOf(curr), reflect.TypeOf(next)

		fmt.Println(currType, nextType)
		if currType.Out(0) != nextType.In(0) {
			panic(fmt.Sprintf("output of pipeFunc[%v] does not match the input of pipeFunc[%v]", i, i+1))
		}
	}

	// Validation complete
	cleanup := NewCleanup()

	errors := make([]<-chan error, 0)
	lastOutput := generator
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

	return nil, err
}
