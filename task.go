package gtask

import (
	"context"
	"time"
)

type TaskOptions struct {
	ZombieMode bool
}

type TaskGraph struct {
	Head     *Task
	indgrees map[string]int
	nodes    map[string]struct{}
}

type Task struct {
	RunFunc    func(ctx context.Context) error
	Identifier string
	Timeout    time.Duration
	state      int //0: init 1: success 2: fail
	Children   []*Task
}

type RunResult struct {
	Identifier string
	Err        error
	Buf        string
}

func (t *Task) AddChildren(children ...*Task) {
	for _, child := range children {
		t.Children = append(t.Children, child)
	}
}

func (t *Task) Depends(parents ...*Task) {
	for _, parent := range parents {
		parent.AddChildren(t)
	}
}

func (t *Task) walk(walked map[string]struct{}, f func(*Task)) {
	if _, ok := walked[t.Identifier]; ok {
		return
	}
	f(t)
	walked[t.Identifier] = struct{}{}
	for _, child := range t.Children {
		child.walk(walked, f)
	}
}

func (t *Task) Walk(f func(*Task)) {
	walked := make(map[string]struct{})
	t.walk(walked, f)
}

func (t *Task) ToGraph() *TaskGraph {
	graph := &TaskGraph{
		Head:     t,
		indgrees: make(map[string]int),
		nodes:    make(map[string]struct{}),
	}
	t.Walk(func(task *Task) {
		graph.nodes[task.Identifier] = struct{}{}
		for i, _ := range task.Children {
			id := task.Children[i].Identifier
			graph.indgrees[id]++
		}
	})
	return graph
}
