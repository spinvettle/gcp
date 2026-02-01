package gcp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task func()

type Pool struct {
	isClosed         atomic.Bool
	capacity         int32
	numWorkers       int32
	timeout          time.Duration
	availableWorkers []*Worker
	panicHandlerFn   func(any)
	mu               sync.Mutex
	wg               sync.WaitGroup
}

func New(capacity int32) *Pool {
	if capacity <= 0 {
		return nil
	}
	return &Pool{
		capacity:         capacity,
		availableWorkers: make([]*Worker, capacity),
	}
}

func (p *Pool) putWorker(w *Worker) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isClosed.Load() {
		return PoolClosedError
	}

	p.availableWorkers = append(p.availableWorkers, w)

	return nil
}
func (p *Pool) getWorker() *Worker {
	p.mu.Lock()
	defer p.mu.Unlock()
	size := len(p.availableWorkers)
	if size >= 1 {
		item := p.availableWorkers[size-1]
		p.availableWorkers[size-1] = nil
		p.availableWorkers = p.availableWorkers[:size-1]
		return item

	}
	return nil
}

func (p *Pool) Submit(task Task) error {
	if p.isClosed.Load() {
		return PoolClosedError
	}
	w := p.getWorker()
	if w == nil {
		if atomic.LoadInt32(&p.numWorkers) >= p.capacity {
			return PoolFullError
		}

		w = &Worker{
			pool:      p,
			taskQueue: make(chan Task, 1)}
		p.wg.Add(1)
		go w.run()
		w.taskQueue <- task
		atomic.AddInt32(&p.numWorkers, 1)
		return nil
	} else {
		w.taskQueue <- task
		return nil
	}
}

func (p *Pool) ShutDown() {
	p.isClosed.Store(true)
	p.wg.Wait()

}

type Options func(p *Pool)

func Capacity(capacity int32) Options {
	return func(p *Pool) {
		p.capacity = capacity
	}
}

type Worker struct {
	pool      *Pool
	taskQueue chan Task
	lastTime  time.Time
}

func (w *Worker) run() {
	defer func() {
		w.pool.wg.Done()
		atomic.AddInt32(&w.pool.numWorkers, -1)
		if r := recover(); r != nil {
			if w.pool.panicHandlerFn != nil {
				w.pool.panicHandlerFn(r)
			} else {
				fmt.Println("worker error:" + r.(error).Error())
			}
		}

	}()

	for {
		select {
		case task := <-w.taskQueue:
			task()
			w.lastTime = time.Now()
			err := w.putBackPool() //用完放回
			if err != nil {        //放回失败
				return
			}

		case <-time.After(w.pool.timeout):
			//超时销毁
			return
		}
	}
}

func (w *Worker) putBackPool() error {

	return w.pool.putWorker(w)

}
