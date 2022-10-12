package cmd

import (
	"bufio"
	"fmt"
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
	installCmd.Flags().StringVarP(&odooVer, "odoo", "o", "", "Odoo version to install")
	installCmd.Flags().StringVarP(&dbName, "db-name", "d", "", "Database name to create or use")
	installCmd.Flags().StringVarP(&pythonVer, "python", "p", "", "Python version to use")

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [flags] directory_name",
	Short: "Install and configure odoo",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("too many arguments")
		}
		if len(args) == 1 && !utils.IsValidDirName(args[0]) {
			return fmt.Errorf("invalid directory name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if isEnterprise {
			// when enterprise is checked, ask for github username and token
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter github username: ")
			ghUsername, _ = reader.ReadString('\n')
			ghUsername = utils.RemoveNewLine(ghUsername)

			fmt.Print("Enter github token: ")
			ghToken, _ = reader.ReadString('\n')
			ghToken = utils.RemoveNewLine(ghToken)
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
		dirName := utils.DirName(odooVer, isEnterprise)

		if len(args) == 1 {
			dirName = args[0]
		}

		installConf := NewInstallConf(isEnterprise, odooVer, pythonVer, dbName, ghUsername, ghToken, dirName)
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
	dirName      string // where odoo want to be installed
}

func NewInstallConf(isEnterprise bool, odooVer, pythonVer, dbName, ghUser, ghToken, dirName string) *InstallConf {
	if isEnterprise && (ghUser == "" || ghToken == "") {
		fmt.Println("Please provide github username and token to clone odoo enterprise.")
		os.Exit(1)
	}
	return &InstallConf{
		isEnterprise: isEnterprise,
		odooVer:      odooVer,
		pythonVer:    pythonVer,
		dbName:       dbName,
		ghUsername:   ghUser,
		ghToken:      ghToken,
		dirName:      dirName,
	}
}

func (ic InstallConf) InstallOdoo() {
	Logger.Println("Installing Odoo")
	if !ic.isEnterprise {
		fmt.Printf("Cloning odoo %s community\n", ic.odooVer)
	} else {
		fmt.Printf("Cloning odoo %s enterprise\n", ic.odooVer)
	}
	err := ic.cloneOdooCommunity()
	if err != nil {
		Logger.Println("Clone odoo community: ", err)
	}

	if ic.isEnterprise {
		err = ic.cloneOdooEnterprise()
		if err != nil {
			Logger.Println("Clone odoo enterprise: ", err)
		}
	}

	err = ic.initPyenv()
	if err != nil {
		Logger.Println("Initialize pyenv: ", err)
	}

	err = ic.installOdooDeps()
	if err != nil {
		Logger.Println("Install odoo dependencies: ", err)
	}

	err = exec.Command("createdb", ic.dbName).Run()
	if err != nil {
		Logger.Println("CreateDB odoo: ", err)
	}

	err = ic.createOdooConf()
	if err != nil {
		Logger.Println("Create odoo conf: ", err)
	}

	dirName := utils.DirName(ic.odooVer, ic.isEnterprise)
	fmt.Printf("Your odoo %s is ready to use at\n%s\n", dirName, config.OdooDir()+"/"+ic.dirName)
	fmt.Println("for the first run, you need to initialize the db with base addons by add -i base")
}

func (ic InstallConf) cloneOdooCommunity() error {
	Logger.Println("Cloning Odoo Community")
	err := os.Chdir(config.OdooDir())
	if err != nil {
		Logger.Println("Change directory: ", err)
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println("Please run `odoo-one-click init` first.")
			os.Exit(1)
		}
		return err
	}

	err = exec.Command("git", "clone", "https://github.com/odoo/odoo", "--branch", ic.odooVer, "--depth", "1", ic.dirName).Run()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 128") {
			return err
		}
	}

	err = os.Chdir(config.OdooDir() + "/" + ic.dirName)
	if err != nil {
		return err
	}

	return nil
}

func (ic InstallConf) cloneOdooEnterprise() error {
	Logger.Println("Cloning Odoo Enterprise")
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
	Logger.Println("Initializing pyenv")
	isPyVerInstalled, err := utils.CheckPythonInstalled(ic.pythonVer)
	if err != nil {
		Logger.Println(err)
	}

	if !isPyVerInstalled {
		err := exec.Command("pyenv", "install", ic.pythonVer).Run()
		if err != nil {
			return err
		}
	}

	isVenvCreated, err := utils.CheckVenvCreated(ic.dirName)
	if err != nil {
		Logger.Println(err)
	}

	if !isVenvCreated {
		err = exec.Command("pyenv", "virtualenv", ic.pythonVer, ic.dirName).Run()
		if err != nil {
			Logger.Println("Error on create venv: ", err)
			return err
		}
	}

	err = exec.Command("pyenv", "local", ic.dirName).Run()
	if err != nil {
		return err
	}

	return nil
}

func (ic InstallConf) installOdooDeps() error {
	Logger.Println("Installing Odoo Dependencies")
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
	Logger.Println("Creating Odoo Configuration")
	confFile := fmt.Sprintf(`
[options]
admin_passwd = admin
db_host = localhost
db_port = 5432
db_user = %s
db_password = %s
db_name = %s
addons_path = ./addons, ./odoo/addons`, config.DBUsername(), config.DB_PASSWORD, ic.dbName)

	if ic.isEnterprise {
		// When it's enterprise, add enterprise addons path to odoo.conf
		confFile = confFile + ", ./enterprise\n"
	} else {
		confFile = confFile + "\n"
	}

	err := os.WriteFile("odoo.conf", []byte(confFile), 0644)
	if err != nil {
		return err
	}

	return nil
}
