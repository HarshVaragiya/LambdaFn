package golambda

import (
	"fmt"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
)

type ContainerExecutor struct {
	functionName string
	runtime      string
	name         string
	rpcPort      uint16
	codeExecutor BasicCodeExecutor
}

func NewContainerExecutor(function Function) ContainerExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: defaultTimeout}
	containerExecutor := ContainerExecutor{codeExecutor: executor, runtime: function.Runtime, functionName: function.Name}
	return containerExecutor
}

func (executor ContainerExecutor) execute(event *lambda.Event) (response *lambda.Response, err error) {
	log.Debugf("Invoking Lambda [%s] in container.", executor.functionName)
	return &lambda.Response{}, fmt.Errorf("not implemented")
}

func (executor ContainerExecutor) startContainer() (err error) {
	return fmt.Errorf("not implemented")
}

func (executor ContainerExecutor) stopContainer() (err error) {
	return fmt.Errorf("not implemeted")
}
