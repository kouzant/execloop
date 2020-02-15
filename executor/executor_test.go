/*
This file is part of execloop.

execloop is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

execloop is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with execloop.  If not, see <https://www.gnu.org/licenses/>.
*/
package executor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kouzant/execloop"
	"github.com/stretchr/testify/require"
)

type Action int

const (
	Pre Action = iota
	PerformAction
	Post
)

const Attempt = Post + 1

type ParentTask struct {
	DummyTask
}

func (p *ParentTask) Pre() error {
	p.tasksLog[0][Pre] += 1
	return nil
}

func (p *ParentTask) PerformAction() ([]Task, error) {
	p.tasksLog[0][Attempt] += 1
	p.tasksLog[0][PerformAction] += 1
	var kids []Task
	kid0 := &KidTask{DummyTask{taskName: "KidTask0", tasksLog: p.tasksLog}}
	kids = append(kids, kid0)
	kid1 := &KidTask{DummyTask{taskName: "KidTask1", tasksLog: p.tasksLog}}
	kids = append(kids, kid1)

	return kids, nil
}

func (p *ParentTask) Post() error {
	p.tasksLog[0][Post] += 1
	return nil
}

type KidTask struct {
	DummyTask
}

func (d *KidTask) Pre() error {
	d.updateCounter(Pre)
	return nil
}

func (d *KidTask) PerformAction() ([]Task, error) {
	if d.taskName == "KidTask0" {
		d.tasksLog[1][Attempt] += 1
	} else if d.taskName == "KidTask1" {
		d.tasksLog[2][Attempt] += 1
	}
	d.updateCounter(PerformAction)
	var kids []Task
	return kids, nil
}

func (d *KidTask) Post() error {
	d.updateCounter(Post)
	return nil
}

func (d *KidTask) updateCounter(action Action) {
	if d.taskName == "KidTask0" {
		d.tasksLog[1][action] += 1
	} else if d.taskName == "KidTask1" {
		d.tasksLog[2][action] += 1
	}
}

type DummyTask struct {
	taskName    string
	tasksLog    [][]int
	shouldIFail int
}

func (d *DummyTask) Pre() error {
	d.updateCounter(Pre)
	return nil
}

func (d *DummyTask) PerformAction() ([]Task, error) {
	if d.taskName == "DummyTask0" {
		d.tasksLog[0][Attempt] += 1
	} else if d.taskName == "DummyTask1" {
		d.tasksLog[1][Attempt] += 1
	} else if d.taskName == "DummyTask2" {
		d.tasksLog[2][Attempt] += 1
	}

	if d.shouldIFail == 1 {
		return nil, &FatalError{"Something terrible has happened", nil}
	}

	if d.shouldIFail == 2 {
		return nil, errors.New("A small tiny error")
	}

	d.updateCounter(PerformAction)
	var children []Task
	return children, nil
}

func (d *DummyTask) Post() error {
	d.updateCounter(Post)
	return nil
}

func (d *DummyTask) Name() string {
	return d.taskName
}

func (d *DummyTask) updateCounter(action Action) {
	if d.taskName == "DummyTask0" {
		d.tasksLog[0][action] += 1
	} else if d.taskName == "DummyTask1" {
		d.tasksLog[1][action] += 1
	} else if d.taskName == "DummyTask2" {
		d.tasksLog[2][action] += 1
	}
}

type SuccessPlan struct {
	tasksLog [][]int
}

