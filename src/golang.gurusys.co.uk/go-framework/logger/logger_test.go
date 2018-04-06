package logger

import (
	"flag"
	"testing"
	"time"
)

func TestSendLog(t *testing.T) {

	flag.Parse()

	q, err := NewAsyncLogQueue("buildrepo", "?", "?", "?", "?")
	if err != nil {
		t.Error(err)
		return
	}

	err = q.LogCommandStdout("this is a test searchforme", "")
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(3 * time.Second)
}
