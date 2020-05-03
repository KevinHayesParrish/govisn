package main

import (
	"fmt"
	//"gosnmp"
	"log"
	"strconv"
	"time"

	g "github.com/soniah/gosnmp"
)

/*
* TODO:
 */

//DiscoverVersion is the file version number
const DiscoverVersion = "0.1.0"

func discover(debugFlag bool, snmpTarget string, community string, maxHopsStr string) {

	type ifTable struct {
		ifEntry struct {
			//	ifIndexRow []struct {
			ifIndexOID    string
			ifIndexType   byte
			ifIndex       int
			ifIndexLogger string
			//	}

			//	ifDescrRow []struct {
			ifDescrOID    string
			ifDescrType   byte
			ifDescr       string
			ifDescrLogger string
			//	}

			//	ifTypeRow []struct {
			ifTypeOID    string
			ifTypeType   byte
			ifType       int
			ifTypeLogger string
			//	}

			//	ifMtuRow []struct {
			ifMtuOID    string
			ifMtuType   byte
			ifMtu       int
			ifMTULogger string
			//	}

			//	ifSpeedRow []struct {
			ifSpeedOID    string
			ifSpeedType   byte
			ifSpeed       uint
			ifSpeedLogger string
			//	}

			//	ifPhysAddressRow []struct {
			ifPhysAddressOID  string
			ifPhysAddressType byte
			//ifPhysAddress       []byte
			ifPhysAddress       string
			ifPhysAddressLogger string
			//	}

			//	ifAdminStatusRow []struct {
			ifAdminStatusOID    string
			ifAdminStatusType   byte
			ifAdminStatus       string
			ifAdminStatusLogger string
			//	}

			//	ifOperStatusRow []struct {
			ifOperStatusOID    string
			ifOperStatusType   byte
			ifOperStatus       string
			ifOperStatusLogger string
			//	}

			//	ifLastChangeRow []struct {
			ifLastChangeOID    string
			ifLastChangeType   byte
			ifLastChange       uint32
			ifLastChangeLogger string
			//	}

			//	ifInOctetsRow []struct {
			ifInOctetsOID    string
			ifInOctetsType   byte
			ifInOctets       uint
			ifInOctetsLogger string
			//	}

			//	ifInUcastPktsRow []struct {
			ifInUcastPktsOID    string
			ifInUcastPktsType   byte
			ifInUcastPkts       uint
			ifInUcastPktsLogger string
			//	}

			//	ifInNUcastPktsRow []struct {
			ifInNUcastPktsOID    string // deprecated
			ifInNUcastPktsType   byte   // deprecated
			ifInNUcastPkts       uint   // deprecated
			ifInNUcastPktsLogger string
			//	}

			//	ifInDiscardsRow []struct {
			ifInDiscardsOID    string
			ifInDiscardsType   byte
			ifInDiscards       uint
			ifInDiscardsLogger string
			//	}

			//	ifInErrorsRow []struct {
			ifInErrorsOID    string
			ifInErrorsType   byte
			ifInErrors       uint
			ifInErrorsLogger string
			//	}

			//	ifInUnknownProtosRow []struct {
			ifInUnknownProtosOID    string
			ifInUnknownProtosType   byte
			ifInUnknownProtos       uint
			ifInUnknownProtosLogger string
			//	}

			//	ifOutOctetsRow []struct {
			ifOutOctetsOID    string
			ifOutOctetsType   byte
			ifOutOctets       uint
			ifOutOctetsLogger string
			//	}

			//	ifOutUcastPktsRow []struct {
			ifOutUcastPktsOID    string
			ifOutUcastPktsType   byte
			ifOutUcastPkts       uint
			ifOutUcastPktsLogger string
			//	}

			//	ifOutNUcastPktsRow []struct {
			ifOutNUcastPktsOID    string // deprecated
			ifOutNUcastPktsType   byte   // deprecated
			ifOutNUcastPkts       uint   //deprecated
			ifOutNUcastPktsLogger string
			//	}

			//	ifOutDiscardsRow []struct {
			ifOutDiscardsOID    string
			ifOutDiscardsType   byte
			ifOutDiscards       uint
			ifOutDiscardsLogger string
			//	}

			//	ifOutErrorsRow []struct {
			ifOutErrorsOID    string
			ifOutErrorsType   byte
			ifOutErrors       uint
			ifOutErrorsLogger string
			//	}

			//	ifOutQLenRow []struct {
			ifOutQLenOID    string
			ifOutQLenType   byte
			ifOutQLen       uint // deprecated
			ifOutQLenLogger string
			//	}

			//	ifSpecificRow []struct {
			ifSpecificOID    string
			ifSpecificType   byte
			ifSpecific       string // deprecated
			ifSpecificLogger string
			//	}
		}
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

	fmt.Println("\nfunc discover version", DiscoverVersion, "started.\ndebugFlag=", debugFlag)

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

	maxHops, _ := strconv.Atoi(maxHopsStr)
	// Discover network, constrained by input parm maximum hops away from snmpTarget node
	for i := 0; i <= maxHops; i++ {
		if debugFlag {
			fmt.Println("Discover iteration")
		}
	}

	// GoSNMP struct
	params := &g.GoSNMP{
		Target:    snmpTarget,
		Port:      uint16(port),
		Community: community,
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Logger:    nil,
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
	result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	//	if debugFlag {
	//		fmt.Println("\ngovisn.discover.results=", result)
	//	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case g.OctetString:

			fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}

	// TODO: Write Router row to database

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
		log.Fatalf("Get() err: %v", err2)
	}
	if debugFlag {
		//		fmt.Println("\nifTable PDU=", walkPDU)
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
			interfaceTable.ifEntry.ifType = walkPDU[i].Value.(int)
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
			//			interfaceTable.ifEntry.ifPhysAddress = walkPDU[i].Value.([]byte)
			//var physAddrInt [12]byte
			physAddrUint := walkPDU[i].Value.([]byte)
			//physAddrUint8 := physAddrUint[0]
			//var physAddrInt []int
			//for j := 0; j < 12; j++ {
			//	physAddrInt[j] = int(physAddrUint[j])
			//}
			var physAddrHex [6]string

			//physAddrHex[0] = fmt.Sprintf("%x", physAddrUint8)

			for l := 0; l < 6; l++ {
				physAddrUint8 := physAddrUint[l]
				physAddrHex[l] = fmt.Sprintf("%x", physAddrUint8)
			}
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

			if debugFlag {
				fmt.Println("ifPhysAddress=", interfaceTable.ifEntry.ifPhysAddress)
			}

			i++

			// TODO: Write RouterMac row to database
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifAdminStatusOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifAdminStatusType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifAdminStatus = string(walkPDU[i].Value.(int))
			if debugFlag {
				fmt.Println("ifAdminStatus=", interfaceTable.ifEntry.ifAdminStatus)
			}
			i++
		}
		for k := 0; k < nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifOperStatusOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifOperStatusType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifOperStatus = string(walkPDU[i].Value.(int))
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

	if debugFlag { //  TROUBLESHOOTING ONLY. REMOVE AFTER TROUBLESHOOTING
		goto end //  TROUBLESHOOTING ONLY. REMOVE AFTER TROUBLESHOOTING
	} //  TROUBLESHOOTING ONLY. REMOVE AFTER TROUBLESHOOTING

	// get ipAddrTable
	//var addressTable ipAddrTable
	walkPDU, walkError = params.WalkAll(ipAddrTableOID)
	if walkError != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	if debugFlag {
		fmt.Println("\nipAddrTable PDU=", walkPDU)
		fmt.Println("interfaceTable=", interfaceTable)
		return // TESTING ONLY, REMOVE AFTER TEST
	}

	// TODO: Write RouterIp Row to database

	//var ipAddrTableResult ipAddrTable

	// get ipRouteTable
	walkPDU, walkError = params.WalkAll(ipRouteTableOID)
	if walkError != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	if debugFlag {
		fmt.Println("\nipRouteTable PDU=", walkPDU)
	}

	// TODO: Write RouterIp Row to database

	//var ipRouteTableResult ipRouteTable

	// TODO: Write Links row to database

end:
	return
}
