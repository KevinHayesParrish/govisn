package main

import (
	"fmt"
)

/*
 * TODO:
 	*
*/

// SCANNETVERSION is the file version number
const SCANNETVERSION = "0.0.1"

func scanNet(debugFlag bool, snmpTarget string, community string) ([]string, []string) {

	fmt.Println("\nfunc scanNet version", SCANNETVERSION, "started.")

	var RouterNames, RouterIPs []string

	return RouterNames, RouterIPs
}
