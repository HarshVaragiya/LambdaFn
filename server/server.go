package main

import (
	log "github.com/sirupsen/logrus"
	viper2 "github.com/spf13/viper"
)

func main() {
	log.Println("Starting LambdaFn Server")
	vi := viper2.New()
	setApplicationDefaults(vi)

	if err := vi.ReadInConfig(); err != nil {
		if _, ok := err.(viper2.ConfigFileNotFoundError); ok {
			log.Warn("config file not found. creating one using the default values.")
			if err = vi.SafeWriteConfig(); err != nil {
				log.Errorf("error saving default config file. error = %v", err)
			} else {
				log.Info("saved default config file.")
			}
		} else {
			log.Fatalf("error loading config file. error = %v", err)
		}
	}

	var appConfig ApplicationConfig
	if err := vi.Unmarshal(&appConfig); err != nil {
		log.Fatalf("error unmarshalling config file. error = %v", err)
	}
	log.Printf("Application Config Loaded : %v \n", appConfig)
}
