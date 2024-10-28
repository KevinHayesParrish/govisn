// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"

	"net"
	"strconv"
	"strings"

	"github.com/g3n/engine/util/logger"
	g "github.com/gosnmp/gosnmp"
)

// DISCOVERY_VERSION is the file version number
const DISCOVERY_VERSION = "0.3.9"

/*
 * func discover  discovers the network, constrained by input parm maximum hops away from snmpTarget node.
 * It returns the database.
 */
func discover(log *logger.Logger, snmpTarget string, params *g.GoSNMP, maxHopsStr string, database *sql.DB) *sql.DB {

	log.Info("func discover version %s started.", DISCOVERY_VERSION)

	// Discover network, constrained by input parm maximum hops away from snmpTarget node
	maxHops, _ := strconv.Atoi(maxHopsStr)
	for i := 0; i <= maxHops; i++ {
		log.Debug("Discover iteration")
	}

	// Get Router attributes
	var router Router

	log.Debug("params= %v", params)

	err := params.Connect()
	if err != nil {
		database.Close()
		log.Fatal("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	// Initialize the database
	database = initDB(log, database)

	router = getRouterInfo(log, snmpTarget, params, router, database)

	log.Debug("func discover version %s", DISCOVERY_VERSION+" ended.")

	return database
}

/*
 * func initDB initializes the database with its tables.
 */
func initDB(log *logger.Logger, database *sql.DB) *sql.DB {

	initDbVersion := "0.0.4"
	log.Debug("initDB version %s", initDbVersion)

	/*
	 *	Add Routers table to DB
	 */
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY UNIQUE, Name TEXT, Description TEXT, UpTime TEXT, Contact TEXT, Location TEXT, Services INTEGER, GpsLat REAL, GpsLong REAL, GpsAlt REAL)")
	if err != nil {
		database.Close()
		log.Fatal("Router Table Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add RouteTable table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, IPRouteIfIndex TEXT, NextHop TEXT)")
	if err != nil {
		database.Close()
		log.Fatal("RouteTable Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add RouterIP table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER NOT NULL, IpAddr TEXT, IfIndex TEXT)")
	if err != nil {
		database.Close()
		log.Fatal("RouterIP Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add RouterMac table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
	if err != nil {
		database.Close()
		log.Fatal("RouterMac Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add Links table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER NOT NULL UNIQUE, FromRouterID INTEGER, FromRouterName TEXT, FromRouterIP TEXT, FromRouterIfIndex Text, ToRouterID INTEGER, ToRouterName TEXT, ToRouterIP TEXT)")
	if err != nil {
		database.Close()
		log.Fatal("Links Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	return database
}

func getRtrName(ipAddr string) []string {
	names, err := net.LookupAddr(ipAddr)
	if err != nil {
		log.Warn("No reverse lookup found for %s", ipAddr)
	}

	if len(names) > 0 {
		return names
	} else {
		//unknown := []string{"Unknown"}
		unknown := []string{ipAddr}
		return unknown
	}
}

/*
 * func getRouterInfo uses SNMP to retrieve the router's system information
 * and writes it to the database.
 */
func getRouterInfo(log *logger.Logger, snmpTarget string, params *g.GoSNMP, router Router, database *sql.DB) Router {
	log.Debug("Starting discover.getRouterInfo.")

	oids := []string{
		SYS_NAME_OID + ".0",     // sysName
		SYS_DESCR_OID + ".0",    // sysDescr
		SYS_UPTIME_OID + ".0",   // sysUpTime
		SYS_CONTACT_OID + ".0",  // sysContact
		SYS_LOCATION_OID + ".0", // sysLocation
		SYS_SERVICES_OID + ".0", // sysServices
	}

	// get FQDN with IP Address
	fqdn := getRtrName(snmpTarget)
	log.Debug("Return from getRtrName=%s", fqdn)

	var routerSupportsSNMP bool
	result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "refused") {
			log.Warn("Get() err %s", err.Error()+
				"\n Router "+snmpTarget+" not responding to SNMP Get. Continuing with network discovery.")

			router.System.Name = fqdn[0]
			router.System.Description = ""
			router.System.UpTime = 0
			router.System.Contact = ""
			router.System.Location = ""
			router.System.Services = 0
			routerSupportsSNMP = false
		} else {
			database.Close()
			log.Fatal("Get() err %s", err.Error())
		}
	} else {
		routerSupportsSNMP = true
		//router.System.Name = fqdn[0]
		router.System.Name = string(result.Variables[0].Value.([]byte))
		router.System.Description = string(result.Variables[1].Value.([]byte))
		router.System.UpTime = result.Variables[2].Value.(uint32)
		router.System.Contact = string(result.Variables[3].Value.([]byte))
		router.System.Location = string(result.Variables[4].Value.([]byte))
		router.System.Services = result.Variables[5].Value.(int)

	}

	// get FQDN with IP Address
	//fqdn := getRtrName(snmpTarget)

	//	router.System.Name = fqdn[0]
	//	router.System.Description = string(result.Variables[1].Value.([]byte))
	//	router.System.UpTime = result.Variables[2].Value.(uint32)
	//	router.System.Contact = string(result.Variables[3].Value.([]byte))
	//	router.System.Location = string(result.Variables[4].Value.([]byte))
	//	router.System.Services = result.Variables[5].Value.(int)

	/*
		// Retrieve GPS data from DNS
	*/

	// get GPS data from DNS
	router.System.GPS.Latitude = "0.0"  // initialze with float data to allow for missing GPS on DB
	router.System.GPS.Longitude = "0.0" // initialze with float data to allow for missing GPS on DB
	router.System.GPS.Altitude = "0.0"  // initialze with float data to allow for missing GPS on DB

	if len(fqdn) > 0 {
		/*
		 * Use router's hostname for DNS query, instead of fqdn[0].
		 * This provides a consistent IP Address, fqdn[0] can be
		 * an interface address relating to a different DNS name than
		 * hostname.
		 */
		//gpsDNS := getGPS(fqdn[0])
		gpsDNS := getGPS(router.System.Name)

		for n := 0; n < len(gpsDNS); n++ {
			s := gpsDNS[n]
			// Split TXT record into prefix and value
			sr := strings.Split(s, "=")
			if sr[0] == "Long" {
				router.System.GPS.Longitude = sr[1]
			}
			if sr[0] == "Lat" {
				router.System.GPS.Latitude = sr[1]
			}
			if sr[0] == "Alt" {
				router.System.GPS.Altitude = sr[1]
			}
		}
	}

	log.Debug("router.System.Name= %s", router.System.Name)
	log.Debug("router.System.Description= %s", router.System.Description)
	log.Debug("router.System.UpTime= %d", router.System.UpTime)
	log.Debug("router.System.Contact= %s", router.System.Contact)
	log.Debug("router.System.Location= %s", router.System.Location)
	log.Debug("router.System.Services= %d", router.System.Services)
	log.Debug("router.System.GPS= %v", router.System.GPS)

	// Write Router row to database
	statement, _ := database.Prepare("INSERT INTO Routers (RouterID, Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	statement.Exec()
	defer statement.Close()

	Name := router.System.Name

	ifPhysAddress1, err := getIfPhysAddress(log, snmpTarget, params)
	var RouterIDUint32 uint32
	if err != nil {
		log.Warn("Router %s has no ifPhysAddress.1", snmpTarget)
		RouterIDUint32 = crc32.ChecksumIEEE([]byte(Name))
		log.Debug("Calculated RouterID using %s as: %d", Name, RouterIDUint32)
	} else {
		RouterIDUint32 = crc32.ChecksumIEEE([]byte(ifPhysAddress1))
		log.Debug("Calculated RouterID using %s as: %d", ifPhysAddress1, RouterIDUint32)
	}

	router.System.RouterID = int(RouterIDUint32)
	Description := router.System.Description
	UpTime := router.System.UpTime
	Contact := router.System.Contact
	Location := router.System.Location
	Services := router.System.Services
	GpsLat := router.System.GPS.Latitude
	GpsLong := router.System.GPS.Longitude
	GpsAlt := router.System.GPS.Altitude

	routerIsInDB := false
	//	_, err = statement.Exec(strconv.Itoa(int(RouterIDUint32)), Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt) // Add router
	_, err = statement.Exec(strconv.Itoa(router.System.RouterID), Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt) // Add router
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			log.Warn("Router %s", Name+" already exists in database. Continuing discovery.")
			routerIsInDB = true
		} else {
			database.Close()
			log.Fatal("RouterMac Insert Exec err: %v", err)
		}
	}
	defer statement.Close()

	if !routerIsInDB && routerSupportsSNMP {
		getInterfaces(log, params, router, database)

		getIPAddresses(log, params, router, database)

		getIPRouteTable(log, params, router, database)
	}

	log.Debug("Ended discover.getRouterInfo.")
	return router
}

/*
 * func getInterfaces uses SNMP to retrieve the router's interfaces.
 */
func getInterfaces(log *logger.Logger, params *g.GoSNMP, router Router, database *sql.DB) {

	// get Number of Interfaces
	ifNumberArray := []string{IF_NUMBER_OID + ".0"}
	getPDU, getError := params.Get(ifNumberArray)
	if getError != nil {
		database.Close()
		log.Fatal("Get() err: %s", getError.Error())
	}
	log.Debug("ifNumber walkPDU= %v", getPDU)

	nbrOfInterfaces := getPDU.Variables[0].Value.(int)
	log.Debug("nbrOfInterfaces= %v", nbrOfInterfaces)

	// get ifTable
	walkPDU, walkError := params.WalkAll(IF_TABLE_OID)
	if walkError != nil {
		database.Close()
		log.Fatal("Get() err: %v", walkError)
	}
	log.Debug("\nifTable PDU= %v", walkPDU)

	var interfaceTable ifTable

	log.Debug("len(walkPDU)= %d", len(walkPDU))

	for i := 0; i < len(walkPDU); i++ { // skip ifIndex array within walkPDU
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifIndexOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifIndexType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifIndex = walkPDU[i].Value.(int)
			log.Debug("ifIndex= %d", interfaceTable.ifEntry.ifIndex)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifDescrOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifDescrType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifDescr = string(walkPDU[i].Value.([]uint8))
			log.Debug("ifDescr= %s", interfaceTable.ifEntry.ifDescr)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifTypeOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifTypeType = byte(walkPDU[i].Type)
			ifTypeInt := walkPDU[i].Value.(int)
			switch ifTypeInt {
			case 1:
				interfaceTable.ifEntry.ifType = "other(1)"
			case 6:
				interfaceTable.ifEntry.ifType = "ethernetCsmacd(6)"
			default:
				interfaceTable.ifEntry.ifType = "other"
			}
			log.Debug("ifType= %s", interfaceTable.ifEntry.ifType)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifMtuOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifMtuType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifMtu = walkPDU[i].Value.(int)
			log.Debug("ifMtu= %d", interfaceTable.ifEntry.ifMtu)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifSpeedOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifSpeedType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifSpeed = walkPDU[i].Value.(uint)
			log.Debug("ifSpeed= %v", interfaceTable.ifEntry.ifSpeed)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifPhysAddressOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifPhysAddressType = byte(walkPDU[i].Type)
			physAddrUint := walkPDU[i].Value.([]byte)
			var physAddrHex [6]string

			for l := 0; l < 6; l++ {
				if len(physAddrUint) == 0 {
					interfaceTable.ifEntry.ifPhysAddress = "Null0"
					break
				}
				physAddrUint8 := physAddrUint[l]
				physAddrHex[l] = fmt.Sprintf("%x", physAddrUint8)
			}
			if len(physAddrUint) > 0 {
				interfaceTable.ifEntry.ifPhysAddress = physAddrHex[0] +
					":" +
					physAddrHex[1] +
					":" +
					physAddrHex[2] +
					":" +
					physAddrHex[3] +
					":" +
					physAddrHex[4] +
					":" +
					physAddrHex[5]
			} else {
				interfaceTable.ifEntry.ifPhysAddress = "Null0"
			}

			log.Debug("ifPhysAddress= %s", interfaceTable.ifEntry.ifPhysAddress)

			writeMacToDB(log, router, interfaceTable, database)

			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifAdminStatusOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifAdminStatusType = byte(walkPDU[i].Type)
			ifAdminStatusInt := walkPDU[i].Value.(int)
			switch ifAdminStatusInt {
			case 1:
				interfaceTable.ifEntry.ifAdminStatus = "up(1)"
			case 2:
				interfaceTable.ifEntry.ifAdminStatus = "down(2)"
			case 3:
				interfaceTable.ifEntry.ifAdminStatus = "testing(3)"
			}
			log.Debug("ifAdminStatus= %s", interfaceTable.ifEntry.ifAdminStatus)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOperStatusOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOperStatusType = byte(walkPDU[i].Type)
			ifAdminStatusInt := walkPDU[i].Value.(int)
			switch ifAdminStatusInt {
			case 1:
				interfaceTable.ifEntry.ifOperStatus = "up(1)"
			case 2:
				interfaceTable.ifEntry.ifOperStatus = "down(2)"
			case 3:
				interfaceTable.ifEntry.ifOperStatus = "testing(3)"
			}
			log.Debug("ifOperStatus= %s", interfaceTable.ifEntry.ifOperStatus)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifLastChangeOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifLastChangeType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifLastChange = walkPDU[i].Value.(uint32)
			log.Debug("ifLastChange= %d", interfaceTable.ifEntry.ifLastChange)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInOctetsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInOctetsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInOctets = walkPDU[i].Value.(uint)
			log.Debug("ifInOctets= %d", interfaceTable.ifEntry.ifInOctets)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInUcastPkts = walkPDU[i].Value.(uint)
			log.Debug("ifInucastPkts= %d", interfaceTable.ifEntry.ifInUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInNUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInNUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInNUcastPkts = walkPDU[i].Value.(uint)
			log.Debug("ifINUcastPkts= %d", interfaceTable.ifEntry.ifInUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInDiscardsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInDiscardsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInDiscards = walkPDU[i].Value.(uint)
			log.Debug("ifDiscards= %d", interfaceTable.ifEntry.ifInDiscards)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInErrorsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInErrorsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInErrors = walkPDU[i].Value.(uint)
			log.Debug("ifInErrors= %d", interfaceTable.ifEntry.ifInErrors)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInUnknownProtosOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInUnknownProtosType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInUnknownProtos = walkPDU[i].Value.(uint)
			log.Debug("ifInUnknownProtos= %d", interfaceTable.ifEntry.ifInUnknownProtos)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutOctetsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutOctetsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutOctets = walkPDU[i].Value.(uint)
			log.Debug("ifOutOctets= %d", interfaceTable.ifEntry.ifOutOctets)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutUcastPkts = walkPDU[i].Value.(uint)
			log.Debug("ifOutUcastPkts= %d", interfaceTable.ifEntry.ifOutUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutNUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutNUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutNUcastPkts = walkPDU[i].Value.(uint)
			log.Debug("ifOutNUcastPkts= %d", interfaceTable.ifEntry.ifOutNUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutDiscardsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutDiscardsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutDiscards = walkPDU[i].Value.(uint)
			log.Debug("ifOutDiscards= %d", interfaceTable.ifEntry.ifOutDiscards)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutErrorsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutErrorsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutErrors = walkPDU[i].Value.(uint)
			log.Debug("ifOutErrors= %d", interfaceTable.ifEntry.ifOutErrors)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutQLenOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutQLenType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutQLen = walkPDU[i].Value.(uint)
			log.Debug("ifOutQLen= %d", interfaceTable.ifEntry.ifOutQLen)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifSpecificOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifSpecificType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifSpecific = walkPDU[i].Value.(string)
			log.Debug("ifSpecific= %s", interfaceTable.ifEntry.ifSpecific)
			i++
		}

	} // End of Interfaces code

}

func getHostIP(routerName string) []string {
	addrs, err := net.LookupHost(routerName)
	if err != nil {
		log.Warn("No Hostname lookup found for %s", routerName)
	}
	if len(addrs) == 0 {
		log.Warn("No Hostname records for %s", routerName)
		return addrs
	}
	return addrs
}

func getGPS(sysName string) []string {
	txts, err := net.LookupTXT(sysName)
	if err != nil {
		log.Debug("No TXT records for %s", sysName)
	}
	if len(txts) == 0 {
		//		fmt.Println("No DNS TXT records for", sysName)
		log.Debug("No DNS TXT records for %s", sysName)
	}
	return txts
}

/*
 * func getIPAddresses uses SNMP to retrieve the router's IP Address information
 * and writes it to the database
 */
func getIPAddresses(log *logger.Logger, params *g.GoSNMP, router Router, database *sql.DB) {

	// get ipAddrTable
	walkPDU, err := params.WalkAll(IP_AD_ENT_ADDR_OID)
	if err != nil {
		database.Close()
		log.Fatal("Get(walkPDU) err: %v", err)
	}
	ifIndexPDU, err := params.WalkAll(IP_AD_ENT_IF_INDEX)
	if err != nil {
		database.Close()
		log.Fatal("Get(ifIndexPDU) err: %v", err)
	}

	log.Debug("\nipAdEntAddr PDU= %v", walkPDU)
	log.Debug("\nifIndex PDU= %v", ifIndexPDU)

	var ipTable ipAddrTable

	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntAddr = walkPDU[i].Value.(string)
		ipTable.ipAddrEntry.ipAdEntIfIndex = ifIndexPDU[i].Value.(int)

		log.Debug("ipAdEntAddr= %s", ipTable.ipAddrEntry.ipAdEntAddr)

		// Add row to RouterIp table
		statement, err := database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr, IfIndex) VALUES (?, ?, ?)")
		if err != nil {
			database.Close()
			log.Fatal("RouterIp Prepare Insert Exec err: %v", err)
		}
		//RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
		//		_, err = statement.Exec(RouterID, ipTable.ipAddrEntry.ipAdEntAddr, ipTable.ipAddrEntry.ipAdEntIfIndex)
		_, err = statement.Exec(router.System.RouterID, ipTable.ipAddrEntry.ipAdEntAddr, ipTable.ipAddrEntry.ipAdEntIfIndex)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				// Continue executing if this is a duplicate IP Address. This assume this router is being processed again.
				// In case this is a duplicate MAC Address within the network, print error output to stdoutput.
				//				fmt.Println("\n****\n Non-Unique IP Address", ipTable.ipAddrEntry.ipAdEntAddr, "\n This may be because this router is being re-discovered.\n If not, then this is a serious network violation condition.\n****")
				log.Warn("\n****\n Non-Unique IP Address %s", ipTable.ipAddrEntry.ipAdEntAddr+
					"\n This may be because this router is being re-discovered."+
					"\n If not, then this is a serious network violation condition.\n****")
			} else {
				database.Close()
				log.Fatal("RouterIp Exec Insert Exec err: %v", err)
			}
		}
		defer statement.Close()

	}

	walkPDU, err = params.WalkAll(IP_AD_ENT_NET_MASK)
	if err != nil {
		log.Debug("WalkAll failed. Err= %s", err)
	}

	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntNetMask = walkPDU[i].Value.(string)
		log.Debug("ipAdEntNetMask= %s", ipTable.ipAddrEntry.ipAdEntNetMask)
	}

	walkPDU, err = params.WalkAll(IP_AD_ENT_BCAST_ADDR)
	if err != nil {
		log.Debug("WalkAll failed. Err= %s", err)
	}

	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntBcastAddr = walkPDU[i].Value.(int)
		log.Debug("ipAdEntBcastAddr= %d", ipTable.ipAddrEntry.ipAdEntBcastAddr)
	}

	walkPDU, err = params.WalkAll(IP_AD_ENT_REASM_MAX_SIZE)
	if err != nil {
		log.Debug("WalkAll failed. Err= %s", err)
	}

	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntReasmMaxSize = walkPDU[i].Value.(int)
		log.Debug("ipAdEntReasmMaxSize= %d", ipTable.ipAddrEntry.ipAdEntReasmMaxSize)
	}
}

func getIfPhysAddress(log *logger.Logger, snmpTarget string, params *g.GoSNMP) (string, error) {
	log.Debug("getIfPhysAddress for %s", snmpTarget)

	oids := []string{
		IF_PHYS_ADDRESS_OID + ".1",
	}
	result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "refused") {
			log.Warn("Get() err %s", err.Error()+
				"\n Router "+snmpTarget+" not responding to SNMP Get. Continuing with network discovery.")

			return "", err
		} else {
			log.Fatal("Router %s has no ifPhysAddress.1", snmpTarget)
		}
	}

	ifPhysAddress1 := string(result.Variables[0].Value.([]byte))

	return ifPhysAddress1, err

}

