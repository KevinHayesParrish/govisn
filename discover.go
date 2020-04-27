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
const DiscoverVersion = "0.0.1"

func discover(debugFlag bool, snmpTarget string, community string, maxHopsStr string) {

	//type ifTable struct {
	//	ifEntry struct {
	//		ifIndex           []int
	//		ifDescr           []string
	//		ifType            []string
	//		ifMtu             []int32
	//		ifSpeed           []int32
	//		ifPhysAddress     []string
	//		ifAdminStatus     []string
	//		ifOperStatus      []string
	//		ifLastChange      []uint32
	//		ifInOctets        []uint32
	//		ifInNUcastPkts    []uint32
	//		ifInDiscards      []uint32
	//		ifInErrors        []uint32
	//		ifInUnknownProtos []uint32
	//		ifOutOctets       []uint32
	//		ifOutNUcastPkts   []uint32
	//		ifOutDiscards     []uint32
	//		ifOutErrors       []uint32
	//		ifOutQLen         []uint32 // deprecated
	//		ifSpecific        []string // deprecated
	//	}
	//}
	type ifTable struct {
		ifEntry struct {
			//ifIndexRow []struct {
			ifIndexRow struct {
				ifIndexOID  string
				ifIndexType string
				ifIndex     int
				Logger      string
			}

			//ifDescrRow []struct {
			ifDescrRow struct {
				ifDescrOID  string
				ifDescrType byte
				ifDescr     string
				Logger      string
			}

			//ifTypeRow []struct {
			ifTypeRow struct {
				ifTypeOID  string
				ifTypeType string
				ifType     int
				Logger     string
			}

			//ifMtuRow []struct {
			ifMtuRow struct {
				ifMtuOID  string
				ifMtuType string
				ifMtu     int32
				Logger    string
			}

			//ifSpeedRow []struct {
			ifSpeedRow struct {
				ifSpeedOID  string
				ifSpeedType string
				ifSpeed     uint32
				Logger      string
			}

			//ifPhysAddressRow []struct {
			ifPhysAddressRow []struct {
				ifPhysAddressOID  string
				ifPhysAddressType string
				ifPhysAddress     string
				Logger            string
			}

			//ifAdminStatusRow []struct {
			ifAdminStatusRow struct {
				ifAdminStatusOID  string
				ifAdminStatusType string
				ifAdminStatus     string
				Logger            string
			}

			//ifOperStatusRow []struct {
			ifOperStatusRow struct {
				ifOperStatusOID  string
				ifOperStatusType string
				ifOperStatus     string
				Logger           string
			}

			//ifLastChangeRow []struct {
			ifLastChangeRow struct {
				ifLastChangeOID  string
				ifLastChangeType string
				ifLastChange     uint32
				Logger           string
			}

			//ifInOctetsRow []struct {
			ifInOctetsRow struct {
				ifInOctetsOID  string
				ifInOctetsType string
				ifInOctets     uint32
				Logger         string
			}

			//ifInUcastPktsRow []struct {
			ifInUcastPktsRow struct {
				ifInUcastPktsOID  string
				ifInUcastPktsType string
				ifInUcastPkts     uint32
				Logger            string
			}

			//ifInNUcastPktsRow []struct {
			ifInNUcastPktsRow struct {
				ifInNUcastPktsOID  string // deprecated
				ifInNUcastPktsType string // deprecated
				ifInNUcastPkts     uint32 // deprecated
				Logger             string
			}

			//ifInDiscardsRow []struct {
			ifInDiscardsRow struct {
				ifInDiscardsOID  string
				ifInDiscardsType string
				ifInDiscards     uint32
				Logger           string
			}

			//ifInErrorsRow []struct {
			ifInErrorsRow struct {
				ifInErrorsOID  string
				ifInErrorsType string
				ifInErrors     uint32
				Logger         string
			}

			//ifInUnknownProtosRow []struct {
			ifInUnknownProtosRow struct {
				ifInUnknownProtosOID  string
				ifInUnknownProtosType string
				ifInUnknownProtos     uint32
				Logger                string
			}

			//ifOutOctetsRow []struct {
			ifOutOctetsRow struct {
				ifOutOctetsOID  string
				ifOutOctetsType string
				ifOutOctets     uint32
				Logger          string
			}

			//ifOutUcastPktsRow []struct {
			ifOutUcastPktsRow struct {
				ifOutUcastPktsOID  string
				ifOutUcastPktsType string
				ifOutUcastPkts     uint32
				Logger             string
			}

			//ifOutNUcastPktsRow []struct {
			ifOutNUcastPktsRow struct {
				ifOutNucastPktsOID  string // deprecated
				ifOutNucastPktsType string // deprecated
				ifOutNUcastPkts     uint32 //deprecated
				Logger              string
			}

			//ifOutDiscardsRow []struct {
			ifOutDiscardsRow struct {
				ifOutDiscardsOID  string
				ifOutDiscardsType string
				ifOutDiscards     uint32
				Logger            string
			}

			//ifOutErrorsRow []struct {
			ifOutErrorsRow struct {
				ifOutErrorsOID  string
				ifOutErrorsType string
				ifOutErrors     uint32
				Logger          string
			}

			//ifOutQLenRow []struct {
			ifOutQLenRow struct {
				ifOutQLenOID  string
				ifOutQLenType string
				ifOutQLen     uint32 // deprecated
				Logger        string
			}

			//IfOutSpecificRow []struct {
			IfOutSpecificRow struct {
				ifOutSpecficOID    string
				ifOutSpecificType  string
				ifSpecificSpecific string // deprecated
				Logger             string
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

	fmt.Println("\nfunc discover started.\ndebugFlag=", debugFlag)

	// get Target and Port from environment
	//	envTarget := os.Getenv("GOSNMP_TARGET")
	//envTarget := snmpTarget
	//	envPort := os.Getenv("GOSNMP_PORT")

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

	// Build our own GoSNMP struct, rather than using g.Default.
	// Do verbose logging of packets.
	//params := &g.GoSNMP{
	//	Target: envTarget,
	//	Port:   uint16(port),
	//	Community: "public",
	//	Version:   g.Version2c,
	//	Timeout:   time.Duration(2) * time.Second,
	//	Logger:    log.New(os.Stdout, "govisn.discover: ", 0), // TESTING ONLY
	//}

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

	// Retrive sysName, sysDescr, sysContact, sysLocation
	//oids := []string{
	//	"1.3.6.1.2.1.1.5.0", // sysName
	//	"1.3.6.1.2.1.1.1.0", // sysDescr
	//	"1.3.6.1.2.1.1.4.0", // sysContact
	//	"1.3.6.1.2.1.1.6.0", // sysLocation
	//}
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
	//ifTable, err3 := params.WalkAll("1.3.6.1.2.1.2.2")
	//ifTable, walkError := params.WalkAll(ifTableOID)
	walkPDU, walkError := params.WalkAll(ifTableOID)
	if walkError != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	if debugFlag {
		fmt.Println("\nifTable PDU=", walkPDU)
	}

	//var interfaceTable ifTable
	var interfaceTable ifTable

	//nbrOfInterfacesInt, strErr := strconv.ParseInt(nbrOfInterfaces.(string), 10, 32)
	//if strErr != nil {
	//	log.Fatalf("Get() err: %v", strErr)
	//}
	//var i int32 = 0
	//for i, variable := range walkPDU {
	//fmt.Println("Interface Descr", i, "=", variable.Value)
	//var strErr error
	//var ifIndexInt int64 = 0
	//var ifIndexStr string = ""
	//ifIndexInt, strErr = strconv.ParseInt(variable.Value.(string), 10, 16)
	//ifIndexInt, strErr = variable.Value
	//	fmt.Println("i=", "walkPDU[i]=", walkPDU[i])   // TESTING ONLY
	//	fmt.Println("variable=", variable)             // TESTING ONLY
	//	fmt.Println("variable.Name=", variable.Name)   // TESTING ONLY
	//	fmt.Println("variable.Value=", variable.Value) // TESTING ONLY
	//ifIndexStr := variable.Value.(string)
	//if strErr != nil {
	//	log.Fatalf("Get() err: %v", strErr)
	//}
	//interfaceTable.ifEntry.ifIndex[i] = int(ifIndexInt)
	//interfaceTable.ifEntry.ifIndex[i] = ifIndexInt
	//	interfaceTable.ifEntry.ifIndex[i] = ifIndexStr
	//	fmt.Println("interfaceTable.ifEntry.ifindex[i]", interfaceTable.ifEntry.ifIndex[i])
	//for k := 0; k < int(nbrOfInterfacesInt); k++ {
	//	fmt.Println("walkPDU=", walkPDU)
	//addressTable.ipAddrEntry.ipAdEntAddr[i] = walkPDU.Variables[k].Value
	//	fmt.Println("addressTable.ipAddrEntry.ipAdEntAddr[i]", addressTable.ipAddrEntry.ipAdEntAddr[i])
	//}
	//}

	//for i := 0; i < nbrOfInterfaces; i++ {

	if debugFlag {
		fmt.Println("len(walkPDU)=", len(walkPDU))
	}
	for i := nbrOfInterfaces; i < len(walkPDU); i++ { // skip ifIndex array within walkPDU

		//for k := nbrOfInterfaces; k < len(walkPDU); k++ { // skip ifIndex array within walkPDU
		//for k := 0; k < nbrOfInterfaces; k++ { // skip ifIndex array within walkPDU
		//interfaceTable.ifEntry.ifIndexRow[0].ifIndex = walkPDU[i].Value.(int)
		//ifDescr := string(walkPDU[i].Value.([]uint8))
		//fmt.Println("ifDesc=", ifDescr)
		//interfaceTable.ifEntry.ifDescrRow[k].ifDescr = string(walkPDU[k].Value.([]uint8))
		//interfaceTable.ifEntry.ifDescrRow.ifDescr = ifDescr
		interfaceTable.ifEntry.ifDescrRow.ifDescrOID = walkPDU[i].Name
		interfaceTable.ifEntry.ifDescrRow.ifDescrType = byte(walkPDU[i].Type)
		interfaceTable.ifEntry.ifDescrRow.ifDescr = string(walkPDU[i].Value.([]uint8))
		if debugFlag {
			//println("ifDesc(", k, ")=", interfaceTable.ifEntry.ifDescrRow.ifDescr)
			println("ifDesc(", i, ")=", interfaceTable.ifEntry.ifDescrRow.ifDescr)
		}
		// TODO
		// write ifDesc to database
		// add ifType to interfaceTable
		// add ifMTU to interfaceTable
		// Add ifSpeed to interfaceTable
		// Add ifPhyAddress to interfaceTable
		// Add ifOutOctets to interfaceTable

		//}
		fmt.Println("i=", i) // TROUBLESHOOTING ONLY. REMOVE AFTER TROUBLESHOOTING
	}

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
}
