package golambda

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
	"time"
)

func TestBuildLambdaImage(t *testing.T) {
	t.Run("ContainerImageBuilderTest", func(t *testing.T) {
		log.SetLevel(logrus.TraceLevel)
		environ := map[string]string{}
		codeUri := "images/python/test/test-example.zip"
		testHandler := "example.lambda_handler"
		log.Println(os.Getwd())
		builder := DockerImageBuilder{}
		err := builder.Init()
		if err != nil{
			t.Error(err)
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute * 5)
		defer cancelFunc()
		err = builder.BuildLambdaImage(ctx, "container_test_lambda",codeUri,"python",testHandler, environ)
		if err != nil {
			t.Errorf("error building lambda docker image. error = %v", err)
		}
	})
}

func TestPrepareEnvironmentVariables(t *testing.T) {
	t.Run("PrepareEnvTest", func(t *testing.T) {
		environ := map[string]string{}
		functionName := "secret-lambda-testing-function"
		testHandler := "example.lambda_handler"
		prepareEnvironmentVariables(functionName,testHandler,environ)
		name := environ["LAMBDA_FUNCTION_NAME"]
		handler := environ["LAMBDA_HANDLER_FUNCTION"]
		if name != functionName || testHandler != handler {
			t.Error("unexpected environment variables.")
		}
	})
}

func TestEnvironmentVariablesToString(t *testing.T) {
	t.Run("EnvToStringTest", func(t *testing.T) {
		environ := map[string]string{
			"LAMBDA_FUNCTION_NAME": "secret-lambda-testing-function",
			"LAMBDA_HANDLER_FUNCTION": "example.lambda_handler",
		}
		envString := environmentVariablesToString(environ)
		t1 := strings.Contains(envString, "ENV LAMBDA_FUNCTION_NAME=secret-lambda-testing-function")
		t2 := strings.Contains(envString, "ENV LAMBDA_HANDLER_FUNCTION=example.lambda_handler")
		if !t1 || !t2 {
			t.Error("envString did not contain required environment variables. ")
		}

	})
}