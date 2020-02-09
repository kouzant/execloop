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
