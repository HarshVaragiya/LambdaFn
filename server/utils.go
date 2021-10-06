package main

import (
	viper2 "github.com/spf13/viper"
)

func setupConfiguration(vi *viper2.Viper) *ApplicationConfig {
	log.Tracef("setting application defaults")
	for key, value := range defaults {
		vi.SetDefault(key, value)
	}
	for _, path := range defaultConfigPaths {
		vi.AddConfigPath(path)
	}

	vi.SetConfigType("toml")
	vi.SetEnvPrefix(defaultEnvApplicationPrefix)
	vi.SetConfigName(defaultConfigFileName)

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
	return &appConfig
}
