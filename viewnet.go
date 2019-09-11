package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	_ "github.com/mattn/go-sqlite3"
)

/*
* TODO:
*	1. Add the router's 3D coordinates to the Router table.
*		This is for use when adding the links to the 3D scene.
 */

//ViewnetVersion is the file version number
const ViewnetVersion = "0.4.4"

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("d", false, "Print Debug statements.")
var sampleNetworkDB = flag.Bool("c", false, "Create a sample database.")

//DbName is the name of the discovered network database file
var DbName = flag.String("f", "discoverednetwork.db", "Name of the discovered network database")

//testArangodb is the startup option to test accessing an ArangoDB database
var testArangoDb = flag.Bool("a", false, "Test opening an ArangoDB database")

//routerRadius is the radius of the 3D object representing a network router
const routerRadius float64 = 0.5

//globeRadius is the radius of the 3D object representing the earth
//const globeRadius float64 = 1.5
const globeRadius float64 = 65.0

// Router is the structure representing a network router
type Router struct {
	System struct {
		RouterID    int
		Name        string
		Description string
		UpTime      string
		Contact     string
		Location    string
		GPS         struct {
			Latitude  string
			Longitude string
			Altitude  string
		}
		Coordinates struct {
			X float32
			Y float32
			Z float32
		}
	}
	Addresses struct {
		NetworkAddresses struct {
			IPAddress []string
		}
		MediaAddresses struct {
			MediaAddress string
		}
	}
	Neighbors struct {
		Neighbor []struct {
			DestinationAddress string
			NextHop            string
		}
	}
}

