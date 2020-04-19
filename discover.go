package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	g "github.com/soniah/gosnmp"
)

func discover(debugFlag bool, snmpTarget string, community string, maxHopsStr string) {

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

	// get ifTable
	//ifTable, err3 := params.WalkAll("1.3.6.1.2.1.2.2")
	ifTable, err3 := params.WalkAll(ifTableOID)
	if err3 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	fmt.Println("\nifTable PDU=", ifTable)

	// get ipAddrTable
	resultPDU, err3 := params.WalkAll(ipAddrTableOID)
	if err3 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	fmt.Println("\nipAddrTable PDU=", resultPDU)
	//var ipAddrTableResult ipAddrTable

	// get ipRouteTable
	resultPDU, err3 = params.WalkAll(ipRouteTableOID)
	if err3 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	fmt.Println("\nipRouteTable PDU=", resultPDU)
	//var ipRouteTableResult ipRouteTable
}
