package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	g "github.com/soniah/gosnmp"
)

/*
 * TODO:
 	* Add code to walk route table to discover all routers
 	* Write Links row to database
*/

// DISCOVERYVERSION is the file version number
const DISCOVERYVERSION = "0.3.2"

//func discover(debugFlag bool, dbName string, snmpTarget string, community string, maxHopsStr string) {
func discover(debugFlag bool, dbName string, snmpTarget string, community string, maxHopsStr string) *sql.DB {

	fmt.Println("\nfunc discover version", DISCOVERYVERSION, "started.")

	// Discover network, constrained by input parm maximum hops away from snmpTarget node
	maxHops, _ := strconv.Atoi(maxHopsStr)
	for i := 0; i <= maxHops; i++ {
		if debugFlag {
			fmt.Println("Discover iteration")
		}
	}

	// Get Router attributes
	var router Router

	snmpPort := "161"
	if len(snmpTarget) <= 0 {
		log.Fatalf("environment variable not set: GOSNMP_TARGET")
	} else {
		if debugFlag {
			fmt.Println("snmpTarget=", snmpTarget)
		}
	}
	if len(snmpPort) <= 0 {
		log.Fatalf("environment variable not set: GOSNMP_PORT")
	}
	port, _ := strconv.ParseUint(snmpPort, 10, 16)

	// GoSNMP struct
	params := &g.GoSNMP{
		Target:    snmpTarget,
		Port:      uint16(port),
		Community: community,
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Logger:    nil,
		//		MaxOids:   5,
		MaxOids: 6,
	}

	if debugFlag {
		fmt.Println("params=", params)
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	// Initialize the database
	database, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("sql.Open() err: %v", err)
	}
	database = initDB(debugFlag, database)

	getRouterInfo(debugFlag, snmpTarget, community, maxHopsStr, params, router, database)

	// TODO: Write Links row to database

	// Close database
	database.Close()

	if debugFlag {
		fmt.Println("func discovery version", DISCOVERYVERSION, "ended.")
	}

	return database
}

