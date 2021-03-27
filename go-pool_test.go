package gtask

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestZombie(t *testing.T) {
	pool := NewGoPool(100000)
	pool.TaskOptions.ZombieMode = true
	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			fmt.Println("i am task A")
			panic(1111)
			return nil
		},
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(3 * time.Second)
			fmt.Println("i am task B")
			return nil
		},
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i am task C")
			return nil
		},
		Identifier: "TaskC",
	}
	TaskA.AddChildren(TaskB)

	results := pool.Do(context.Background(), TaskA, TaskC)
	for res := range results {
		fmt.Printf("task done for %s\n", res.Identifier)
		if res.Err != nil {
			fmt.Printf("err from %s: %s\n", res.Identifier, res.Err)
		}
		if res.Buf != "" {
			fmt.Printf("panic from %s: %s\n", res.Identifier, res.Buf)
		}
	}
}

func TestDo(t *testing.T) {
	pool := NewGoPool(100000)
	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i am task A")
			return nil
		},
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(3 * time.Second)
			fmt.Println("i am task B")
			return nil
		},
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i am task C")
			return nil
		},
		Identifier: "TaskC",
	}
	TaskA.AddChildren(TaskB)

	results := pool.Do(context.Background(), TaskA, TaskC)
	for res := range results {
		fmt.Printf("task done for %s\n", res.Identifier)
	}
}

func TestDoParallel(t *testing.T) {
	pool := NewGoPool(100000)
	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task A")
			return nil
		},
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(3 * time.Second)
			fmt.Println("i m task B")
			return nil
		},
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task C")
			return nil
		},
		Identifier: "TaskC",
	}

	pool.Do(context.Background(), TaskA)
	pool.Do(context.Background(), TaskB)
	pool.Do(context.Background(), TaskC)
	time.Sleep(10 * time.Second)
}

func TestDoGraph(t *testing.T) {
	pool := NewGoPool(100000)

	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task A")
			return nil
		},
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task B")
			return nil
		},
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task C")
			return nil
		},
		Identifier: "TaskC",
	}
	TaskD := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task D")
			return nil
		},
		Identifier: "TaskD",
	}
	TaskE := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task E")
			return nil
		},
		Identifier: "TaskE",
	}
	TaskF := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task F")
			return nil
		},
		Identifier: "TaskF",
	}
	TaskG := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task G")
			return nil
		},
		Identifier: "TaskG",
	}
	TaskH := &Task{
		RunFunc: func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("i m task H")
			return nil
		},
		Identifier: "TaskH",
	}

	TaskA.AddChildren(TaskB, TaskC, TaskD)
	TaskE.Depends(TaskB, TaskC, TaskD)
	TaskE.AddChildren(TaskF)
	TaskF.AddChildren(TaskG)
	TaskG.AddChildren(TaskH)
	for res := range pool.Do(context.Background(), TaskA) {
		fmt.Println(res)
	}
}

func TestRange(t *testing.T) {
	pool := NewGoPool(100000)
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
	}
	var tasks []*Task
	for _, v := range m {
		vv := v
		tasks = append(tasks, &Task{
			RunFunc: func(ctx context.Context) error {
				fmt.Println(vv)
				return nil
			},
		})
	}
	for range pool.Do(context.Background(), tasks...) {
	}
}

func TestTimeout(t *testing.T) {
	pool := NewGoPool(100000)
	var ErrCtxDone = errors.New("context canceled")
	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ErrCtxDone
				default:
					time.Sleep(time.Second)
					fmt.Println("i m task A")
				}
			}
		},
		Timeout:    3 * time.Second,
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ErrCtxDone
				default:
					time.Sleep(time.Second)
					fmt.Println("i m task B")
				}
			}
		},
		Timeout:    3 * time.Second,
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ErrCtxDone
				default:
					time.Sleep(time.Second)
					fmt.Println("i m task C")
				}
			}
		},
		Timeout:    3 * time.Second,
		Identifier: "TaskC",
	}
	TaskA.AddChildren(TaskB)
	TaskB.AddChildren(TaskC)
	results := pool.Do(context.Background(), TaskA)
	for r := range results {
		fmt.Printf("%+v\n", r)
	}

	fmt.Println("==========test zombie================")
	pool.TaskOptions.ZombieMode = true
	results1 := pool.Do(context.Background(), TaskA)
	for r := range results1 {
		fmt.Printf("%+v\n", r)
	}
}

func TestCancel(t *testing.T) {
	pool := NewGoPool(100000)
	var ErrCtxDone = errors.New("context canceled")
	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ErrCtxDone
				default:
					time.Sleep(time.Second)
					fmt.Println("i m task A")
				}
			}
		},
		Timeout:    2 * time.Second,
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ErrCtxDone
				default:
					time.Sleep(time.Second)
					fmt.Println("i m task B")
				}
			}
		},
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ErrCtxDone
				default:
					time.Sleep(time.Second)
					fmt.Println("i m task C")
				}
			}
		},
		Identifier: "TaskC",
	}
	TaskA.AddChildren(TaskB)
	TaskB.AddChildren(TaskC)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	results := pool.Do(ctx, TaskA)
	for r := range results {
		fmt.Printf("%+v\n", r)
	}
}

func TestFastFail(t *testing.T) {
	pool := NewGoPool(2)
	pool.FastFail = true
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
	}
	var tasks []*Task
	for _, v := range m {
		vv := v
		tasks = append(tasks, &Task{
			RunFunc: func(ctx context.Context) error {
				time.Sleep(1 * time.Second)
				fmt.Println(vv)
				return nil
			},
		})
	}
	for re := range pool.Do(context.Background(), tasks...) {
		fmt.Printf("%+v\n", re)
	}
}

func TestClosure(t *testing.T) {
	pool := NewGoPool(100000)
	var i int
	TaskA := &Task{
		RunFunc: func(ctx context.Context) error {
			i++
			fmt.Println(i)
			return nil
		},
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc: func(ctx context.Context) error {
			i++
			fmt.Println(i)
			return nil
		},
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc: func(ctx context.Context) error {
			i++
			fmt.Println(i)
			return nil
		},
		Identifier: "TaskC",
	}
	TaskA.AddChildren(TaskB)
	TaskB.AddChildren(TaskC)

	results := pool.Do(context.Background(), TaskA)
	for res := range results {
		fmt.Printf("task done for %s\n", res.Identifier)
	}
}
