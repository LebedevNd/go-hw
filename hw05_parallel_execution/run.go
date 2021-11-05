package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	m = fixErrorsCount(m)

	errorsChannel := make(chan int, m)
	tChan := make(chan Task, len(tasks))

	for _, task := range tasks {
		tChan <- task
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case t := <-tChan:
					runChannelTask(t, errorsChannel, m)
				default:
					return
				}
			}
		}()
	}

	wg.Wait()

	if m > 0 && len(errorsChannel) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func fixErrorsCount(m int) int {
	if m < 0 {
		return 0
	}
	return m
}

func runChannelTask(t Task, errCh chan int, errorsLimit int) {
	if errorsLimit > 0 && len(errCh) >= errorsLimit {
		return
	}

	if t() != nil && errorsLimit > 0 {
		select {
		case errCh <- 1:
		default:
		}
	}
}
