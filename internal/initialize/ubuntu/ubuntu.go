package ubuntu

import (
	"fmt"
	"odoo-one-click/config"
	"odoo-one-click/internal/initialize"
	"odoo-one-click/utils"
	"os"
	"os/exec"
)

type UbuntuInitializer struct {
	*initialize.DefaultInitializer
}

func NewUbuntuInitializer() *UbuntuInitializer {
	return &UbuntuInitializer{
		DefaultInitializer: &initialize.DefaultInitializer{
			Shell: utils.CurrentShell(),
		},
	}
}

func (u *UbuntuInitializer) CheckAdminAccess() error {
	err := exec.Command("sudo", "whoami").Run()
	if err != nil {
		return err
	}

	return nil
}

func (u *UbuntuInitializer) CheckRequirement() ([]string, error) {
	ubuntuDeps := []string{"build-essential", "postgresql", "postgresql-client", "libxml2-dev", "libssl-dev", "libffi-dev", "libxslt1-dev", "libldap2-dev", "libsasl2-dev", "libtiff5-dev", "libjpeg8-dev", "libopenjp2-7-dev", "zlib1g-dev", "libfreetype6-dev", "liblcms2-dev", "libwebp-dev", "libharfbuzz-dev", "libpq-dev", "git", "libsqlite3-dev", "libreadline-dev", "libbz2-dev", "tk-dev"}

	fmt.Println("Checking requirement...")
	fmt.Println("Check if dependencies already installed in the system")

	notInstalledDeps := make([]string, 0)
	for _, dep := range ubuntuDeps {
		err := exec.Command("dpkg", "--status", dep).Run()
		if err != nil {
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	return notInstalledDeps, nil
}

func (u *UbuntuInitializer) InstallDeps(deps []string) error {
	fmt.Println("Installing dependencies, please wait...")
	cmdAptInstall := []string{"apt-get", "install", "-y"}
	cmdAptInstall = utils.PrependCommand(deps, cmdAptInstall)
	err := exec.Command("sudo", cmdAptInstall...).Run()
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

func (u *UbuntuInitializer) CheckDB() error {
	err := utils.CheckCmdExist("psql")
	if err != nil {
		return err
	}

	// Set PGPASSWORD env so it can connect to postgresql without any prompt
	os.Setenv("PGPASSWORD", config.DbPassword())

	psqlCmd := fmt.Sprintf("psql -h %s -p %s -U %s postgres -c 'SELECT 1'", config.DbHost(), config.DbPort(), config.DbUsername())
	err = exec.Command(u.Shell, "-c", psqlCmd).Run()
	if err != nil {
		return fmt.Errorf("can't connect to postgresql: %w", err)
	}

	return nil
}

func (u *UbuntuInitializer) ConfigureDB() error {
	err := u.CheckDB()
	if err == nil {
		return nil
	}
	os.Setenv("PGPASSWORD", config.DbPassword())
	psqlScript := fmt.Sprintf(`psql -c "CREATE ROLE %s SUPERUSER LOGIN PASSWORD '%s';"`, config.DbUsername(), config.DbPassword())
	err = exec.Command("sudo", "su", "-", "postgres", "-c", psqlScript).Run()
	if err != nil {
		return fmt.Errorf("failed to create postgre role: %w", err)
	}

	err = u.CheckDB()
	if err != nil {
		return err
	}

	return nil
}

func (u *UbuntuInitializer) CheckPyenv() error {
	err := utils.CheckCmdExist("pyenv")
	if err != nil {
		return err
	}

	err = exec.Command("pyenv", "--version").Run()
	if err != nil {
		utils.PyenvInfo()
		return err
	}

	return nil
}

func (u *UbuntuInitializer) InstallPyenv() error {
	pyenvScript := exec.Command("curl", "https://pyenv.run")
	runBash := exec.Command(u.Shell)
	runBash.Stdin, _ = pyenvScript.StdoutPipe()
	runBash.Stdout = os.Stdout
	_ = runBash.Start()
	_ = pyenvScript.Run()
	err := runBash.Wait()
	if err != nil {
		return err
	}

	err = exec.Command(u.Shell, "-c", "exec $SHELL").Run()
	if err != nil {
		return err
	}

	return nil
}

func (u *UbuntuInitializer) ConfigurePyenv() error {
	if err := u.CheckPyenv(); err != nil {
		err = u.InstallPyenv()
		if err != nil {
			utils.PyenvInfo()
		}
	}

	return nil
}
