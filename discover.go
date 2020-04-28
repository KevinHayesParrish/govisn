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
			//			ifIndexRow struct {
			ifIndexOID    string
			ifIndexType   byte
			ifIndex       int
			ifIndexLogger string
			//			}

			//			ifDescrRow struct {
			ifDescrOID    string
			ifDescrType   byte
			ifDescr       string
			ifDescrLogger string
			//			}

			//			ifTypeRow struct {
			ifTypeOID    string
			ifTypeType   byte
			ifType       int
			ifTypeLogger string
			//			}

			ifMtuRow struct {
				ifMtuOID    string
				ifMtuType   byte
				ifMtu       int32
				ifMTULogger string
			}

			ifSpeedRow struct {
				ifSpeedOID    string
				ifSpeedType   byte
				ifSpeed       uint32
				ifSpeedLogger string
			}

			ifPhysAddressRow []struct {
				ifPhysAddressOID    string
				ifPhysAddressType   byte
				ifPhysAddress       string
				ifPhysAddressLogger string
			}

			ifAdminStatusRow struct {
				ifAdminStatusOID    string
				ifAdminStatusType   byte
				ifAdminStatus       string
				ifAdminStatusLogger string
			}

			ifOperStatusRow struct {
				ifOperStatusOID    string
				ifOperStatusType   byte
				ifOperStatus       string
				ifOperStatusLogger string
			}

			ifLastChangeRow struct {
				ifLastChangeOID    string
				ifLastChangeType   byte
				ifLastChange       uint32
				ifLastChangeLogger string
			}

			ifInOctetsRow struct {
				ifInOctetsOID    string
				ifInOctetsType   byte
				ifInOctets       uint32
				ifInOctetsLogger string
			}

			ifInUcastPktsRow struct {
				ifInUcastPktsOID    string
				ifInUcastPktsType   byte
				ifInUcastPkts       uint32
				ifInUcastPktsLogger string
			}

			ifInNUcastPktsRow struct {
				ifInNUcastPktsOID    string // deprecated
				ifInNUcastPktsType   byte   // deprecated
				ifInNUcastPkts       uint32 // deprecated
				ifInNUcastPktsLogger string
			}

			ifInDiscardsRow struct {
				ifInDiscardsOID    string
				ifInDiscardsType   byte
				ifInDiscards       uint32
				ifInDiscardsLogger string
			}

			ifInErrorsRow struct {
				ifInErrorsOID    string
				ifInErrorsType   byte
				ifInErrors       uint32
				ifInErrorsLogger string
			}

			ifInUnknownProtosRow struct {
				ifInUnknownProtosOID    string
				ifInUnknownProtosType   byte
				ifInUnknownProtos       uint32
				ifInUnknownProtosLogger string
			}

			ifOutOctetsRow struct {
				ifOutOctetsOID    string
				ifOutOctetsType   byte
				ifOutOctets       uint32
				ifOutOctetsLogger string
			}

			ifOutUcastPktsRow struct {
				ifOutUcastPktsOID    string
				ifOutUcastPktsType   byte
				ifOutUcastPkts       uint32
				ifOutUcastPktsLogger string
			}

			ifOutNUcastPktsRow struct {
				ifOutNucastPktsOID    string // deprecated
				ifOutNucastPktsType   byte   // deprecated
				ifOutNUcastPkts       uint32 //deprecated
				ifOutNUcastPktsLogger string
			}

			ifOutDiscardsRow struct {
				ifOutDiscardsOID    string
				ifOutDiscardsType   byte
				ifOutDiscards       uint32
				ifOutDiscardsLogger string
			}

			ifOutErrorsRow struct {
				ifOutErrorsOID    string
				ifOutErrorsType   byte
				ifOutErrors       uint32
				ifOutErrorsLogger string
			}

			ifOutQLenRow struct {
				ifOutQLenOID    string
				ifOutQLenType   byte
				ifOutQLen       uint32 // deprecated
				ifOutQLenLogger string
			}

			IfOutSpecificRow struct {
				ifOutSpecficOID          string
				ifOutSpecificType        byte
				ifSpecificSpecific       string // deprecated
				ifSpecificSpecificLogger string
			}
		}
	}

	type ipAddrTable struct {
		ipAddrEntry struct {
			ipAdEntAddr         []string
			ipAdEntIfIndex      []int32
			ipAdEntNetMask      []string
			ipAdEntBcastAddr    []int32
			ipAdEntReasmMaxSize []int32
		}
	}

	type ipRouteTable struct {
		ipRouteEntry struct {
			ipRouteDest    []string
			ipRouteIfIndex []int32
			ipRouteMetric1 []int32
			ipRouteMetric2 []int32
			ipRouteMetric3 []int32
			ipRouteMetric4 []int32
			ipRouteNextHop []string
			ipRouteType    []string
			ipRouteProto   []string
			ipRouteAge     []int32
			ipRouteMask    []string
			ipRouteMetric5 []int32
			ipRouteInfo    []string
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
	if debugFlag {
		fmt.Println("\ngovisn.discover.results=", result)
	}

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
		fmt.Println("\nifTable PDU=", walkPDU)
	}

	//var interfaceTable ifTable
	var interfaceTable ifTable

	if debugFlag {
		fmt.Println("len(walkPDU)=", len(walkPDU))
	}
	for i := nbrOfInterfaces; i < len(walkPDU); i++ { // skip ifIndex array within walkPDU
		for k := i; k < k+nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifIndexOID = walkPDU[k].Name
			interfaceTable.ifEntry.ifIndexType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifIndex = walkPDU[i].Value.(int)
		}
		for k := i; k < k+nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifDescrOID = walkPDU[i].Name
			interfaceTable.ifEntry.ifDescrType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifDescr = string(walkPDU[i].Value.([]uint8))
		}
		for k := i; k < k+nbrOfInterfaces; k++ {
			interfaceTable.ifEntry.ifTypeOID = walkPDU[k].Name
			interfaceTable.ifEntry.ifTypeType = byte(walkPDU[i].Type)
			interfaceTable.ifEntry.ifType = walkPDU[i].Value.(int)
		}

		if debugFlag {
			//println("ifDesc(", k, ")=", interfaceTable.ifEntry.ifDescrRow.ifDescr)
			println("ifDesc(", i, ")=", interfaceTable.ifEntry.ifDescr)
		}
		// TODO
		// add ifType to interfaceTable
		// add ifMTU to interfaceTable
		// Add ifSpeed to interfaceTable
		// Add ifPhyAddress to interfaceTable
		// Add ifOutOctets to interfaceTable
		// write Router table row to database
		// write RouterMac table row to database

		//}
		//fmt.Println("i=", i) // TROUBLESHOOTING ONLY. REMOVE AFTER TROUBLESHOOTING
	}

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
	//var ipAddrTableResult ipAddrTable

	// get ipRouteTable
	walkPDU, walkError = params.WalkAll(ipRouteTableOID)
	if walkError != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	if debugFlag {
		fmt.Println("\nipRouteTable PDU=", walkPDU)
	}
	//var ipRouteTableResult ipRouteTable

end:
	return
}
