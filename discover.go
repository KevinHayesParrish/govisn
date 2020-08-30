// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"

	//"log"
	"net"
	"strconv"
	"strings"

	"github.com/g3n/engine/util/logger"
	g "github.com/soniah/gosnmp"
)

/*
 * TODO:
 	* Add code to walk route table to discover all routers
*/

// DISCOVERYVERSION is the file version number
const DISCOVERYVERSION = "0.3.6"

func discover(debugFlag bool, log *logger.Logger, dbName string, snmpTarget string, community string, params *g.GoSNMP, maxHopsStr string, database *sql.DB) *sql.DB {

	//	fmt.Println("\nfunc discover version", DISCOVERYVERSION, "started.")
	log.Info("\nfunc discover version %s", DISCOVERYVERSION+" started.")

	// Discover network, constrained by input parm maximum hops away from snmpTarget node
	maxHops, _ := strconv.Atoi(maxHopsStr)
	for i := 0; i <= maxHops; i++ {
		//		if debugFlag {
		//			fmt.Println("Discover iteration")
		//		}
		log.Debug("Discover iteration")
	}

	// Get Router attributes
	var router Router

	log.Debug("params= %v", params)

	err := params.Connect()
	if err != nil {
		log.Fatal("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	// Initialize the database
	database = initDB(debugFlag, log, database)

	getRouterInfo(debugFlag, log, snmpTarget, community, maxHopsStr, params, router, database)

	log.Debug("func discovery version %s", DISCOVERYVERSION+" ended.")

	return database
}

func getInterfaces(debugFlag bool, log *logger.Logger, snmpTarget string, community string, maxHopsStr string, params *g.GoSNMP, router Router, database *sql.DB) {

	// get Number of Interfaces
	ifNumberArray := []string{ifNumberOID + ".0"}
	getPDU, getError := params.Get(ifNumberArray)
	if getError != nil {
		//		log.Fatalf("Get() err: %v", getError)
		log.Fatal("Get() err: %v", getError)
	}
	//	if debugFlag {
	//		fmt.Println("ifNumber walkPDU=", getPDU)
	//	}
	log.Debug("ifNumber walkPDU= %v", getPDU)

	nbrOfInterfaces := getPDU.Variables[0].Value.(int)
	//	if debugFlag {
	//		fmt.Println("nbrOfInterfaces =", nbrOfInterfaces)
	//	}
	log.Debug("nbrOfInterfaces= %v", nbrOfInterfaces)

	// get ifTable
	walkPDU, walkError := params.WalkAll(ifTableOID)
	if walkError != nil {
		//		log.Fatalf("Get() err: %v", walkError)
		log.Fatal("Get() err: %v", walkError)
	}
	//	if debugFlag {
	//		fmt.Println("\nifTable PDU=", walkPDU)
	//	}
	log.Debug("\nifTable PDU= %v", walkPDU)

	var interfaceTable ifTable

	//	if debugFlag {
	//		fmt.Println("len(walkPDU)=", len(walkPDU))
	//	}
	log.Debug("len(walkPDU)= %d", len(walkPDU))

	for i := 0; i < len(walkPDU); i++ { // skip ifIndex array within walkPDU
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifIndexOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifIndexType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifIndex = walkPDU[i].Value.(int)
			//			if debugFlag {
			//				fmt.Println("ifIndex=", interfaceTable.ifEntry.ifIndex)
			//			}
			log.Debug("ifIndex= %d", interfaceTable.ifEntry.ifIndex)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifDescrOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifDescrType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifDescr = string(walkPDU[i].Value.([]uint8))
			//			if debugFlag {
			//				fmt.Println("ifDescr=", interfaceTable.ifEntry.ifDescr)
			//			}
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
			//			if debugFlag {
			//				fmt.Println("ifType=", interfaceTable.ifEntry.ifType)
			//			}
			log.Debug("ifType= %s", interfaceTable.ifEntry.ifType)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifMtuOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifMtuType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifMtu = walkPDU[i].Value.(int)
			//			if debugFlag {
			//				fmt.Println("ifMtu=", interfaceTable.ifEntry.ifMtu) // TESTING ONLY
			//			}
			log.Debug("ifMtu= %d", interfaceTable.ifEntry.ifMtu)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifSpeedOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifSpeedType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifSpeed = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifSpeed=", interfaceTable.ifEntry.ifSpeed)
			//			}
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

			//			if debugFlag {
			//				fmt.Println("ifPhysAddress=", interfaceTable.ifEntry.ifPhysAddress)
			//			}
			log.Debug("ifPhysAddress= %s", interfaceTable.ifEntry.ifPhysAddress)

			//			writeMacToDB(debugFlag, router, interfaceTable, database)
			writeMacToDB(debugFlag, log, router, interfaceTable, database)

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
			//			if debugFlag {
			//				fmt.Println("ifAdminStatus=", interfaceTable.ifEntry.ifAdminStatus)
			//			}
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
			//			if debugFlag {
			//				fmt.Println("ifOperStatus=", interfaceTable.ifEntry.ifOperStatus)
			//			}
			log.Debug("ifOperStatus= %s", interfaceTable.ifEntry.ifOperStatus)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifLastChangeOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifLastChangeType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifLastChange = walkPDU[i].Value.(uint32)
			//			if debugFlag {
			//				fmt.Println("ifLastChange=", interfaceTable.ifEntry.ifLastChange)
			//			}
			log.Debug("ifLastChange= %d", interfaceTable.ifEntry.ifLastChange)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInOctetsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInOctetsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInOctets = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifInOctets=", interfaceTable.ifEntry.ifInOctets)
			//			}
			log.Debug("ifInOctets= %d", interfaceTable.ifEntry.ifInOctets)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInUcastPkts = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifInucastPkts=", interfaceTable.ifEntry.ifInUcastPkts)
			//			}
			log.Debug("ifInucastPkts= %d", interfaceTable.ifEntry.ifInUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInNUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInNUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInNUcastPkts = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifINUcastPkts=", interfaceTable.ifEntry.ifInUcastPkts)
			//			}
			log.Debug("ifINUcastPkts= %d", interfaceTable.ifEntry.ifInUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInDiscardsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInDiscardsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInDiscards = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifDiscards=", interfaceTable.ifEntry.ifInDiscards)
			//			}
			log.Debug("ifDiscards= %d", interfaceTable.ifEntry.ifInDiscards)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInErrorsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInErrorsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInErrors = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifInErrors=", interfaceTable.ifEntry.ifInErrors)
			//			}
			log.Debug("ifInErrors= %d", interfaceTable.ifEntry.ifInErrors)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInUnknownProtosOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInUnknownProtosType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInUnknownProtos = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifInUnknownProtos=", interfaceTable.ifEntry.ifInUnknownProtos)
			//			}
			log.Debug("ifInUnknownProtos= %d", interfaceTable.ifEntry.ifInUnknownProtos)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutOctetsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutOctetsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutOctets = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifOutOctets=", interfaceTable.ifEntry.ifOutOctets)
			//			}
			log.Debug("ifOutOctets= %d", interfaceTable.ifEntry.ifOutOctets)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutUcastPkts = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifOutUcastPkts=", interfaceTable.ifEntry.ifOutUcastPkts)
			//			}
			log.Debug("ifOutUcastPkts= %d", interfaceTable.ifEntry.ifOutUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutNUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutNUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutNUcastPkts = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifOutNUcastPkts=", interfaceTable.ifEntry.ifOutNUcastPkts)
			//			}
			log.Debug("ifOutNUcastPkts= %d", interfaceTable.ifEntry.ifOutNUcastPkts)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutDiscardsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutDiscardsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutDiscards = walkPDU[i].Value.(uint)
			//			if debugFlag
			//				fmt.Println("ifOutDiscards=", interfaceTable.ifEntry.ifOutDiscards)
			//			}
			log.Debug("ifOutDiscards= %d", interfaceTable.ifEntry.ifOutDiscards)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutErrorsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutErrorsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutErrors = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifOutErrors=", interfaceTable.ifEntry.ifOutErrors)
			//			}
			log.Debug("ifOutErrors= %d", interfaceTable.ifEntry.ifOutErrors)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutQLenOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutQLenType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutQLen = walkPDU[i].Value.(uint)
			//			if debugFlag {
			//				fmt.Println("ifOutQLen=", interfaceTable.ifEntry.ifOutQLen)
			//			}
			log.Debug("ifOutQLen= %d", interfaceTable.ifEntry.ifOutQLen)
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifSpecificOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifSpecificType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifSpecific = walkPDU[i].Value.(string)
			//			if debugFlag {
			//				fmt.Println("ifSpecific=", interfaceTable.ifEntry.ifSpecific)
			//			}
			log.Debug("ifSpecific= %s", interfaceTable.ifEntry.ifSpecific)
			i++
		}

	} // End of Interfaces code

}

func initDB(debugFlag bool, log *logger.Logger, database *sql.DB) *sql.DB {

	/*
	* TODO:
	 */

	initDbVersion := "0.0.3"
	//	if debugFlag {
	//		fmt.Println("initDB version:", initDbVersion)
	//	}
	log.Debug("initDB version %s", initDbVersion)

	/*
	 *	Add Routers table to DB
	 */
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY UNIQUE, Name TEXT, Description TEXT, UpTime TEXT, Contact TEXT, Location TEXT, Services INTEGER, GpsLat REAL, GpsLong REAL, GpsAlt REAL)")
	if err != nil {
		//		log.Fatalf("Router Table Create err: %v", err)
		log.Fatal("Router Table Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add RouteTable table to DB
	 */
	//	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, NextHop TEXT)")
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, IPRouteIfIndex TEXT, NextHop TEXT)")
	if err != nil {
		//		log.Fatalf("RouteTable Create err: %v", err)
		log.Fatal("RouteTable Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add RouterIP table to DB
	 */
	//	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER NOT NULL, IpAddr TEXT)")
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER NOT NULL, IpAddr TEXT, IfIndex TEXT)")
	if err != nil {
		//		log.Fatalf("RouterIP Create err: %v", err)
		log.Fatal("RouterIP Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add RouterMac table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
	if err != nil {
		//		log.Fatalf("RouterMac Create err: %v", err)
		log.Fatal("RouterMac Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	/*
	 *	Add Links table to DB
	 */
	//	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER NOT NULL UNIQUE, FromRouterName TEXT, FromRouterIP TEXT, ToRouterName TEXT, ToRouterIP TEXT)")
	//	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER NOT NULL UNIQUE, FromRouterID INTEGER, FromRouterName TEXT, FromRouterIP TEXT, ToRouterID INTEGER, ToRouterName TEXT, ToRouterIP TEXT)")
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER NOT NULL UNIQUE, FromRouterID INTEGER, FromRouterName TEXT, FromRouterIP TEXT, FromRouterIfIndex Text, ToRouterID INTEGER, ToRouterName TEXT, ToRouterIP TEXT)")
	if err != nil {
		//		log.Fatalf("Links Create err: %v", err)
		log.Fatal("Links Create err: %v", err)
	}
	defer statement.Close()
	statement.Exec()

	return database
}

func getRtrName(ipAddr string) []string {
	names, err := net.LookupAddr(ipAddr)
	if err != nil {
		//		fmt.Println("No reverse lookup found for", ipAddr)
		log.Warn("No reverse lookup found for %s", ipAddr)
	}
	if len(names) == 0 {
		//		fmt.Println("No FQDN records for", ipAddr)
		log.Warn("No FQDN records for %s", ipAddr)
		return names
	}
	return names
}

func getHostIP(routerName string) []string {
	addrs, err := net.LookupHost(routerName)
	if err != nil {
		//		fmt.Println("No Hostname lookup found for", routerName)
		log.Warn("No Hostname lookup found for %s", routerName)
	}
	if len(addrs) == 0 {
		//		fmt.Println("No Hostname records for", routerName)
		log.Warn("No Hostname records for %s", routerName)
		return addrs
	}
	return addrs
}

func getGPS(sysName string) []string {
	txts, err := net.LookupTXT(sysName)
	if err != nil {
		//		panic(err)
		//		fmt.Println("No TXT records for", sysName)
		log.Debug("No TXT records for %s", sysName)
	}
	if len(txts) == 0 {
		//		fmt.Println("No DNS TXT records for", sysName)
		log.Debug("No DNS TXT records for %s", sysName)
	}
	return txts
}

func writeMacToDB(debugFlag bool, log *logger.Logger, router Router, interfaceTable ifTable, database *sql.DB) {

	statement, err := database.Prepare("INSERT INTO RouterMac (RouterID, MacAddr) VALUES (?, ?)")
	if err != nil {
		//		log.Fatalf("RouterMac Insert Prepare err: %v", err)
		//		log.Fatal(err)
		log.Fatal("RouterMac Insert Prepare err: %v", err)
	}
	defer statement.Close()

	RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
	_, err = statement.Exec(strconv.Itoa(int(RouterID)), interfaceTable.ifEntry.ifPhysAddress)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// Continue executing if this is a duplicate MAC Address. This assume this router is being processed again.
			// In case this is a duplicate MAC Address within the network, print error output to stdoutput.
			//			fmt.Println("\n****\n Non-Unique MAC Address", interfaceTable.ifEntry.ifPhysAddress, "\n This may be because this router is being re-discovered.\n If not, then this is a serious network violation condition.\n****")
			log.Warn("\n****\n Non-Unique MAC Address %s", interfaceTable.ifEntry.ifPhysAddress+
				"\n This may be because this router is being re-discovered."+
				"\n If not, then this is a serious network violation condition.\n****")
		} else {
			//			fmt.Printf("RouterMac Insert Exec err: %v", err)
			//			log.Fatal(err)
			log.Fatal("RouterMac Insert Exec err: %v", err)
		}
	}
	defer statement.Close()
}

func getIPAddresses(debugFlag bool, log *logger.Logger, params *g.GoSNMP, router Router, database *sql.DB) {
	// get ipAddrTable
	walkPDU, err := params.WalkAll(ipAdEntAddrOID)
	if err != nil {
		//		log.Fatalf("Get(walkPDU) err: %v", err)
		log.Fatal("Get(walkPDU) err: %v", err)
	}
	ifIndexPDU, err := params.WalkAll(ipAdEntIfIndex)
	if err != nil {
		//		log.Fatalf("Get(ifIndexPDU) err: %v", err)
		log.Fatal("Get(ifIndexPDU) err: %v", err)
	}
	//	if debugFlag {
	//		fmt.Println("\nipAdEntAddr PDU=", walkPDU)
	//		fmt.Println("\nifIndex PDU=", ifIndexPDU)
	//	}
	log.Debug("\nipAdEntAddr PDU= %v", walkPDU)
	log.Debug("\nifIndex PDU= %v", ifIndexPDU)

	var ipTable ipAddrTable

	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntAddr = walkPDU[i].Value.(string)
		ipTable.ipAddrEntry.ipAdEntIfIndex = ifIndexPDU[i].Value.(int)
		//		if debugFlag {
		//			fmt.Println("ipAdEntAddr=", ipTable.ipAddrEntry.ipAdEntAddr)
		//		}
		log.Debug("ipAdEntAddr= %s", ipTable.ipAddrEntry.ipAdEntAddr)

		// Add row to RouterIp table
		//		statement, err := database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr) VALUES (?, ?)")
		statement, err := database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr, IfIndex) VALUES (?, ?, ?)")
		if err != nil {
			//			fmt.Printf("RouterIp Prepare Insert Exec err: %v", err)
			//			log.Fatal(err)
			log.Fatal("RouterIp Prepare Insert Exec err: %v", err)
		}
		RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
		_, err = statement.Exec(RouterID, ipTable.ipAddrEntry.ipAdEntAddr, ipTable.ipAddrEntry.ipAdEntIfIndex)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				// Continue executing if this is a duplicate IP Address. This assume this router is being processed again.
				// In case this is a duplicate MAC Address within the network, print error output to stdoutput.
				//				fmt.Println("\n****\n Non-Unique IP Address", ipTable.ipAddrEntry.ipAdEntAddr, "\n This may be because this router is being re-discovered.\n If not, then this is a serious network violation condition.\n****")
				log.Warn("\n****\n Non-Unique IP Address %s", ipTable.ipAddrEntry.ipAdEntAddr+
					"\n This may be because this router is being re-discovered."+
					"\n If not, then this is a serious network violation condition.\n****")
			} else {
				//				fmt.Printf("RouterIp Exec Insert Exec err: %v", err)
				//				log.Fatal(err)
				log.Fatal("RouterIp Exec Insert Exec err: %v", err)
			}
		}
		defer statement.Close()

	}

	walkPDU, err = params.WalkAll(ipAdEntNetMask)
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntNetMask = walkPDU[i].Value.(string)
		//		if debugFlag {
		//			fmt.Println("ipAdEntNetMask=", ipTable.ipAddrEntry.ipAdEntNetMask)
		//		}
		log.Debug("ipAdEntNetMask= %s", ipTable.ipAddrEntry.ipAdEntNetMask)
	}
	walkPDU, err = params.WalkAll(ipAdEntBcastAddr)
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntBcastAddr = walkPDU[i].Value.(int)
		//		if debugFlag {
		//			fmt.Println("ipAdEntBcastAddr=", ipTable.ipAddrEntry.ipAdEntBcastAddr)
		//		}
		log.Debug("ipAdEntBcastAddr= %d", ipTable.ipAddrEntry.ipAdEntBcastAddr)
	}
	walkPDU, err = params.WalkAll(ipAdEntReasmMaxSize)
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntReasmMaxSize = walkPDU[i].Value.(int)
		//		if debugFlag {
		//			fmt.Println("ipAdEntReasmMaxSize=", ipTable.ipAddrEntry.ipAdEntReasmMaxSize)
		//		}
		log.Debug("ipAdEntReasmMaxSize= %d", ipTable.ipAddrEntry.ipAdEntReasmMaxSize)
	}
}

