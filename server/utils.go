package main

import (
	log "github.com/sirupsen/logrus"
	viper2 "github.com/spf13/viper"
)

func setApplicationDefaults(vi *viper2.Viper){
	log.Tracef("setting application defaults")
	for key, value := range defaults{
		vi.SetDefault(key, value)
	}
	for _, path := range defaultConfigPaths {
		vi.AddConfigPath(path)
	}
	vi.SetConfigType("toml")
	vi.SetEnvPrefix(defaultEnvApplicationPrefix)
	vi.SetConfigName(defaultConfigFileName)
}
