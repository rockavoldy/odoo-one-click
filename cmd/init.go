package cmd

import (
	"fmt"
	"io"
	"log"
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
	Short: "First initialization, install pyenv, configure postgresql and clone odoo",
	Long:  "First initialization for installing pyenv, configure postgresql, and clone odoo",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Verbose {
			log.SetOutput(io.Discard)
		}
		CheckRequirement()
		checkPyenv()
	},
}

func CheckRequirement() {
	// first, need to confirm if it is ubuntu or derivatives

	listOfDeps := []string{"postgresql", "postgresql-client", "libxml2-dev", "libxslt1-dev", "libldap2-dev", "libsasl2-dev", "libtiff5-dev", "libjpeg8-dev", "libopenjp2-7-dev", "zlib1g-dev", "libfreetype6-dev", "liblcms2-dev", "libwebp-dev", "libharfbuzz-dev", "libpq-dev", "git", "libsqlite3-dev", "libreadline-dev", "libbz2-dev", "tk-dev"}

	notInstalledDeps := make([]string, 0)

	dpkgStatus := utils.PrependCommand(listOfDeps, []string{"--status"})

	for _, dep := range listOfDeps {
		err := exec.Command("dpkg", dpkgStatus...).Run()
		// when the package already installed, err will be nil
		if err != nil {
			log.Println("Check dependencies: ", err)
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	if len(notInstalledDeps) > 0 {
		// When there is still missing dependencies, install them
		log.Printf("Install missing dependencies: %s", strings.Join(notInstalledDeps, " "))
		err := utils.CheckSudoAccess()
		if err != nil {
			log.Fatalln("Wrong password: ", err)
		}

		cmdAptInstall := []string{"apt-get", "install", "-y"}
		cmdAptInstall = utils.PrependCommand(notInstalledDeps, cmdAptInstall)

		err = exec.Command("sudo", cmdAptInstall...).Run()
		if err != nil {
			log.Println("Install dependencies: ", err)
		}
	}

	dbAccess, err := checkDBAccess()
	if err != nil {
		log.Println("Can't access DB: ", err)
	}

	if !dbAccess {
		err := configureDB()
		if err != nil {
			log.Fatalln("Error when configure DB: ", err)
		}
		log.Println("Database successfully configured")
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

	log.Println("Database can be accessed")

	return true, nil
}

func configureDB() error {
	// TODO: Add validation first, to make sure no create command again executed
	os.Setenv("PGPASSWORD", config.DB_PASSWORD)
	psqlScript := fmt.Sprintf(`psql -c "CREATE ROLE %s SUPERUSER LOGIN PASSWORD '%s';"`, config.DBUsername(), config.DB_PASSWORD)
	err := exec.Command("sudo", "su", "-", "postgres", "-c", psqlScript).Run()
	if err != nil {
		return err
	}

	// if db for the user already exist, there is no need to new one, so this will return error
	_ = exec.Command("createdb", config.DBUsername()).Run()
	if err != nil {
		return nil
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
