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
	parameterLogFormat         = "logformat"
	parameterDBConnectionURI   = "dbconnectionuri"
	parameterDBInMemory        = "dbinmemory"
	parameterDBName            = "dbname"
	parameterPort              = "port"
)

var (
	defaultLogLevel        = logrus.WarnLevel.String()
	defaultLogFormat       = utils.LogFormatText
	defaultDBConnectionURI = ""
	defaultDBName          = ""
	defaultPort            = 8080
)

var rootCmd = &cobra.Command{
	Use:   "go-api-skeleton",
	Short: "go-api-skeleton",
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitLogger(config.LogLevel, config.LogFormat)

		logrus.
			WithField(parameterConfigurationFile, cfgFile).
			WithField(parameterLogLevel, config.LogLevel).
			WithField(parameterLogFormat, config.LogFormat).
			WithField(parameterPort, config.Port).
			WithField(parameterDBInMemory, config.Mock).
			WithField(parameterDBConnectionURI, config.DBConnectionURI).
			WithField(parameterDBName, config.DBName).
			Warn("Configuration")

		router := handlers.NewRouter(config)
		err := router.Run(fmt.Sprintf(":%d", config.Port))
		if err != nil {
			utils.GetLogger().WithError(err).Error("error while starting app")
		}
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
	_ = viper.BindPFlag(parameterLogLevel, rootCmd.Flags().Lookup(parameterLogLevel))

	rootCmd.Flags().String(parameterLogFormat, defaultLogFormat, "Use this flag to set the logging format")
	_ = viper.BindPFlag(parameterLogFormat, rootCmd.Flags().Lookup(parameterLogFormat))

	rootCmd.Flags().Int(parameterPort, defaultPort, "Use this flag to set the listening port of the api")
	_ = viper.BindPFlag(parameterPort, rootCmd.Flags().Lookup(parameterPort))

	rootCmd.Flags().String(parameterDBConnectionURI, defaultDBConnectionURI, "Use this flag to set the db connection URI")
	_ = viper.BindPFlag(parameterDBConnectionURI, rootCmd.Flags().Lookup(parameterDBConnectionURI))

	rootCmd.Flags().String(parameterDBName, defaultDBName, "Use this flag to set the db name. This parameter is used when using a MongoDB database")
	_ = viper.BindPFlag(parameterDBName, rootCmd.Flags().Lookup(parameterDBName))

	rootCmd.Flags().Bool(parameterDBInMemory, false, "Use this flag to enable the db in memory mode")
	_ = viper.BindPFlag(parameterDBInMemory, rootCmd.Flags().Lookup(parameterDBInMemory))
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

	config.LogLevel = viper.GetString(parameterLogLevel)
	config.LogFormat = viper.GetString(parameterLogFormat)
	config.Port = viper.GetInt(parameterPort)
	config.DBConnectionURI = viper.GetString(parameterDBConnectionURI)
	config.DBName = viper.GetString(parameterDBName)
	config.DBInMemory = viper.GetBool(parameterDBInMemory)
}
