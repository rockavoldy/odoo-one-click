package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "odoo-one-click",
	Short: "Odoo-one-click is wrapper to install and run odoo easily",
	Long:  `Odoo-one-click is wrapper to install and run odoo easily.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hei hoi!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
