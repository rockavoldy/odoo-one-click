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
	withSudo bool
}

func NewUbuntuInitializer() *UbuntuInitializer {
	var withSudo bool
	if err := utils.CheckCmdExist("sudo"); err != nil {
		withSudo = false
	} else {
		withSudo = true
	}

	return &UbuntuInitializer{
		DefaultInitializer: &initialize.DefaultInitializer{
			Shell: utils.CurrentShell(),
		},
		withSudo: withSudo,
	}
}

func (u *UbuntuInitializer) CheckAdminAccess() error {
	var err error
	if u.withSudo {
		err = exec.Command("sudo", "whoami").Run()
	} else {
		err = exec.Command("whoami").Run()
	}

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
		var err error
		if u.withSudo {
			fmt.Println("sudo dpkg --status", dep)
			err = exec.Command("sudo", "dpkg", "--status", dep).Run()
		} else {
			fmt.Println("dpkg --status", dep)
			err = exec.Command("dpkg", "--status", dep).Run()
		}

		if err != nil {
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	if len(notInstalledDeps) > 0 {
		// apt-get update
		if u.withSudo {
			fmt.Println("sudo apt-get update")
			exec.Command("sudo", "apt-get", "update").Run()
		} else {
			fmt.Println("apt-get update")
			exec.Command("apt-get", "update").Run()
		}
	}

	return notInstalledDeps, nil
}

func (u *UbuntuInitializer) InstallDeps(deps []string) error {
	fmt.Println("Installing dependencies, please wait...")
	var err error
	if u.withSudo {
		cmdAptInstall := []string{"apt-get", "install", "-y"}
		cmdAptInstall = utils.PrependCommand(deps, cmdAptInstall)
		exec.Command("sudo", cmdAptInstall...).Run()
	} else {
		cmdAptInstall := []string{"install", "-y"}
		cmdAptInstall = utils.PrependCommand(deps, cmdAptInstall)
		exec.Command("apt-get", cmdAptInstall...).Run()
	}
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
	err = nil
	os.Setenv("PGPASSWORD", config.DbPassword())
	if u.withSudo {
		psqlScript := fmt.Sprintf(`psql -c "CREATE ROLE %s SUPERUSER LOGIN PASSWORD '%s';"`, config.DbUsername(), config.DbPassword())
		err = exec.Command("sudo", "su", "-", "postgres", "-c", psqlScript).Run()
	} else {
		psqlScript := fmt.Sprintf(`psql -c "CREATE ROLE %s SUPERUSER LOGIN PASSWORD '%s';"`, config.DbUsername(), config.DbPassword())
		err = exec.Command("su", "-", "postgres", "-c", psqlScript).Run()
	}
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
	pyenvScript := exec.Command("curl", "https://pyenv.run", "|", u.Shell)
	err := pyenvScript.Run()
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
