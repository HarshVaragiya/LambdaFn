package golambda

import (
	"time"
)

type Function struct {
	Name        string
	Description string
	Executor    LambdaExecutor
	CodeUri     string
	Handler     string
	Timeout     time.Duration
	Runtime     string
}

func (function Function) Invoke(event Event) (response Response, err error) {
	log.Infof("Invoking Lambda [%s]", function.Name)
	log.Tracef("Event: %s", event.EventData)
	response, err = function.Executor.execute(event)
	if err != nil {
		log.Errorf("Error Invoking Lambda Function. error = %v", err)
	}
	log.Tracef("Response Data: [%s]", response.Data)
	log.Tracef("Response StatusCode: [%v]", response.StatusCode)
	log.Tracef("Response Log: [%s]", response.Stderr)
	return
}

func NewLocalOsLambdaFunction(name, codeUri, argument string) Function {
	log.Infof("Creating New Lambda Function [%s] with defaults.", name)
	function := Function{Name: name, CodeUri: codeUri, Handler: argument, Timeout: defaultTimeout}
	function.Executor = NewLocalOsExecutor(function)
	return function
}
