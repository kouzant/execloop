package executor

import (
	"errors"
	"fmt"
	"time"

	"github.com/kouzant/execloop"
)

var fatalError *FatalError

type Executor struct {
	options        *execloop.Options
	numberOfErrors int
}

func New(options *execloop.Options) *Executor {
	return &Executor{
		options:        options,
		numberOfErrors: 0,
	}
}

func (e *Executor) Run(plan Plan) error {
	for {
		tasks, err := plan.Create()
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			e.options.Infof("No more tasks to execute\n")
			return nil
		}
		e.options.Debugf("Tasks remaining: %d\n", len(tasks))
		err = e.execute(tasks)
		if err != nil {
			e.options.Errorf("%s. Reason: %s", err, errors.Unwrap(err))
			return err
		}

		time.Sleep(e.options.SleepBetweenRuns)
	}
}

func (e *Executor) execute(tasks []Task) error {
	for _, task := range tasks {
		e.options.Infof("Executing Task: %s\n", task.Name())
		e.options.Debugf("Executing Pre of Task: %s\n", task.Name())
		prerr := task.Pre()
		ferr := e.handleTaskError(prerr)
		if ferr != nil {
			return ferr
		}

		if prerr == nil {
			e.options.Debugf("Executing PerfomAction of Task: %s\n", task.Name())
			paerr := task.PerformAction()
			ferr = e.handleTaskError(paerr)
			if ferr != nil {
				return ferr
			}

			if paerr == nil {
				e.options.Debugf("Executing Post of Task: %s\n", task.Name())
				poerr := task.Post()
				ferr = e.handleTaskError(poerr)
				if ferr != nil {
					return ferr
				}
				e.options.Infof("Finished executing Task: %s\n", task.Name())
			}
		}
	}
	return nil
}

func (e *Executor) handleTaskError(err error) error {
	if err == nil {
		return nil
	}
	e.options.Warningf("%s\n", err)
	e.numberOfErrors++
	if e.numberOfErrors > e.options.ErrorsToTolerate {
		return &FatalError{fmt.Sprintf("Reached maximum number of errors to tolerate %d", e.options.ErrorsToTolerate),
			err}
	}
	if errors.As(err, &fatalError) {
		return err
	}
	return nil
}
