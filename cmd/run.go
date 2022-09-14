package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run odoo instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run odoo instance")
	},
}
