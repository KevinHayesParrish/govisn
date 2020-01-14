package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//loaddbVersion is the file version number
const loadbVersion = "0.2.4"

// The V15NDiscoveredNetwork struct contains the discovered network, and it's sub-structs;
// essentially, the XML input file.
type V15NDiscoveredNetwork struct {
	XMLName xml.Name `xml:"V15N_Discovered_Network"`
	Text    string   `xml:",chardata"`
	// The Router struct, this contains
	// the router's Name, Descriptions, UpTime
	// Contact, Location and GPS Coordinates.
	// It also contains nested structs for Addresses
	// and Neighbor routers.
	Router []struct {
		Text   string `xml:",chardata"`
		System struct {
			Text        string `xml:",chardata"`
			Name        string `xml:"Name"`
			Description string `xml:"Description"`
			UpTime      string `xml:"Up_Time"`
			Contact     string `xml:"Contact"`
			Location    string `xml:"Location"`
			GPS         struct {
				Text      string `xml:",chardata"`
				Latitude  string `xml:"Latitude"`
				Longitude string `xml:"Longitude"`
				Altitude  string `xml:"Altitude"`
			} `xml:"GPS"`
		} `xml:"System"`
		Addresses struct {
			Text             string `xml:",chardata"`
			NetworkAddresses struct {
				Text      string   `xml:",chardata"`
				IPAddress []string `xml:"IP_Address"`
			} `xml:"Network_Addresses"`
			MediaAddresses struct {
				Text         string   `xml:",chardata"`
				MediaAddress []string `xml:"Media_Address"`
			} `xml:"Media_Addresses"`
		} `xml:"Addresses"`
		Neighbors struct {
			Text     string `xml:",chardata"`
			Neighbor []struct {
				Text               string `xml:",chardata"`
				DestinationAddress string `xml:"Destination_Address"`
				NextHop            string `xml:"Next_Hop"`
			} `xml:"Neighbor"`
		} `xml:"Neighbors"`
	} `xml:"Router"`
}

