package golambda

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestContainerStartStop(t *testing.T) {
	t.Run("ContainerTest", func(t *testing.T) {
		builder := DockerContainerManager{}
		err := builder.Init()
		if err != nil {
			t.Fatalf("error initializing container manager. error = %v", err)
		}
		ctx := context.Background()
		srcPath, _ := os.MkdirTemp(os.TempDir(), "lambdaFn-test")
		defer os.RemoveAll(srcPath)
		env := make(map[string]string)
		containerId, err := builder.startContainer(ctx, "amazonlinux-python", "9000", srcPath, "/lambda/archive/", env)
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		log.Printf("started container with id [%s]", containerId)
		time.Sleep(time.Second * 2)
		log.Println("stopping the container.")
		if err = builder.stopContainer(ctx, containerId); err != nil {
			t.Errorf("error stopping the container.")
		}

	})
}

func TestPrepareEnvironmentVariables(t *testing.T) {
	t.Run("PrepareEnvTest", func(t *testing.T) {
		environ := map[string]string{}
		functionName := "secret-lambda-testing-function"
		testHandler := "example.lambda_handler"
		prepareEnvironmentVariables(functionName, testHandler, environ)
		name := environ["LAMBDA_FUNCTION_NAME"]
		handler := environ["LAMBDA_HANDLER_FUNCTION"]
		if name != functionName || testHandler != handler {
			t.Error("unexpected environment variables.")
		}
	})
}
