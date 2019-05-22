package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//ViewnetVersion is the file version number
const ViewnetVersion = "0.1.0"

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("d", false, "Print Debug statements.")

//DbName is the name of the discovered network database file
var DbName = flag.String("f", "kpnetviz.db", "Name of the discovered network database")

// DiscoveredNetwork is the network that was discovered and the subject of the visualization.
type DiscoveredNetwork struct {
	Routers []Router
}

// Router is the structure representing a network router
type Router struct {
	System struct {
		Text        string
		Name        string
		Description string
		UpTime      string
		Contact     string
		Location    string
		GPS         struct {
			Text      string
			Latitude  string
			Longitude string
			Altitude  string
		}
	}
	Addresses struct {
		Text             string
		NetworkAddresses struct {
			Text      string
			IPAddress []string
		}
		MediaAddresses struct {
			Text         string
			MediaAddress string
		}
	}
	Neighbors struct {
		Text     string
		Neighbor []struct {
			Text               string
			DestinationAddress string
			NextHop            string
		}
	}
}

func main() {
	flag.Parse() // Scan the arguments list
	fmt.Println("viewnet version:", ViewnetVersion)
	if *versionFlag {
		return
	}
	if *debugFlag {
		fmt.Println("Debug option selected")
	}

	// Open the database containing the discovered network
	database, _ := sql.Open("sqlite3", *DbName)
	routers, _ := database.Query("SELECT RouterID, SystemName, SystemDesc FROM Routers")

	var RouterID int
	var SystemName string
	var SystemDesc string
	if *debugFlag {
		for routers.Next() {
			routers.Scan(&RouterID, &SystemName, &SystemDesc)
			fmt.Println(strconv.Itoa(RouterID) + ": " + SystemName + " " + SystemDesc)
		}
	}

}
