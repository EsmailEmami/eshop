package cmd

import (
	"github.com/esmailemami/eshop/dbseed"
	"github.com/spf13/cobra"
)

var DBCmd = &cobra.Command{
	Use: "db",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	DBCmd.AddCommand(&cobra.Command{
		Use: "seed",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbseed.Run()
		},
	})
}
