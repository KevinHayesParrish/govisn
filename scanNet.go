// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"

	"strings"

	"github.com/g3n/engine/util/logger"
	_ "github.com/mattn/go-sqlite3"

	g "github.com/gosnmp/gosnmp"
)

// SCAN_NET_VERSION is the file version number
const SCAN_NET_VERSION = "0.0.4"

/*
 * func scanNet accepts a seed cidr subnet address and discovers the network by iterating through IP Addresses.
 */
func scanNet(log *logger.Logger, cidr string, community string, params *g.GoSNMP) []ScannedRouter {

	log.Info("func scanNet version %s ", SCAN_NET_VERSION+" started.")
	log.Debug(" seed=%s", seed+" community="+community)

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
		log.Fatal(err.Error())
	}
	log.Debug(strconv.Itoa(len(subnetIPAddrs)) + " Host IP Addresses to be scanned= " + strings.Join(subnetIPAddrs, " "))

	// Query all IP Addresses in the requested CIDR subnet
	for i := 0; i < len(subnetIPAddrs); i++ {
		// get sysName and sysServices
		oids := []string{
			SYS_NAME_OID + ".0",     // sysName
			SYS_SERVICES_OID + ".0", // sysServices
		}

		fqdn := getRtrName(subnetIPAddrs[i])
		params.Target = subnetIPAddrs[i]

		// Build SNMP connection to Router
		err = params.Connect()
		if err != nil {
			log.Warn("Router not SNMP Enabled, or SNMP parameters incorrect. Continuing to scan CIDR.")
			params.Conn.Close()
			continue
		}

		result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err != nil {
			if strings.Contains(err.Error(), "request timeout") || strings.Contains(err.Error(), "connection refused") {
				log.Warn(subnetIPAddrs[i] + " not answering SNMP get. Continuing network scan.")
				continue
			} else {
				log.Fatal("Get() err: %v", err)
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

	log.Debug("Returning, scannedRouters= %v", scannedRouters)
	log.Info("func scanNet %s", SCAN_NET_VERSION+" ended.")
	return scannedRouters
}
