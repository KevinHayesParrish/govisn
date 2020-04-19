package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	g "github.com/soniah/gosnmp"
)

func discover(debugFlag bool, snmpTarget string, community string, maxHopsStr string) {

	type ifTable struct {
		ifEntry struct {
			ifIndex           []int
			ifDescr           []string
			ifType            []string
			ifMtu             []int32
			ifSpeed           []int32
			ifPhysAddress     []string
			ifAdminStatus     []string
			ifOperStatus      []string
			ifLastChange      []uint32
			ifInOctets        []uint32
			ifInNUcastPkts    []uint32
			ifInDiscards      []uint32
			ifInErrors        []uint32
			ifInUnknownProtos []uint32
			ifOutOctets       []uint32
			ifOutNUcastPkts   []uint32
			ifOutDiscards     []uint32
			ifOutErrors       []uint32
			ifOutQLen         []uint32 // deprecated
			ifSpecific        []string // deprecated
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
	} // TESTING ONLY
	nbrOfInterfaces := getPDU.Variables[0].Value
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

	var interfaceTable ifTable
	//nbrOfInterfacesInt, strErr := strconv.ParseInt(nbrOfInterfaces.(string), 10, 32)
	//if strErr != nil {
	//	log.Fatalf("Get() err: %v", strErr)
	//}
	//var i int32 = 0
	for i, variable := range walkPDU {
		//fmt.Println("Interface Descr", i, "=", variable.Value)
		var strErr error
		interfaceTable.ifEntry.ifIndex[i], strErr = strconv.ParseInt(variable.Value.(string), 10, 16)
		fmt.Println("interfaceTable.ifEntry.ifindex[i]", interfaceTable.ifEntry.ifIndex[i])
		//for k := 0; k < int(nbrOfInterfacesInt); k++ {
		//	fmt.Println("walkPDU=", walkPDU)
		//addressTable.ipAddrEntry.ipAdEntAddr[i] = walkPDU.Variables[k].Value
		//	fmt.Println("addressTable.ipAddrEntry.ipAdEntAddr[i]", addressTable.ipAddrEntry.ipAdEntAddr[i])
		//}
	}

	// get ipAddrTable
	//var addressTable ipAddrTable
	walkPDU, walkError = params.WalkAll(ipAddrTableOID)
	if walkError != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	if debugFlag {
		fmt.Println("\nipAddrTable PDU=", walkPDU)
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
