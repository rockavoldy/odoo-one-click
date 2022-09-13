package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
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

func CheckPythonInstalled(pythonVer string) (bool, error) {
	log.Println("Check python version if it already innstalled by pyenv")
	bashCommand := fmt.Sprintf("pyenv versions | grep %s", pythonVer)
	out, err := exec.Command("bash", "-c", bashCommand).Output()
	if err != nil {
		return false, err
	}

	log.Println(string(out))
	return true, nil
}

func CheckVenvCreated(venv string) (bool, error) {
	log.Println("Check if virtualenv with the same name already created")
	bashCommand := fmt.Sprintf("pyenv virtualenvs | grep %s", venv)
	out, err := exec.Command("bash", "-c", bashCommand).Output()
	if err != nil {
		return false, err
	}

	log.Println(string(out))
	return true, nil
}

func GetPythonBasedOdooVer(odooVer string) string {
	ver, _ := strconv.Atoi(strings.Split(odooVer, ".")[0])
	if ver < 13 {
		return "3.7.13"
	}

	return "3.8.13"
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
