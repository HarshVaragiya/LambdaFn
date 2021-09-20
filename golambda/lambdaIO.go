package golambda

type Event struct {
	EventData string
	Context   string
}

type Response struct {
	Data       string
	Stderr     string
	StatusCode int
	Message    string
}

func NewSimpleLambdaEvent(event string) Event {
	return Event{EventData: event}
}

func NewSimpleLambdaResponse(stdout, stderr string, err error) (response Response) {
	statusCode := 200
	if err != nil {
		if knownErrorStatusCode, exists := statusCodeMap[err.Error()]; exists {
			statusCode = knownErrorStatusCode
		} else {
			statusCode = 404
		}
	}
	stderr += err.Error()
	return Response{Data: stdout, Stderr: stderr, StatusCode: statusCode}
}
