# execloop

Simple task execution engine.

## Installation

To install `execloop` simply run `go get -u github.com/kouzant/execloop`

## Usage

`execloop` is a very simple and generic execution engine with retryable tasks.
In the very basics it's an oversimplified control loop constantly trying to reach
a final state by executing a `Plan` consisting of some `Tasks`.

### Task

A `Task` is an isolated unit of work. A Task can have sub-tasks which will be
executed after the parent task has finished successfully. The main workload is
executed in the `PerformAction` function surrounded by a `Pre` and `Post`
action.

Your `Task` implementation should comply with the following interface:

    type Task interface {
	    Pre() error
	    PerformAction() ([]Task, error)
	    Post() error
	    Name() string
    }

### Plan

One or more `Tasks` form a `Plan` and this is what is going to be executed
by the scheduler. In every loop, the plan will be asked to return a set of
tasks to be executed in order to reach a final state. **Only** when the plan
will return an **empty** set of tasks the execution will stop.

Each time the scheduler will invoke the `Create` function of a plan to get
a set of tasks to execute. Eventually the `Create` function should return an
empy slice of `Tasks` meaning the final state has been reached.

A `Plan` should implement the follwing interface:

    type Plan interface {
	    Create() ([]Task, error)
    }

### Executor

Finally a plan will be executed by the scheduler. You can invoke the `Run(plan Plan) error`
or the `RunWithContext(ctx context.Context, plan Plan) error` function to
execute a plan. The latter will timeout after a configurable period of time.

Call `executor.New(options *execloop.Options) *Executor` to create a new
scheduler. The `Options` are the following:

    type Options struct {
	    Logger           Logger
	    SleepBetweenRuns time.Duration
	    ErrorsToTolerate int
	    ExecutionTimeout time.Duration
    }

Use the `With*` functions to override the default options obtained by `execloop.DefaultOptions()`

## Development

`make` to build and test

`make test` to run the tests

`test-no-cache` to run all tests with no cache

`make build` to build the library