/*
 * fun getIPRouteTable uses SNMP to retrieve the router's route table information
 * and writes it to the database
 */
func getIPRouteTable(log *logger.Logger, params *g.GoSNMP, router Router, database *sql.DB) {

	// get ipRouteTable
	ipRouteDestPDU, err := params.WalkAll(IP_ROUTE_DEST_OID)
	if err != nil {
		database.Close()
		log.Fatal("Get(ipRouteDestPDU) err")
	}
	log.Debug("\nipRouteDestPDU PDU= %v", ipRouteDestPDU)

	ipRouteIfIndexPDU, err := params.WalkAll(IP_ROUTE_IF_INDEX_OID)
	if err != nil {
		database.Close()
		log.Fatal("Get(ipRouteIfIndexPDU) err")
	}
	log.Debug("\nipRouteIfIndexPDU PDU= %v", ipRouteIfIndexPDU)

	ipRouteNextHopPDU, err := params.WalkAll(IP_ROUTE_NEXT_HOP_OID)
	if err != nil {
		database.Close()
		log.Fatal("Get(ipRouteNextHopPDU) err")
	}
	log.Debug("\nipRouteNextHopPDU PDU= %v", ipRouteNextHopPDU)

	// Parse Dest and NextHop PDUs, adding row to ipRouteTable for each PDU element.
	var ipRouteTab ipRouteTable

	for i := 0; i < (len(ipRouteDestPDU)); i++ {
		ipRouteTab.ipRouteEntry.ipRouteDest = ipRouteDestPDU[i].Value.(string)
		ipRouteTab.ipRouteEntry.ipRouteIfIndex = ipRouteIfIndexPDU[i].Value.(int)
		ipRouteTab.ipRouteEntry.ipRouteNextHop = ipRouteNextHopPDU[i].Value.(string)

		log.Debug("ipRouteDest= %s", ipRouteTab.ipRouteEntry.ipRouteDest)
		log.Debug("ipRouteIfIndex= %d", ipRouteTab.ipRouteEntry.ipRouteIfIndex)
		log.Debug("ipRouteNextHop= %s", ipRouteTab.ipRouteEntry.ipRouteNextHop)

		// Add row to RouteTable table
		statement, _ := database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, IPRouteIfIndex, NextHop) VALUES (?, ?, ?, ?)")
		if err != nil {
			database.Close()
			log.Fatal("RouterTable Prepare Insert Exec err")
		}

		ifPhysAddress1, err := getIfPhysAddress(log, params.Target, params)
		var RouterID uint32
		if err != nil {
			log.Warn("Router %s has no ifPhysAddress.1", params.Target)
			RouterID = crc32.ChecksumIEEE([]byte(router.System.Name))
		} else {
			RouterID = crc32.ChecksumIEEE([]byte(ifPhysAddress1))
		}

		//RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))

		_, err = statement.Exec(RouterID, ipRouteTab.ipRouteEntry.ipRouteDest, ipRouteTab.ipRouteEntry.ipRouteIfIndex, ipRouteTab.ipRouteEntry.ipRouteNextHop)
		if err != nil {
			database.Close()
			log.Fatal("RouteTable Insert err")
		}
		defer statement.Close()
	}

}

