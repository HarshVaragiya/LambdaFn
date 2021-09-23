package golambda

import (
	"bytes"
	"context"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
	"io"
	"os/exec"
)

type LocalOsExecutor struct {
	codeExecutor BasicCodeExecutor
}

func NewLocalOsExecutor(function Function) LocalOsExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: function.Timeout}
	return LocalOsExecutor{codeExecutor: executor}
}

func (localOsExecutor LocalOsExecutor) execute(event *lambda.Event) (response *lambda.Response, err error) {
	log.Debug("Invoking Local OS Executor.")
	stdout, stderr, err := localOsExecutor.runLocalBinary(event.EventData)
	return NewSimpleLambdaResponse(stdout, stderr, err), nil
}

func (localOsExecutor LocalOsExecutor) runLocalBinary(event string) (string, string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), localOsExecutor.codeExecutor.functionTimeout)
	defer cancelFunc()
	cmd := exec.CommandContext(ctx, localOsExecutor.codeExecutor.codeUri, localOsExecutor.codeExecutor.functionHandler)
	log.Debugf("Executing local binary: '%s', with args: '%s' ", localOsExecutor.codeExecutor.codeUri, localOsExecutor.codeExecutor.functionHandler)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Warn("Unable to Open STDIN.")
		log.Warn(err)
		return "", "", err
	}
	_, err = io.WriteString(stdin, event)
	if err != nil {
		log.Warn("Unable to write to STDIN.")
		log.Warn(err)
		return "", "", err
	}
	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	err = cmd.Run()
	if err != nil {
		log.Warn("Error running the required binary")
		log.Warn(err)
	}
	return stdoutBytes.String(), stderrBytes.String(), err
}
