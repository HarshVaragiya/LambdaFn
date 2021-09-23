package golambda

import (
	"context"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	log           = logrus.New()
	statusCodeMap = map[string]int32{
		"default":        200,
		"signal: killed": 201,
		"error":          208,
	}
	defaultTimeout = time.Second * 3
)

type LambdaExecutor interface {
	execute(event *lambda.Event)(*lambda.Response, error)
}

type BasicCodeExecutor struct {
	codeUri         string
	functionHandler string
	context         context.Context
	functionTimeout time.Duration
}


func NewSimpleLambdaEvent(event string) *lambda.Event {
	return &lambda.Event{EventData: event}
}

func NewSimpleLambdaResponse(stdout, stderr string, err error) (response *lambda.Response) {
	statusCode := int32(200)
	if err != nil {
		if knownErrorStatusCode, exists := statusCodeMap[err.Error()]; exists {
			statusCode = knownErrorStatusCode
		} else {
			statusCode = 404
		}
		stderr += err.Error()
	}
	return &lambda.Response{Data: stdout, Stderr: stderr, StatusCode: statusCode}
}
