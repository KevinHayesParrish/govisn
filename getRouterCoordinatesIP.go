package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// Get Router Coordinates from routerArray
func getRouterCoordinatesIP(debugFlag bool, database *sql.DB, ToRouterIPIn string) (float32, float32, float32) {
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

	/*
		for i := 0; i < len(routers); i++ {
			//		if routerArray[i].System.Name == "" {
			//			break // end of routerArray entries
			//		}
			//		if routerArray[i].System.Name == routerName {
				if routers[i].System.Name == routerName {
				//			x = routerArray[i].System.Coordinates.X
				//			x1, err := strconv.ParseFloat(routerArray[i].System.GPS.Longitude, 32)
				x1, err := strconv.ParseFloat(routers[i].System.GPS.Longitude, 32)
				if err != nil {
					panic(err)
				}
				x = (float32)(x1)
				//			y = routerArray[i].System.Coordinates.Y
				//			y1, err := strconv.ParseFloat(routerArray[i].System.GPS.Latitude, 32)
				y1, err := strconv.ParseFloat(routers[i].System.GPS.Latitude, 32)
				if err != nil {
					panic(err)
				}
				y = (float32)(y1)
				//			z = routerArray[i].System.Coordinates.Z
				//			z1, err := strconv.ParseFloat(routerArray[i].System.GPS.Altitude, 32)
				z1, err := strconv.ParseFloat(routers[i].System.GPS.Altitude, 32)
				if err != nil {
					panic(err)
				}
				z = (float32)(z1)
				break
			}
		}
	*/
	return x, y, z
}
