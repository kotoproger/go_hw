package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currentIn := in

	for _, stageItem := range stages {
		currentIn = func(in In) (out Out) {
			middlewareChannel := make(Bi, 1)

			stageOutput := stageItem(in)

			go func(in In, out Bi) {
				for {
					select {
					case <-done:
						close(out)
						for range in { //nolint:revive
						}
						return
					case data, ok := <-in:
						if !ok {
							close(out)
							return
						}
						out <- data
					}
				}
			}(stageOutput, middlewareChannel)
			return middlewareChannel
		}(currentIn)
	}

	return currentIn
}
