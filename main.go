// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math"
	"strconv"
	"time"

	//	"github.com/g3n/g3nd/material"
	"github.com/g3n/engine/util/logger"
	_ "github.com/mattn/go-sqlite3"

	g "github.com/gosnmp/gosnmp"
)

// GOVISN_VERSION is the file version number
const GOVISN_VERSION = "0.22.7"

// ROUTER_RADIUS is the radius of the 3D object representing a network router
const ROUTER_RADIUS float64 = 0.5

// GLOBE_RADIUS is the radius of the 3D object representing the earth
const GLOBE_RADIUS float64 = 63.7

const CONST_X = math.Pi / 180

var log *logger.Logger

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("de", false, "Print Debug statements.")

var loadDBFlag = flag.Bool("l", false, "Load a database from an XML document.")

var dbName = "govisnDiscoveredNet.db"

// DbName is the name of the discovered network database file or name of XML input file
var DbName = flag.String("f", "govisnDiscoveredNet.db", "Name of the discovered network database -or-\nName of the XML input file, if combined with -l option.")

// discoverFlag is the option to discover a network
var discoverFlag = flag.String("di", "", "Discover a network using a seed IP Address")

var kmlFlag = flag.String("k", "", "Export GoVisn database to KML format file")

var seed = "127.0.0.1"

var community = flag.String("co", "public", "SNMP Community ReadOnly String")
var maxHops = flag.String("m", "10", "Scope of discovery. Maximum number of Hops from seed.")
var visualizeFlag = flag.Bool("vi", false, "Visualize the Network.")

// scanNetFlag is the startup option to scan the network for SNMP capable routers.
var scanNetFlag = flag.String("s", "", "Scan the CIDR network for SNMP capable routers.\nCIDR format = x.x.x.x/n. ex: 192.168.1.0/24\nOnce the network is scanned, the list of found routers\nwill be queried and their information added to the database.")

// Rad converts degrees to radians
func Rad(d float64) float64 { return d * CONST_X }

// Deg converts radians to degrees
func Deg(r float64) float64 { return r / CONST_X }

