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
			middlewareChannel := make(Bi)
			go func() {
				defer close(middlewareChannel)
				for {
					select {
					case <-done:
						return
					case data, ok := <-in:
						if !ok {
							return
						}
						select {
						case <-done:
						case middlewareChannel <- data:
						}
					}
				}
			}()

			return stageItem(middlewareChannel)
		}(currentIn)
	}

	return currentIn
}
