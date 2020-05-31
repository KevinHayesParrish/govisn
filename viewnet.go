package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"

	//	"github.com/g3n/g3nd/material"

	_ "github.com/mattn/go-sqlite3"
)

/*
* TODO:
 */

//ViewnetVersion is the file version number
const ViewnetVersion = "0.8.10"

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
var discoverFlag = flag.String("di", "", "Discover a network using seed IP Address")

var seed = "127.0.0.1"

var community = flag.String("co", "public", "SNMP Community ReadOnly String")
var maxHops = flag.String("m", "0", "Scope of discovery. Maximum number of Hops from seed. (Default:10)")
var visualizeFlag = flag.Bool("vi", false, "Visualize the Network.")

//scanNetFlag is the startup option to scan the network for SNMP capable routers.
var scanNetFlag = flag.String("s", "", "Scan the network for SNMP capable routers.\nOnce the network is scanned, the list of found routers\nwill be queried and their information added to the database.")

//routerRadius is the radius of the 3D object representing a network router
const routerRadius float64 = 0.5

//globeRadius is the radius of the 3D object representing the earth
const globeRadius float64 = 63.7

func main() {
	flag.Parse() // Scan the arguments list
	fmt.Println("viewnet version:", ViewnetVersion)
	if *versionFlag {
		return
	}
	if *debugFlag {
		fmt.Println("Debug option selected")
	}
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
		if *debugFlag {
			fmt.Println("seed=", seed, "community=", *community)
		}

		// Open the database connection
		database, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatalf("sql.Open() err: %v", err)
		}

		// Discover the network
		database = discover(*debugFlag, dbName, seed, *community, *maxHops, database)

		// Close database. Completed initialization and update of all tables, except Links.
		database.Close()

		// Open database. buildLinks joins Router and RouteTable tables.
		database, err = sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatalf("sql.Open() err: %v", err)
		}

		// Build Links
		database = buildLinks(*debugFlag, database)
		database.Close()

	}
	/*
		// Open the database containing the discovered network
		databaseForRead, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			fmt.Println("Error opening databaseForRead", *DbName)
			log.Fatal(openErr)
		}
		defer databaseForRead.Close()

		databaseForUpdate, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			fmt.Println("Error opening databaseForUpdate", *DbName)
			log.Fatal(openErr)
		}
		defer databaseForUpdate.Close()
	*/

	var scannedRouters []ScannedRouter
	if *scanNetFlag != "" {
		seed = *scanNetFlag
		// Open the database connection
		database, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			fmt.Println("Error opening database", *DbName)
			log.Fatal(openErr)
		}
		defer database.Close()

		scannedRouters = scanNet(*debugFlag, seed, *community, database)
		if *debugFlag {
			fmt.Println("scnnedRouters=", scannedRouters)
		}
	}

	if *visualizeFlag {
		// Open the database containing the discovered network
		databaseForRead, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			fmt.Println("Error opening databaseForRead", *DbName)
			log.Fatal(openErr)
		}
		defer databaseForRead.Close()

		databaseForUpdate, openErr := sql.Open("sqlite3", *DbName)
		if openErr != nil {
			fmt.Println("Error opening databaseForUpdate", *DbName)
			log.Fatal(openErr)
		}
		defer databaseForUpdate.Close()
		databaseForRead = visualizeNetwork(*debugFlag, databaseForRead)
	}
}

const constX = math.Pi / 180

// Rad converts degrees to radians
func Rad(d float64) float64 { return d * constX }

// Deg converts radians to degrees
func Deg(r float64) float64 { return r / constX }
