package cmd

import (
	"log/slog"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Test client to connect to the game-srv",
	Long: `Use this command to connect to the game-srv.
You can pass a few arguments like host, workers and jobs.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("client started",
			slog.String(version, viper.GetString(version)),
			slog.String(host, viper.GetString(host)),
		)

		client.Dial(client.Config{
			Version:   viper.GetString("version"),
			Host:      viper.GetString(host),
			NbWorkers: viper.GetInt(workers),
			NbWrites:  viper.GetInt(jobs),
			JWTSecret: viper.GetString(jwtsecret),
		})
	},
}

// flags/configs keys
const (
	host    = "host"
	workers = "workers"
	jobs    = "jobs"
)

// default values
const (
	defautWorkers = 5
	defaultJobs   = 20
	defaultHost   = "127.0.0.1:80"
)

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String(host, defaultHost,
		"pass the host of the server to connect to")
	clientCmd.Flags().IntP("workers", "w", defautWorkers,
		"how many workers should run")
	clientCmd.Flags().IntP(jobs, "j", defaultJobs,
		"how many jobs each worker should do")

	viper.BindPFlag(host, clientCmd.Flags().Lookup(host))
	viper.BindPFlag(workers, clientCmd.Flags().Lookup(workers))
	viper.BindPFlag(jobs, clientCmd.Flags().Lookup(jobs))

	viper.SetDefault(host, defaultHost)
	viper.SetDefault(workers, defautWorkers)
	viper.SetDefault(jobs, defaultJobs)
}
