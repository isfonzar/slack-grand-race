package logs

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	var tests = []struct {
		isDebug bool
	}{
		{true},
		{false},
	}

	for _, test := range tests {
		_, err := New(test.isDebug)
		if err != nil {
			t.Errorf("could not initialize logger with debug = %t", test.isDebug)
		}
	}
}
