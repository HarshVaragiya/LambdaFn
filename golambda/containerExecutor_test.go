package golambda

import (
	"fmt"
	"testing"
)

func TestContainerExecutor_execute(t *testing.T) {
	t.Run("ContainerExecutorTest", func(t *testing.T) {
		lambda := NewContainerLambdaFunction("my-awesome-function","","python3","example.lambda_handler")
		event := NewSimpleLambdaEvent("Hello from the integration test!")
		log.Printf("Sending EventId [%s]: %v", event.EventId, event.EventData)
		response, err := lambda.Invoke(event)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("--------- LAMBDA RESPONSE ---------")
		fmt.Print("Data : ",response.Data)
		fmt.Println("InvocationStatusCode : ", response.StatusCode)
		fmt.Print("Message : ", response.Message)
		fmt.Println("Stderr : ",response.Stderr)
		fmt.Println("--------- LAMBDA RESPONSE ---------")
	})
}
