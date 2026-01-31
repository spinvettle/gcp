package gcp

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Task func()

var (
	PoolClosedError = errors.New("Pool has closed")
	PoolRejectError = errors.New("Pool ")
)

type Pool struct {
	capacity      int32
	runningWorker int32
	queue         chan Task
	state         atomic.Bool
	mu            sync.Locker
}

func New(cap int32) (*Pool, error) {
	if cap <= 0 {
		return nil, errors.New("Errror must >0")
	}
	return &Pool{capacity: cap,
		mu: &sync.Mutex{}}, nil
}
func (p *Pool) IncRuningNum() {
	atomic.AddInt32(&p.runningWorker, 1)
}
func (p *Pool) DecRuningNum() {
	atomic.AddInt32(&p.runningWorker, -1)
}

func (p *Pool) Submit(t Task) error {
	if !p.state.Load() {
		return PoolClosedError
	}
	if atomic.LoadInt32(&p.runningWorker) < p.capacity {
		p.IncRuningNum()
		go p.worker(t)
		return nil
	}
	p.queue <- t

	return nil
}
func (p *Pool) worker(t Task) {

}
func (p *Pool) ShutDown() {

}
