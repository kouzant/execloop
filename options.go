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

package execloop

import "time"

type Options struct {
	Logger           Logger
	SleepBetweenRuns time.Duration
	ErrorsToTolerate int
	ExecutionTimeout time.Duration
}

func DefaultOptions() Options {
	return Options{
		Logger:           defaultLog,
		SleepBetweenRuns: time.Second,
		ErrorsToTolerate: 5,
		ExecutionTimeout: 20 * time.Minute,
	}
}

func (o Options) WithLogger(logger Logger) Options {
	o.Logger = logger
	return o
}

func (o Options) WithSleepBetweenRuns(sleep time.Duration) Options {
	o.SleepBetweenRuns = sleep
	return o
}

func (o Options) WithErrorsToTolerate(numOfErrors int) Options {
	o.ErrorsToTolerate = numOfErrors
	return o
}

func (o Options) WithExecutionTimeout(timeout time.Duration) Options {
	o.ExecutionTimeout = timeout
	return o
}