func getIPRouteTable(debugFlag bool, log *logger.Logger, params *g.GoSNMP, router Router, database *sql.DB) {

	// get ipRouteTable
	ipRouteDestPDU, err := params.WalkAll(ipRouteDestOID)
	if err != nil {
		//		log.Fatalf("Get(ipRouteDestPDU) err: %v", err)
		log.Fatal("Get(ipRouteDestPDU) err")
	}
	//	if debugFlag {
	//		fmt.Println("\nipRouteDestPDU PDU=", ipRouteDestPDU)
	//	}
	log.Debug("\nipRouteDestPDU PDU= %v", ipRouteDestPDU)

	ipRouteIfIndexPDU, err := params.WalkAll(ipRouteIfIndexOID)
	if err != nil {
		//		log.Fatalf("Get(ipRouteIfIndexPDU) err: %v", err)
		log.Fatal("Get(ipRouteIfIndexPDU) err")
	}
	//	if debugFlag {
	//		fmt.Println("\nipRouteIfIndexPDU PDU=", ipRouteIfIndexPDU)
	//	}
	log.Debug("\nipRouteIfIndexPDU PDU= %v", ipRouteIfIndexPDU)

	ipRouteNextHopPDU, err := params.WalkAll(ipRouteNextHopOID)
	if err != nil {
		//		log.Fatalf("Get(ipRouteNextHopPDU) err: %v", err)
		log.Fatal("Get(ipRouteNextHopPDU) err")
	}
	//	if debugFlag {
	//		fmt.Println("\nipRouteNextHopPDU PDU=", ipRouteNextHopPDU)
	//	}
	log.Debug("\nipRouteNextHopPDU PDU= %v", ipRouteNextHopPDU)

	// Parse Dest and NextHop PDUs, adding row to ipRouteTable for each PDU element.
	var ipRouteTab ipRouteTable

	for i := 0; i < (len(ipRouteDestPDU)); i++ {
		ipRouteTab.ipRouteEntry.ipRouteDest = ipRouteDestPDU[i].Value.(string)
		ipRouteTab.ipRouteEntry.ipRouteIfIndex = ipRouteIfIndexPDU[i].Value.(int)
		ipRouteTab.ipRouteEntry.ipRouteNextHop = ipRouteNextHopPDU[i].Value.(string)
		//		if debugFlag {
		//			fmt.Println("ipRouteDest=", ipRouteTab.ipRouteEntry.ipRouteDest)
		//			fmt.Println("ipRouteIfIndex=", ipRouteTab.ipRouteEntry.ipRouteIfIndex)
		//			fmt.Println("ipRouteNextHop=", ipRouteTab.ipRouteEntry.ipRouteNextHop)
		//		}
		log.Debug("ipRouteDest= %s", ipRouteTab.ipRouteEntry.ipRouteDest)
		log.Debug("ipRouteIfIndex= %d", ipRouteTab.ipRouteEntry.ipRouteIfIndex)
		log.Debug("ipRouteNextHop= %s", ipRouteTab.ipRouteEntry.ipRouteNextHop)

		// Add row to RouteTable table
		statement, _ := database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, IPRouteIfIndex, NextHop) VALUES (?, ?, ?, ?)")
		if err != nil {
			//			fmt.Printf("RouterTable Prepare Insert Exec err: %v", err)
			//			log.Fatal(err)
			log.Fatal("RouterTable Prepare Insert Exec err")
		}

		RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
		_, err = statement.Exec(RouterID, ipRouteTab.ipRouteEntry.ipRouteDest, ipRouteTab.ipRouteEntry.ipRouteIfIndex, ipRouteTab.ipRouteEntry.ipRouteNextHop)
		if err != nil {
			//			log.Fatalf("RouteTable Insert err: %v", err)
			log.Fatal("RouteTable Insert err")
		}
		defer statement.Close()
	}

}