func getInterfaces(debugFlag bool, snmpTarget string, community string, maxHopsStr string, params *g.GoSNMP, router Router, database *sql.DB) {

	// get Number of Interfaces
	ifNumberArray := []string{ifNumberOID + ".0"}
	getPDU, getError := params.Get(ifNumberArray)
	if getError != nil {
		log.Fatalf("Get() err: %v", getError)
	}
	if debugFlag {
		fmt.Println("ifNumber walkPDU=", getPDU)
	}
	nbrOfInterfaces := getPDU.Variables[0].Value.(int)
	if debugFlag {
		fmt.Println("nbrOfInterfaces =", nbrOfInterfaces)
	}
	// get ifTable
	walkPDU, walkError := params.WalkAll(ifTableOID)
	if walkError != nil {
		log.Fatalf("Get() err: %v", walkError)
	}
	if debugFlag {
		fmt.Println("\nifTable PDU=", walkPDU)
	}

	var interfaceTable ifTable

	if debugFlag {
		fmt.Println("len(walkPDU)=", len(walkPDU))
	}

	for i := 0; i < len(walkPDU); i++ { // skip ifIndex array within walkPDU
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifIndexOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifIndexType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifIndex = walkPDU[i].Value.(int)
			if debugFlag {
				fmt.Println("ifIndex=", interfaceTable.ifEntry.ifIndex)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifDescrOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifDescrType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifDescr = string(walkPDU[i].Value.([]uint8))
			if debugFlag {
				fmt.Println("ifDescr=", interfaceTable.ifEntry.ifDescr)
			}
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
			if debugFlag {
				fmt.Println("ifType=", interfaceTable.ifEntry.ifType)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifMtuOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifMtuType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifMtu = walkPDU[i].Value.(int)
			if debugFlag {
				fmt.Println("ifMtu=", interfaceTable.ifEntry.ifMtu) // TESTING ONLY
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifSpeedOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifSpeedType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifSpeed = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifSpeed=", interfaceTable.ifEntry.ifSpeed)
			}
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

			if debugFlag {
				fmt.Println("ifPhysAddress=", interfaceTable.ifEntry.ifPhysAddress)
			}

			writeMacToDB(debugFlag, router, interfaceTable, database)

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
			if debugFlag {
				fmt.Println("ifAdminStatus=", interfaceTable.ifEntry.ifAdminStatus)
			}
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
			if debugFlag {
				fmt.Println("ifOperStatus=", interfaceTable.ifEntry.ifOperStatus)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifLastChangeOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifLastChangeType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifLastChange = walkPDU[i].Value.(uint32)
			if debugFlag {
				fmt.Println("ifLastChange=", interfaceTable.ifEntry.ifLastChange)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInOctetsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInOctetsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInOctets = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifInOctets=", interfaceTable.ifEntry.ifInOctets)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInUcastPkts = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifInucastPkts=", interfaceTable.ifEntry.ifInUcastPkts)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInNUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInNUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInNUcastPkts = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifINUcastPkts=", interfaceTable.ifEntry.ifInUcastPkts)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInDiscardsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInDiscardsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInDiscards = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifDiscards=", interfaceTable.ifEntry.ifInDiscards)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInErrorsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInErrorsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInErrors = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifInErrors=", interfaceTable.ifEntry.ifInErrors)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifInUnknownProtosOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifInUnknownProtosType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifInUnknownProtos = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifInUnknownProtos=", interfaceTable.ifEntry.ifInUnknownProtos)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutOctetsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutOctetsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutOctets = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifOutOctets=", interfaceTable.ifEntry.ifOutOctets)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutUcastPkts = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifOutUcastPkts=", interfaceTable.ifEntry.ifOutUcastPkts)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutNUcastPktsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutNUcastPktsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutNUcastPkts = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifOutNUcastPkts=", interfaceTable.ifEntry.ifOutNUcastPkts)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutDiscardsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutDiscardsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutDiscards = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifOutDiscards=", interfaceTable.ifEntry.ifOutDiscards)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutErrorsOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutErrorsType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutErrors = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifOutErrors=", interfaceTable.ifEntry.ifOutErrors)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOutQLenOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOutQLenType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOutQLen = walkPDU[i].Value.(uint)
			if debugFlag {
				fmt.Println("ifOutQLen=", interfaceTable.ifEntry.ifOutQLen)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifSpecificOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifSpecificType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifSpecific = walkPDU[i].Value.(string)
			if debugFlag {
				fmt.Println("ifSpecific=", interfaceTable.ifEntry.ifSpecific)
			}
			i++
		}

	} // End of Interfaces code

}

func initDB(debugFlag bool, database *sql.DB) *sql.DB {
	initDbVersion := "0.0.2"
	if debugFlag {
		fmt.Println("initDB version:", initDbVersion)
	}

	/*
	 *	Add Routers table to DB
	 */
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY UNIQUE, Name TEXT, Description TEXT, UpTime TEXT, Contact TEXT, Location TEXT, Services INTEGER, GpsLat REAL, GPSLong REAL, GpsAlt REAL)")
	if err != nil {
		log.Fatalf("Router Table Create err: %v", err)
	}
	statement.Exec()

	/*
	 *	Add RouteTable table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, NextHop TEXT)")
	if err != nil {
		log.Fatalf("RouteTable Create err: %v", err)
	}
	statement.Exec()

	/*
	 *	Add RouterIP table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER NOT NULL, IpAddr TEXT)")
	if err != nil {
		log.Fatalf("RouterIP Create err: %v", err)
	}
	statement.Exec()

	/*
	 *	Add RouterMac table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
	if err != nil {
		log.Fatalf("RouterMac Create err: %v", err)
	}
	statement.Exec()

	/*
	 *	Add Links table to DB
	 */
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER NOT NULL, RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP TEXT)")
	if err != nil {
		log.Fatalf("Links Create err: %v", err)
	}
	statement.Exec()

	return database
}

func getIPADDR(ipAddr string) []string {
	names, err := net.LookupAddr(ipAddr)
	if err != nil {
		fmt.Println("No reverse lookup found for", ipAddr)
	}
	if len(names) == 0 {
		fmt.Println("No FQDN records for", ipAddr)
		return names
	}
	//	for _, name := range names {
	//		fmt.Printf("%s\n", name)
	//	}
	return names
}

