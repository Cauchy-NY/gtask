package gtask

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

var (
	ErrTaskTimeout = errors.New("task timeout")
	ErrParentFail  = errors.New("parent task fail and zombie mode on")
)

type GoPool struct {
	TokenBucket *TokenBucket
	TaskOptions TaskOptions
	FastFail    bool
}

func NewGoPool(maxGoroutineSize int) *GoPool {
	return &GoPool{
		NewTokenBucket(maxGoroutineSize),
		TaskOptions{},
		false,
	}
}

func (p *GoPool) Go(ctx context.Context, f func()) error {
	err := p.TokenBucket.Get(p.FastFail)
	if err != nil {
		return fmt.Errorf("GoPool get TokenBucket err: %s", err)
	}
	go func() {
		defer p.TokenBucket.Put()
		f()
	}()
	return nil
}

func (p *GoPool) do(ctx context.Context, task *Task) *RunResult {
	ret := &RunResult{
		Identifier: task.Identifier,
	}
	if p.TaskOptions.ZombieMode && task.state == 2 {
		ret.Err = ErrParentFail
		return ret
	}
	err := p.TokenBucket.Get(p.FastFail)
	if err != nil {
		ret.Err = err
		return ret
	}
	defer p.TokenBucket.Put()

	done := make(chan struct{}, 1)
	if task.Timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, task.Timeout)
		ctx = timeoutCtx
		defer cancel()
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				ret.Buf = string(debug.Stack())
				ret.Err = fmt.Errorf("run task %v error: %v", ret.Identifier, err)
				done <- struct{}{}
			}
		}()
		cancleCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		ret.Err = task.RunFunc(cancleCtx)
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		ret.Err = ErrTaskTimeout
	case <-done:
	}
	return ret
}

func (p *GoPool) Do(ctx context.Context, tasks ...*Task) <-chan *RunResult {
	var graphs []*TaskGraph
	var length int
	for _, task := range tasks {
		//todo ring detect
		graph := task.ToGraph()
		graphs = append(graphs, graph)
		length += len(graph.nodes)
	}
	res := make(chan *RunResult, length)
	var wg sync.WaitGroup
	wg.Add(length)
	var recursive func(graph *TaskGraph, t *Task)
	recursive = func(graph *TaskGraph, t *Task) {
		go func() {
			result := p.do(ctx, t)
			res <- result
			wg.Done()
			if result.Err != nil || result.Buf != "" {
				t.state = 2
			}
			for _, child := range t.Children {
				graph.indgrees[child.Identifier]--
				if graph.indgrees[child.Identifier] <= 0 {
					if p.TaskOptions.ZombieMode && t.state == 2 {
						child.state = 2
					}
					recursive(graph, child)
				}
			}
		}()
	}
	for i, task := range tasks {
		recursive(graphs[i], task)
	}
	go func() {
		wg.Wait()
		close(res)
	}()
	return res
}

var DefaultGoPool = NewGoPool(100000)

func Go(ctx context.Context, f func()) {
	DefaultGoPool.Go(ctx, f)
}

func Do(ctx context.Context, tasks ...*Task) <-chan *RunResult {
	return DefaultGoPool.Do(ctx, tasks...)
}
