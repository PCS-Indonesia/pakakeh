package workerpool_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/PCS-Indonesia/pakakeh/concurrency/workerpool"
)

func TestExamplePool(t *testing.T) {
	done := make(chan struct{})
	mu := &sync.RWMutex{}
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sum := 0
	jobHandlerFunc := func() workerpool.JobFunc {
		return func(j workerpool.Job) error {
			mu.Lock()
			defer mu.Unlock()
			sum += j.Data.(int)
			return nil
		}
	}

	size := 2

	opt := func(c *workerpool.Config) error {
		c.InitDispatcherNum = size
		c.WorkerNum = 5
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
