package mac

import (
	"fmt"
	"odoo-one-click/internal/initialize"
	"odoo-one-click/utils"
	"os/exec"
)

type MacInitializer struct {
	*initialize.DefaultInitializer
}

func (m *MacInitializer) CheckAdminAccess() error {
	err := exec.Command("sudo", "whoami").Run()
	if err != nil {
		return err
	}

	return nil
}

func (m *MacInitializer) CheckRequirement() ([]string, error) {
	macDeps := []string{"zlib", "libjpeg", "openssl@1.1"}
	fmt.Println("Checking requirement...")
	fmt.Println("Check if homebrew already installed in the system")
	err := exec.Command("brew", "--version").Run()
	if err != nil {
		return nil, fmt.Errorf("homebrew not installed, please install homebrew first from https://brew.sh/")
	}

	fmt.Println("Check if the dependencies already installed in the system")
	notInstalledDeps := make([]string, 0)
	for _, dep := range macDeps {
		err := exec.Command("brew", "ls", "--versions", dep).Run()
		if err != nil {
			notInstalledDeps = append(notInstalledDeps, dep)
		}
	}

	return notInstalledDeps, nil
}

func (m *MacInitializer) InstallDeps(deps []string) error {
	fmt.Println("Installing dependencies, please wait...")
	utils.PrependCommand(deps, []string{"install"})

	err := exec.Command("brew", deps...).Run()
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	if err := m.macQuirk(); err != nil {
		return err
	}

	return nil
}

// macQuirk will be called in the end of installation process,
// to fix some issues when compile python libraries on Mac M1
func (m *MacInitializer) macQuirk() error {
	fmt.Println("Installing mac quirk, please wait...")
	err := exec.Command("brew", "link", "zlib", "--force").Run()
	if err != nil {
		return fmt.Errorf("failed to link zlib: %w", err)
	}

	err = exec.Command("brew", "link", "libjpeg", "--force").Run()
	if err != nil {
		return fmt.Errorf("failed to link libjpeg: %w", err)
	}

	return nil
}

func (m *MacInitializer) CheckDB() error {
	return nil
}

func (m *MacInitializer) ConfigureDB() error {
	return nil
}