func getGPS(sysName string) []string {
	txts, err := net.LookupTXT(sysName)
	if err != nil {
		panic(err)
	}
	if len(txts) == 0 {
		fmt.Println("No DNS TXT records for", sysName)
	}
	//	for _, txt := range txts {
	//		fmt.Printf("%s\n", txt)
	//	}
	return txts
}

func writeMacToDB(debugFlag bool, router Router, interfaceTable ifTable, database *sql.DB) {

	statement, err := database.Prepare("INSERT INTO RouterMac (RouterID, MacAddr) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("RouterMac Insert Prepare err: %v", err)
		log.Fatal(err)
	}

	RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
	_, err = statement.Exec(strconv.Itoa(int(RouterID)), interfaceTable.ifEntry.ifPhysAddress)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// Continue executing if this is a duplicate MAC Address. This assume this router is being processed again.
			// In case this is a duplicate MAC Address within the network, print error output to stdoutput.
			fmt.Println("\n****\n Non-Unique MAC Address", interfaceTable.ifEntry.ifPhysAddress, "\n This may be because this router is being re-discovered.\n If not, then this is a serious network violation condition.\n****")
		} else {
			fmt.Printf("RouterMac Insert Exec err: %v", err)
			log.Fatal(err)
		}
	}
}

func getIPAddresses(debugFlag bool, params *g.GoSNMP, router Router, database *sql.DB) {
	// get ipAddrTable
	walkPDU, err := params.WalkAll(ipAdEntAddrOID)
	if err != nil {
		log.Fatalf("Get(walkPDU) err: %v", err)
	}
	if debugFlag {
		fmt.Println("\nipAdEntAddr PDU=", walkPDU)
	}

	var ipTable ipAddrTable

	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntAddr = walkPDU[i].Value.(string)
		if debugFlag {
			fmt.Println("ipAdEntAddr=", ipTable.ipAddrEntry.ipAdEntAddr)
		}

		// Add row to RouterIp table
		statement, err := database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr) VALUES (?, ?)")
		if err != nil {
			fmt.Printf("RouterIp Prepare Insert Exec err: %v", err)
			log.Fatal(err)
		}
		RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
		_, err = statement.Exec(RouterID, ipTable.ipAddrEntry.ipAdEntAddr)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				// Continue executing if this is a duplicate MAC Address. This assume this router is being processed again.
				// In case this is a duplicate MAC Address within the network, print error output to stdoutput.
				fmt.Println("\n****\n Non-Unique IP Address", ipTable.ipAddrEntry.ipAdEntAddr, "\n This may be because this router is being re-discovered.\n If not, then this is a serious network violation condition.\n****")
			} else {
				fmt.Printf("RouterIp Exec Insert Exec err: %v", err)
				log.Fatal(err)
			}
		}
	}

	walkPDU, err = params.WalkAll(ipAdEntIfIndex)
	if err != nil {
		log.Fatalf("Get(walkPDU) err: %v", err)
	}
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntIfIndex = walkPDU[i].Value.(int)
		if debugFlag {
			fmt.Println("ipAdEntIfIndex=", ipTable.ipAddrEntry.ipAdEntIfIndex)
		}
	}
	walkPDU, err = params.WalkAll(ipAdEntNetMask)
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntNetMask = walkPDU[i].Value.(string)
		if debugFlag {
			fmt.Println("ipAdEntNetMask=", ipTable.ipAddrEntry.ipAdEntNetMask)
		}
	}
	walkPDU, err = params.WalkAll(ipAdEntBcastAddr)
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntBcastAddr = walkPDU[i].Value.(int)
		if debugFlag {
			fmt.Println("ipAdEntBcastAddr=", ipTable.ipAddrEntry.ipAdEntBcastAddr)
		}
	}
	walkPDU, err = params.WalkAll(ipAdEntReasmMaxSize)
	for i := 0; i < (len(walkPDU)); i++ {
		ipTable.ipAddrEntry.ipAdEntReasmMaxSize = walkPDU[i].Value.(int)
		if debugFlag {
			fmt.Println("ipAdEntReasmMaxSize=", ipTable.ipAddrEntry.ipAdEntReasmMaxSize)
		}
	}
}

