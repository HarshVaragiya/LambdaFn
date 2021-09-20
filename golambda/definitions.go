package golambda

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	log           = logrus.New()
	statusCodeMap = map[string]int{
		"default":        200,
		"signal: killed": 201,
		"error":          208,
	}
	defaultTimeout = time.Second * 3
)

type LambdaExecutor interface {
	execute(Event) (Response, error)
}

type BasicCodeExecutor struct {
	codeUri         string
	functionHandler string
	context         context.Context
	functionTimeout time.Duration
}
