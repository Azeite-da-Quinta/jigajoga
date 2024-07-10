// Package cmd holds all game-srv Cobra commands
package cmd

import (
	"log/slog"
	"os"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/slogt"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/token"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// filled by rootCmd flag
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "game-srv",
	Short: "The Jigajoga game server app",
	Long: `This app may run a game server. 
It supports the creation of game rooms that players may join.
Handles communication through websockets`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets
// flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic("root cmd failed to execute")
	}
}

// flags/configs keys
const (
	version   = "version"
	level     = "level"
	jwtsecret = "jwtsecret"
)

// default values
const (
	defaultVersion = "v0.1.0"
	defaultLevel   = "INFO"
	defaultSecret  = token.DefaultSecret
)

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.game-srv.yaml)")

	rootCmd.PersistentFlags().StringP(version, "v", defaultVersion,
		"defines the version of the app")
	viper.BindPFlag(version, rootCmd.PersistentFlags().Lookup(version))
	viper.SetDefault(version, defaultVersion)

	rootCmd.PersistentFlags().StringP(level, "l", defaultLevel,
		"defines the log level of the app")
	viper.BindPFlag(level, rootCmd.PersistentFlags().Lookup(level))
	viper.SetDefault(level, defaultLevel)

	rootCmd.PersistentFlags().StringP(jwtsecret, "s", defaultSecret,
		"defines the version of the app")
	viper.BindPFlag(jwtsecret, rootCmd.PersistentFlags().Lookup(jwtsecret))
	viper.SetDefault(jwtsecret, defaultSecret)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		// ❗ at this stage, it might not exist
		viper.SetConfigFile(cfgFile)
	} else {
		viperFromHome()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("using config file", "file", viper.ConfigFileUsed())
	}

	viperWatch()
}

func viperFromHome() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name
	//  ".game-srv.yml".
	viper.AddConfigPath(home)

	viper.SetConfigType("yaml")
	viper.SetConfigName(".game-srv.yml")
}

func viperWatch() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("config file changed", "file", e.Name)
		// Careful not displaying senstive info here
		setLogLevel()
	})
	viper.WatchConfig()
}

func setLogLevel() {
	var lvl slog.Level
	err := lvl.UnmarshalText([]byte(viper.GetString(level)))
	if err != nil {
		slog.Error("slog: could not unmarshal level", slogt.Error(err))
		return
	}

	slog.SetLogLoggerLevel(lvl)
}
