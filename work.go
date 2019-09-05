package future

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type workgroup struct {
	wg    sync.WaitGroup
	ctx   context.Context
	limit chan struct{}
}

type Future struct {
	group  *workgroup
	err    error
	result interface{}
	cost   time.Duration
}

func (f *Future) Get() (result interface{}, err error) {
	if f.Wait() {
		return nil, fmt.Errorf("Get Future timeout ")
	}
	return f.result, f.err
}

func (f *Future) Wait() (timeout bool) {
	return f.group.Wait()
}

func NewGroupWithContext(ctx context.Context) *workgroup {
	return &workgroup{ctx: ctx, limit: make(chan struct{}, runtime.GOMAXPROCS(0))}
}

func NewGroup() *workgroup {
	return &workgroup{ctx: context.Background(), limit: make(chan struct{}, runtime.GOMAXPROCS(0))}
}

func NewGroupWithTimeout(duration time.Duration) *workgroup {
	ctx, _ := context.WithTimeout(context.Background(), duration)
	return &workgroup{ctx: ctx, limit: make(chan struct{}, runtime.GOMAXPROCS(0))}
}

func (w *workgroup) WithSize(size int) *workgroup {
	w.limit = make(chan struct{}, size)
	return w
}

func (w *workgroup) Wait() (timeout bool) {
	done := make(chan bool, 1)
	go func() {
		w.wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		return false
	case <-w.ctx.Done():
		return true
	}
}

func (w *workgroup) Submit(f func() (interface{}, error)) *Future {
	w.wg.Add(1)
	future := &Future{group: w}
	go func() {
		defer func() { <-w.limit }()
		defer w.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case error:
					future.err = err.(error)
				}
			}
		}()
		w.limit <- struct{}{}
		bg := time.Now()
		r, e := f()
		future.cost = time.Since(bg)
		future.result = r
		future.err = e
	}()
	return future
}
