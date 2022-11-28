package main

import (
	"fmt"
	"odoo-one-click/cmd"
	"odoo-one-click/utils"
	"os"
	"runtime"
)

func main() {
	// OS Check, for now it's only for ubuntu and derivatives
	if runtime.GOOS != "linux" {
		fmt.Println("This program only works on Ubuntu or derivatives")
		os.Exit(1)
	} else {
		// Check if UBUNTU_CODENAME is on allowedOS
		if err := utils.CheckUbuntuVersion(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd.Execute()
}
