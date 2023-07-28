package cmd

import (
	"errors"

	"github.com/esmailemami/eshop/migrations"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use: "migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	MigrateCmd.AddCommand(&cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrations.Migrate()
		},
	})

	MigrateCmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrations.Rollback()
		},
	})

	MigrateCmd.AddCommand(&cobra.Command{
		Use:   "make",
		Short: "Making a new migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("migration name is required")
			}
			name := args[0]
			return migrations.MakeMigration(&name)
		},
	})

	MigrateCmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Rollback all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrations.RollbackAll()
		},
	})
}
