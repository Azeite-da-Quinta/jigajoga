/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Test client to connect to the game-srv",
	Long: `TODO EDIT
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("client started",
			slog.String("version", viper.GetString("version")),
			slog.Int("port", viper.GetInt("port")),
		)

		client.Dial(client.Config{
			Version: viper.GetString("version"),
			Host:    viper.GetString("host"),
		})
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String("host", "http://127.0.0.1:8080", "pass the host of the server to connect to")
	viper.BindPFlag("host", clientCmd.Flags().Lookup("host"))
	viper.SetDefault("host", "http://127.0.0.1:80")
}
