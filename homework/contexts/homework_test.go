package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Group struct {
	wg     *sync.WaitGroup
	err    chan error
	ctx    context.Context
	cancel context.CancelFunc
}

func NewErrGroup(parentCtx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	return &Group{
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
		err:    make(chan error, 1),
	}, ctx
}

func (g *Group) Go(action func() error) {
	if action == nil {
		panic("nil action")
	}
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := action(); err != nil {
			select {
			case g.err <- err:
				g.cancel()
			case <-g.ctx.Done():
			}
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	close(g.err)

	select {
	case err := <-g.err:
		return err
	default:
		return nil
	}
}

func TestErrGroupWithoutError(t *testing.T) {
	var counter atomic.Int32
	group, _ := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			time.Sleep(time.Second)
			counter.Add(1)
			return nil
		})
	}

	err := group.Wait()
	assert.Equal(t, int32(5), counter.Load())
	assert.NoError(t, err)
}

func TestErrGroupWithError(t *testing.T) {
	var counter atomic.Int32
	group, ctx := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				counter.Add(1)
				return nil
			}
		})
	}

	group.Go(func() error {
		return errors.New("error")
	})

	err := group.Wait()
	assert.Equal(t, int32(0), counter.Load())
	assert.Error(t, err)
}
