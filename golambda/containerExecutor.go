package golambda

import (
	"context"
	"fmt"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type ContainerExecutor struct {
	functionName string
	runtime      string
	name         string
	rpcPort      string
	codeExecutor BasicCodeExecutor
	client       liblambda.LambdaClient
	Env          map[string]string
	State        string
	ContainerId  string
}

func NewContainerExecutor(function *Function) ContainerExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: defaultTimeout}
	containerExecutor := ContainerExecutor{codeExecutor: executor, runtime: function.Runtime, functionName: function.Name, Env: function.Env}
	return containerExecutor
}

func (executor ContainerExecutor) execute(event *liblambda.Event) (*liblambda.Response, error) {
	log.Infof("Invoking function [%s] in container executor.", executor.functionName)
	tempDir, _ := os.MkdirTemp(os.TempDir(), "LambdaFn")
	defer os.RemoveAll(tempDir)
	if err := copySrcArchiveToDir(executor.codeExecutor.codeUri, tempDir); err != nil {
		return NewErrorLambdaResponse(event.EventId, err, ""), err
	}
	err := executor.StartContainer(tempDir)
	if err != nil {
		log.Errorf("error starting docker container. exiting.")
		return NewErrorLambdaResponse(event.EventId, err, ""), err
	}

	defer executor.StopContainer()

	target := fmt.Sprintf("127.0.0.1:%s", executor.rpcPort)

	time.Sleep(time.Millisecond * 100) // give time for the gRPC server to start before connecting to it

	executor.client, err = NewLambdaClient(target)
	if err != nil {
		log.Warnf("error connecting to container executor. error = %v", err)
		return NewErrorLambdaResponse(event.EventId, err, ""), err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), executor.codeExecutor.functionTimeout)
	defer cancelFunc()
	response, err := executor.client.Invoke(ctx, event)
	if err != nil {
		log.Warnf("error executing function [%s]", executor.functionName)
	}
	return response, err
}

func copySrcArchiveToDir(codeUri, tempDir string) error {
	srcBytes, err := ioutil.ReadFile(codeUri)
	if err != nil {
		log.Errorf("error reading source from codeuri. error = %v", err)
		return err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s/source.zip", tempDir), srcBytes, 0666); err != nil {
		log.Errorf("error writing to required dir. error = %v", err)
		return err
	}
	return nil
}

func (executor *ContainerExecutor) StartContainer(tempDir string) error {
	startCtx, startCancelFunc := context.WithTimeout(context.Background(), containerTimeout)
	defer startCancelFunc()
	log.Debugf("attempting to start docker container for [%s]", executor.functionName)
	env := prepareEnvironmentVariables(executor.functionName, executor.codeExecutor.functionHandler, executor.Env)
	dockerImage, exists := runtimeToDockerImageMap[executor.runtime]
	if !exists {
		return fmt.Errorf("lambda runtime [%v] not defined", executor.runtime)
	}
	executor.rpcPort = fmt.Sprintf("%d", 2000+rand.Int()%20000)
	containerId, err := dockerContainerManager.startContainer(startCtx, dockerImage, executor.rpcPort, tempDir, containerArchivePath, env)
	if err != nil {
		log.Errorf("error starting docker container. error = %v", err)
		executor.State = "ERROR"
	} else {
		log.Debugf("started container %s with rpc port mapped to %s", containerId, executor.rpcPort)
		executor.State = "RUNNING"
	}
	executor.ContainerId = containerId
	return err
}

func (executor *ContainerExecutor) StopContainer() {
	go func() {
		log.Debugf("attempting to stop container for function [%s]", executor.functionName)
		ctx, cancelFunc := context.WithTimeout(context.Background(), containerTimeout)
		defer cancelFunc()
		err := dockerContainerManager.stopContainer(ctx, executor.ContainerId)
		if err != nil {
			log.Warnf("error stopping docker container. error = %v", err)
			executor.State = "ERROR-STOPPING"
		} else {
			log.Debugf("stopped docker container for function [%s]", executor.functionName)
			executor.State = "READY"
		}
	}()
}
