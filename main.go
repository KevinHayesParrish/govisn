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
	g "github.com/soniah/gosnmp"
)

/*
* TODO:
 */

//GOVISNVERSION is the file version number
const GOVISNVERSION = "0.10.1"

var log *logger.Logger

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("de", false, "Print Debug statements.")
var sampleNetworkDB = flag.Bool("cr", false, "Create a sample database.")
var loadDBFlag = flag.Bool("l", false, "Load a database from an XML document.")

var dbName = "govisnDiscoveredNet.db"

// DbName is the name of the discovered network database file or name of XML input file
var DbName = flag.String("f", "govisnDiscoveredNet.db", "Name of the discovered network database -or-\nName of the XML input file, if combined with -l option.")

//testArangodb is the startup option to test accessing an ArangoDB database
var testArangoDb = flag.Bool("a", false, "Test opening an ArangoDB database")

//discoverFlag is the option to discover a network
var discoverFlag = flag.String("di", "", "Discover a network using a seed IP Address")

var seed = "127.0.0.1"

var community = flag.String("co", "public", "SNMP Community ReadOnly String")
var maxHops = flag.String("m", "0", "Scope of discovery. Maximum number of Hops from seed.")
var visualizeFlag = flag.Bool("vi", false, "Visualize the Network.")

//scanNetFlag is the startup option to scan the network for SNMP capable routers.
var scanNetFlag = flag.String("s", "", "Scan the CIDR network for SNMP capable routers.\nCIDR format = x.x.x.x/n. ex: 192.168.1.0/24\nOnce the network is scanned, the list of found routers\nwill be queried and their information added to the database.")

//routerRadius is the radius of the 3D object representing a network router
const routerRadius float64 = 0.5

//globeRadius is the radius of the 3D object representing the earth
const globeRadius float64 = 63.7

func main() {

	// Create logger
	log = logger.New("GoVisn", nil)
	log.AddWriter(logger.NewConsole(false))
	log.SetFormat(logger.FTIME | logger.FMICROS)

	flag.Parse() // Scan the arguments list
	if *versionFlag {
		fmt.Println("GoVision version", GOVISNVERSION)
		return
	}

	if *debugFlag {
		log.SetLevel(logger.DEBUG)
	} else {
		log.SetLevel(logger.INFO)
	}
	log.Debug("Log Level set to DEBUG")

	log.Info("GoVision version %s", GOVISNVERSION+
		" started")

	if *sampleNetworkDB {
		createsampledb()
	}
	if *testArangoDb {
		testarango()
	}
	if *DbName != "govisnDiscoveredNet.db" {
		dbName = *DbName
	}
	if *loadDBFlag {
		loaddb(*debugFlag, dbName)
		return
	}

	if *discoverFlag != "" {
		seed = *discoverFlag
		//		if *debugFlag {
		log.Debug("seed= %s", seed+" community= %v"+*community)
		//		}

		snmpPort := "161"
		snmpTarget := seed
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
			Logger:    nil,
			MaxOids:   6,
		}

		// Open the database connection
		database, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatal("sql.Open() err: %v", err)
		}

		// Discover the network
		database = discover(*debugFlag, log, dbName, seed, *community, params, *maxHops, database)

		// Close database. Completed initialization and update of all tables, except Links.
		database.Close()

		// Open database. buildLinks joins Router and RouteTable tables.
		database, err = sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatal("sql.Open() err: %v", err)
		}

		// Build Links
		database = buildLinks(*debugFlag, log, database)
		database.Close()

	}

	snmpPort := "161"
	snmpTarget := seed

	var scannedRouters []ScannedRouter

	if *scanNetFlag != "" {
		seed = *scanNetFlag
		// Open the database connection
		database, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			log.Fatal("Error opening database: %s" + *DbName + "err %s" + openErr.Error())
		}
		defer database.Close()

		//		snmpPort := "161"
		//		snmpTarget := seed
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
			Logger:    nil,
			MaxOids:   6,
		}

		// Scan the requested network for Router hosts
		//		scannedRouters = scanNet(*debugFlag, seed, *community, *params)
		scannedRouters = scanNet(*debugFlag, log, seed, *community, *params)
		//		if *debugFlag {
		//			fmt.Println("scnnedRouters=", scannedRouters)
		//		}
		log.Debug("scnnedRouters= %v", scannedRouters)

		// Discover router information from list of scanned routers.
		// Open the database connection
		database, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatal("sql.Open() err: %v", err)
		}
		for i := 0; i < len(scannedRouters); i++ {

			// Discover the router's information and add to database
			params.Target = scannedRouters[i].IPAddress
			database = discover(*debugFlag, log, dbName, scannedRouters[i].IPAddress, *community, params, *maxHops, database)

			// Close database. Completed initialization and update of all tables, except Links.
			database.Close()

			// Open database. buildLinks joins Router and RouteTable tables.
			database, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal("sql.Open() err %v", err)
			}

			// Build Links
			database = buildLinks(*debugFlag, log, database)

		}
		database.Close()

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
			Logger:    nil,
			MaxOids:   6,
		}
		databaseForRead = visualizeNetwork(*debugFlag, log, databaseForRead, snmpTarget, *community, params)
	}
	log.Info("GoVisn version %s", GOVISNVERSION+" ending.")
}

const constX = math.Pi / 180

// Rad converts degrees to radians
func Rad(d float64) float64 { return d * constX }

// Deg converts radians to degrees
func Deg(r float64) float64 { return r / constX }
