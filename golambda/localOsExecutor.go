package golambda

import (
	"context"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
)

type LocalOsExecutor struct {
	codeExecutor BasicCodeExecutor
}

func NewLocalOsExecutor(function Function) LocalOsExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: function.Timeout}
	return LocalOsExecutor{codeExecutor: executor}
}

func (localOsExecutor LocalOsExecutor) execute(event *lambda.Event) (response *lambda.Response, err error) {
	log.Debug("Invoking Local Binary.")
	stdout, stderr, err := localOsExecutor.runLocalBinary(event.EventData)
	return NewSimpleLambdaResponse(event.EventId, stdout, stderr, err), nil
}

func (localOsExecutor LocalOsExecutor) runLocalBinary(event string) (string, string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), localOsExecutor.codeExecutor.functionTimeout)
	defer cancelFunc()
	return lambda.RunLocalBinary(event, ctx, localOsExecutor.codeExecutor.codeUri, localOsExecutor.codeExecutor.functionHandler)
}
