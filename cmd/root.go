package cmd

import (
	"fmt"
	"io"
	"log"
	"odoo-one-click/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "odoo-one-click",
	Short:            "Odoo-one-click is wrapper to install and run odoo easily",
	Long:             `Odoo-one-click is wrapper to install and run odoo easily.`,
	Version:          config.VERSION,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Verbose {
			// When verbose flag called, print logs to stdout
			log.SetOutput(io.Discard)
		}

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
