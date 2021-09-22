package main

import "LambdaFn/golambda"

type LambdaRuntime struct {
	runtime      string
	handlerProxy string
}

// handlerProxy is a runtime language specific function that imports and invokes the actual lambda Handler function
// ex. for a python code present in /lambda of container, gRPC requests to container is received by the lambdaRuntime (go code)
// the request is then unmarshalled and the event, context are sent via STDIN to the "handlerProxy" which is a small python code
// the python code then imports the actual request handler for the lambda function and passes the event, context and gets response

func (obj LambdaRuntime) Invoke(event golambda.Event) (response golambda.Response, err error) {
	log.Debugf("Sending event '%s' to lambda Handler Proxy ", event.EventData)
	return
}
