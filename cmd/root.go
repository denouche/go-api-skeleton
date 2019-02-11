package cmd

import (
	"fmt"
	"os"

	"github.com/denouche/go-api-skeleton/handlers"
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config  = &handlers.Config{}
	cfgFile string
)

const (
	parameterConfigurationFile = "config"
	parameterMock              = "mock"
	parameterEnvironment       = "environment"
	parameterDBConnectionURI   = "dbconnectionuri"
	parameterPort              = "port"
)

var (
	defaultEnvironment     = "production"
	defaultDBConnectionURI = ""
	defaultPort            = 80
)

var rootCmd = &cobra.Command{
	Use:   "go-api-skeleton",
	Short: "go-api-skeleton",
	Run: func(cmd *cobra.Command, args []string) {
		utils.GetLoggerForEnvironment(config.Environment).Sugar().
			Infow("Configuration",
				parameterConfigurationFile, cfgFile,
				parameterEnvironment, config.Environment,
				parameterMock, config.Mock,
				parameterPort, config.Port,
				parameterDBConnectionURI, config.DBConnectionURI)

		router := handlers.NewRouter(config)
		router.Run(fmt.Sprintf(":%d", config.Port))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, parameterConfigurationFile, "", "Config file. All flags given in command line will override the values from this file.")

	rootCmd.Flags().String(parameterEnvironment, defaultEnvironment, "Use this flag to set the current environemnt")
	viper.BindPFlag(parameterEnvironment, rootCmd.Flags().Lookup(parameterEnvironment))

	rootCmd.Flags().String(parameterDBConnectionURI, defaultDBConnectionURI, "Use this flag to set the db connection URI")
	viper.BindPFlag(parameterDBConnectionURI, rootCmd.Flags().Lookup(parameterDBConnectionURI))

	rootCmd.Flags().Int(parameterPort, defaultPort, "Use this flag to set the listening port of the api")
	viper.BindPFlag(parameterPort, rootCmd.Flags().Lookup(parameterPort))

	rootCmd.Flags().Bool(parameterMock, false, "Use this flag to enable the mock mode")
	viper.BindPFlag(parameterMock, rootCmd.Flags().Lookup(parameterMock))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config.Mock = viper.GetBool(parameterMock)
	config.Port = viper.GetInt(parameterPort)
	config.Environment = viper.GetString(parameterEnvironment)
	config.DBConnectionURI = viper.GetString(parameterDBConnectionURI)
}
