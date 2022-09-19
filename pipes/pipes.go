package pipes

func Pipe[I any, O any](inputChannel <-chan I, fn func(I) (O, error), cleanup *Cleanup) (<-chan O, <-chan error) {
	resc := make(chan O)
	errc := make(chan error)

	cleanup.wait.Add(1)
	go func() {
		defer func() {
			defer cleanup.wait.Done()
			close(resc)
			close(errc)
		}()

		for {
			select {
			case <-cleanup.E:
				return
			case input := <-inputChannel:
				res, err := fn(input)
				if err != nil {
					errc <- err
				} else {
					resc <- res
				}
			}
		}
	}()

	return resc, errc
}

func PipeWithFanout[I any, O any](inputChannel <-chan I, fn func(I) (O, error), nFanout int, cleanup *Cleanup) (<-chan O, <-chan error) {
	rs, es := make([]<-chan O, 0, nFanout), make([]<-chan error, 0, nFanout)
	for i := 0; i < nFanout; i++ {
		r, e := Pipe(inputChannel, fn, cleanup)
		rs = append(rs, r)
		es = append(es, e)
	}

	return Merge(cleanup, rs...), Merge(cleanup, es...)
}

func Merge[T any](cleanup *Cleanup, channels ...<-chan T) <-chan T {
	out := make(chan T)

	for i := range channels {
		c := channels[i]
		go func() {
			select {
			case <-cleanup.E:
				return
			case v := <-c:
				out <- v
			}
		}()
	}

	return out
}
