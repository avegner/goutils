package workers

import (
	"context"
	"errors"
)

var ErrTokenClosed = errors.New("token closed")

type Tokener interface {
	WaitWithContext(ctx context.Context) error
	Cancel()
}

type token struct {
	errc   chan error
	ctx    context.Context
	cancel context.CancelFunc
}

func newToken(ctx context.Context) *token {
	ctx, cancel := context.WithCancel(ctx)
	return &token{
		errc:   make(chan error, 1),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *token) close(err error) {
	t.cancel()
	t.errc <- err
	close(t.errc)
}

func (t *token) WaitWithContext(ctx context.Context) error {
	select {
	case err, ok := <-t.errc:
		if !ok {
			err = ErrTokenClosed
		}
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *token) Cancel() {
	t.cancel()
}
