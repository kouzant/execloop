package executor

import (
	"errors"
	"testing"
	"time"

	"github.com/kouzant/execloop"
	"github.com/stretchr/testify/require"
)

type DummyTask struct {
	taskName    string
	tasksLog    [][]int
	shouldIFail bool
}

func (d *DummyTask) Pre() error {
	return nil
}

func (d *DummyTask) PerformAction() error {
	if d.taskName == "DummyTask0" {
		d.tasksLog[0][0] += 1
	} else if d.taskName == "DummyTask1" {
		d.tasksLog[1][0] += 1
	} else if d.taskName == "DummyTask2" {
		d.tasksLog[2][0] += 1
	}

	if d.shouldIFail {
		return &FatalError{"Something terrible has happened", nil}
	}

	if d.taskName == "DummyTask0" {
		d.tasksLog[0][1] += 1
	} else if d.taskName == "DummyTask1" {
		d.tasksLog[1][1] += 1
	} else if d.taskName == "DummyTask2" {
		d.tasksLog[2][1] += 1
	}
	return nil
}

func (d *DummyTask) Post() error {
	return nil
}

func (d *DummyTask) Name() string {
	return d.taskName
}

type SuccessPlan struct {
	tasksLog [][]int
}

func (p *SuccessPlan) Create() ([]Task, error) {
	var tasks []Task
	if p.tasksLog[0][1] == 0 {
		t := &DummyTask{
			taskName: "DummyTask0",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	if p.tasksLog[1][1] == 0 {
		t := &DummyTask{
			taskName: "DummyTask1",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	if p.tasksLog[2][1] == 0 {
		t := &DummyTask{
			taskName: "DummyTask2",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

type FailFatalPlan struct {
	tasksLog [][]int
}

func (p *FailFatalPlan) Create() ([]Task, error) {
	var tasks []Task
	if p.tasksLog[0][1] == 0 {
		t := &DummyTask{
			taskName: "DummyTask0",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}

	// Second task should fail with Fatal error
	if p.tasksLog[1][1] == 0 {
		t := &DummyTask{
			taskName:    "DummyTask1",
			tasksLog:    p.tasksLog,
			shouldIFail: true,
		}
		tasks = append(tasks, t)
	}
	if p.tasksLog[2][1] == 0 {
		t := &DummyTask{
			taskName: "DummyTask2",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func TestSuccessExecutor(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0},
		{0, 0},
		{0, 0},
	}
	plan := &SuccessPlan{tasksLog: tasksLog}
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond)
	exec := New(&opts)
	err := exec.Run(plan)
	if err != nil {
		t.Errorf("Did not expect error but gotten :%v\n", err)
	}

	// All tasks should have exactly one attempt
	require.Equal(t, 1, plan.tasksLog[0][0])
	require.Equal(t, 1, plan.tasksLog[1][0])
	require.Equal(t, 1, plan.tasksLog[2][0])

	// All tasks should have run exactly once
	require.Equal(t, 1, plan.tasksLog[0][1])
	require.Equal(t, 1, plan.tasksLog[1][1])
	require.Equal(t, 1, plan.tasksLog[2][1])
}

func TestFatalTask(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0},
		{0, 0},
		{0, 0},
	}
	plan := &FailFatalPlan{tasksLog: tasksLog}
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond)
	exec := New(&opts)
	err := exec.Run(plan)
	if err != nil {
		if !errors.As(err, &fatalError) {
			t.Errorf("Expected FatalError, instead got %v\n", err)
		}
	} else {
		t.Errorf("Expected error")
	}

	// All BUT the last task should have exactly one attempt
	require.Equal(t, 1, plan.tasksLog[0][0])
	require.Equal(t, 1, plan.tasksLog[1][0])
	require.Equal(t, 0, plan.tasksLog[2][0])

	// ONLY the first task should have succeeded
	require.Equal(t, 1, plan.tasksLog[0][1])
	require.Equal(t, 0, plan.tasksLog[1][1])
	require.Equal(t, 0, plan.tasksLog[2][1])
}
