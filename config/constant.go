package config

import (
	"log"
	"os"
	"strings"
)

const (
	VERSION         = "0.1.0"
	ODOO_VERSION    = "15.0"
	ODOO_PERMISSION = 0755
	DB_PASSWORD     = "odoo"
	DB_HOST         = "localhost"
	DB_PORT         = "5432"
)

func OdooVersion() string {
	if ODOOVER := os.Getenv("ODOO_VERSION"); ODOOVER != "" {
		if !strings.Contains(ODOOVER, ".") {
			log.Fatalf("It seems ODOO_VERSION = %s is not correct, please change with the branch name of odoo. Ex: 15.0\n", ODOOVER)
		}
		return os.Getenv("ODOO_VERSION")
	}

	return ODOO_VERSION
}

func DBUsername() string {
	return "odoo"
	// return os.Getenv("USER")
}

func OdooDir() string {
	home, _ := os.UserHomeDir()
	return home + "/odoo"
}
