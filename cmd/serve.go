package cmd

import (
	"log"

	"github.com/esmailemami/eshop/server"
	"github.com/esmailemami/eshop/services/token"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the api",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := token.InitJWT(); err != nil {
			log.Fatalln(err)
		}

		server.RunServer()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
