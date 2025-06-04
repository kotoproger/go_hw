package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, gorutineCount, errorsCountLimit int) error {
	if gorutineCount == 0 {
		gorutineCount = len(tasks)
	}

	tasksChannel := make(chan Task, gorutineCount)

	outputs := make([]<-chan error, gorutineCount)
	for i := 0; i < gorutineCount; i++ {
		outputs[i] = func() <-chan error {
			output := make(chan error)
			go func() {
				defer close(output)

				worker(tasksChannel, output)
			}()
			return output
		}()
	}

	taskSendNumber := 0
	for ; taskSendNumber < gorutineCount; taskSendNumber++ {
		tasksChannel <- tasks[taskSendNumber]
	}

	tasksDoneCount, errorsCount := 0, 0
	workStopped := false

	for {
		hasActiveChannel, isError, hasData := processChannels(outputs)
		if isError {
			errorsCount++
		}
		if hasData {
			tasksDoneCount++
		}

		if !workStopped && (errorsCount >= errorsCountLimit || tasksDoneCount >= taskSendNumber) {
			close(tasksChannel)
			workStopped = true
		}
		if hasData && !workStopped && (taskSendNumber+1) <= len(tasks) {
			tasksChannel <- tasks[taskSendNumber]
			taskSendNumber++
		}

		if !hasActiveChannel {
			break
		}
	}

	if errorsCount >= errorsCountLimit {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func processChannels(outputs []<-chan error) (hasActiveChannel bool, isError bool, hasData bool) {
	isError = false
	var data error
	for {
		hasActiveChannel = false
		for _, channel := range outputs {
			select {
			case data, hasData = <-channel:
				if hasData {
					hasActiveChannel = true
				} else {
					continue
				}
				if data != nil {
					isError = true
				}
				return
			default:
				hasActiveChannel = true
			}
		}
		if !hasActiveChannel {
			break
		}
	}
	return
}

func worker(input <-chan Task, output chan<- error) {
	for task := range input {
		output <- task()
	}
}
