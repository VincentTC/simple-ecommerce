package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "app",
		Short: "Simple Ecommerce",
		Run: func(*cobra.Command, []string) {
			server()
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "run-server",
		Short: "Run Server",
		Run: func(*cobra.Command, []string) {
			server()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "run-script-report",
		Short: "Run Script Report",
		Run: func(*cobra.Command, []string) {
			scriptReport()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "run-cron",
		Short: "Run Cron",
		Run: func(*cobra.Command, []string) {
			cron()
		},
	})
	cmd.Execute()
}
