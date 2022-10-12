package config

import (
	"os"
)

var (
	VERSION = "0.1.0"
	Verbose = false
)

const (
	ODOO_VERSION    = "15.0"
	ODOO_PERMISSION = 0755
	DB_PASSWORD     = "odoo"
	DB_HOST         = "localhost"
	DB_PORT         = "5432"
)

func OdooVersion() string {
	return ODOO_VERSION
}

func DBUsername() string {
	return "odoo"
}

func OdooDir() string {
	home, _ := os.UserHomeDir()
	return home + "/odoo"
}

func PyenvBin() string {
	home, _ := os.UserHomeDir()
	return home + "/.pyenv"
}

var allowedOS = map[string]bool{
	"focal":   true,
	"groovy":  true,
	"hirsute": true,
	"impish":  true,
	"jammy":   true,
}

func IsAllowedOS(os string) bool {
	return allowedOS[os]
}
