package main

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/twpayne/go-kml"

	"github.com/g3n/engine/util/logger"
)

// EXPORTKMLVERSION is the version number of this function's source code file.
const EXPORTKMLVERSION = "0.1.0"

func exportKML(log *logger.Logger, kmlFilename string, DbName string) {
	log.Info("exportKML version %s ", EXPORTKMLVERSION+" started")
	log.Debug("kmlFilename = %s", kmlFilename)
	log.Debug("DbName = %s", DbName)

	databaseForRead, openErr := sql.Open("sqlite3", DbName)
	if openErr != nil {
		log.Fatal("Error opening databaseForRead %v", DbName)
	}
	defer databaseForRead.Close()

	// Retrieve the Routers table
	routerRows, queryErr := databaseForRead.Query("SELECT Name, Description, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		log.Fatal("databaseForRead Query error %v; %v does not contain Router table entries.", queryErr, DbName)
	}
	log.Debug("Successful Routers table Select")

	var Name string
	var Description string
	var Location string
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	doc := kml.Document(
		kml.Name("GoVisn"),
		kml.Description("KML representation of a network."),
	)
	kmlFile, err := os.Create(kmlFilename)
	if err != nil {
		log.Fatal("Cannot Create %s", kmlFilename+" error="+err.Error())
	}

	var GpsLongFloat float64
	var GpsLatFloat float64
	var GpsAltFloat float64

	// Add the routers to the KML document
	for routerRows.Next() {
		routerRows.Scan(&Name, &Description, &Location, &GpsLat, &GpsLong, &GpsAlt)

		s, _ := strconv.ParseFloat(GpsLong, 64)
		GpsLongFloat = s
		s, _ = strconv.ParseFloat(GpsLat, 64)
		GpsLatFloat = s
		s, _ = strconv.ParseFloat(GpsAlt, 64)
		GpsAltFloat = s

		doc.Add(
			kml.Placemark(
				kml.Name(Name),
				kml.Description(Description),
				kml.Point(
					kml.Coordinates(kml.Coordinate{Lon: GpsLongFloat, Lat: GpsLatFloat, Alt: GpsAltFloat}),
				),
			),
		)
	}

	// Add the links to the KML document
	// Retrieve the Links table
	linkRows, queryErr := databaseForRead.Query("SELECT LinkID, FromRouterName, ToRouterName FROM Links")
	if queryErr != nil {
		log.Fatal("databaseForRead Query error %v", queryErr)
	}
	log.Debug("Successful Links table Select")

	var FromRouterName, ToRouterName string
	var LinkID int
	for linkRows.Next() {
		linkRows.Scan(&LinkID, &FromRouterName, &ToRouterName)

		var routerName string
		// retrieve From Router coordinates
		routerName = FromRouterName
		fromRouterRows, queryErr := databaseForRead.Query("SELECT Name, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name=?", routerName)
		if queryErr != nil {
			log.Fatal("databaseForRead Query error %v", queryErr)
		}
		var FromGpsLongFloat float64
		var FromGpsLatFloat float64
		var FromGpsAltFloat float64
		for fromRouterRows.Next() {
			fromRouterRows.Scan(&Name, &GpsLat, &GpsLong, &GpsAlt)
			s, _ := strconv.ParseFloat(GpsLong, 64)
			FromGpsLongFloat = s
			s, _ = strconv.ParseFloat(GpsLat, 64)
			FromGpsLatFloat = s
			s, _ = strconv.ParseFloat(GpsAlt, 64)
			FromGpsAltFloat = s
		}

		// retrieve ToRouter coordinates
		var ToGpsLongFloat float64
		var ToGpsLatFloat float64
		var ToGpsAltFloat float64
		routerName = ToRouterName
		toRouterRows, queryErr := databaseForRead.Query("SELECT Name, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name=?", routerName)
		if queryErr != nil {
			log.Fatal("databaseForRead Query error %v", queryErr)
		}
		for toRouterRows.Next() {
			toRouterRows.Scan(&Name, &GpsLat, &GpsLong, &GpsAlt)
			s, _ := strconv.ParseFloat(GpsLong, 64)
			ToGpsLongFloat = s
			s, _ = strconv.ParseFloat(GpsLat, 64)
			ToGpsLatFloat = s
			s, _ = strconv.ParseFloat(GpsAlt, 64)
			ToGpsAltFloat = s
		}

		doc.Add(
			kml.Placemark(
				kml.Name(strconv.Itoa(LinkID)),
				kml.LineString(
					kml.Coordinates(kml.Coordinate{Lon: FromGpsLongFloat, Lat: FromGpsLatFloat, Alt: FromGpsAltFloat},
						kml.Coordinate{Lon: ToGpsLongFloat, Lat: ToGpsLatFloat, Alt: ToGpsAltFloat}),
				),
			),
		)

	}

	// Write the KML document to the file
	if err := doc.WriteIndent(kmlFile, "", "  "); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("exportKML version %s ", EXPORTKMLVERSION+" ending")
	//	return
}
