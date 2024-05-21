package cmd

import (
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the game endpoint",
	Long: `TODO EDIT
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(server.Config{
			Port:    viper.GetInt("port"),
			Version: viper.GetString("version"),
		})
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().Int32P("port", "p", 8080, "")
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.SetDefault("port", 80)
}
