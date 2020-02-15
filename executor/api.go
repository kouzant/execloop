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

import "fmt"

type Task interface {
	Pre() error
	PerformAction() ([]Task, error)
	Post() error
	Name() string
}

type Plan interface {
	Create() ([]Task, error)
}

type FatalError struct {
	msg string
	err error
}

func (e *FatalError) Error() string {
	return fmt.Sprintf("FatalError: %s", e.msg)
}

func (e *FatalError) Unwrap() error {
	return e.err
}
