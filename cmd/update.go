package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"odoo-one-click/config"
	"os"
	"runtime"
	"strings"

	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

type Asset struct {
	Url                string `json:"url"`
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	HtmlUrl    string  `json:"html_url"`
	Id         int     `json:"id"`
	TagName    string  `json:"tag_name"`
	Prerelease bool    `json:"prerelease"`
	Assets     []Asset `json:"assets"`
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your odoo-one-click app",
	Run: func(cmd *cobra.Command, args []string) {
		newVer, assets := CheckNewVersion()
		if !newVer {
			fmt.Println("Your installed odoo-one-click already on latest version.")
			return
		}

		// log.Println("You can update the binary by run again the installer script.")
		// TODO: implement self-update binary
		var url string
		for _, asset := range assets {
			if strings.Contains(asset.Name, fmt.Sprintf("odoo-one-click_%s_%s.zip", runtime.GOOS, runtime.GOARCH)) {
				url = asset.BrowserDownloadUrl
				break
			}
		}

		if url == "" {
			fmt.Println("No update for your current system, please make sure if your OS is supported.")
			os.Exit(1)
		}
		DoUpdate(url)
	},
}

func CheckNewVersion() (bool, []Asset) {
	installedVersion := config.VERSION

	resp, err := http.Get("https://api.github.com/repos/rockavoldy/odoo-one-click/releases/latest")
	if err != nil {
		fmt.Println("Can't check github. Please check again your internet connection", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var jsonData Release
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	currentVersion := jsonData.TagName
	if currentVersion != installedVersion {
		fmt.Printf("Your currently installed version is: %s. The latest one is: %s.\n", installedVersion, currentVersion)
		return true, jsonData.Assets
	}

	return false, nil
}

func DoUpdate(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Updated to latest version!")
}
