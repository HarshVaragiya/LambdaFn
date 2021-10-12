package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	log                         = logrus.New()
	defaultEnvApplicationPrefix = "LAMBDA_SERVER"
	defaults                    = map[string]interface{}{
		"ServerPort":       8080,
		"AwsAccountNumber": 123412341234,
		"Debug":            true,
	}
	defaultConfigPaths = []string{"."}
	statusCodeMap      = map[string]int32{
		"default":        200,
		"signal: killed": 201,
		"error":          208,
		"does-not-exist": 404,
	}
)

type ApplicationConfig struct {
	ServerPort       int
	AwsAccountNumber int64
	Debug            bool
}

func setLogLevel() {
	if appConfig.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
