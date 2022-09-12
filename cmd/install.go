package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"odoo-one-click/config"
	"odoo-one-click/utils"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var isEnterprise bool
var odooVer string
var dbName string
var pythonVer string
var ghUsername string
var ghToken string

func init() {
	installCmd.Flags().BoolVarP(&isEnterprise, "enterprise", "e", false, "Install odoo enterprise")
	installCmd.Flags().StringVarP(&odooVer, "odoo-version", "o", "", "Odoo version to install")
	installCmd.Flags().StringVarP(&dbName, "db-name", "d", "", "Database name to create or use")
	installCmd.Flags().StringVarP(&pythonVer, "python-version", "p", "", "Python version to use")

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure odoo",
	Long:  "Install and configure odoo with demo data",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Verbose {
			// TODO: create 1 log file for the project to use; can be extended to log to a file
			log.SetOutput(io.Discard)
		}

		if isEnterprise {
			// when enterprise is checked, ask for github username and token
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter github username: ")
			ghUsername, _ = reader.ReadString('\n')
			ghUsername = strings.Replace(ghUsername, "\n", "", -1)

			fmt.Print("Enter github token: ")
			ghToken, _ = reader.ReadString('\n')
			ghToken = strings.Replace(ghToken, "\n", "", -1)
		}

		if odooVer == "" {
			odooVer = config.OdooVersion()
		}

		if dbName == "" {
			dbName = utils.DirName(odooVer, isEnterprise)
		}
		if pythonVer == "" {
			pythonVer = utils.GetPythonBasedOdooVer(odooVer)
		}

		installConf := NewInstallConf(isEnterprise, odooVer, pythonVer, dbName, ghUsername, ghToken)
		installConf.InstallOdoo()
	},
}

type InstallConf struct {
	isEnterprise bool   // is it enterprise or community
	odooVer      string // which version of odoo want to be installed
	pythonVer    string // version of python that going to run this odoo instance
	dbName       string // Database name for odoo
	ghUsername   string // when is_enterprise is true, need to fill this username
	ghToken      string // when is_enterprise is true, need to fill this token
}

func NewInstallConf(isEnterprise bool, odooVer, pythonVer, dbName, ghUser, ghToken string) *InstallConf {
	if isEnterprise && (ghUser == "" || ghToken == "") {
		log.Fatalln("Please provide github username and token to clone odoo enterprise.")
	}
	return &InstallConf{
		isEnterprise: isEnterprise,
		odooVer:      odooVer,
		pythonVer:    pythonVer,
		dbName:       dbName,
		ghUsername:   ghUser,
		ghToken:      ghToken,
	}
}

func (ic InstallConf) InstallOdoo() {
	log.Println("Installing Odoo")
	err := ic.cloneOdooCommunity()
	if err != nil {
		log.Println("Clone odoo community: ", err)
	}

	if ic.isEnterprise {
		err = ic.cloneOdooEnterprise()
		if err != nil {
			log.Println("Clone odoo enterprise: ", err)
		}
	}

	err = ic.initPyenv()
	if err != nil {
		log.Println("Initialize pyenv: ", err)
	}

	err = ic.installOdooDeps()
	if err != nil {
		log.Println("Install odoo dependencies: ", err)
	}

	err = exec.Command("createdb", ic.dbName).Run()
	if err != nil {
		log.Println("CreateDB odoo: ", err)
	}

	err = ic.createOdooConf()
	if err != nil {
		log.Println("Create odoo conf: ", err)
	}
}

func (ic InstallConf) cloneOdooCommunity() error {
	log.Println("Cloning Odoo Community")
	err := os.Chdir(config.OdooDir())
	if err != nil {
		log.Println("Change directory: ", err)
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println("Please run `odoo-one-click init` first.")
			os.Exit(1)
		}
		return err
	}

	dirName := utils.DirName(ic.odooVer, ic.isEnterprise)

	err = exec.Command("git", "clone", "https://github.com/odoo/odoo", "--branch", ic.odooVer, "--depth", "1", dirName).Run()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 128") {
			return err
		}
	}

	err = os.Chdir(config.OdooDir() + "/" + dirName)
	if err != nil {
		return err
	}

	return nil
}

func (ic InstallConf) cloneOdooEnterprise() error {
	log.Println("Cloning Odoo Enterprise")
	enterpriseUrl := fmt.Sprintf("https://%s:%s@github.com/odoo/enterprise", ic.ghUsername, ic.ghToken)
	err := exec.Command("git", "clone", enterpriseUrl, "--branch", ic.odooVer, "--depth", "1").Run()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 128") {
			return err
		}
	}

	return nil
}

func (ic InstallConf) initPyenv() error {
	log.Println("Initializing pyenv")
	isPyVerInstalled, err := utils.CheckPythonInstalled(ic.pythonVer)
	if err != nil {
		log.Println(err)
	}

	if !isPyVerInstalled {
		err := exec.Command("pyenv", "install", ic.pythonVer).Run()
		if err != nil {
			return err
		}
	}

	dirName := utils.DirName(ic.odooVer, ic.isEnterprise)
	isVenvCreated, err := utils.CheckVenvCreated(dirName)
	if err != nil {
		log.Println(err)
	}

	if !isVenvCreated {
		err = exec.Command("pyenv", "virtualenv", ic.pythonVer, dirName).Run()
		if err != nil {
			log.Println("Error on create venv: ", err)
			return err
		}
	}

	err = exec.Command("pyenv", "local", dirName).Run()
	if err != nil {
		return err
	}

	return nil
}

func (ic InstallConf) installOdooDeps() error {
	log.Println("Installing Odoo Dependencies")
	err := exec.Command("pip", "install", "-r", "requirements.txt").Run()
	if err != nil {
		return err
	}

	// handle this issue https://github.com/odoo/odoo/issues/99809, because it seems affect odoo 13 and up
	err = exec.Command("pip", "install", "pyopenssl==22.0.0").Run()
	if err != nil {
		return err
	}

	return nil
}

func (ic InstallConf) createOdooConf() error {
	log.Println("Creating Odoo Configuration")
	confFile := fmt.Sprintf(`
[options]
admin_passwd = admin
db_host = localhost
db_port = 5432
db_user = %s
db_password = %s
db_name = %s
addons_path = ./addons, ./odoo/addons
`, config.DBUsername(), config.DB_PASSWORD, ic.dbName)

	if ic.isEnterprise {
		// When it's enterprise, add enterprise addons path to odoo.conf
		confFile = confFile + ", ./enterprise"
	}

	err := ioutil.WriteFile("odoo.conf", []byte(confFile), 0644)
	if err != nil {
		return err
	}

	return nil
}
