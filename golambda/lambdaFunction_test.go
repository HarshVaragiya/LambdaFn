package golambda

import (
	"testing"
)

func TestFunction_Invoke(t *testing.T) {
	t.Run("LocalExecutionLambda", func(t *testing.T) {
		lambda := NewLocalOsLambdaFunction("user-test-lambda-function-1", "whoami", "--help")
		resp, err := lambda.Invoke(NewSimpleLambdaEvent(""))
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("expected statusCode to be 200")
			t.Error(resp.StatusCode)
		}
	})
}

func TestNewSimpleLocalFunctionExecutor(t *testing.T) {
	t.Run("LocalExecutionLambdaTimeout", func(t *testing.T) {
		lambda := NewLocalOsLambdaFunction("test-lambda-function-2", "sleep", "5")
		response, err := lambda.Invoke(NewSimpleLambdaEvent(""))
		if err != nil {
			t.Error(err)
		}
		if response.StatusCode != statusCodeMap["signal: killed"] {
			t.Errorf("Status code was not as expected. statusCode = %v", response.StatusCode)
			t.Errorf("stderr = %v", response.Stderr)
		}
	})
}
