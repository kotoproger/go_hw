package hw06pipelineexecution

import (
	"sync/atomic"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currentIn := in
	skip := atomic.Bool{}
	go func() {
		<-done
		skip.Store(true)
	}()

	for _, stageItem := range stages {
		currentIn = func(in In) (out Out) {
			middlewareChannel := make(Bi, 1)
			go func() {
				defer close(middlewareChannel)
				for data := range in {
					if !skip.Load() {
						middlewareChannel <- data
					}
				}
			}()

			return stageItem(middlewareChannel)
		}(currentIn)
	}

	return currentIn
}
