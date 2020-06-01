package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// Get Router Coordinates from routerArray
func getRouterCoordinatesIP(debugFlag bool, database *sql.DB, ToRouterIPIn string) (float32, float32, float32) {
	if debugFlag {
		fmt.Println("getRouterCoordiantesIP starting")
	}
	var x float32
	var y float32
	var z float32

	LinksTableRows, err := database.Query("SELECT RouterID  FROM RouterIP WHERE IpAddr = $1", ToRouterIPIn)
	if err != nil {
		log.Fatalln("LinksTableRows error", err.Error())
	}
	if debugFlag {
		fmt.Println("Successful LinksTableRows Query")
	}
	defer LinksTableRows.Close()

	var RouterIDLinks string
	var RouterIDRouter, GpsLat, GpsLong, GpsAlt string
	for LinksTableRows.Next() {
		LinksTableRows.Scan(&RouterIDLinks)

		RouterTableRows, err := database.Query("SELECT RouterID, GpsLat, GpsLong, GpsAlt FROM Routers WHERE RouterID = $1", RouterIDLinks)
		if err != nil {
			log.Fatalln("RouterTableRows error", err.Error())
		}
		defer LinksTableRows.Close()

		if debugFlag {
			fmt.Println("Successful RouterTable Query")
		}
		for RouterTableRows.Next() {
			RouterTableRows.Scan(&RouterIDRouter, &GpsLat, &GpsLong, &GpsAlt)
		}

		x1, parseErr := strconv.ParseFloat(GpsLong, 32)
		if parseErr != nil {
			log.Fatalln("x1 ParseFloat error", parseErr.Error())
		}
		x = (float32)(x1)

		y1, parseErr := strconv.ParseFloat(GpsLat, 32)
		if parseErr != nil {
			log.Fatalln("y1 ParseFloat error", parseErr.Error())
		}
		y = (float32)(y1)

		z1, parseErr := strconv.ParseFloat(GpsAlt, 32)
		if parseErr != nil {
			log.Fatalln("z1 ParseFloat error", parseErr.Error())
		}
		z = (float32)(z1)
	}
	if debugFlag {
		fmt.Println("getRouterCoordiantesIP ending")
	}

	return x, y, z
}
