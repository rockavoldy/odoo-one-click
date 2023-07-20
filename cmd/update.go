package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"odoo-one-click/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your odoo-one-click app",
	Run: func(cmd *cobra.Command, args []string) {
		newVer := CheckNewVersion()
		if !newVer {
			log.Println("Your installed odoo-one-click already on latest version.")
		}

		log.Println("You can update the binary by run again the installer script.")
		// TODO: implement self-update binary
		// check for update
		// then run the update script
	},
}

func CheckNewVersion() bool {
	installedVersion := config.VERSION

	resp, err := http.Get("https://api.github.com/repos/rockavoldy/odoo-one-click/releases/latest")
	if err != nil {
		log.Fatalln("Can't check github. Please check again your internet connection", err)
	}
	defer resp.Body.Close()

	var jsonData map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	currentVersion := jsonData["tag_name"]

	if currentVersion != installedVersion {
		log.Printf("Your currently installed version is: %s. The latest one is: %s.\n", installedVersion, currentVersion)
		return true
	}

	return false
}
