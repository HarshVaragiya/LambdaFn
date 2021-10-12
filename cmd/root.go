/*
Copyright © 2021 Harsh Varagiya

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string
var appConfig *ApplicationConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "LambdaFn",
	Short: "LambdaFn Application Server",
	Long: `LambdaFn application server is responsible for managing lambda functions
	and running them inside docker containers whenever the function is invoked. 
	The Application server allows creation / invocation / deletion of lambda functions.
	The --code-uri parameter of lambdacli should be relative to the application server.`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("starting LambdaFn application server")
		if err := viper.Unmarshal(&appConfig); err != nil {
			log.Fatalf("error unmarshalling config file. error = %v", err)
		}
		setLogLevel()
		log.Debug("application config loaded. initializing gin router")
		router := gin.Default()
		lambdaRestApi := LambdaRestApi{router: router, lambdaManager: NewLambdaManager(appConfig.AwsAccountNumber)}
		lambdaRestApi.Init()

		log.Infof("initialized REST API. serving application on port: [ %v ]", appConfig.ServerPort)

		if err := router.Run(fmt.Sprintf(":%d", appConfig.ServerPort)); err != nil {
			log.Fatalf("error running REST server. error = %v", err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.LambdaFn.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".LambdaFn" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".LambdaFn")
		viper.SetEnvPrefix(defaultEnvApplicationPrefix)
		for _, path := range defaultConfigPaths {
			viper.AddConfigPath(path)
		}
		for key, value := range defaults {
			viper.SetDefault(key, value)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config cmdfile is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %v", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("config file not found. creating one using the default values")
			if err = viper.SafeWriteConfig(); err != nil {
				log.Errorf("error saving default config file. error = %v", err)
			} else {
				log.Info("saved default config file")
			}
		} else {
			log.Fatalf("error loading config file. error = %v", err)
		}
	}
}
