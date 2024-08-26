package main

import (

	//"log"

	"strconv"
	"strings"

	"github.com/g3n/engine/util/logger"
	_ "github.com/mattn/go-sqlite3"

	//g "github.com/soniah/gosnmp"
	g "github.com/gosnmp/gosnmp"
)

//var scannedRouters []ScannedRouter

// WALKROUTETABLEVERSION is the file version number
const WALKROUTETABLEVERSION = "0.1.1"

/*
func walkRouteTable(log *logger.Logger, seed string, community string, params *g.GoSNMP) []ScannedRouter {

	log.Info("\nfunc walkRouteTable version %s", WALKROUTETABLEVERSION+" started.")
	log.Debug("seed= %s", seed+" community= "+community)

	//	var scannedRouters []ScannedRouter

	// get sysName and sysServices
	//	oids := []string{
	//		sysNameOID + ".0",     // sysName
	//		sysServicesOID + ".0", // sysServices
	//	}

	fqdn := getRtrName(seed)
	params.Target = seed

	// Build SNMP connection to Router
	err := params.Connect()
	if err != nil {
		log.Warn("Router not SNMP Enabled, or SNMP parameters incorrect. Continuing to scan CIDR.")
		params.Conn.Close()
	}

	//	result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	//	if err != nil {
	//		if strings.Contains(err.Error(), "Request timeout") || strings.Contains(err.Error(), "connection refused") {
	//			log.Warn(seed + " not answering SNMP get. Continuing network scan.")
	//		} else {
	//			log.Fatal("Get() err: %v", err)
	//		}
	//	}

	var scannedRouter ScannedRouter

	// Add seed router to list of routers
	//	if result.Variables[1].Value.(int) >= 4 {
	//		scannedRouter.Name = fqdn[0]
	//		scannedRouter.IPAddress = seed
	//		scannedRouters = append(scannedRouters, scannedRouter)
	//	}
	scannedRouter.IPAddress = seed
	fqdn = getRtrName(scannedRouter.IPAddress)
	scannedRouter.Name = fqdn[0]
	scannedRouters = append(scannedRouters, scannedRouter)

	// Retrieve the route table and add each Next Hop address to the list of routers
	ipRouteNextHopPDU, err := params.WalkAll(ipRouteNextHopOID)
	if err != nil {
		log.Fatal("Get(ipRouteNextHopPDU) err")
	}
	log.Debug("\nipRouteNextHopPDU PDU= %v", ipRouteNextHopPDU)

	for i := 0; i < len(ipRouteNextHopPDU); i++ {
		scannedRouter.IPAddress = ipRouteNextHopPDU[i].Value.(string)
		fqdn = getRtrName(scannedRouter.IPAddress)
		scannedRouter.Name = fqdn[0]
		scannedRouters = append(scannedRouters, scannedRouter)
	}

	nbrOfRouters := len(scannedRouters)
	for j := 0; j < nbrOfRouters; j++ {
		params.Target = scannedRouters[j].IPAddress
		log.Debug("Recusively calling walkRouteTable." + scannedRouters[j].IPAddress + " params.Target=" + params.Target)
		scannedRouters = walkRouteTable(log, scannedRouters[j].IPAddress, community, params)
	}

	params.Conn.Close()
	log.Info("\nfunc walkRouteTable version %s ", WALKROUTETABLEVERSION+" ended.")
	return scannedRouters
}
*/

func walkRouteTableMap(log *logger.Logger, seed string, community string, params *g.GoSNMP) map[string]string {

	log.Info("func walkRouteTableMap version %s", WALKROUTETABLEVERSION+" started.")
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
	// TODO = Allow for no SNMP agent on router
	ipRouteNextHopPDU, err := params.WalkAll(ipRouteNextHopOID)

	//	if err != nil {
	//		log.Fatal("Get(ipRouteNextHopPDU) err")
	//	}
	if err != nil {
		if strings.Contains(err.Error(), "request timeout") || strings.Contains(err.Error(), "connection refused") {
			log.Warn("walkRouteTable: " + seed + " not answering SNMP get. Continue walking route table.")
			return scannedRouterMap
			//		} else {
			//			log.Fatal("Get() err: %v", err)
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
		strMaxHops, err := strconv.Atoi(*maxHops)
		if err != nil {
			log.Fatal("walkRouteTable: strconv eror.")
		}
		if walkedHops < strMaxHops {
			scannedRouterMap = walkRouteTableMap(log, IPAddresses[j], community, params)
		} else {
			log.Debug("end of recursion.")
			break
		}
	}

	params.Conn.Close()
	log.Info("\nfunc walkRouteTableMap version %s ", WALKROUTETABLEVERSION+" ended.")
	return scannedRouterMap
}
