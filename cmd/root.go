package cmd

import (
	"fmt"
	"os"

	"github.com/denouche/go-api-skeleton/handlers"
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config  = &handlers.Config{}
	cfgFile string
)

const (
	parameterConfigurationFile = "config"
	parameterLogLevel          = "loglevel"
	parameterMock              = "mock"
	parameterLogFormat         = "logformat"
	parameterDBConnectionURI   = "dbconnectionuri"
	parameterPort              = "port"
)

var (
	defaultLogLevel        = logrus.WarnLevel.String()
	defaultLogFormat       = utils.LogFormatText
	defaultDBConnectionURI = ""
	defaultPort            = 80
)

var rootCmd = &cobra.Command{
	Use:   "go-api-skeleton",
	Short: "go-api-skeleton",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel := utils.ParseLogrusLevel(config.LogLevel)
		logrus.SetLevel(logLevel)

		logFormat := utils.ParseLogrusFormat(config.LogFormat)
		logrus.SetFormatter(logFormat)

		logrus.
			WithField(parameterConfigurationFile, cfgFile).
			WithField(parameterMock, config.Mock).
			WithField(parameterLogLevel, config.LogLevel).
			WithField(parameterLogFormat, config.LogFormat).
			WithField(parameterPort, config.Port).
			WithField(parameterDBConnectionURI, config.DBConnectionURI).
			Warn("Configuration")

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

	rootCmd.Flags().String(parameterLogLevel, defaultLogLevel, "Use this flag to set the logging level")
	viper.BindPFlag(parameterLogLevel, rootCmd.Flags().Lookup(parameterLogLevel))

	rootCmd.Flags().String(parameterLogFormat, defaultLogFormat, "Use this flag to set the logging format")
	viper.BindPFlag(parameterLogFormat, rootCmd.Flags().Lookup(parameterLogFormat))

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
	config.DBConnectionURI = viper.GetString(parameterDBConnectionURI)
	config.Port = viper.GetInt(parameterPort)
	config.LogLevel = viper.GetString(parameterLogLevel)
	config.LogFormat = viper.GetString(parameterLogFormat)
}
