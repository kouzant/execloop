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

func (e *Executor) RunWithContext(ctx context.Context, plan Plan) error {
	e.options.Debugf("Running with context")
	execCtx, cancel := context.WithTimeout(ctx, e.options.ExecutionTimeout)
	defer cancel()

	controlChannel := make(chan error)
	go func() {
		controlChannel <- e.run(plan)
	}()

	select {
	case <-execCtx.Done():
		return execCtx.Err()
	case controlResponse := <-controlChannel:
		return controlResponse
	}
}

func (e *Executor) Run(plan Plan) error {
	e.options.Debugf("Running without context")
	return e.run(plan)
}

func (e *Executor) run(plan Plan) error {
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
			childrenTasks, paerr := task.PerformAction()
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
				if poerr == nil && childrenTasks != nil && len(childrenTasks) > 0 {
					e.options.Debugf("Executig children tasks of %s\n", task.Name())
					inerr := e.execute(childrenTasks)
					if ferr = e.handleTaskError(inerr); ferr != nil {
						return ferr
					}
				}
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
