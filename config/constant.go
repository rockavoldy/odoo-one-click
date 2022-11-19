package config

import (
	"io/fs"
	"os"
)

var (
	VERSION        = "0.1.0"
	Verbose        = false
	odooVersion    = "15.0"
	odooPermission = 0755
	odooDir        = ""
	dbUsername     = "odoo"
	dbPassword     = "odoo"
	dbHost         = "127.0.0.1"
	dbPort         = "5432"
)

func init() {
	SetOdooVersion("15.0")
	SetOdooDir("/odoo")
	SetDbUsername("odoo")
	SetDbPassword("odoo")
	SetDbHost("127.0.0.1")
	SetDbPort("5432")
}

// odooVersion
func OdooVersion() string {
	return odooVersion
}

func SetOdooVersion(odooVer string) {
	odooVersion = odooVer
}

// odooPermission
func OdooPermission() os.FileMode {
	return fs.FileMode(odooPermission)
}

// odooDir
func OdooDir() string {
	return odooDir
}

func SetOdooDir(dir string) {
	home, _ := os.UserHomeDir()
	odooDir = home + dir
}

// dbUsername
func DbUsername() string {
	return dbUsername
}

func SetDbUsername(dbUser string) {
	dbUsername = dbUser
}

// dbPassword
func DbPassword() string {
	return dbPassword
}

func SetDbPassword(dbPass string) {
	dbPassword = dbPass
}

// dbHost
func DbHost() string {
	return dbHost
}

func SetDbHost(host string) {
	dbHost = host
}

// dbPort
func DbPort() string {
	return dbPort
}

func SetDbPort(port string) {
	dbPort = port
}

func PyenvDir() string {
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
