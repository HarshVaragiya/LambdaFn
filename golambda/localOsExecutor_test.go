package golambda

import (
	"strings"
	"testing"
	"time"
)

func TestLocalOsExecutor_runLocalBinary(t *testing.T) {
	t.Run("LocalCodeExecution", func(t *testing.T) {
		executor := BasicCodeExecutor{codeUri: "whoami", functionHandler: "--help", functionTimeout: time.Millisecond * 100}
		exec := LocalOsExecutor{codeExecutor: executor}
		_, stderr, _, err := exec.runLocalBinary("")
		if err != nil {
			t.Errorf("runLocalBinary() error = %v", err)
			t.Errorf(stderr)
		}
	})
	t.Run("LocalCodeExecutionTimeout", func(t *testing.T) {
		executor := BasicCodeExecutor{codeUri: "sleep", functionHandler: "1", functionTimeout: time.Millisecond * 100}
		exec := LocalOsExecutor{codeExecutor: executor}
		_, stderr,_, err := exec.runLocalBinary("")
		if err == nil {
			t.Error("runLocalBinary() Expected error.")
		}
		if !strings.Contains(err.Error(), "signal: killed") {
			t.Errorf("Timeout error expected. unknown error encountered. error = %v", err)
			t.Errorf("stderr = %v ", stderr)
		}
	})
}
