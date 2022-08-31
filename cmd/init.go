package cmd

import (
	"fmt"
	"log"
	"odoo-one-click/config"
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
	Short: "First initialization, install pyenv, configure postgresql and clone odoo",
	Long:  "First initialization for installing pyenv, configure postgresql, and clone odoo",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(config.VERSION)
		CheckRequirement()
		checkPyenv()
	},
}

func CheckRequirement() {
	// first, need to confirm if it is ubuntu or derivatives

	listOfDeps := []string{"postgresql", "postgresql-client", "libxml2-dev", "libxslt1-dev", "libldap2-dev", "libsasl2-dev", "libtiff5-dev", "libjpeg8-dev", "libopenjp2-7-dev", "zlib1g-dev", "libfreetype6-dev", "liblcms2-dev", "libwebp-dev", "libharfbuzz-dev", "libpq-dev", "git"}

	notInstalledDeps := make([]string, 0)

	dpkgStatus := []string{"--status"}
	dpkgStatus = append(dpkgStatus, listOfDeps...)

	for _, dep := range listOfDeps {
		err := exec.Command("dpkg", dpkgStatus...).Run()
		// when the package already installed, err will be nil
		if err != nil {
			log.Println("Check dependencies: ", err)
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	if len(notInstalledDeps) > 0 {
		log.Printf("Install missing dependencies: %s", strings.Join(notInstalledDeps, " "))
		_, err := CheckSudoAccess()
		if err != nil {
			log.Fatalln("Wrong password: ", err)
		}

		cmd := exec.Command("bash", "-c", "sudo apt-get install -y "+strings.Join(notInstalledDeps, " "))
		err = cmd.Run()
		if err != nil {
			log.Println("Install dependencies: ", err)
		}
	}

	dbAccess, err := checkDBAccess()
	if err != nil {
		log.Println("DB Access: ", err)
	}

	if !dbAccess {
		isDbConfigured, err := configureDB()
		if err != nil {
			log.Fatalln("Err db configure: ", err)
		}
		log.Println("Database is configured: ", isDbConfigured)
	}

	// TODO: as for now, it should be only work for ubuntu 20.04 and 22.04
}

func checkDBAccess() (bool, error) {
	os.Setenv("PGPASSWORD", config.DB_PASSWORD)
	psqlCmd := fmt.Sprintf("psql -h %s -p %s -U %s -c 'SELECT 1'", config.DB_HOST, config.DB_PORT, config.DBUsername())

	err := exec.Command("bash", "-c", psqlCmd).Run()
	if err != nil {
		return false, err
	}

	log.Println("db can be accessed")

	return true, nil
}

func configureDB() (bool, error) {
	// TODO: Add validation first, to make sure no create command again executed
	os.Setenv("PGPASSWORD", config.DB_PASSWORD)
	psqlScript := fmt.Sprintf(`psql -c "CREATE ROLE %s SUPERUSER PASSWORD '%s';"`, config.DBUsername(), config.DB_PASSWORD)
	err := exec.Command("sudo", "su", "-", "postgres", "-c", psqlScript).Run()
	if err != nil {
		return false, err
	}

	// if db for the user already exist, there is no need to new one, so this will return error
	_ = exec.Command("createdb", config.DBUsername()).Run()
	if err != nil {
		return true, nil
	}

	return true, nil
}

func CheckSudoAccess() (bool, error) {
	err := exec.Command("sudo", "whoami").Run()
	if err != nil {
		return false, err
	}

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
	err := exec.Command("bash", "-c", "curl https://pyenv.run | bash").Run()
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
		log.Println(err)
	}

	if !pyenvInstalled {
		log.Println("Pyenv is not installed, installing pyenv")
		_, err := installPyenv()
		if err != nil {
			log.Fatalf("Failed to install pyenv: %s", err)
		}
	}

	if _, err := os.Stat(config.OdooDir()); os.IsNotExist(err) {
		err = os.MkdirAll(config.OdooDir(), config.ODOO_PERMISSION)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
