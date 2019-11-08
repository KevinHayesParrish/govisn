package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//loaddbVersion is the file version number
const loadbVersion = "0.1.1"

func loaddb(networkXML string) {
	fmt.Println("loaddb version:", loadbVersion)
	fmt.Println("Loading database from XML document", networkXML) // FOR TESTING ONLY

	type V15NDiscoveredNetwork struct {
		XMLName xml.Name `xml:"V15N_Discovered_Network"`
		Text    string   `xml:",chardata"`
		Router  struct {
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

	var discoveredNetworkXML V15NDiscoveredNetwork

	fmt.Println("discoveredNetworkXML=", discoveredNetworkXML)
	// Open our xmlFile
	xmlFile, err := os.Open(networkXML)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened", networkXML)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	discoveredNetworkBytes, _ := ioutil.ReadAll(xmlFile)

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'discoveredNetworkXML' which we defined above
	xml.Unmarshal(discoveredNetworkBytes, &discoveredNetworkXML)

	fmt.Println("discoveredNetworkBytes=", discoveredNetworkBytes) // TESTING ONLY

    for i := 0; i < len(users.Users); i++ {
        fmt.Println("User Type: " + users.Users[i].Type)
    }

	/*
		database, _ := sql.Open("sqlite3", "./networkXML")

		//*	Add Routers table to DB
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat REAL, GPSLong REAL, GpsAlt REAL)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

		// add router to the database
		SystemName := "media"
		RouterIDUint32 := crc32.ChecksumIEEE([]byte(SystemName))
		SystemDesc := "Hardware: x86 Family 15 Model 2 Stepping 7 AT/AT COMPATIBLE - Software: Windows Version 5.2 (Build 3790 Uniprocessor Free)"
		UpTime := "18h 4m 9s 40"
		Contact := "Kevin Parrish"
		Location := "9218 Faxon Place, Elk Grove, CA 95624 USA"
		GpsLat := "38.419471"
		GpsLong := "-121.357212"
		GpsAlt := "10.668"
		statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

		//*	Add RouteTable table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, NextHop TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, NextHop) VALUES (?, ?, ?)")

		//*	Add RouterIP table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER, IpAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr) VALUES (?, ?)")

		//*	Add RouterMac table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterMac (RourterID, MacAddr) VALUES (?, ?)")

		//*	Add Links table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, FromRouter TEXT, ToRouter TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO Links (LinkID, FromRouter, ToRouter) VALUES (?, ?, ?)")

		//* Add a set of Links to the database
		var dest string
		var nextHop string
		dest = "media"
		nextHop = "router"
		destToNextHopLinkStr := dest + nextHop // directional link from dest to nextHop
		nextHopToDestStr := nextHop + dest     // directional link from nextHop to dest

		// add direction link from dest to nextHop to the database
		destToNextHopLinkUint32 := crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
		statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

		// add direction link from nextHop to dest to the database
		nextHopToDestUint32 := crc32.ChecksumIEEE([]byte(nextHopToDestStr))
		statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

		//* Add another set of Links to the database
		dest = "router"
		nextHop = "wan router"
		// add direction link from dest to nextHop to the database
		destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
		nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

		destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
		statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

		// add direction link from nextHop to dest to the database
		nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
		statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

		//* Add another set of Links to the database
		dest = "wan router"
		nextHop = "hub"
		// add direction link from dest to nextHop to the database
		destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
		nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

		destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
		statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

		// add direction link from nextHop to dest to the database
		nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
		statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

		//* Add another set of Links to the database
		dest = "wan router"
		nextHop = "old-country-road"
		// add direction link from dest to nextHop to the database
		destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
		nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

		destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
		statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

		// add direction link from nextHop to dest to the database
		nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
		statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

		//* Add another set of Links to the database
		dest = "wan router"
		nextHop = "fukui"
		// add direction link from dest to nextHop to the database
		destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
		nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

		destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
		statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

		// add direction link from nextHop to dest to the database
		nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
		statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

		//* Add another set of Links to the database
		dest = "wan router"
		nextHop = "amsterdam"
		// add direction link from dest to nextHop to the database
		destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
		nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

		destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
		statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

		// add direction link from nextHop to dest to the database
		nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
		statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

		//* Print contents of the db
		rows, _ := database.Query("SELECT LinkID, FromRouter, ToRouter FROM Links")

		var LinkID int
		var FromRouter string
		var ToRouter string
		for rows.Next() {
			//		fmt.Println(idUint32 + ": " + FromRouter + " " + ToRouter)
			rows.Scan(&LinkID, &FromRouter, &ToRouter)
			fmt.Println(strconv.Itoa(LinkID) + ": " + FromRouter + " " + ToRouter)
		}

		//* Add 3DCoordinates table to the database.
		//* This table is used dynamically througout the visualization to hold the 3D
		//* coordinates of the routers.
		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Coordinates (RouterID INTEGER, SystemName TEXT, X3D REAL, Y3D REAL, Z3D REAL)")
		if err != nil {
			fmt.Println("Error preparing 3DCoordinates table Create statement. Result=", statement)
			log.Fatal(err)
		}
		statement.Exec()

		statement, err = database.Prepare("INSERT INTO Coordinates (RouterID, SystemName, X3D, Y3D, Z3D) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("Error preparing Coordinates insert statement. Result=", statement)
			log.Fatal(err)
		}
		RouterIDUint32 = 589093411
		SystemName = "wan-router"
		X3D := "0.0"
		Y3D := "0.0"
		Z3D := "0.0"
		statement.Exec(X3D, Y3D, Z3D) // Add 3D Coordinates
	*/
}
