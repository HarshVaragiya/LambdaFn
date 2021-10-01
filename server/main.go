package main

import (
	"github.com/sirupsen/logrus"
	viper2 "github.com/spf13/viper"
)

var (
	log = logrus.New()
)

func main() {
	log.Println("Starting LambdaFn Server")
	vi := viper2.New()
	appConfig := setupConfiguration(vi)
	log.Printf("Application Config Loaded : %v \n", appConfig)
}
