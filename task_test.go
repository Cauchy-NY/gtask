package gtask

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTaskGraph(t *testing.T) {
	var runfunc = func(ctx context.Context) error {
		time.Sleep(time.Second)
		fmt.Println("i am a silly task")
		return nil
	}
	TaskA := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskA",
	}
	TaskB := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskB",
	}
	TaskC := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskC",
	}
	TaskD := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskD",
	}
	TaskE := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskE",
	}
	TaskF := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskF",
	}
	TaskG := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskG",
	}
	TaskH := &Task{
		RunFunc:    runfunc,
		Identifier: "TaskH",
	}

	TaskA.AddChildren(TaskB, TaskC, TaskD)
	TaskE.Depends(TaskB, TaskC, TaskD)
	TaskE.AddChildren(TaskF)
	TaskF.AddChildren(TaskG)
	TaskG.AddChildren(TaskH)

	graph := TaskA.ToGraph()
	t.Logf("%+v\n", graph)
}

func TestRange1(t *testing.T) {
	fetchedConnectionInfos := map[int]string{
		1: "1",
		3: "3",
		4: "4",
		6: "6",
	}
	uids := []int{1, 3, 4, 6}
	generalExtraMap := make(map[int]string)
	for _, uid := range uids {
		generalExtra, _ := fetchedConnectionInfos[uid]
		generalExtraMap[uid] = generalExtra + "extra"
	}
	fmt.Println(generalExtraMap)
}

func TestMarshalMapInt(t *testing.T) {
	m := make(map[interface{}]interface{})
	m[1] = "ddd"
	m["k2"] = "ddd"
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
