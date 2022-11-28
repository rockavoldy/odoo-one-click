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

func CheckCmdExist(cmd string) error {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return fmt.Errorf("command %s not found", cmd)
	}

	return nil
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

func CurrentShell() string {
	envShell := strings.Split(os.Getenv("SHELL"), "/")
	return envShell[len(envShell)-1]
}

func PyenvInfo() {
	shell := CurrentShell()
	fmt.Println()
	fmt.Printf("One more thing you need to do, please add this line to your ~/.%src file:\n", shell)
	fmt.Println("(Just copy and paste to your terminal line per line)")
	fmt.Println()
	fmt.Printf("echo 'PYENV_ROOT=\"$HOME/.pyenv\"' >> ~/.%src\n", shell)
	fmt.Printf("echo 'PATH=\"$PYENV_ROOT/bin:$PATH\"' >> ~/.%src\n", shell)
	fmt.Printf("echo 'eval \"$(pyenv init -)\"' >> ~/.%src\n", shell)
	fmt.Printf("echo 'eval \"$(pyenv virtualenv-init -)\"' >> ~/.%src\n", shell)
	fmt.Println()
	fmt.Printf("Then run this command to reload your %src\n", shell)
	fmt.Printf("source ~/.%src\n\n", shell)
}

func CheckUbuntuVersion() error {
	out, err := exec.Command("bash", "-c", "source /etc/os-release; echo $UBUNTU_CODENAME").Output()
	if err != nil {
		return err
	}

	// check if UBUNTU_CODENAME on allowedOS
	codename := RemoveNewLine(string(out))
	if !config.IsAllowedOS(codename) {
		return fmt.Errorf("ubuntu version is not supported")
	}

	return nil
}

func IsPyenvConfigured() bool {
	_, err := exec.LookPath("pyenv")
	return err == nil
}

func IsPyenvInstalled() bool {
	_, err := os.Stat(config.PyenvDir())
	return err == nil
}
