package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	g "github.com/soniah/gosnmp"
)

/*
* TODO:
 */

//DISCOVERYVERSION is the file version number
const DISCOVERYVERSION = "0.1.2"

func discover(debugFlag bool, snmpTarget string, community string, maxHopsStr string) {

	type Router struct {
		sysName     string
		sysDescr    string
		sysUpTime   uint32
		sysContact  string
		sysLocation string
		sysServices *big.Int
		GpsLat      string
		GpsLong     string
		GpsAlt      string
	}

	type ipAddrTable struct {
		ipAddrEntry struct {
			ipAdEntAddr         string
			ipAdEntIfIndex      int32
			ipAdEntNetMask      string
			ipAdEntBcastAddr    int32
			ipAdEntReasmMaxSize int32
		}
	}

	type ipRouteTable struct {
		ipRouteEntry struct {
			ipRouteDest    string
			ipRouteIfIndex int32
			ipRouteMetric1 int32
			ipRouteMetric2 int32
			ipRouteMetric3 int32
			ipRouteMetric4 int32
			ipRouteNextHop string
			ipRouteType    string
			ipRouteProto   string
			ipRouteAge     int32
			ipRouteMask    string
			ipRouteMetric5 int32
			ipRouteInfo    string
		}
	}

	fmt.Println("\nfunc discover version", DISCOVERYVERSION, "started.\ndebugFlag=", debugFlag)

	maxHops, _ := strconv.Atoi(maxHopsStr)
	// Discover network, constrained by input parm maximum hops away from snmpTarget node
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
		MaxOids:   5,
	}

	if debugFlag {
		fmt.Println("params=", params)
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	oids := []string{
		sysNameOID + ".0",     // sysName
		sysDescrOID + ".0",    // sysDescr
		sysContactOID + ".0",  // sysContact
		sysLocationOID + ".0", // sysLocation
		sysServicesOID + ".0", // sysServices
	}

	result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		log.Fatalf("Get() err: %v", err)
	}

	router.sysName = string(result.Variables[0].Value.([]byte))
	router.sysDescr = string(result.Variables[1].Value.([]byte))
	router.sysContact = string(result.Variables[2].Value.([]byte))
	router.sysLocation = string(result.Variables[3].Value.([]byte))
	router.sysServices = g.ToBigInt(result.Variables[4].Value)

	if debugFlag {
		fmt.Println("router.sysName=", router.sysName)
		fmt.Println("router.sysDescr=", router.sysDescr)
		fmt.Println("router.sysName=", router.sysName)
		fmt.Println("router.sysContact=", router.sysContact)
		fmt.Println("router.sysLocation=", router.sysLocation)
		fmt.Println("router.sysServices=", router.sysServices)
	}

	// TODO: Write Router row to database

	//database := initDB()

	//statement, _ := database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	discoverInterfaces(debugFlag, snmpTarget, community, maxHopsStr, params)

	// get ipAddrTable
	//var addressTable ipAddrTable
	//	walkPDU, walkError := params.WalkAll(ipAddrTableOID)
	//	if walkError != nil {
	//		log.Fatalf("Get() err: %v", err2)
	//	}

	// TODO: Write RouterIp Row to database

	//var ipAddrTableResult ipAddrTable

	// get ipRouteTable
	//	walkPDU, walkError = params.WalkAll(ipRouteTableOID)
	//	if walkError != nil {
	//		log.Fatalf("Get() err: %v", err2)
	//	}
	//	if debugFlag {
	//		fmt.Println("\nipRouteTable PDU=", walkPDU)
	//	}

	// TODO: Write RouterIp Row to database

	//var ipRouteTableResult ipRouteTable

	// TODO: Write Links row to database

	if debugFlag {
		fmt.Println("func discovery version", DISCOVERYVERSION, "ended.")
	}
}