func loaddb(debug bool, networkXML string) {
	fmt.Println("loaddb version:", loadbVersion)
	fmt.Println("Loading database from XML document", networkXML) // FOR TESTING ONLY

	if debug {
		fmt.Println("Debug option selected")
	}

	// The struc which contains all the Routers in the XML input file.
	//	type Routers struct {
	//		XMLName     xml.Name `xml:"V15N_Discovered_Network"`
	//		NetworkName string   `xml:"V15N_Discovered_Network"`
	//		Routers     []Router `xml:"Router"`
	//	}

	// Open our xmlFile
	xmlFile, err := os.Open(networkXML)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened", networkXML)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	xmlFileBytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
	}

	// Initialize the network struct
	var network V15NDiscoveredNetwork

	// Unmarshal our byteArray which contains our discovered network
	err = xml.Unmarshal(xmlFileBytes, &network)
	if err != nil {
		fmt.Println(err)
		//		return
		panic(err)
	}

	// Open the database
	databaseName := *DbName + ".db"
	fmt.Println("dabaseName=", databaseName) //TESTING ONLY
	database, _ := sql.Open("sqlite3", databaseName)

	// Create Routers DB table
	//	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat REAL, GPSLong REAL, GpsAlt REAL)")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat TEXT, GPSLong TEXT, GpsAlt TEXT)")
	statement.Exec()

	// Add Routers to the database
	for i := 0; i < len(network.Router); i++ {
		statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

		if debug {
			fmt.Println("Router Name: ", network.Router[i].System.Name)
			fmt.Println("Description=", network.Router[i].System.Description)
			fmt.Println("Up_Time=", network.Router[i].System.UpTime)
			fmt.Println("GPS=", network.Router[i].System.GPS)
			fmt.Println("Addresses=", network.Router[i].Addresses)
			fmt.Println("Neighbors=", network.Router[i].Neighbors)
		}
		fmt.Println("Router Name: ", network.Router[i].System.Name)       // TESTING ONLY
		fmt.Println("Description=", network.Router[i].System.Description) // TESTING ONLY
		fmt.Println("Up_Time=", network.Router[i].System.UpTime)          // TESTING ONLY
		fmt.Println("GPS=", network.Router[i].System.GPS)                 // TESTING ONLY

		SystemName := network.Router[i].System.Name
		RouterIDUint32 := crc32.ChecksumIEEE([]byte(SystemName))
		SystemDesc := network.Router[i].System.Description
		UpTime := network.Router[i].System.UpTime
		Contact := network.Router[i].System.Contact
		Location := network.Router[i].System.Location
		GpsLat := network.Router[i].System.GPS.Latitude
		GpsLong := network.Router[i].System.GPS.Longitude
		GpsAlt := network.Router[i].System.GPS.Altitude

		statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

		// Create RouterIp DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER, IPAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterIp (RouterID, IPAddr) VALUES (?, ?)")

		// Add IP Addresses to current router
		for j := 0; j < len(network.Router[i].Addresses.NetworkAddresses.IPAddress); j++ {
			IPAddr := network.Router[i].Addresses.NetworkAddresses.IPAddress[j]
			statement.Exec(strconv.Itoa(int(RouterIDUint32)), IPAddr) // Add IP Address
		}

		//	Create RouterMac DB table
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER, MacAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterMac (RouterID, MacAddr) VALUES (?, ?)")

		// Add Media Addresses to current router
		for k := 0; k < len(network.Router[i].Addresses.MediaAddresses.MediaAddress); k++ {
			MacAddr := network.Router[i].Addresses.MediaAddresses.MediaAddress[k]
			statement.Exec(strconv.Itoa(int(RouterIDUint32)), MacAddr) // Add Media Address
		}

		//	Create Links DB table
		//		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, FromRouter TEXT, ToRouter TEXT)")
		//		statement, _ = database.Prepare("INSERT INTO Links (LinkID, FromRouter, ToRouter) VALUES (?, ?, ?)")
		statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, RouterName TEXT, DestinationName TEXT, DestinationIP TEXT, NextHopName TEXT, NextHopIP TEXT)")
		if err != nil {
			fmt.Println("Error creating Links table.")
			panic(err)
		}
		statement.Exec()

		statement, err = database.Prepare("INSERT INTO Links (LinkID, RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("Error inserting row into Links table.")
			panic(err)
		}

		// Add Link records to Links table
		var dest string
		var nextHop string
		var DestinationName string
		var NextHopName string

		if debug {
			fmt.Println("Adding link records to Links table.")
			//fmt.Println("network.Router[", i, "]=", network.Router[i])
		}
		for l := 0; l < len(network.Router[i].Neighbors.Neighbor); l++ {
			// Don't add link row for loopback interface
			if network.Router[i].Neighbors.Neighbor[l].DestinationAddress == "127.0.0.0" {
				continue
			}
			if network.Router[i].Neighbors.Neighbor[l].DestinationAddress == "127.0.0.1" {
				continue
			}
			// Don't add link row for Multicast route
			if network.Router[i].Neighbors.Neighbor[l].DestinationAddress == "224.0.0.0" {
				continue
			}
			// Don't add link row for broadcast
			if network.Router[i].Neighbors.Neighbor[l].DestinationAddress == "255.255.255.255" {
				continue
			}

			//* Add a set of Links to the database
			dest = network.Router[i].Neighbors.Neighbor[l].DestinationAddress
			nextHop = network.Router[i].Neighbors.Neighbor[l].NextHop
			destToNextHopLinkStr := dest + nextHop // directional link from dest to nextHop
			nextHopToDestStr := nextHop + dest     // directional link from nextHop to dest

			// add direction link from dest to nextHop to the database
			destToNextHopLinkUint32 := crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))

			// Lookup DestinationName
			if debug {
				fmt.Println("\nCalling getRouterNameUsingIP with Destination", dest)
			}
			DestinationName = getRouterNameUsingIP(debug, dest, network)
			if debug {
				fmt.Println(" Returned DestinationName", DestinationName)
			}
			if DestinationName == "Not Found" {
				DestinationName = "Unknown"
				if debug {
					fmt.Println("router name with destination IP of", dest, " Not Found.")
				}
			}
			// Lookup NextHopName
			if debug {
				fmt.Println("\nCalling getRouterNameUsingIP with NextHop", nextHop)
			}
			NextHopName = getRouterNameUsingIP(debug, nextHop, network)
			if debug {
				fmt.Println(" Returned NextHopName", NextHopName)
			}
			if NextHopName == "Not Found" {
				NextHopName = "Unknown"
				if debug {
					fmt.Println("router name with Next Hop IP of", nextHop, " Not Found.")
				}
			}

			//			statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), DestinationName, dest, NextHopName, nextHop)
			if debug {
				fmt.Println("Adding link row with fields =", SystemName, DestinationName, dest, NextHopName, nextHop)
			}
			statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), SystemName, DestinationName, dest, NextHopName, nextHop)

			// add direction link from nextHop to dest to the database
			nextHopToDestUint32 := crc32.ChecksumIEEE([]byte(nextHopToDestStr))
			statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), SystemName, DestinationName, nextHop, NextHopName, dest)
		}
	}

}

func getRouterNameUsingIP(debug bool, ipAddress string, network V15NDiscoveredNetwork) string {
	var routerName string
	//	var network V15NDiscoveredNetwork
	routerName = "Not Found"

	if debug {
		fmt.Println("getRouterNameUsingIP")
		fmt.Println(" network.Router length is", len(network.Router))
	}

	for i := 0; i < len(network.Router); i++ {
		if debug {
			fmt.Println(" network.Router[i].System.Name=", network.Router[i].System.Name)
			fmt.Println(" i=", i)
			//fmt.Println(" network.Router[i]", network.Router[i])
			fmt.Println(" len(network.Router[i].Addresses.NetworkAddresses.IPAddress=", len(network.Router[i].Addresses.NetworkAddresses.IPAddress))
		}

		for j := 0; j < len(network.Router[i].Addresses.NetworkAddresses.IPAddress); j++ {
			if debug {
				fmt.Println(" j=", j)
				fmt.Println(" ipAddress=", ipAddress)
				fmt.Println(" network.Router[i].Addresses.NetworkAddresses.IPAddress[j]=", network.Router[i].Addresses.NetworkAddresses.IPAddress[j])
				//fmt.Println(" network.Router[i]", network.Router[i])
			}
			if network.Router[i].Addresses.NetworkAddresses.IPAddress[j] == ipAddress {
				routerName = network.Router[i].System.Name
				break
			}
		}
	}
	// router name not found in network
	if debug {
		fmt.Println(" Returning routerName of", routerName)
		fmt.Println()
	}
	return routerName
}