// Function to connect to an SNMP agent
func connectToSNMP(target string, community string) (*g.GoSNMP, error) {
	client := &g.GoSNMP{
		Target:    target,
		Port:      uint16(161),
		Community: community,
		Version:   g.Version2c,
		//		Timeout:   10,
		Timeout: time.Duration(2) * time.Second,
	}
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Function to discover routes from a router
func discoverRoutes(client *g.GoSNMP, oid string) (map[string]string, error) {
	routes := make(map[string]string)
	err := client.Walk(oid, func(pdu g.SnmpPDU) error {
		// Example processing, customize based on OID structure
		routes[pdu.Name] = fmt.Sprintf("%v", pdu.Value)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return routes, nil
}

/*
 * main function
 * Parse startup aurguments
 * Discover the network
 * Build the links between routers
 * Export the network to a KML file, if startup option is set
 *
 */
func main() {

	// Create logger
	log = logger.New("GoVisn", nil)
	log.AddWriter(logger.NewConsole(false))
	log.SetFormat(logger.FTIME | logger.FMICROS)

	flag.Parse() // Scan the arguments list
	if *versionFlag {
		fmt.Println("GoVision version", GOVISN_VERSION)
		return
	}

	if *debugFlag {
		log.SetLevel(logger.DEBUG)
	} else {
		log.SetLevel(logger.INFO)
	}
	log.Debug("Log Level set to DEBUG")

	log.Info("GoVision version %s", GOVISN_VERSION+
		" started")

	if *DbName != "govisnDiscoveredNet.db" {
		dbName = *DbName
	}
	if *loadDBFlag {
		loaddb(dbName)
		return
	}

	snmpPort := "161"
	snmpTarget := seed
	log.Debug("snmpPort=%s", snmpPort)
	log.Debug("snmpTarget=%s", snmpTarget)
	if len(snmpTarget) <= 0 {
		log.Fatal("environment variable not set: GOSNMP_TARGET")
	} else {
		if *debugFlag {
			log.Debug("snmpTarget= %s", snmpTarget)
		}
	}
	if len(snmpPort) <= 0 {
		log.Fatal("environment variable not set: GOSNMP_PORT")
	}
	port, _ := strconv.ParseUint(snmpPort, 10, 16)

	// GoSNMP struct
	params := &g.GoSNMP{
		Target:    snmpTarget,
		Port:      uint16(port),
		Community: *community,
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   2,
		Logger:    g.Default.Logger,
		MaxOids:   6,
	}
	log.Debug("params=%v", params)

	if *discoverFlag != "" {
		seed = *discoverFlag
		log.Debug("seed= %s", seed+" community= "+*community)

		// Open the database connection
		database, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatal("sql.Open() err: %v", err)
		}
		//Initialize the database
		database = initDB(log, database)

		// Discover network starting from the initial router
		discoveredRouters := make(map[string]bool)
		toDiscover := []string{seed}

		var router Router

		for len(toDiscover) > 0 {
			currentRouter := toDiscover[0]
			toDiscover = toDiscover[1:]

			if discoveredRouters[currentRouter] {
				continue
			}

			// Connect to the router
			params, err := connectToSNMP(
				currentRouter,
				*community,
			)
			if err != nil {
				log.Warn("Error connecting to %s: %v", currentRouter, err)
				continue
			}
			defer params.Conn.Close()

			// Get the Router's info
			router = getRouterInfo(
				log,
				params.Target,
				params,
				router,
				database,
			)

			// Discover routes
			routes, err := discoverRoutes(params, IP_ROUTE_NEXT_HOP_OID)
			if err != nil {
				log.Warn("Error discovering routes from %s: %v", currentRouter, err)
				continue
			}

			// Print discovered routes
			log.Debug("Routes from %s", currentRouter)
			for _, route := range routes {
				log.Debug(route)
				// Parse the route to find new router IP addresses if applicable
				toDiscover = append(toDiscover, route)
				//}
			}

			// Mark the current router as discovered
			discoveredRouters[currentRouter] = true
		}

		// Build Links
		database = buildLinks(log, database)

		defer database.Close()

	}

	var scannedRouters []ScannedRouter

	if *scanNetFlag != "" {
		seed = *scanNetFlag
		// Open the database connection
		database, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			log.Fatal("Error opening database: %s" + *DbName + "err:" + openErr.Error())
		}
		defer database.Close()

		if len(snmpTarget) <= 0 {
			log.Fatal("environment variable not set: GOSNMP_TARGET")
		} else {
			log.Debug("snmpTarget= %s", snmpTarget)
		}
		if len(snmpPort) <= 0 {
			log.Fatal("environment variable not set: GOSNMP_PORT")
		}
		port, _ := strconv.ParseUint(snmpPort, 10, 16)

		// GoSNMP struct
		params := &g.GoSNMP{
			Target:    snmpTarget,
			Port:      uint16(port),
			Community: *community,
			Version:   g.Version2c,
			Timeout:   time.Duration(2) * time.Second,
			Logger:    g.Default.Logger,
			MaxOids:   6,
		}

		// Scan the requested network for Router hosts
		scannedRouters = scanNet(log, seed, *community, params)
		log.Debug("scannedRouters= %v", scannedRouters)

		// Discover router information from list of scanned routers.
		// Open the database connection
		database, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatal("sql.Open() err: %v", err)
		}
		for i := 0; i < len(scannedRouters); i++ {

			// Discover the router's information and add to database
			params.Target = scannedRouters[i].IPAddress
			database = discover(log, scannedRouters[i].IPAddress, params, *maxHops, database)

			// Close database. Completed initialization and update of all tables, except Links.
			database.Close()

			// Open database. buildLinks joins Router and RouteTable tables.
			database, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal("sql.Open() err %v", err)
			}

			// Build Links
			database = buildLinks(log, database)

		}
		database.Close()

	}

	if *kmlFlag != "" {
		kmlFilename := *kmlFlag
		exportKML(log, kmlFilename, *DbName)
	}

	if *visualizeFlag {
		// Open the database containing the discovered network
		log.Info("Beginning Network Visualization.")

		databaseForRead, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			log.Fatal("Error opening databaseForRead %v", *DbName)
		}
		defer databaseForRead.Close()

		databaseForUpdate, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			log.Fatal("Error opening databaseForUpdate %v", *DbName)
		}
		defer databaseForUpdate.Close()

		// GoSNMP struct
		port, _ := strconv.ParseUint(snmpPort, 10, 16)
		params := &g.GoSNMP{
			Target:    snmpTarget,
			Port:      uint16(port),
			Community: *community,
			Version:   g.Version2c,
			Timeout:   time.Duration(2) * time.Second,
			Logger:    g.Default.Logger,
			MaxOids:   6,
		}
		databaseForRead = visualizeNetwork(log, databaseForRead, snmpTarget, params)
	}
	log.Info("GoVisn version %s", GOVISN_VERSION+" ending.")
}
