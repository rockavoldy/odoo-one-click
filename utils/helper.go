package utils

import (
	"fmt"
	"log"
	"odoo-one-click/config"
	"os"
	"os/exec"
	"strings"
)

func CheckSudoAccess() error {
	err := exec.Command("sudo", "whoami").Run()
	if err != nil {
		return err
	}

	return nil
}

func PrependCommand(data []string, cmd []string) []string {
	// copy the command first, then append the data
	retData := cmd
	return append(retData, data...)
}

func DirName(odooVer string, isEnterprise bool) string {
	if isEnterprise {
		return strings.Split(odooVer, ".")[0] + "e"
	}

	return strings.Split(odooVer, ".")[0] + "c"
}

func CheckDirExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func CheckPythonInstalled(pythonVer string) (bool, error) {
	bashCommand := fmt.Sprintf("pyenv versions | grep %s", pythonVer)
	err := exec.Command("bash", "-c", bashCommand).Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckVenvCreated(venv string) (bool, error) {
	bashCommand := fmt.Sprintf("pyenv virtualenvs | grep %s", venv)
	out, err := exec.Command("bash", "-c", bashCommand).Output()
	if err != nil {
		return false, err
	}

	log.Println(string(out))
	return true, nil
}

func GetPythonBasedOdooVer(odooVer string) string {
	// Because of some issue with py3.8 (gevent, cython, etc)
	// it's better to use python3.7 for odoo 11.0 to 16.0
	return "3.7.13"
	// ver, _ :=÷ strconv.Atoi(strings.Split(odooVer, ".")[0])
	// if ver < 13 {
	// 	return "3.7.13"
	// }

	// return "3.8.13"
}

func RemoveNewLine(data string) string {
	return strings.Replace(data, "\n", "", -1)
}

func IsValidDirName(dirName string) bool {
	if dirName == "" {
		return false
	}
	if strings.Contains(dirName, " ") {
		return false
	}
	if strings.Contains(dirName, "/") {
		return false
	}
	if strings.HasPrefix(dirName, ".") {
		return false
	}

	return true
}

func OdooConf(isEnterprise bool, dbUser, dbPass, dbName string) string {
	// [options]
	// admin_passwd = admin
	// db_host = localhost
	// db_port = 5432
	// db_user = %s
	// db_password = %s
	// db_name = %s
	// addons_path = ./addons, ./odoo/addons

	confFile := fmt.Sprintf(`
[options]
admin_passwd = admin
db_host = localhost
db_port = 5432
db_user = %s
db_password = %s
db_name = %s
addons_path = ./addons, ./odoo/addons`, dbUser, dbPass, dbName)
	if isEnterprise {
		confFile += ", ./enterprise\n"
	} else {
		confFile += "\n"
	}

	return confFile
}

func PyenvInfoBash() {
	rcfile := "~/.bashrc"
	if strings.Contains(os.Getenv("SHELL"), "zsh") {
		rcfile = "~/.zshrc"
	}
	fmt.Println()
	fmt.Printf("One more thing you need to do, please add this line to your %s file:\n", rcfile)
	fmt.Println("(Just copy and paste to your terminal line per line)")
	fmt.Println()
	fmt.Printf(`echo 'PYENV_ROOT="$HOME/.pyenv"' >> %s\n`, rcfile)
	fmt.Printf(`echo 'PATH="$PYENV_ROOT/bin:$PATH"' >> %s\n`, rcfile)
	fmt.Printf(`echo 'eval "$(pyenv init -)"' >> %s\n`, rcfile)
	fmt.Printf(`echo 'eval "$(pyenv virtualenv-init -)"' >> %s\n`, rcfile)
	fmt.Println()
	fmt.Println("Then run this command to reload your shell")
	fmt.Printf("source %s\n", rcfile)
}

func IsPyenvConfigured() bool {
	_, err := exec.LookPath("pyenv")
	return err == nil
}

func IsPyenvInstalled() bool {
	_, err := os.Stat(config.PyenvDir())
	return err == nil
}

func ValidateOdooVer(odooVer string) bool {
	// TODO: change this validation later by taking the data directly from github branches name, or using regex?
	// Only allow Odoo 10 and plus since it's the versions i have tried myself
	switch odooVer {
	case "10.0":
		return true
	case "11.0":
		return true
	case "12.0":
		return true
	case "13.0":
		return true
	case "14.0":
		return true
	case "15.0":
		return true
	case "16.0":
		return true
	case "17.0":
		return true
	default:
		return false
	}
}
