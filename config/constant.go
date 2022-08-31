package config

import (
	"os"
)

const (
	VERSION         = "0.1.0"
	ODOO_VERSION    = "15.0"
	ODOO_PERMISSION = 0755
	DB_PASSWORD     = "odoo"
	DB_HOST         = "localhost"
	DB_PORT         = "5432"
)

func DBUsername() string {
	return "odoo"
	// return os.Getenv("USER")
}

func OdooDir() string {
	home, _ := os.UserHomeDir()
	return home + "/odoo"
}
