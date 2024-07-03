package cmd

import (
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/srv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the game endpoint",
	Long: `This is the main command of the app.
It'll serve an HTTP and websocket server. Players joining
this server will have their messages forwarded to the correct
game room based on a JWT token`,
	Run: func(cmd *cobra.Command, args []string) {
		setLogLevel()

		s := srv.Server{
			Config: srv.Config{
				Port:      viper.GetInt(port),
				Version:   viper.GetString(version),
				JWTSecret: viper.GetString(jwtsecret),
				Mode:      viper.GetString(mode),
			},
		}

		s.Start()
	},
}

// flags/configs keys
const (
	port = "port"
	mode = "mode"
)

// default values
const (
	defaultPort = 80
	defaultMode = "single"
)

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().Int32P(port, "p", defaultPort, "the http listening port")
	viper.BindPFlag(port, serveCmd.Flags().Lookup(port))
	viper.SetDefault(port, defaultPort)

	serveCmd.Flags().StringP(mode, "m", defaultMode, "the exec mode")
	viper.BindPFlag(mode, serveCmd.Flags().Lookup(mode))
	viper.SetDefault(mode, defaultMode)
}
