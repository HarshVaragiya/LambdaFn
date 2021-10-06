package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	viper2 "github.com/spf13/viper"
)

var (
	log = logrus.New()
)

func main() {
	log.Println("starting LambdaFn application deployment server")
	vi := viper2.New()
	appConfig := setupConfiguration(vi)

	if appConfig.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Infof("application config loaded. initializing gin router")

	router := gin.Default()
	lambdaRestApi := LambdaRestApi{router: router, lambdaManager: NewLambdaManager(appConfig.AwsAccountNumber)}
	lambdaRestApi.Init()

	log.Infof("serving REST application on port: %v", appConfig.ServerPort)

	if err := router.Run(fmt.Sprintf(":%d", appConfig.ServerPort)); err != nil {
		log.Fatalf("error running REST server. error = %v", err)
	}

}
