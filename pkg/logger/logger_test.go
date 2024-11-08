package logger

import (
	"os"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestCreateLogger(t *testing.T) {
	logFile := "./testlog.log"

	testCases := []struct {
		desc  string
		level Level
	}{
		{
			desc:  "new DEVELOPMENT level logger",
			level: DEVELOPMENT,
		},
		{
			desc:  "new PRODUCTION level logger",
			level: PRODUCTION,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := New(tC.level, logFile)
			if err != nil {
				t.Errorf("creating logger instance %s level with error: %s", tC.level, err)
			}
			os.Remove(logFile)
		})
	}

}
