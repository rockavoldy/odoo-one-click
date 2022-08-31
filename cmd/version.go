package cmd

import (
	"fmt"
	"odoo-one-click/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current installed version",
	Long:  "Show current installed version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.VERSION)
		fmt.Println(args)
	},
}
