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

// WALK_ROUTE_TABLE_VERSION is the file version number
const WALK_ROUTE_TABLE_VERSION = "0.1.2"

/*
 * func walkRouteTableMap walks the router's ipRouteTable and returns a map of the results.
 */
func walkRouteTableMap(log *logger.Logger, seed string, community string, params *g.GoSNMP) map[string]string {

	var walkedHops = 0

	log.Info("func walkRouteTableMap version %s", WALK_ROUTE_TABLE_VERSION+" started.")
	log.Debug("seed=%s", seed)
	log.Debug("community=%s", community)
	log.Debug("params.snmpTarget=%s", params.Target)

	fqdn := getRtrName(seed)
	params.Target = seed

	// Build SNMP connection to Router
	err := params.Connect()
	if err != nil {
		log.Warn("Router not SNMP Enabled, or SNMP parameters incorrect. Continuing to scan CIDR.")
		params.Conn.Close()
	}

	scannedRouterMap := make(map[string]string)

	// Add seed router to list of routers
	fqdn = getRtrName(seed)
	log.Debug("fqdn=%v", fqdn)
	scannedRouterMap[fqdn[0]] = seed

	// Retrieve the route table and add each Next Hop address to the list of routers
	ipRouteNextHopPDU, err := params.WalkAll(IP_ROUTE_NEXT_HOP_OID)
	if err != nil {
		if strings.Contains(err.Error(), "request timeout") || strings.Contains(err.Error(), "connection refused") {
			log.Warn("walkRouteTable: " + seed + " not answering SNMP get. Continue walking route table.")
			return scannedRouterMap
		}
		log.Fatal("Get() err: %v", err)
	}
	log.Debug("\nipRouteNextHopPDU PDU= %v", ipRouteNextHopPDU)

	for i := 0; i < len(ipRouteNextHopPDU); i++ {
		fqdn = getRtrName(ipRouteNextHopPDU[i].Value.(string))
		scannedRouterMap[fqdn[0]] = ipRouteNextHopPDU[i].Value.(string)
	}

	// retrieve IP Addresses of scanned routers
	var IPAddresses []string
	for _, IPAddress := range scannedRouterMap {
		IPAddresses = append(IPAddresses, IPAddress)
	}

	nbrOfRouters := len(scannedRouterMap)
	for j := 0; j < nbrOfRouters; j++ {
		params.Target = IPAddresses[j]
		log.Debug("Recusively calling walkRouteTable." + IPAddresses[j] + " params.Target=" + params.Target)
		walkedHops++
		intMaxHops, err := strconv.Atoi(*maxHops)
		if err != nil {
			log.Fatal("walkRouteTable: strconv eror.")
		}
		if walkedHops < intMaxHops {
			scannedRouterMap = walkRouteTableMap(log, IPAddresses[j], community, params)
		} else {
			log.Debug("end of recursion.")
			break
		}
	}

	params.Conn.Close()
	log.Info("\nfunc walkRouteTableMap version %s ", WALK_ROUTE_TABLE_VERSION+" ended.")
	return scannedRouterMap
}
