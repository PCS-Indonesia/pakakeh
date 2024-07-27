package workerpool

import (
	"log"
	"os"
	"sync"
)

// PoolImplementor is the Pool interface.
type PoolImplementor interface {
	Start()
	Closed() bool
	GetSize() int
	StopDispatch(...int)
	SetMaxPoolNum(int)
}

type (
	Option         = func(*Config) error
	JobHandlerFunc = func() JobFunc
)

var (
	logger = log.New(os.Stdout, "workerpool:", log.LstdFlags)

	// set as global
	DefaultConfig = Config{
		InitDispatcherNum:  1,
		MaxDispatcherNum:   3,
		WorkerNum:          20,
		JobQueueBufferSize: 1000,
	}
)

// New creates a pool.
func New(done <-chan struct{}, jobHandlerFunc JobHandlerFunc, options ...Option) *Pool {
	pConfig := DefaultConfig
	setOption(&pConfig, options...)

	// validate config
	if pConfig.InitDispatcherNum < 1 {
		logger.Panicln("InitDispatcherNum must greater than 0")
	}

	if pConfig.MaxDispatcherNum < pConfig.InitDispatcherNum {
		pConfig.MaxDispatcherNum = pConfig.InitDispatcherNum
	}

	if pConfig.WorkerNum < 1 {
		logger.Panicln("WorkerNum must greater than 0")
	}

	p := &Pool{
		JobQueue:       make(chan Job, pConfig.JobQueueBufferSize),
		config:         pConfig,
		JobHandlerFunc: jobHandlerFunc,
		done:           done,
		doneDispatcher: make(chan struct{}, pConfig.InitDispatcherNum),
		mu:             &sync.Mutex{},
		wg:             &sync.WaitGroup{},
	}

	if pConfig.Errors {
		p.Errors = make(chan error, 1)
	}

	return p
}

// Start run dispatchers in the pool.
func (p *Pool) Start() {
	dispatcherCounter := make(chan []struct{}, p.config.InitDispatcherNum)

	for range dispatcherCounter {
		p.newDispatcher()
	}

	p.mu.Lock()
	p.poolNum = p.config.InitDispatcherNum
	p.mu.Unlock()

	go p.listen() // listen for pool done signal concurrently
}

func (p *Pool) newDispatcher() {
	j := p.JobHandlerFunc()
	p.wg.Add(1)
	d := NewDispatcher(p.doneDispatcher, p.wg, p.config.WorkerNum,
		p.JobQueue, j, p.Errors)
	d.Dispatch()
}

func (p *Pool) listen() {
	for { // loop until pool is closed or done
		select {
		case _, open := <-p.done:
			if !open {
				close(p.JobQueue)
				p.wg.Wait()

				close(p.doneDispatcher)
				p.mu.Lock()

				if !p.closed {
					p.closed = true
				}

				p.poolNum = 0
				p.mu.Unlock()

				if p.config.Errors {
					close(p.Errors)
				}

				return
			}
		}
	}
}

// SetMaxPoolNum applies MaxPoolNum to Pool Config.
func (p *Pool) SetMaxPoolNum(maxPoolNum int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	opt := func(c *Config) error {
		c.MaxDispatcherNum = maxPoolNum
		return nil
	}
	setOption(&p.config, opt)
}

// StopDispatch signals dispatcher to stop
func (p *Pool) StopDispatch(num ...int) {
	n := Min(1, p.poolNum)
	if len(num) > 0 && num[0] > 1 {
		n = Min(num[0], p.poolNum)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.closed && p.poolNum-n > 0 {
		poolCounter := make(chan []struct{}, n)

		for range poolCounter {
			p.doneDispatcher <- struct{}{}
			p.poolNum--
			p.wg.Done()
		}
	}
}

func (p *Pool) Closed() bool {
	p.mu.Lock()
	// p.closed = false
	defer p.mu.Unlock()

	return p.closed
}

// GetSize returns current number of dispatcher.
func (p *Pool) GetSize() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.poolNum
}