// Link is the structure representing a network link between two routers
type Link struct {
	LinkID     int
	FromRouter string
	ToRouter   string
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
	if *sampleNetworkDB {
		createsampledb()
	}
	if *testArangoDb {
		testarango()
	}

	// Open the database containing the discovered network
	//	dbOpenString := *DbName + "?mode=wal"
	//	database, openErr := sql.Open("sqlite3", dbOpenString)
	databaseForRead, openErr := sql.Open("sqlite3", *DbName)
	if openErr != nil {
		//		fmt.Println("Error opening database", dbOpenString)
		fmt.Println("Error opening databaseForRead", *DbName)
		log.Fatal(openErr)
	}
	/*	// Set database to WAL Mode
		dbResult, PRAGMAErr := database.Exec("PRAGMA journal_mode=WAL;")
		if PRAGMAErr != nil {
			log.Fatal(openErr)
		}
		if *debugFlag {
			fmt.Println("PRAMGA=", dbResult)
		}
		database.SetMaxOpenConns(0)
		if *debugFlag {
			fmt.Println("After SetMaxOpenConns function call")
		}
	*/
	defer databaseForRead.Close()

	databaseForUpdate, openErr := sql.Open("sqlite3", *DbName)
	if openErr != nil {
		fmt.Println("Error opening databaseForUpdate", *DbName)
		log.Fatal(openErr)
	}
	defer databaseForUpdate.Close()

	// Retrieve the Routers table
	//	routerRows, queryErr := database.Query("SELECT RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	//	routerRows, queryErr := database.Query("SELECT RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt, X3D, Y3D, Z3D FROM Routers")
	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		fmt.Println("databaseForRead Query error", queryErr)
		log.Fatal(openErr)
	}
	if *debugFlag {
		fmt.Println("Successful Routers table Select")
	}

	// Retrieve the Links table
	links, queryErr := databaseForRead.Query("SELECT LinkID, FromRouter, ToRouter FROM Links")
	if queryErr != nil {
		fmt.Println("databaseForRead Query error", queryErr)
		log.Fatal(openErr)
	}
	if *debugFlag {
		fmt.Println("Successful Links table Select")
	}

	// Initialize the 3D space
	app, appErr := application.Create(application.Options{
		Title:  "GoVisn - 3D Network Visualization",
		Width:  1200,
		Height: 1000,
	})
	if appErr != nil {
		fmt.Println("Error Creating 3D g3n app", *DbName)
		log.Fatal(openErr)
	}

	var RouterID int
	var SystemName string
	var SystemDesc string
	var UpTime string
	var Contact string
	var Location string
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	//	var X3D float32
	//	var Y3D float32
	//	var Z3D float32
	var router Router
	var link Link
	var LinkID int
	var FromRouter string
	var ToRouter string
	var x float32
	var y float32 = 1.0
	var z float32 = 1.0

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	//	pointLight.SetPosition(1, 0, 2)
	//	pointLight.SetPosition(10, 10, 20)
	pointLight.SetPosition((float32)(globeRadius+10), (float32)(globeRadius+10), (float32)(globeRadius+20))
	app.Scene().Add(pointLight)
	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	// Set initial camera position, i.e. viewing point
	//	app.CameraPersp().SetPosition(4, 0, 15)
	//	app.CameraPersp().SetPosition(0, 0, (float32)(globeRadius+40))
	app.CameraPersp().SetPosition(0, 0, (float32)(globeRadius*3.0))

	// Create a sphere representing the globe
	globe3D := geometry.NewSphere(globeRadius, 16, 16, 0, math.Pi*2, 0, math.Pi)
	globeMat := material.NewPhong(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetPosition(-1, -1, -1)
	app.Scene().Add(globeMesh)

	if *debugFlag {
		fmt.Println("Beginning routerRows.Next loop; adding routers to 3D scene.")
	}
	/*
	* Add the routers to the 3D scene
	 */
	for routerRows.Next() {
		routerRows.Scan(&RouterID, &SystemName, &SystemDesc, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt)
		//		routerRows.Scan(&RouterID, &SystemName, &SystemDesc, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt, &X3D, &Y3D, &Z3D)

		// Load router struct from DB fields
		router.System.RouterID = RouterID
		router.System.Name = SystemName
		router.System.UpTime = UpTime
		router.System.Contact = Contact
		router.System.Location = Location
		router.System.GPS.Latitude = GpsLat
		router.System.GPS.Longitude = GpsLong
		router.System.GPS.Altitude = GpsAlt

		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		cylinderMesh := graphic.NewMesh(rtr3D, mat)
		/*
		 * Set coordinates and altitude
		 */
		GpsLatFloat64, parseErr := strconv.ParseFloat(GpsLat, 64)
		if parseErr != nil {
			fmt.Println("Error parsing GpsLat =", GpsLat)
			log.Fatal(openErr)
		}
		xRadianLat := Rad(GpsLatFloat64)
		GpsLongFloat64, parseErr := strconv.ParseFloat(GpsLong, 64)
		if parseErr != nil {
			fmt.Println("Error parsing GpsLong", GpsLong)
			log.Fatal(openErr)
		}
		xRadianLong := Rad(GpsLongFloat64)
		x = (float32)(globeRadius * math.Sin(xRadianLat) * math.Cos(xRadianLong))
		router.System.Coordinates.X = x // update router struc with x coordinate
		yRadianLat := Rad(GpsLatFloat64)
		yRadianLong := Rad(GpsLongFloat64)
		y = (float32)(globeRadius * math.Sin(yRadianLat) * math.Sin(yRadianLong))
		router.System.Coordinates.Y = y // update router struc with y coordinate

		GpsAltFloat64, parseErr := strconv.ParseFloat(GpsAlt, 64)
		z = (float32)(globeRadius*(math.Cos(yRadianLat)) + GpsAltFloat64)
		router.System.Coordinates.Z = z // update router struc with z coordinate

		if *debugFlag {
			fmt.Println("x =", x, "y =", y, "z", z)
			fmt.Println("router =", router)
			fmt.Println("router.System.Coordinates =", router.System.Coordinates)
			fmt.Println("RouterID=", RouterID, "SystemName=", SystemName)
		}

		// Write 3D coordinates to Coordinates table row for later retrieval

		coordStatement, coordErr := databaseForUpdate.Prepare("UPDATE Coordinates SET X3D = ?, Y3D = ?, Z3D = ? WHERE RouterID = ?")
		if coordErr != nil {
			fmt.Println("Error preparing Coordinates Update statement:", coordErr)
			fmt.Println("coordStatement=", coordStatement)
			log.Fatal(coordErr)
		}
		//coordStatement.Exec(router.System.RouterID, router.System.Coordinates.X, router.System.Coordinates.Y, router.System.Coordinates.Z)
		_, coordErr = coordStatement.Exec(router.System.RouterID, router.System.Coordinates.X, router.System.Coordinates.Y, router.System.Coordinates.Z)
		if coordErr != nil {
			//			fmt.Println("Error executing Update to table Coordinates X3D, Y3D and Z3D coordinates.")
			//			fmt.Println("coordStatement=", coordStatement)
			//			log.Fatal(coordErr)
			//			log.Fatalln("Error executing Update to table Coordinates X3D, Y3D and Z3D coordinates.", "\ncoordStatement =", coordStatement, "\nError = ", coordErr)
			log.Fatalln("Error executing UPDATE to table Coordinates; X3D, Y3D and Z3D.",
				"\n System.RouterID =", router.System.RouterID,
				"\n System.Coordeinates.X =", router.System.Coordinates.X,
				"\n System.Coordeinates.Y =", router.System.Coordinates.Y,
				"\n System.Coordeinates.Z =", router.System.Coordinates.Z,
				"\n Error = ", coordErr)
		}
		defer coordStatement.Close()

		//		cylinderMesh.SetPosition(x, y, z)
		cylinderMesh.SetPosition(router.System.Coordinates.X, router.System.Coordinates.Y, router.System.Coordinates.Z)
		app.Scene().Add(cylinderMesh)

		//x = x + 2.0
		queryErr = routerRows.Err()
		if queryErr != nil {
			log.Fatal(queryErr)
		}
	}
	defer routerRows.Close()

	if *debugFlag {
		fmt.Println("Beginning links.Next loop; adding links to the 3D scene")
	}
	/*
	* Add the links to the 3D scene
	 */
	for links.Next() {
		err := links.Scan(&LinkID, &FromRouter, &ToRouter)
		if err != nil {
			log.Fatal(err)
		}
		// Load link struct from DB fields
		link.LinkID = LinkID
		link.FromRouter = FromRouter
		link.ToRouter = ToRouter
		if *debugFlag {
			fmt.Println("link =", link)
		}

		// retrieve FromRouter coordinates from router struc
		//		routerRows, queryErr = database.Query("SELECT RouterID, SystemName FROM Routers WHERE SystemName =", FromRouter)

		//		routerRows.Scan(&RouterID, &SystemName, &SystemDesc, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt, X3D, Y3D, Z3D)
		x = router.System.Coordinates.X
		y = router.System.Coordinates.Y
		z = router.System.Coordinates.Z
		if *debugFlag {
			fmt.Println("router coordinates =", router.System.Coordinates)
		}

		err = links.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Run()
}

const constX = math.Pi / 180

// Rad converts degrees to radians
func Rad(d float64) float64 { return d * constX }

// Deg converts radians to degrees
func Deg(r float64) float64 { return r / constX }
