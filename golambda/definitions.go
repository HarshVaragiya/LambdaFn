package golambda

import (
	"context"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

func init() {
	if err := dockerContainerManager.Init(); err != nil{
		log.Fatalf("error initializing docker container manager. error = %v", err)
	}
}

var (
	log           = logrus.New()
	statusCodeMap = map[string]int32{
		"default":        200,
		"signal: killed": 201,
		"error":          208,
		"does-not-exist": 404,
	}
	defaultTimeout = time.Second * 10
	containerRpcPort = uint16(8888)
	containerTimeout = time.Second * 2
	dockerContainerManager = DockerContainerManager{}
	runtimeToDockerImageMap = map[string]string{
		"python3": "amazonlinux-python",
	}
	containerArchivePath = "/lambda/archive/"
)

type LambdaExecutor interface {
	execute(event *liblambda.Event)(*liblambda.Response, error)
}

type BasicCodeExecutor struct {
	codeUri         string
	functionHandler string
	context         context.Context
	functionTimeout time.Duration
}


func NewSimpleLambdaEvent(event string) *liblambda.Event {
	return &liblambda.Event{EventData: event, EventId: uuid.New().String()}
}

func NewSimpleLambdaResponse(eventId, stdout, stderr, responseString string, err error) (response *liblambda.Response) {
	statusCode := int32(200)
	if err != nil {
		if knownErrorStatusCode, exists := statusCodeMap[err.Error()]; exists {
			statusCode = knownErrorStatusCode
		} else {
			statusCode = 404
		}
		stderr += err.Error()
	}
	return &liblambda.Response{Data: responseString, Stderr: stderr, StatusCode: statusCode, EventId: eventId, Message: stdout}
}

func NewErrorLambdaResponse(eventId string, err error, stderr string) *liblambda.Response {
	return &liblambda.Response{EventId: eventId, Message: err.Error(), Stderr: stderr}
}