package main

import (
	"odoo-one-click/cmd"
)

func main() {
	// TODO: add validation to check if the OS is ubuntu or derivatives
	// on elementary, there is UPSTREAM lsb-release, but it seemms not on all derivatives
	cmd.Execute()
}
