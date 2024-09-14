// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// LOAD_DB_VERSION is the file version number
const LOAD_DB_VERSION = "0.2.5"

/*
 * The V15NDiscoveredNetwork struct contains the discovered network, and it's sub-structs;
 * essentially, the XML input file.
 */
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

func loaddb(networkXML string) {
	/*
	 * TODO
	 *    replace Println with log.DEBUG, etc.
	 */

	log.Info("func loaddb version %s started.", LOAD_DB_VERSION)
	log.Info("Loading database from XML document %s", networkXML)

	log.Debug("Debug option selected")

	// Open our xmlFile
	xmlFile, err := os.Open(networkXML)
	if err != nil {
		log.Warn("Open XML file error %v", err.Error())
	}
	log.Info("Successfully Opened %s", networkXML)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	xmlFileBytes, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
	}

	// Initialize the network struct
	var network V15NDiscoveredNetwork

	// Unmarshal our byteArray which contains our discovered network
	err = xml.Unmarshal(xmlFileBytes, &network)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Open the database
	databaseName := *DbName + ".db"
	log.Debug("databaseName=%s", databaseName)
	database, _ := sql.Open("sqlite3", databaseName)

	// Create Routers DB table
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat TEXT, GPSLong TEXT, GpsAlt TEXT)")
	statement.Exec()

	// Add Routers to the database
	for i := 0; i < len(network.Router); i++ {
		statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

		log.Debug("Router Name: %s", network.Router[i].System.Name)
		log.Debug("Description: %s", network.Router[i].System.Description)
		log.Debug("Up_Time: %s", network.Router[i].System.UpTime)
		log.Debug("GPS: %s", network.Router[i].System.GPS)
		log.Debug("Addresses: %s", network.Router[i].Addresses)
		log.Debug("Neighbors: %s", network.Router[i].Neighbors)

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
		statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, RouterName TEXT, DestinationName TEXT, DestinationIP TEXT, NextHopName TEXT, NextHopIP TEXT)")
		if err != nil {
			log.Fatal("Error creating Links table.\n%v", err.Error())
		}
		statement.Exec()

		statement, err = database.Prepare("INSERT INTO Links (LinkID, RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Error inserting row into Links table.\n%v", err.Error())
		}

		// Add Link records to Links table
		var dest string
		var nextHop string
		var DestinationName string
		var NextHopName string

		log.Debug("Adding link records to Links table.")
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
			log.Debug("\nCalling getRouterNameUsingIP with Desitination %s", dest)

			DestinationName = getRouterNameUsingIP(dest, network)
			log.Debug("Returned DestinationName %s", DestinationName)

			if DestinationName == "Not Found" {
				DestinationName = "Unknown"
				log.Debug("router name with destination IP of %s Not Found.", dest)
			}
			// Lookup NextHopName
			log.Debug("\nCalling getRouterNameUsingIP with NextHop %s", nextHop)

			NextHopName = getRouterNameUsingIP(nextHop, network)
			log.Debug(" Returned NextHopName %s", NextHopName)

			if NextHopName == "Not Found" {
				NextHopName = "Unknown"
				log.Debug("router name with Next Hop IP of %s Not Found.", nextHop)
			}

			log.Debug("Adding link row with fields: %s, %s, %s, %s, %s", SystemName, DestinationName, dest, NextHopName, nextHop)

			statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), SystemName, DestinationName, dest, NextHopName, nextHop)

			// add direction link from nextHop to dest to the database
			nextHopToDestUint32 := crc32.ChecksumIEEE([]byte(nextHopToDestStr))
			statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), SystemName, DestinationName, nextHop, NextHopName, dest)
		}
	}

}

func getRouterNameUsingIP(ipAddress string, network V15NDiscoveredNetwork) string {
	var routerName string
	routerName = "Not Found"

	log.Debug("func getRouterNameUsingIP started")
	log.Debug(" network.Router length is %d", len(network.Router))

	for i := 0; i < len(network.Router); i++ {
		log.Debug(" network.Router[i].System.Name= %s", network.Router[i].System.Name)
		log.Debug(
			" i= %d", i,
		)
		log.Debug(
			" network.Router[i].Addresses.NetworkAddresses.IPAddress= %d",
			len(network.Router[i].Addresses.NetworkAddresses.IPAddress),
		)

		for j := 0; j < len(network.Router[i].Addresses.NetworkAddresses.IPAddress); j++ {
			log.Debug(" j=%d", j)
			log.Debug(" ipAddress=%s", ipAddress)
			log.Debug(" network.Router[i].Addresses.NetworkAddresses.IPAddress[j]=%s", network.Router[i].Addresses.NetworkAddresses.IPAddress[j])

			if network.Router[i].Addresses.NetworkAddresses.IPAddress[j] == ipAddress {
				routerName = network.Router[i].System.Name
				break
			}
		}
	}
	// router name not found in network
	log.Info(" Returning routerName of %s", routerName)

	return routerName
}
