package hw06_pipeline_execution //nolint:golint,stylecheck,revive

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		ch := make(Bi)
		go func(ch Bi, in In) {
			defer close(ch)

			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					ch <- val
				}
			}
		}(ch, in)
		in = stage(ch)
	}

	return in
}