func getRouterInfo(debugFlag bool, log *logger.Logger, snmpTarget string, community string, maxHopsStr string, params *g.GoSNMP, router Router, database *sql.DB) {
	oids := []string{
		sysNameOID + ".0",     // sysName
		sysDescrOID + ".0",    // sysDescr
		sysUpTimeOID + ".0",   // sysUpTime
		sysContactOID + ".0",  // sysContact
		sysLocationOID + ".0", // sysLocation
		sysServicesOID + ".0", // sysServices
	}

	result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		//		log.Fatalf("Get() err: %v", err)
		//		log.Fatal("Get() err %s", err.Error())
		log.Fatal("Get() err %s", err.Error()+"\nRouter not responding to SNMP Get. Continuing with network discovery.")
	}

	// get FQDN with IP Address
	fqdn := getRtrName(snmpTarget)

	router.System.Name = fqdn[0]
	router.System.Description = string(result.Variables[1].Value.([]byte))
	router.System.UpTime = result.Variables[2].Value.(uint32)
	router.System.Contact = string(result.Variables[3].Value.([]byte))
	router.System.Location = string(result.Variables[4].Value.([]byte))
	router.System.Services = result.Variables[5].Value.(int)

	/*
		// Retrieve GPS data from DNS
	*/

	// get GPS data from DNS
	router.System.GPS.Latitude = "0.0"  // initialze with float data to allow for missing GPS on DB
	router.System.GPS.Longitude = "0.0" // initialze with float data to allow for missing GPS on DB
	router.System.GPS.Altitude = "0.0"  // initialze with float data to allow for missing GPS on DB

	if len(fqdn) > 0 {
		gpsDNS := getGPS(fqdn[0])
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

	//	if debugFlag {
	//		fmt.Println("router.System.Name=", router.System.Name)
	//		fmt.Println("router.System.Description=", router.System.Description)
	//		fmt.Println("router.System.UpTime=", router.System.UpTime)
	//		fmt.Println("router.System.Contact=", router.System.Contact)
	//		fmt.Println("router.System.Location=", router.System.Location)
	//		fmt.Println("router.System.Services=", router.System.Services)
	//		fmt.Println("router.System.GPS=", router.System.GPS)
	//	}
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
	RouterIDUint32 := crc32.ChecksumIEEE([]byte(Name))
	Description := router.System.Description
	UpTime := router.System.UpTime
	Contact := router.System.Contact
	Location := router.System.Location
	Services := router.System.Services
	GpsLat := router.System.GPS.Latitude
	GpsLong := router.System.GPS.Longitude
	GpsAlt := router.System.GPS.Altitude

	routerIsInDB := false
	_, err = statement.Exec(strconv.Itoa(int(RouterIDUint32)), Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt) // Add router
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			//			fmt.Println("Router", Name, "is already exists in database. Continuing discovery.")
			log.Warn("Router %s", Name+" already exists in database. Continuing discovery.")
			routerIsInDB = true
		} else {
			//			fmt.Printf("RouterMac Insert Exec err: %v", err)
			log.Fatal("RouterMac Insert Exec err: %v", err)
		}
	}
	defer statement.Close()

	if !routerIsInDB {
		//		getInterfaces(debugFlag, snmpTarget, community, maxHopsStr, params, router, database)
		getInterfaces(debugFlag, log, snmpTarget, community, maxHopsStr, params, router, database)

		getIPAddresses(debugFlag, log, params, router, database)

		getIPRouteTable(debugFlag, log, params, router, database)
	}
}
