package golambda

import (
	"context"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
)

type LocalOsExecutor struct {
	codeExecutor BasicCodeExecutor
}

func NewLocalOsExecutor(function Function) *LocalOsExecutor {
	executor := BasicCodeExecutor{codeUri: function.CodeUri, functionHandler: function.Handler, functionTimeout: function.Timeout}
	return &LocalOsExecutor{codeExecutor: executor}
}

func (localOsExecutor *LocalOsExecutor) execute(event *liblambda.Event) (response *liblambda.Response, err error) {
	log.Debug("Invoking Local Binary.")
	stdout, stderr, responseString, err := localOsExecutor.runLocalBinary(event.EventData)
	return NewSimpleLambdaResponse(event.EventId, stdout, stderr, responseString, err), nil
}

func (localOsExecutor *LocalOsExecutor) runLocalBinary(event string) (string, string, string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), localOsExecutor.codeExecutor.functionTimeout)
	defer cancelFunc()
	return liblambda.RunLocalBinary(event, ctx, localOsExecutor.codeExecutor.codeUri, localOsExecutor.codeExecutor.functionHandler)
}
