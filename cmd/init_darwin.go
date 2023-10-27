package cmd

import (
	"fmt"
	"odoo-one-click/config"
	"odoo-one-click/utils"
	"os"
	"os/exec"
	"strings"

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
		CheckRequirement()
		checkPyenv()
	},
}

func CheckRequirement() {
	// Check requirement for ubuntu and derivatives
	fmt.Println("Checking requirement for macOS")

	err := exec.Command("which", "brew").Run()
	if err != nil {
		Logger.Println("Check brew command: ", err)
	}

	listOfDeps := []string{"jpeg", "postgresql", "zlib"}

	notInstalledDeps := make([]string, 0)

	for _, dep := range listOfDeps {
		err := exec.Command("brew", "list", dep).Run()
		// when the package already installed, err will be nil
		if err != nil {
			Logger.Println("Check dependencies: ", err)
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	if len(notInstalledDeps) > 0 {
		fmt.Println("Seems like you don't have all the dependencies installed, installing dependencies")

		// brew didn't need to use sudo like on linux, so can just directly run the brew command
		fmt.Println("Update apt repositories, please wait...")
		err = exec.Command("brew", "update").Run()
		if err != nil {
			Logger.Fatalln("Failed to update brew: ", err)
		}

		fmt.Println("Installing dependencies, please wait...")
		cmdAptInstall := utils.PrependCommand(notInstalledDeps, []string{"brew", "install"})

		err = exec.Command("sudo", cmdAptInstall...).Run()
		if err != nil {
			Logger.Fatalln("Install dependencies: ", err)
		}

		fmt.Println("Dependencies installed")
		err = exec.Command("brew", "services", "start", "postgresql").Run()
		if err != nil {
			Logger.Fatalln("Failed to start postgresql: ", err)
		}

	}

	_, err = checkDBAccess()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 127") {
			Logger.Fatalln("Postgresql is not installed, please install it first")
		}

		Logger.Fatalln("Can't access DB: ", err)
	}

	fmt.Printf("Database user '%s' with password '%s' created with superuser access\n", config.DBUsername(), config.DB_PASSWORD)

	err = addBuildPrerequisite()
	if err != nil {
		Logger.Fatalln("Failed to setup some requirements: ", err)
	}

	utils.PyenvInfoBash()
}

func checkDBAccess() (bool, error) {
	os.Setenv("PGPASSWORD", config.DB_PASSWORD)
	psqlCmd := fmt.Sprintf("psql -h %s -p %s -U %s -c 'SELECT 1'", config.DB_HOST, config.DB_PORT, config.DBUsername())

	err := exec.Command("zsh", "-c", psqlCmd).Run()
	if err != nil {
		return false, err
	}

	Logger.Println("Database can be accessed")

	return true, nil
}

func isPyenvInstalled() (bool, error) {
	err := exec.Command("pyenv", "--version").Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

func installPyenv() (bool, error) {
	pyenvScript := exec.Command("curl", "https://pyenv.run")
	runBash := exec.Command("zsh")
	runBash.Stdin, _ = pyenvScript.StdoutPipe()
	runBash.Stdout = os.Stdout
	_ = runBash.Start()
	_ = pyenvScript.Run()
	err := runBash.Wait()
	if err != nil {
		return false, err
	}

	err = exec.Command("zsh", "-c", "exec $SHELL").Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

func checkPyenv() {
	pyenvInstalled, err := isPyenvInstalled()
	if err != nil {
		Logger.Println(err)
	}

	if !pyenvInstalled {
		Logger.Println("Pyenv is not installed, installing pyenv")
		_, err := installPyenv()
		if err != nil {
			Logger.Fatalf("Failed to install pyenv: %s", err)
		}
	}

	if _, err := os.Stat(config.OdooDir()); os.IsNotExist(err) {
		err = os.MkdirAll(config.OdooDir(), config.ODOO_PERMISSION)
		if err != nil {
			Logger.Fatalln(err)
		}
	}
}

func addBuildPrerequisite() error {
	// some build flags and some brew packages need to be linked forcefully

	// link jpeg and zlib first
	err := exec.Command("brew", "link", "jpeg", "--force").Run()
	if err != nil {
		return err
	}

	err = exec.Command("brew", "link", "zlib", "--force").Run()
	if err != nil {
		return err
	}

	// add build flags to cppflags
	// LDFLAGS=-L$(brew --prefix openssl@1.1)/lib
	err = exec.Command("export", "LDFLAGS=-L$(brew --prefix openssl@1.1)/lib").Run()
	if err != nil {
		return err
	}
	// CPPFLAGS=-I$(brew --prefix openssl@1.1)/include
	err = exec.Command("export", "CPPFLAGS=-I$(brew --prefix openssl@1.1)/include").Run()
	if err != nil {
		return err
	}
	// CFLAGS="-Wno-error=implicit-function-declaration"
	err = exec.Command("export", "CFLAGS=\"-Wno-error=implicit-function-declaration\"").Run()
	if err != nil {
		return err
	}

	return nil
}
