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
	fmt.Println("Checking requirement")

	listOfDeps := []string{"build-essential", "postgresql", "postgresql-client", "libxml2-dev", "libssl-dev", "libffi-dev", "libxslt1-dev", "libldap2-dev", "libsasl2-dev", "libtiff5-dev", "libjpeg8-dev", "libopenjp2-7-dev", "zlib1g-dev", "libfreetype6-dev", "liblcms2-dev", "libwebp-dev", "libharfbuzz-dev", "libpq-dev", "git", "libsqlite3-dev", "libreadline-dev", "libbz2-dev", "tk-dev"}

	notInstalledDeps := make([]string, 0)

	dpkgStatus := utils.PrependCommand(listOfDeps, []string{"--status"})

	for _, dep := range listOfDeps {
		err := exec.Command("dpkg", dpkgStatus...).Run()
		// when the package already installed, err will be nil
		if err != nil {
			Logger.Println("Check dependencies: ", err)
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	if len(notInstalledDeps) > 0 {
		fmt.Println("Seems like you don't have all the dependencies installed, installing dependencies")

		// When there is still missing dependencies, install them
		err := utils.CheckSudoAccess()
		if err != nil {
			Logger.Fatalln("Wrong password: ", err)
		}

		fmt.Println("Update apt repositories, please wait...")
		err = exec.Command("sudo", "apt-get", "update").Run()
		if err != nil {
			Logger.Fatalln("Failed to update repositories: ", err)
		}

		fmt.Println("Installing dependencies, please wait...")
		cmdAptInstall := []string{"apt-get", "install", "-y"}
		cmdAptInstall = utils.PrependCommand(notInstalledDeps, cmdAptInstall)

		err = exec.Command("sudo", cmdAptInstall...).Run()
		if err != nil {
			Logger.Fatalln("Install dependencies: ", err)
		}

		fmt.Println("Dependencies installed")
		err = exec.Command("sudo", "service", "postgresql", "start").Run()
		if err != nil {
			Logger.Fatalln("Failed to start postgresql: ", err)
		}

	}

	dbAccess, err := checkDBAccess()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 127") {
			Logger.Fatalln("Postgresql is not installed, please install it first")
		}

		Logger.Println("Can't access DB: ", err)
	}

	if !dbAccess {
		err := configureDB()
		if err != nil {
			Logger.Fatalln("Error when configure DB: ", err)
		}
		fmt.Println("Database successfully configured")
	}

	fmt.Printf("Database user '%s' with password '%s' created with superuser access\n", config.DBUsername(), config.DB_PASSWORD)

	utils.PyenvInfoBash()
}

func checkDBAccess() (bool, error) {
	os.Setenv("PGPASSWORD", config.DB_PASSWORD)
	psqlCmd := fmt.Sprintf("psql -h %s -p %s -U %s -c 'SELECT 1'", config.DB_HOST, config.DB_PORT, config.DBUsername())

	err := exec.Command("bash", "-c", psqlCmd).Run()
	if err != nil {
		return false, err
	}

	Logger.Println("Database can be accessed")

	return true, nil
}

func configureDB() error {
	// TODO: Add validation first, to make sure no create command again executed
	os.Setenv("PGPASSWORD", config.DB_PASSWORD)
	psqlScript := fmt.Sprintf(`psql -c "CREATE ROLE %s SUPERUSER LOGIN PASSWORD '%s';"`, config.DBUsername(), config.DB_PASSWORD)
	err := exec.Command("sudo", "su", "-", "postgres", "-c", psqlScript).Run()
	if err != nil {
		Logger.Println("Create role: ", err.Error())
	}

	// if db for the user already exist, there is no need to new one, so this will return error
	err = exec.Command("createdb", "-h", config.DB_HOST, "-U", config.DBUsername(), config.DBUsername()).Run()
	if err != nil {
		Logger.Println("Create database: ", err.Error())
	}

	return nil
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
	runBash := exec.Command("bash")
	runBash.Stdin, _ = pyenvScript.StdoutPipe()
	runBash.Stdout = os.Stdout
	_ = runBash.Start()
	_ = pyenvScript.Run()
	err := runBash.Wait()
	if err != nil {
		return false, err
	}

	err = exec.Command("bash", "-c", "exec $SHELL").Run()
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
