package cmd

import (
	"fmt"
	"io"
	"log"
	"odoo-one-click/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run odoo instance",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Verbose {
			// TODO: create 1 log file for the project to use; can be extended to log to a file
			log.SetOutput(io.Discard)
		}
		fmt.Println("Run odoo instance")
	},
}
