package workers

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/avegner/utils/channels"
)

var ErrPoolClosed = errors.New("pool closed")

type Pool struct {
	workc  chan *work
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type work struct {
	task  Task
	token *token
}

type Task func(ctx context.Context) error

func NewPool(size ...int) *Pool {
	sz := runtime.GOMAXPROCS(0)
	if len(size) > 0 && size[0] > 0 {
		sz = size[0]
	}
	ctx, cancel := context.WithCancel(context.Background())

	p := Pool{
		workc:  make(chan *work, sz),
		ctx:    ctx,
		cancel: cancel,
	}

	p.wg.Add(sz)
	for i := 0; i < sz; i++ {
		go p.worker()
	}

	return &p
}

func (p *Pool) Dispatch(task Task) (Tokener, error) {
	select {
	case <-p.ctx.Done():
		return nil, ErrPoolClosed
	default:
	}

	w := &work{
		task:  task,
		token: newToken(p.ctx),
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		select {
		case p.workc <- w:
		case <-p.ctx.Done():
			w.token.close(ErrPoolClosed)
		}
	}()

	return w.token, nil
}

func (p *Pool) Close() error {
	select {
	case <-p.ctx.Done():
		return ErrPoolClosed
	default:
	}

	p.cancel()
	p.wg.Wait()

	channels.Drain(p.workc, func(v interface{}) {
		w := v.(*work)
		w.token.close(ErrPoolClosed)
	})
	return nil
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case w := <-p.workc:
			err := w.task(w.token.ctx)
			w.token.close(err)
		case <-p.ctx.Done():
			return
		}
	}
}
