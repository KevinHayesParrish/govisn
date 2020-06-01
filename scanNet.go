package main

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	g "github.com/soniah/gosnmp"
)

/*
 * TODO:
 	*
*/

// SCANNETVERSION is the file version number
const SCANNETVERSION = "0.0.1"

func scanNet(debugFlag bool, cidr string, community string, params g.GoSNMP) []ScannedRouter {

	fmt.Println("\nfunc scanNet version", SCANNETVERSION, "started.")
	if debugFlag {
		fmt.Println("seed=", seed, "community=", community)
	}

	var scannedRouters []ScannedRouter

	/*
	 * TODO:
	 * Add scanning loop.
	 * SNMP Get for each IP address
	 * If Router System.Services is => IP services, then add Router FQDN and IP Address to scannedRouters

	 */

	// get all the addresses within the cidr subnet, given the input parameter.
	subnetIPAddrs, err := getHosts(cidr)
	if err != nil {
		log.Fatal(err)
	}
	if debugFlag {
		fmt.Println(len(subnetIPAddrs), "Host IP Addresses to be scanned=", subnetIPAddrs)
	}

	// Query all IP Addresses in the requested CIDR subnet
	//	err = params.Connect()
	//	if err != nil {
	//		log.Fatalf("Connect() err: %v", err)
	//	}
	//	defer params.Conn.Close()

	for i := 0; i < len(subnetIPAddrs); i++ {
		// get sysName and sysServices
		oids := []string{
			sysNameOID + ".0",     // sysName
			sysServicesOID + ".0", // sysServices
		}

		fqdn := getRtrName(subnetIPAddrs[i])
		params.Target = subnetIPAddrs[i]
		// Build SNMP connection to Router
		err = params.Connect()
		if err != nil {
			//			log.Fatalf("Connect() err: %v", err)
			fmt.Println("Router not SNMP Enabled, or SNMP parameters incorrect. Continuing to scann CIDR.")
			params.Conn.Close()
			continue
		}

		result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err != nil {
			if strings.Contains(err.Error(), "Request timeout") {
				fmt.Println(subnetIPAddrs[i], "not answering SNMP get. Continuing network scan.")
				continue
			} else {
				log.Fatalf("Get() err: %v", err)
			}
		}
		var scannedRouter ScannedRouter
		if result.Variables[1].Value.(int) >= 4 {
			scannedRouter.Name = fqdn[0]
			scannedRouter.IPAddress = subnetIPAddrs[i]
			scannedRouters = append(scannedRouters, scannedRouter)
		}
		params.Conn.Close()
	}

	if debugFlag {
		fmt.Println("Returning, scannedRouters=", scannedRouters)
	}
	fmt.Println("func scanNet", SCANNETVERSION, "ended.")
	return scannedRouters
}
