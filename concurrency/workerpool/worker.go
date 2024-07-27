package workerpool

import (
	"sync"
)

type WorkerImplementor interface {
	Start(JobFunc)
	Closed() bool
}

type JobFunc func(Job) error

// NewWorker creates a worker.
func NewWorker(done <-chan struct{}, workerPool chan<- chan Job, wg *sync.WaitGroup,
	jobPool <-chan struct{}, errors chan error) *Worker {

	bucket := make(chan Job, 1)

	return &Worker{
		pool:    workerPool,
		jobPool: jobPool,
		bucket:  bucket,
		done:    done,
		errors:  errors,
		wg:      wg,
		mu:      &sync.Mutex{},
	}
}

// Start will pushes the worker into workerqueue, listens stop sinyal.
func (w *Worker) Start(jobFunc JobFunc) {
	go func() {
		for {
			select {
			case <-w.jobPool:
				// worker has received a job request
				w.pool <- w.bucket
				job := <-w.bucket
				if err := jobFunc(job); err != nil {
					logger.Printf("Error when process job -> %s \n", err.Error())
					if w.errors != nil {
						w.errors <- err
					}
				}
			case <-w.done:
				// worker has received stop sinyal
				w.done = nil
				close(w.bucket)
				w.wg.Done()

				w.mu.Lock()
				if !w.closed {
					w.closed = true
				}
				w.mu.Unlock()

				return
			}
		}
	}()
}

// Closed worker received a signal to stop
func (w *Worker) Closed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closed
}
