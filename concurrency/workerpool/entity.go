package workerpool

import "sync"

type (

	// JobFunc completes the job.
	JobFunc func(Job) error

	// Pool represents a pool with dispatcher.
	Pool struct {
		JobQueue       chan Job // JobQueue channel for incoming job request
		JobHandlerFunc JobHandlerFunc
		config         Config
		poolNum        int // current number of dispatcher pool
		Errors         chan error
		done           <-chan struct{} // done channel signals the pool to stop,
		doneDispatcher chan struct{}
		wg             *sync.WaitGroup
		mu             *sync.Mutex
		closed         bool
	}

	// Config
	Config struct {
		InitDispatcherNum  int
		MaxDispatcherNum   int
		WorkerNum          int
		JobQueueBufferSize int
		Errors             bool // if true, pool will send errors to Errors channel
	}

	Worker struct {
		pool    chan<- chan Job
		jobPool <-chan struct{}
		bucket  chan Job
		done    <-chan struct{}
		errors  chan error
		wg      *sync.WaitGroup
		closed  bool
		mu      *sync.Mutex
	}

	// Dispatcher handling dispatch the job to the worker.
	Dispatcher struct {
		// a pool of workers bucket
		workerPool chan chan Job
		jobPool    chan struct{}
		jobQueue   <-chan Job
		doneWorker chan struct{}
		once       sync.Once
		errors     chan error
		done       <-chan struct{}
		wg         *sync.WaitGroup
		wgJob      *sync.WaitGroup
		wgPool     *sync.WaitGroup
		closed     bool
		mu         *sync.Mutex
		numWorkers int
		jobFunc    JobFunc
	}
)