func (p *SuccessPlan) Create() ([]Task, error) {
	var tasks []Task
	if p.tasksLog[0][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask0",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	if p.tasksLog[1][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask1",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	if p.tasksLog[2][PerformAction] == 0 {
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
	if p.tasksLog[0][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask0",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}

	// Second task should fail with Fatal error
	if p.tasksLog[1][PerformAction] == 0 {
		t := &DummyTask{
			taskName:    "DummyTask1",
			tasksLog:    p.tasksLog,
			shouldIFail: 1,
		}
		tasks = append(tasks, t)
	}
	if p.tasksLog[2][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask2",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

type FailOneTaskPlan struct {
	tasksLog     [][]int
	succeedAfter int
	failures     int
}

func (p *FailOneTaskPlan) Create() ([]Task, error) {
	var tasks []Task
	if p.tasksLog[0][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask0",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}

	// Second task should be failing but not fatal for succeedAfter times
	if p.tasksLog[1][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask1",
			tasksLog: p.tasksLog,
		}
		if p.failures <= p.succeedAfter {
			t.shouldIFail = 2
			p.failures++
		}
		tasks = append(tasks, t)
	}

	if p.tasksLog[2][PerformAction] == 0 {
		t := &DummyTask{
			taskName: "DummyTask2",
			tasksLog: p.tasksLog,
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

type ChildrenPlan struct {
	tasksLog [][]int
}

func (p *ChildrenPlan) Create() ([]Task, error) {
	var tasks []Task
	if p.tasksLog[0][PerformAction] == 0 {
		t := &ParentTask{
			DummyTask{taskName: "ParentTask0",
				tasksLog: p.tasksLog}}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

type SleepyTask struct {
	Sleep time.Duration
	DummyTask
}

func (s *SleepyTask) PerformAction() ([]Task, error) {
	time.Sleep(s.Sleep)
	var tasks []Task
	return tasks, nil
}

type SleepyPlan struct {
	TaskSleep time.Duration
}

func (s *SleepyPlan) Create() ([]Task, error) {
	var tasks []Task
	// We should always check if task has been executed
	// and skip adding it to the tasks.
	// In this case it's OK as the executor will timeout
	t := &SleepyTask{Sleep: s.TaskSleep}
	tasks = append(tasks, t)
	return tasks, nil
}

func TestSuccessExecutor(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	plan := &SuccessPlan{tasksLog: tasksLog}
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond)
	exec := New(&opts)
	err := exec.RunWithContext(context.Background(), plan)
	if err != nil {
		t.Errorf("Did not expect error but gotten :%v\n", err)
	}

	// All tasks should have exactly one attempt
	require.Equal(t, 1, plan.tasksLog[0][Attempt])
	require.Equal(t, 1, plan.tasksLog[1][Attempt])
	require.Equal(t, 1, plan.tasksLog[2][Attempt])

	// All tasks should have exactly one Pre, PerformAction and Post
	require.Equal(t, 1, plan.tasksLog[0][Pre])
	require.Equal(t, 1, plan.tasksLog[1][Pre])
	require.Equal(t, 1, plan.tasksLog[2][Pre])

	require.Equal(t, 1, plan.tasksLog[0][PerformAction])
	require.Equal(t, 1, plan.tasksLog[1][PerformAction])
	require.Equal(t, 1, plan.tasksLog[2][PerformAction])

	require.Equal(t, 1, plan.tasksLog[0][Post])
	require.Equal(t, 1, plan.tasksLog[1][Post])
	require.Equal(t, 1, plan.tasksLog[2][Post])
}

func TestFatalTask(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
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
	require.Equal(t, 1, plan.tasksLog[0][Attempt])
	require.Equal(t, 1, plan.tasksLog[1][Attempt])
	require.Equal(t, 0, plan.tasksLog[2][Attempt])

	// Pre should have run only on the first two
	require.Equal(t, 1, plan.tasksLog[0][Pre])
	require.Equal(t, 1, plan.tasksLog[1][Pre])
	require.Equal(t, 0, plan.tasksLog[2][Pre])

	// PerformAction should have run ONLY on the first one
	require.Equal(t, 1, plan.tasksLog[0][PerformAction])
	require.Equal(t, 0, plan.tasksLog[1][PerformAction])
	require.Equal(t, 0, plan.tasksLog[2][PerformAction])

	// Post should have run only on the first one
	require.Equal(t, 1, plan.tasksLog[0][Post])
	require.Equal(t, 0, plan.tasksLog[1][Post])
	require.Equal(t, 0, plan.tasksLog[2][Post])
}

func TestMultipleAttemptsPlan(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	plan := &FailOneTaskPlan{tasksLog: tasksLog, succeedAfter: 3}
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond)
	exec := New(&opts)
	err := exec.Run(plan)
	if err != nil {
		t.Errorf("Did not expect error but gotten %v\n", err)
	}

	// All BUT the second task should have exactly one attempt
	// Second attempt should have 5 attempts
	require.Equal(t, 1, plan.tasksLog[0][Attempt])
	require.Equal(t, 5, plan.tasksLog[1][Attempt])
	require.Equal(t, 1, plan.tasksLog[2][Attempt])

	// First and third task should have succeded Pre once
	// Second should be 5
	require.Equal(t, 1, plan.tasksLog[0][Pre])
	require.Equal(t, 5, plan.tasksLog[1][Pre])
	require.Equal(t, 1, plan.tasksLog[2][Pre])

	// PerformAction should have succeeded only once
	require.Equal(t, 1, plan.tasksLog[0][PerformAction])
	require.Equal(t, 1, plan.tasksLog[1][PerformAction])
	require.Equal(t, 1, plan.tasksLog[2][PerformAction])

	// Post should have been executed only once for every task
	require.Equal(t, 1, plan.tasksLog[0][Post])
	require.Equal(t, 1, plan.tasksLog[1][Post])
	require.Equal(t, 1, plan.tasksLog[2][Post])

	require.Equal(t, 4, plan.failures)
}

func TestFailMultipleTimesGivingUpPlan(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	plan := &FailOneTaskPlan{tasksLog: tasksLog, succeedAfter: 100}
	var errors2tolerate = 3
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond).WithErrorsToTolerate(errors2tolerate)
	exec := New(&opts)
	err := exec.Run(plan)
	if err != nil {
		if !errors.As(err, &fatalError) {
			t.Errorf("Expected FatalError, instead got %v\n", err)
		}
	} else {
		t.Errorf("Expected error")
	}

	// All BUT the second task should have exactly one attempt
	// Second task should have MAX ERRORS TO TOLERATE + 1
	require.Equal(t, 1, plan.tasksLog[0][Attempt])
	require.Equal(t, errors2tolerate+1, plan.tasksLog[1][Attempt])
	require.Equal(t, 1, plan.tasksLog[2][Attempt])

	// First and third task should have succeded Pre once
	// Second should be MAX ERRORS TO TOLERATE + 1
	require.Equal(t, 1, plan.tasksLog[0][Pre])
	require.Equal(t, errors2tolerate+1, plan.tasksLog[1][Pre])
	require.Equal(t, 1, plan.tasksLog[2][Pre])

	// PerformAction should have succeeded only once
	// except for the second
	require.Equal(t, 1, plan.tasksLog[0][PerformAction])
	require.Equal(t, 0, plan.tasksLog[1][PerformAction])
	require.Equal(t, 1, plan.tasksLog[2][PerformAction])

	// Post should have succedded only once
	// except for the second
	require.Equal(t, 1, plan.tasksLog[0][Post])
	require.Equal(t, 0, plan.tasksLog[1][Post])
	require.Equal(t, 1, plan.tasksLog[2][Post])

	require.Equal(t, 4, plan.failures)
}

func TestChildrenPlan(t *testing.T) {
	var tasksLog = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	plan := &ChildrenPlan{tasksLog: tasksLog}
	opts := execloop.DefaultOptions().WithSleepBetweenRuns(50 * time.Millisecond)
	exec := New(&opts)
	err := exec.Run(plan)
	if err != nil {
		t.Errorf("Did not expect error but gotten :%v\n", err)
	}

	// All tasks should have exactly one attempt
	require.Equal(t, 1, plan.tasksLog[0][Attempt])
	require.Equal(t, 1, plan.tasksLog[1][Attempt])
	require.Equal(t, 1, plan.tasksLog[2][Attempt])

	// All tasks should have exactly one Pre, PerformAction and Post
	require.Equal(t, 1, plan.tasksLog[0][Pre])
	require.Equal(t, 1, plan.tasksLog[1][Pre])
	require.Equal(t, 1, plan.tasksLog[2][Pre])

	require.Equal(t, 1, plan.tasksLog[0][PerformAction])
	require.Equal(t, 1, plan.tasksLog[1][PerformAction])
	require.Equal(t, 1, plan.tasksLog[2][PerformAction])

	require.Equal(t, 1, plan.tasksLog[0][Post])
	require.Equal(t, 1, plan.tasksLog[1][Post])
	require.Equal(t, 1, plan.tasksLog[2][Post])
}

func TestExecutorTimeout(t *testing.T) {
	plan := &SleepyPlan{1 * time.Second}
	opts := execloop.DefaultOptions().WithExecutionTimeout(700 * time.Millisecond)
	exec := New(&opts)
	err := exec.RunWithContext(context.Background(), plan)
	require.NotNil(t, err)
	require.Equal(t, err, context.DeadlineExceeded)
}
