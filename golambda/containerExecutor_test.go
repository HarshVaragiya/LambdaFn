package golambda

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestContainerExecutor_execute(t *testing.T) {
	t.Run("ContainerExecutorDockerIntegrationTest", func(t *testing.T) {
		codeUri := "../images/python/test/test-example.zip"
		log.SetLevel(logrus.DebugLevel)
		lambda := NewContainerLambdaFunction("test-function", codeUri, "python3", "example.lambda_handler")
		event := NewSimpleLambdaEvent("Hi there!")
		resp, err := lambda.Invoke(event)
		if err != nil {
			t.Fatalf("error invoking function. error = %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("response not as expected. ")
		}
		log.Printf("response : %v", resp)
		log.Printf("response Data : %s", resp.Data)
		log.Printf("response EventId : %s", resp.EventId)
		log.Printf("response Stderr : %s", resp.Stderr)
	})
}
