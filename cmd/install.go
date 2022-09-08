package cmd

import (
	"log"
	"odoo-one-click/config"
	"odoo-one-click/utils"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure odoo",
	Long:  "Install and configure odoo with demo data",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(config.OdooVersion())
		installConf := NewInstallConf(false, config.OdooVersion(), "3.8.13", "", "")
		installConf.InstallOdoo()
	},
}

type InstallConf struct {
	isEnterprise bool   // is it enterprise or community
	odooVer      string // which version of odoo want to be installed
	pythonVer    string // version of python that going to run this odoo instance
	ghUsername   string // when is_enterprise is true, need to fill this username
	ghToken      string // when is_enterprise is true, need to fill this token
}

func NewInstallConf(isEnterprise bool, odooVer, pythonVer, ghUser, ghToken string) *InstallConf {
	if isEnterprise && (ghUser == "" || ghToken == "") {
		log.Fatalln("Please provide github username and token to clone odoo enterprise.")
	}
	return &InstallConf{
		isEnterprise: isEnterprise,
		odooVer:      odooVer,
		pythonVer:    pythonVer,
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

	err = ic.initPyenv()
	if err != nil {
		log.Println("Initialize pyenv: ", err)
	}
}

func (ic InstallConf) cloneOdooCommunity() error {
	log.Println("Cloning Odoo Community")
	err := os.Chdir(config.OdooDir())
	if err != nil {
		log.Println("Change directory: ", err)
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
