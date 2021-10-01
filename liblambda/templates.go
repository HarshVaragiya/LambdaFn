package liblambda

type LambdaRestApiResponse struct {
	Message				string	`json:"message"`
	ErrorMessage		string	`json:"error-message"`
	FunctionName		string	`json:"function-name"`
	ModifiedResource	string	`json:"modified-resource"`
	DebugMessage		string	`json:"debug-message"`
}

func NewErrorRestResponse(errorMessage, functionName, debugMessage string) *LambdaRestApiResponse {
	return &LambdaRestApiResponse{ErrorMessage: errorMessage, FunctionName: functionName, DebugMessage: debugMessage}
}

func NewOkRestResponse(message, functionName, modifiedArn, debugMsg string) *LambdaRestApiResponse {
	return &LambdaRestApiResponse{Message: message, FunctionName: functionName,ModifiedResource: modifiedArn, DebugMessage: debugMsg}
}