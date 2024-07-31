package workerpool

import (
	"sync"
)

type (
	// IDispatcher is the Dispatcher interface.
	DispatcherImplementor interface {
		Run()
		Closed() bool
		DeWorker(...int)
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
		jobFunc:    jobFunc,
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

// Run creates the workers pool and dispatches available jobs.
func (d *Dispatcher) Run() {

	d.wg.Add(d.numWorkers)
	// starting all workers in the dispatcher
	for i := 0; i < d.numWorkers; i++ {
		worker := NewWorker(d.doneWorker, d.workerPool, d.wg, d.jobPool, d.errors)
		worker.Start(d.jobFunc)
	}

	go d.dispatch()
}

func (d *Dispatcher) closeWorkerDoneCh() {
	d.once.Do(func() {
		d.wgJob.Wait()
		close(d.doneWorker)
		d.wg.Wait()
	})
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job, open := <-d.jobQueue:
			if !open {
				d.jobQueue = nil
				d.closeWorkerDoneCh()
				d.wgPool.Done()
				break
			}
			// a job request has been received
			d.wgJob.Add(1)
			go func(j Job) {
				// try to obtain an available worker
				// this will block until a worker is idle
				d.jobPool <- struct{}{}
				w := <-d.workerPool
				// dispatch job to worker's job channel
				w <- j
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

// Closed returns true if dispatcher received a signal to stop.
func (d *Dispatcher) Closed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.closed
}

// DeWorker signals worker to stop,
// num is the number of workers to stop, default to 1.
func (d *Dispatcher) DeWorker(num ...int) {
	n := Min(1, d.numWorkers)
	if len(num) > 0 && num[0] > 1 {
		n = Min(num[0], d.numWorkers)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.closed && d.numWorkers-n > 1 {
		for range Range(n) {
			d.doneWorker <- struct{}{}
			d.numWorkers--
		}
	}
}
