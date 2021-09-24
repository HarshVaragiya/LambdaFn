package golambda

import (
	"context"
	"fmt"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
	"github.com/docker/docker/client"
)

type ContainerExecutor struct {
	functionName string
	runtime      string
	name         string
	rpcPort      uint16
	codeExecutor BasicCodeExecutor
	client		 lambda.LambdaClient
}

func NewContainerExecutor(function Function) ContainerExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: defaultTimeout}
	containerExecutor := ContainerExecutor{codeExecutor: executor, runtime: function.Runtime, functionName: function.Name}
	return containerExecutor
}

func (executor ContainerExecutor) execute(event *lambda.Event) (*lambda.Response, error) {
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

func (executor ContainerExecutor) startContainer() (err error) {
	_, err = client.NewClientWithOpts()
	if err != nil {
		log.Warnf("Error communicating with docker daemon. error = %v", err)
		return
	}

	//ctx := context.Background()
	//resp, err := cli.ContainerCreate(ctx, &container.Config{
	//	Image:        "mongo",
	//	ExposedPorts: nat.PortSet{"8080": struct{}{}},
	//}, &container.HostConfig{
	//	PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("8080"): {{HostIP: "127.0.0.1", HostPort: "8080"}}},
	//}, nil, "mongo-go-cli")
	//if err != nil {
	//	panic(err)
	//}
	//
	//if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	//	panic(err)
	//}
	return fmt.Errorf("not implemented")
}

func (executor ContainerExecutor) stopContainer() (err error) {
	return fmt.Errorf("not implemeted")
}
