package app

import (
	"go-start-template/internal/app"
	"log"

	"github.com/spf13/cobra"
)

func AppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Run main application",
		Run: func(cmd *cobra.Command, args []string) {
			addr, err := cmd.Flags().GetString("http_addr")
			if err != nil {
				log.Fatalf("failed to get addr: %v", err)
			}

			app.Run(addr)

			// TODO: Build and run application
		},
	}
	cmd.PersistentFlags().String("http_addr",
		"0.0.0.0:8080", "An address that is used to run http server")

	return cmd
}
