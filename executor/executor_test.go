package executor

import (
	"testing"
	"time"

	"github.com/kouzant/execloop"
	"github.com/stretchr/testify/require"
)

type SuccessTask struct {
	TaskName string
	taskLogs [][]int
}

func (d *SuccessTask) Pre() error {
	return nil
}

func (d *SuccessTask) PerformAction() error {
	if d.TaskName == "DummyTask0" {
		d.taskLogs[0][0] += 1
		d.taskLogs[0][1] += 1
	} else if d.TaskName == "DummyTask1" {
		d.taskLogs[1][0] += 1
		d.taskLogs[1][1] += 1
	} else if d.TaskName == "DummyTask2" {
		d.taskLogs[2][0] += 1
		d.taskLogs[2][1] += 1
	}
	return nil
}

func (d *SuccessTask) Post() error {
	return nil
}

func (d *SuccessTask) Name() string {
	return d.TaskName
}

type SuccessPlan struct {
	taskLogs [][]int
}

func (p *SuccessPlan) Create() ([]Task, error) {
	var tasks []Task
	if p.taskLogs[0][1] == 0 {
		t := &SuccessTask{"DummyTask0", p.taskLogs}
		tasks = append(tasks, t)
	}
	if p.taskLogs[1][1] == 0 {
		t := &SuccessTask{"DummyTask1", p.taskLogs}
		tasks = append(tasks, t)
	}
	if p.taskLogs[2][1] == 0 {
		t := &SuccessTask{"DummyTask2", p.taskLogs}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
func TestExecutor(t *testing.T) {
	var taskLogs = [][]int{
		{0, 0},
		{0, 0},
		{0, 0},
	}
	plan := &SuccessPlan{taskLogs: taskLogs}
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond)
	exec := New(&opts)
	err := exec.Run(plan)
	if err != nil {
		t.Errorf("Did not expect error but gotten :%v\n", err)
	}

	// All tasks should have exactly one attempt
	require.Equal(t, 1, plan.taskLogs[0][0])
	require.Equal(t, 1, plan.taskLogs[1][0])
	require.Equal(t, 1, plan.taskLogs[2][0])

	// All tasks should have run exactly once
	require.Equal(t, 1, plan.taskLogs[0][1])
	require.Equal(t, 1, plan.taskLogs[1][1])
	require.Equal(t, 1, plan.taskLogs[2][1])
}
