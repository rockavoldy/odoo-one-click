package cmd

import (
	"fmt"
	"odoo-one-click/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure odoo",
	Long:  "Install and configure odoo with demo data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.ODOO_VERSION)
		InstallOdoo()
	},
}

func InstallOdoo() {
	fmt.Println("Installing Odoo")
}
