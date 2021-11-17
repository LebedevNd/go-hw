package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, f := range stages {
		Out := make(Bi)
		go func(in In) {
			defer close(Out)
			for {
				select {
				case <-done:
					return
				case a, ok := <-in:
					if !ok {
						return
					}
					Out <- a
				}
			}
		}(in)
		in = f(Out)
	}
	return in
}
