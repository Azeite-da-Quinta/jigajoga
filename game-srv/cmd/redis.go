package cmd

import (
	"github.com/spf13/cobra"
)

// redisCmd
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "x",
	Long:  `x`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(redisCmd)
}
