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
		log.Fatal("databaseForRead Query error %v", queryErr)
	}
	log.Debug("Successful Routers table Select")

	var Name string
	var Description string
	var Location string
	var GpsLat string
	var GpsLong string
	var GpsAlt string

	k := kml.KML(
		kml.Placemark(
			kml.Name("GoVision"),
			kml.Description("KML representation of a network."),
		),
	)
	kmlFile, err := os.Create(kmlFilename)
	if err != nil {
		log.Fatal("Cannot Create %s", kmlFilename+" error="+err.Error())
	}

	var GpsLongFloat float64
	var GpsLatFloat float64
	var GpsAltFloat float64

	for routerRows.Next() {
		routerRows.Scan(&Name, &Description, &Location, &GpsLat, &GpsLong, &GpsAlt)

		s, _ := strconv.ParseFloat(GpsLong, 64)
		GpsLongFloat = s
		s, _ = strconv.ParseFloat(GpsLat, 64)
		GpsLatFloat = s
		s, _ = strconv.ParseFloat(GpsAlt, 64)
		GpsAltFloat = s

		k.Add(
			kml.Placemark(
				kml.Name(Name),
				kml.Description(Description),
				kml.Point(
					kml.Coordinates(kml.Coordinate{Lon: GpsLongFloat, Lat: GpsLatFloat, Alt: GpsAltFloat}),
				),
			),
		)
	}
	if err := k.WriteIndent(kmlFile, "", "  "); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("exportKML version %s ", EXPORTKMLVERSION+" ending")
	return
}
