package cmd

import (
	"fmt"
	"log"
	"odoo-one-click/internal/initialize/ubuntu"
	"odoo-one-click/utils"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "First initialization to install pyenv, and configure postgresql",
	Long:  "First initialization to install pyenv, and configure postgresql for local development",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "linux" {
			if err := utils.CheckUbuntuVersion(); err != nil {
				log.Fatalln(err)
			}
			RunScript("ubuntu")
		} else {
			fmt.Println("This app currently only support Ubuntu 20.04 or later")
		}
	},
}

func RunScript(osPlatform string) error {
	if osPlatform == "ubuntu" {
		u := ubuntu.NewUbuntuInitializer()
		u.CheckAdminAccess()
		missingDeps, err := u.CheckRequirement()
		if err != nil {
			Logger.Fatalln(err)
		}
		err = u.InstallDeps(missingDeps)
		if err != nil {
			Logger.Fatalln(err)
		}

		err = u.ConfigureDB()
		if err != nil {
			Logger.Fatalln(err)
		}

		err = u.ConfigurePyenv()
		if err != nil {
			Logger.Fatalln(err)
		}

		utils.PyenvInfo()
	}
	return nil
}
