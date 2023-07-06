package main

import (
	"fmt"
	"odoo-one-click/cmd"
	"odoo-one-click/config"
	"odoo-one-click/utils"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// OS Check, for now it's only for ubuntu and derivatives
	if runtime.GOOS != "linux" {
		fmt.Println("This application only works on Ubuntu or derivatives")
		os.Exit(1)
	} else {
		// Check if UBUNTU_CODENAME is on allowedOS
		if err := checkOSVersion(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd.Execute()
}

func checkOSVersion() error {
	skipUbuntuCheck := os.Getenv("SKIP_UBUNTU_CHECK")
	if skipUbuntuCheck != "" {
		return nil
	}

	// Continue to check the OS list when there is no skip ubuntu check
	out, err := exec.Command("bash", "-c", "source /etc/os-release; echo $UBUNTU_CODENAME").Output()
	if err != nil {
		return err
	}

	// check if UBUNTU_CODENAME is n allowedOS
	codename := utils.RemoveNewLine(string(out))
	if !config.IsAllowedOS(codename) {
		return fmt.Errorf("ubuntu version is not supported\nIf it's a mistakes, try to run the app again with\nenv SKIP_UBUNTU_CHECK=YES odoo-one-click")
	}

	return nil
}