/*
 * func writeMactoDB writes the router's MAC address information to the RouterMac table
 */
func writeMacToDB(log *logger.Logger, router Router, interfaceTable ifTable, database *sql.DB) {

	statement, err := database.Prepare("INSERT INTO RouterMac (RouterID, MacAddr) VALUES (?, ?)")
	if err != nil {
		database.Close()
		log.Fatal("RouterMac Insert Prepare err: %v", err)
	}
	defer statement.Close()

	// RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
	// _, err = statement.Exec(strconv.Itoa(int(RouterID)), interfaceTable.ifEntry.ifPhysAddress)
	_, err = statement.Exec(router.System.RouterID, interfaceTable.ifEntry.ifPhysAddress)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// Continue executing if this is a duplicate MAC Address. This assume this router is being processed again.
			// In case this is a duplicate MAC Address within the network, print error output to stdoutput.
			//			fmt.Println("\n****\n Non-Unique MAC Address", interfaceTable.ifEntry.ifPhysAddress, "\n This may be because this router is being re-discovered.\n If not, then this is a serious network violation condition.\n****")
			log.Warn("\n****\n Non-Unique MAC Address %s", interfaceTable.ifEntry.ifPhysAddress+
				"\n This may be because this router is being re-discovered."+
				"\n If not, then this is a serious network violation condition.\n****")
		} else {
			database.Close()
			log.Fatal("RouterMac Insert Exec err: %v", err)
		}
	}
	defer statement.Close()
}
