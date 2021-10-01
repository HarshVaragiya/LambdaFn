package golambda

import (
	"fmt"
	"github.com/HarshVaragiya/LambdaFn/liblambda"
	"time"
)

type Function struct {
	Name        	string				`json:"name"`
	Description 	string				`json:"description"`
	Executor    	LambdaExecutor
	CodeUri     	string				`json:"codeuri"`
	Handler     	string				`json:"handler"`
	Timeout     	time.Duration
	TimeoutSeconds	string				`json:"timeout-seconds"`
	Runtime     	string				`json:"runtime"`
	Arn				string				`json:"arn"`
	Tags			map[string]string	`json:"tags"`
	Env				map[string]string	`json:"env"`
}

func (function *Function) Invoke(event *liblambda.Event) (response *liblambda.Response, err error) {
	log.Infof("Invoking Lambda [%s]", function.Name)
	log.Tracef("Event: %s", event.EventData)
	response, err = function.Executor.execute(event)
	if err != nil {
		log.Errorf("Error Invoking Lambda Function. error = %v", err)
		return
	}
	log.Tracef("Response Data: [%s]", response.Data)
	log.Tracef("Response StatusCode: [%v]", response.StatusCode)
	log.Tracef("Response Log: [%s]", response.Stderr)
	return
}

func NewLocalOsLambdaFunction(name, codeUri, argument string) Function {
	log.Infof("Creating New Default Local OS Lambda Function [%s]", name)
	function := Function{Name: name, CodeUri: codeUri, Handler: argument, Timeout: defaultTimeout}
	function.Executor = NewLocalOsExecutor(function)
	return function
}

func NewContainerLambdaFunction(name, codeUri, runtime, handler string) *Function {
	log.Infof("Creating New Default Container Lambda Function [%s]", name)
	env := make(map[string]string)
	function := Function{Name: name, CodeUri: codeUri, Handler: handler, Timeout: defaultTimeout, Runtime: runtime, Env: env}
	function.Executor = NewContainerExecutor(&function)
	return &function
}

func (function *Function) CalculateArn(awsAccountNumber int64) string {
	function.Arn = fmt.Sprintf("arn:aws:lambda:us-east-1:%d:function:%s",awsAccountNumber, function.Name)
	return function.Arn
}