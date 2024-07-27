package workerpool

import (
	"sync"
)

type (
	DispatcherImplementor interface {
		Dispatch()
		Closed() bool
		StopWorker(numWorkers ...int)
	}
)

// NewDispatcher creates a dispatcher.
func NewDispatcher(done <-chan struct{}, wgPool *sync.WaitGroup, numWorkers int, jobQueue <-chan Job,
	jobFunc JobFunc, errors chan error) *Dispatcher {
	wp := make(chan chan Job, numWorkers)
	return &Dispatcher{
		workerPool: wp,
		numWorkers: numWorkers,
		jobQueue:   jobQueue,
		jobHandler: jobFunc,
		wg:         &sync.WaitGroup{},
		wgJob:      &sync.WaitGroup{},
		wgPool:     wgPool,
		errors:     errors,
		done:       done,
		doneWorker: make(chan struct{}, numWorkers),
		jobPool:    make(chan struct{}, numWorkers),
		mu:         &sync.Mutex{},
	}
}

// Dispatch creates the workers pool and dispatches available jobs.
func (d *Dispatcher) Dispatch() {
	// starting the workers
	d.wg.Add(d.numWorkers)
	// starting all workers in the dispatcher
	for i := 0; i < d.numWorkers; i++ {
		worker := NewWorker(d.doneWorker, d.workerPool, d.wg, d.jobPool, d.errors)
		worker.Start(d.jobHandler)
	}

	go d.startDispatch()
}

func (d *Dispatcher) Closed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.closed
}

// StopWorker signals worker to stop. Default worker is 1
func (d *Dispatcher) StopWorker(numWorkers ...int) {
	n := Min(1, d.numWorkers)
	if len(numWorkers) > 0 && numWorkers[0] > 1 {
		n = Min(numWorkers[0], d.numWorkers)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.closed && d.numWorkers-n > 1 {
		workerCounter := make(chan []struct{}, n)

		for range workerCounter {
			d.doneWorker <- struct{}{}
			d.numWorkers--
		}
	}
}

func (d *Dispatcher) startDispatch() {
	for {
		select {
		case job, open := <-d.jobQueue:
			if !open {
				d.jobQueue = nil
				d.closeWorkerDoneCh()
				d.wgPool.Done()
				break
			}

			d.wgJob.Add(1)

			go func(j Job) {
				d.jobPool <- struct{}{} // will block until worker is available
				w := <-d.workerPool
				w <- j // dispatch job to worker's job channel
				d.wgJob.Done()
			}(job)
		case <-d.done:
			d.done = nil

			d.closeWorkerDoneCh()
			close(d.workerPool)
			close(d.jobPool)

			d.mu.Lock()
			if !d.closed {
				d.closed = true
			}
			d.mu.Unlock()

			return
		}
	}
}

func (d *Dispatcher) closeWorkerDoneCh() {
	d.once.Do(func() {
		d.wgJob.Wait()
		close(d.doneWorker)
		d.wg.Wait()
	})
}
