package workerpool

import "sync"

type (
	// Job represents the job to be run.
	Job struct {
		Data interface{}
	}
)

// NewWorker creates a worker.
func NewWorker(done <-chan struct{}, workerPool chan<- chan Job, wg *sync.WaitGroup,
	jobPool <-chan struct{}, errors chan error) *Worker {
	return &Worker{
		pool:    workerPool,
		jobPool: jobPool,
		bucket:  make(chan Job, 1),
		done:    done,
		errors:  errors,
		wg:      wg,
		mu:      &sync.Mutex{},
	}
}

// Start pushes the worker into worker queue, listens for signal to stop.
func (w *Worker) Start(handler JobFunc) {

	go func() {
		for {
			select {
			case <-w.jobPool:

				// worker has received a job request
				w.pool <- w.bucket
				job := <-w.bucket
				if err := handler(job); err != nil {
					logger.Printf("Error handling job -> %s \n", err.Error())
					if w.errors != nil {
						w.errors <- err
					}
				}

			case <-w.done:
				// worker has received a signal to stop
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

// Closed returns true if worker received a signal to stop.
func (w *Worker) Closed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closed
}
