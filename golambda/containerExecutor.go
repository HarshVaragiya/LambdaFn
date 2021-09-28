package golambda

import (
	"context"
	"fmt"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
)

type ContainerExecutor struct {
	functionName string
	runtime      string
	name         string
	rpcPort      uint16
	codeExecutor BasicCodeExecutor
	client       liblambda.LambdaClient
}

func NewContainerExecutor(function Function) ContainerExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: defaultTimeout}
	containerExecutor := ContainerExecutor{codeExecutor: executor, runtime: function.Runtime, functionName: function.Name}
	return containerExecutor
}

func (executor ContainerExecutor) execute(event *liblambda.Event) (*liblambda.Response, error) {
	log.Debugf("Invoking function [%s] in container executor.", executor.functionName)
	executor.rpcPort = 8888
	target := fmt.Sprintf("127.0.0.1:%v", executor.rpcPort)
	var err error
	executor.client, err = NewLambdaClient(target)
	if err != nil {
		log.Warnf("Error connecting to container executor. error = %v", err)
		return NewErrorLambdaResponse(event.EventId, err, ""), err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), executor.codeExecutor.functionTimeout)
	defer cancelFunc()
	response, err := executor.client.Invoke(ctx, event)
	if err != nil {
		log.Warnf("Error executing function [%s]", executor.functionName)
	}
	return response, err
}
