package main

var (
	defaultEnvApplicationPrefix = "LAMBDA_SERVER"
	defaults = map[string]interface{}{
		"ServerPort": 8080,
		"AwsAccountNumber": 123412341234,
	}
	defaultConfigFileName = "config"
	defaultConfigPaths = []string{"."}
	statusCodeMap = map[string]int32{
		"default":        200,
		"signal: killed": 201,
		"error":          208,
		"does-not-exist": 404,
	}
)

type ApplicationConfig struct {
	ServerPort			int
	AwsAccountNumber	int64
}
