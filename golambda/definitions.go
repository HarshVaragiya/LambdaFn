package golambda

import (
	"context"
	lambda "github.com/HarshVaragiya/LambdaFn/liblambda"
	"github.com/google/uuid"
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
	runtimeToDockerfileMap = map[string]string{
		"python": "../templates/python.dockerfile.template",
	}
	defaultTimeout = time.Second * 10
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
	return &lambda.Event{EventData: event, EventId: uuid.New().String()}
}

func NewSimpleLambdaResponse(eventId, stdout, stderr, responseString string, err error) (response *lambda.Response) {
	statusCode := int32(200)
	if err != nil {
		if knownErrorStatusCode, exists := statusCodeMap[err.Error()]; exists {
			statusCode = knownErrorStatusCode
		} else {
			statusCode = 404
		}
		stderr += err.Error()
	}
	return &lambda.Response{Data: responseString, Stderr: stderr, StatusCode: statusCode, EventId: eventId, Message: stdout}
}

func NewErrorLambdaResponse(eventId string, err error, stderr string) *lambda.Response {
	return &lambda.Response{EventId: eventId, Message: err.Error(), Stderr: stderr}
}