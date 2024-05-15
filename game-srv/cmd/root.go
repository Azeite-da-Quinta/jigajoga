package cmd

import (
	"log/slog"
	"os"

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.game-srv.yaml)")

	rootCmd.PersistentFlags().StringP("version", "v", "v0.0.0", "defines the version of the app")
	viper.BindPFlag("version", serveCmd.Flags().Lookup("version"))
	viper.SetDefault("version", "v0.1.0")

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
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".game-srv" (without extension).
		viper.AddConfigPath(home)

		viper.SetConfigType("yaml")
		viper.SetConfigName(".game-srv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("using config file", "file", viper.ConfigFileUsed())
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("config file changed", "file", e.Name)
	})
	viper.WatchConfig()
}