func discoverInterfaces(debugFlag bool, snmpTarget string, community string, maxHopsStr string, params *g.GoSNMP) {
	type ifTable struct {
		ifEntry struct {
			ifIndexOID    string
			ifIndexType   byte
			ifIndex       int
			ifIndexLogger string

			ifDescrOID    string
			ifDescrType   byte
			ifDescr       string
			ifDescrLogger string

			ifTypeOID    string
			ifTypeType   byte
			ifType       string
			ifTypeLogger string

			ifMtuOID    string
			ifMtuType   byte
			ifMtu       int
			ifMTULogger string

			ifSpeedOID    string
			ifSpeedType   byte
			ifSpeed       uint
			ifSpeedLogger string

			ifPhysAddressOID    string
			ifPhysAddressType   byte
			ifPhysAddress       string
			ifPhysAddressLogger string

			ifAdminStatusOID    string
			ifAdminStatusType   byte
			ifAdminStatus       string
			ifAdminStatusLogger string

			ifOperStatusOID    string
			ifOperStatusType   byte
			ifOperStatus       string
			ifOperStatusLogger string

			ifLastChangeOID    string
			ifLastChangeType   byte
			ifLastChange       uint32
			ifLastChangeLogger string

			ifInOctetsOID    string
			ifInOctetsType   byte
			ifInOctets       uint
			ifInOctetsLogger string

			ifInUcastPktsOID    string
			ifInUcastPktsType   byte
			ifInUcastPkts       uint
			ifInUcastPktsLogger string

			ifInNUcastPktsOID    string // deprecated
			ifInNUcastPktsType   byte   // deprecated
			ifInNUcastPkts       uint   // deprecated
			ifInNUcastPktsLogger string

			ifInDiscardsOID    string
			ifInDiscardsType   byte
			ifInDiscards       uint
			ifInDiscardsLogger string

			ifInErrorsOID    string
			ifInErrorsType   byte
			ifInErrors       uint
			ifInErrorsLogger string

			ifInUnknownProtosOID    string
			ifInUnknownProtosType   byte
			ifInUnknownProtos       uint
			ifInUnknownProtosLogger string

			ifOutOctetsOID    string
			ifOutOctetsType   byte
			ifOutOctets       uint
			ifOutOctetsLogger string

			ifOutUcastPktsOID    string
			ifOutUcastPktsType   byte
			ifOutUcastPkts       uint
			ifOutUcastPktsLogger string

			ifOutNUcastPktsOID    string // deprecated
			ifOutNUcastPktsType   byte   // deprecated
			ifOutNUcastPkts       uint   //deprecated
			ifOutNUcastPktsLogger string

			ifOutDiscardsOID    string
			ifOutDiscardsType   byte
			ifOutDiscards       uint
			ifOutDiscardsLogger string

			ifOutErrorsOID    string
			ifOutErrorsType   byte
			ifOutErrors       uint
			ifOutErrorsLogger string

			ifOutQLenOID    string
			ifOutQLenType   byte
			ifOutQLen       uint // deprecated
			ifOutQLenLogger string

			ifSpecificOID    string
			ifSpecificType   byte
			ifSpecific       string // deprecated
			ifSpecificLogger string
		}
	}

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
			//			interfaceTable.ifEntry.ifType = walkPDU[i].Value.(int)
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
			fmt.Println("ifMtu=", interfaceTable.ifEntry.ifMtu) // TESTING ONLY
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

			i++

			// TODO: Write RouterMac row to database
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
			//			interfaceTable.ifEntry.ifOperStatus = string(walkPDU[i].Value.(int))
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

func initDB() *sql.DB {
	initDbVersion := "0.0.1"
	fmt.Println("initDB version:", initDbVersion)
	database, _ := sql.Open("sqlite3", "./govisionDiscoveredDb.db")

	/*
	 *	Add Routers table to DB
	 */
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat REAL, GPSLong REAL, GpsAlt REAL)")
	statement.Exec()
	//	 statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	/*
	 *	Add RouteTable table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, NextHop TEXT)")
	statement.Exec()
	// 	statement, _ = database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, NextHop) VALUES (?, ?, ?)")

	/*
	 *	Add RouterIP table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER, IpAddr TEXT)")
	statement.Exec()
	//	statement, _ = database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr) VALUES (?, ?)")

	/*
	 *	Add RouterMac table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
	statement.Exec()
	//	statement, _ = database.Prepare("INSERT INTO RouterMac (RourterID, MacAddr) VALUES (?, ?)")

	/*
	 *	Add Links table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, FromRouter TEXT, ToRouter TEXT)")
	statement.Exec()
	//	statement, _ = database.Prepare("INSERT INTO Links (LinkID, FromRouter, ToRouter) VALUES (?, ?, ?)")

	return database
}
