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

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockLogger struct {
	output string
}

func (m *mockLogger) Debugf(f string, v ...interface{}) {
	m.output = fmt.Sprintf("DEBUG: "+f, v...)
}

func (m *mockLogger) Infof(f string, v ...interface{}) {
	m.output = fmt.Sprintf("INFO: "+f, v...)
}

func (m *mockLogger) Warningf(f string, v ...interface{}) {
	m.output = fmt.Sprintf("WARN: "+f, v...)
}

func (m *mockLogger) Errorf(f string, v ...interface{}) {
	m.output = fmt.Sprintf("ERROR: "+f, v...)
}

func TestLog(t *testing.T) {
	logger := &mockLogger{}
	opts := Options{}.WithLogger(logger)

	opts.Debugf("test")
	require.Equal(t, "DEBUG: test", logger.output)

	opts.Infof("test")
	require.Equal(t, "INFO: test", logger.output)

	opts.Warningf("test")
	require.Equal(t, "WARN: test", logger.output)

	opts.Errorf("test")
	require.Equal(t, "ERROR: test", logger.output)
}
