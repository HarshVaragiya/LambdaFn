package golambda

type LambdaManager struct {
	awsAccountNumber		int64			// used to generate function ARN
	lambdaFunctions			map[string]Function
}
