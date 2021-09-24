package main

import (
	"context"
	"github.com/HarshVaragiya/LambdaFn/golambda"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
)

// handlerProxy is a runtime language specific function that imports and invokes the actual lambda Handler function
// ex. for a python code present in /lambda of container, gRPC requests to container is received by the lambdaRuntime (go code)
// the request is then unmarshalled and the event, context are sent via STDIN to the "handlerProxy" which is a small python code
// the python code then imports the actual request handler for the lambda function and passes the event, context and gets response

type lambdaRpcServer struct {
	runtime			string
	handler			string
	bootstrap		string
	functionName	string
}

func (rpcServer lambdaRpcServer) Invoke(ctx context.Context, event *lambda.Event) (*lambda.Response,error){
	log.Printf("Function [%s] invoked with eventId [%s]", rpcServer.functionName, event.EventId)
	stdinInput, err := lambda.EventToStdinInput(event)
	if err != nil {
		log.Warnf("Function [%s] encountered error while processing eventId: [%s]", rpcServer.functionName, event.EventId)
		log.Warnf("Error invoking bootstrap binary. Unable to Marshal event. error = %v", err)
		return &lambda.Response{EventId: event.EventId, Message: err.Error()}, err
	}
	stdout, stderr, err := lambda.RunLocalBinary(stdinInput,ctx,bootstrap,handler)
	return golambda.NewSimpleLambdaResponse(event.EventId, stdout, stderr, err), nil
}
