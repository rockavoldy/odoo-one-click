package cmd

import (
	"fmt"
	"log"
	"odoo-one-click/config"
	"odoo-one-click/utils"
	"os"

	"github.com/spf13/cobra"
)

var Logger *log.Logger

var rootCmd = &cobra.Command{
	Use:              "odoo-one-click",
	Short:            "Odoo-one-click is wrapper to install and run odoo easily",
	Long:             `Odoo-one-click is wrapper to install and run odoo easily.`,
	Version:          config.VERSION,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		Logger = utils.Logger(config.Verbose)

		fmt.Println("Hei hoi!")
		fmt.Println("Use --help to see available commands")
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Print logs to stdout")
	rootCmd.SetVersionTemplate(fmt.Sprintf("Odoo one-click v%s\n", config.VERSION))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