func getIPRouteTable(debugFlag bool, params *g.GoSNMP, router Router, database *sql.DB) {

	// TODO: DB record must be written within the for loop that parses the PDU arrays.
	//       Remove the second nexthop for loop. loop contains dest and nexthop parses, then DB row add.

	// get ipRouteTable
	ipRouteDestPDU, err := params.WalkAll(ipRouteDestOID)
	if err != nil {
		log.Fatalf("Get(ipRouteDestPDU) err: %v", err)
	}
	if debugFlag {
		fmt.Println("\nipRouteDestPDU PDU=", ipRouteDestPDU)
	}

	ipRouteNextHopPDU, err := params.WalkAll(ipRouteNextHopOID)
	if err != nil {
		log.Fatalf("Get(ipRouteNextHopPDU) err: %v", err)
	}
	if debugFlag {
		fmt.Println("\nipRouteNextHopPDU PDU=", ipRouteNextHopPDU)
	}

	// Parse Dest and NextHop PDUs, adding row to ipRouteTable for each PDU element.
	var ipRouteTab ipRouteTable

	for i := 0; i < (len(ipRouteDestPDU)); i++ {
		ipRouteTab.ipRouteEntry.ipRouteDest = ipRouteDestPDU[i].Value.(string)
		if debugFlag {
			fmt.Println("ipRouteDest=", ipRouteTab.ipRouteEntry.ipRouteDest)
		}

		ipRouteTab.ipRouteEntry.ipRouteNextHop = ipRouteNextHopPDU[i].Value.(string)
		if debugFlag {
			fmt.Println("ipRouteNextHop=", ipRouteTab.ipRouteEntry.ipRouteNextHop)
		}

		// Add row to RouteTable table
		statement, _ := database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, NextHop) VALUES (?, ?, ?)")
		if err != nil {
			fmt.Printf("RouterTable Prepare Insert Exec err: %v", err)
			log.Fatal(err)
		}

		RouterID := crc32.ChecksumIEEE([]byte(router.System.Name))
		_, err = statement.Exec(RouterID, ipRouteTab.ipRouteEntry.ipRouteDest, ipRouteTab.ipRouteEntry.ipRouteNextHop)
		if err != nil {
			log.Fatalf("RouteTable Insert err: %v", err)
		}
	}

}

func getRouterInfo(debugFlag bool, snmpTarget string, community string, maxHopsStr string, params *g.GoSNMP, router Router, database *sql.DB) {
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
		log.Fatalf("Get() err: %v", err)
	}

	router.System.Name = string(result.Variables[0].Value.([]byte))
	router.System.Description = string(result.Variables[1].Value.([]byte))
	router.System.UpTime = result.Variables[2].Value.(uint32)
	router.System.Contact = string(result.Variables[3].Value.([]byte))
	router.System.Location = string(result.Variables[4].Value.([]byte))
	router.System.Services = result.Variables[5].Value.(int)

	/*
		// Retrieve GPS data from DNS
	*/

	// get FQDN with IP Address
	fqdn := getIPADDR(snmpTarget)

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

	if debugFlag {
		fmt.Println("router.System.Name=", router.System.Name)
		fmt.Println("router.System.Description=", router.System.Description)
		fmt.Println("router.System.UpTime=", router.System.UpTime)
		fmt.Println("router.System.Contact=", router.System.Contact)
		fmt.Println("router.System.Location=", router.System.Location)
		fmt.Println("router.System.Services=", router.System.Services)
		fmt.Println("router.System.GPS=", router.System.GPS)
	}

	// Write Router row to database
	statement, _ := database.Prepare("INSERT INTO Routers (RouterID, Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	statement.Exec()

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
			fmt.Println("Router", Name, "is already exists in database. Continuing discovery.")
			routerIsInDB = true
		} else {
			fmt.Printf("RouterMac Insert Exec err: %v", err)
			log.Fatal(err)
		}
	}

	if !routerIsInDB {
		getInterfaces(debugFlag, snmpTarget, community, maxHopsStr, params, router, database)

		getIPAddresses(debugFlag, params, router, database)

		getIPRouteTable(debugFlag, params, router, database)
	}

	/*
		// TODO: Write Links row to database

		// Close database
		database.Close()

		if debugFlag {
			fmt.Println("func discovery version", DISCOVERYVERSION, "ended.")
		}
	*/
}
