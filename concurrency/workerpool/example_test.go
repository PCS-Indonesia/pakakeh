package workerpool_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/PCS-Indonesia/pakakeh/concurrency/workerpool"
)

func ExamplePool(t *testing.T) {
	done := make(chan struct{})
	mu := &sync.RWMutex{}
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	sum := 0
	jobHandlerFunc := func() workerpool.JobFunc {
		return func(j workerpool.Job) error {
			mu.Lock()
			defer mu.Unlock()
			sum += j.Data.(int)
			return nil
		}
	}

	size := 10

	opt := func(c *workerpool.Config) error {
		c.InitDispatcherNum = size
		c.WorkerNum = 1000
		return nil
	}

	p := workerpool.New(done, jobHandlerFunc, opt)
	p.Start()

	for i := range data {
		p.JobQueue <- workerpool.Job{
			Data: data[i],
		}
	}

	close(done)

	// wait for jobs to finish
	for {
		// time.Sleep(1 * time.Second)
		if p.Closed() {
			break
		}
	}

	mu.RLock()
	fmt.Println(sum)
	mu.RUnlock()

}
