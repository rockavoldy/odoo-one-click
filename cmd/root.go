package cmd

import (
	"fmt"
	"odoo-one-click/config"
	"odoo-one-click/utils"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "odoo-one-click",
	Short: "Odoo-one-click is wrapper to install and run odoo easily",
	Long:  `Odoo-one-click is wrapper to install and run odoo easily.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hei hoi!")
		fmt.Printf("Odoo one-click v%s\n", config.VERSION)
		odooVersion := config.OdooVersion()
		fmt.Println("Current odoo version: ", odooVersion)
		fmt.Println("Odoo instances directory: ", config.OdooDir())
		fmt.Println(utils.DirName(odooVersion, true))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
