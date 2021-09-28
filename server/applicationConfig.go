package main

var (
	defaultEnvApplicationPrefix = "LAMBDA_SERVER"
	defaults = map[string]interface{}{
		"ServerPort": 8080,
		"AwsAccountNumber": 123412341234,
	}
	defaultConfigFileName = "config"
	defaultConfigPaths = []string{"."}
)

type ApplicationConfig struct {
	ServerPort			int
	AwsAccountNumber	int64
}
