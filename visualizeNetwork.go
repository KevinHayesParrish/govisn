package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/application"
	_ "github.com/mattn/go-sqlite3"
)

func visualizeNetwork(debugFlag bool, databaseForRead *sql.DB) *sql.DB {
	const VISUALIZENETWORKVERSION = "0.0.1"
	if debugFlag {
		fmt.Println("visualizeNetwork", VISUALIZENETWORKVERSION, "func started")
	}

	// Retrieve the Routers table
	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		fmt.Println("databaseForRead Query error", queryErr)
		log.Fatal(queryErr)
	}
	if debugFlag {
		fmt.Println("Successful Routers table Select")
	}

	// Retrieve the Links table
	linkRows, queryErr := databaseForRead.Query("SELECT LinkID, FromRouterName, FromRouterIP, ToRouterName, ToRouterIP FROM Links")
	if queryErr != nil {
		fmt.Println("databaseForRead Query error", queryErr)
		log.Fatal(queryErr)
	}
	if debugFlag {
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
		log.Fatal(appErr)
	}

	var RouterID int
	var Name string
	var Description string
	var UpTime uint32
	var Contact string
	var Location string
	var Services int
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	var link Link
	var LinkID int
	var FromRouterName, FromRouterIP, ToRouterName, ToRouterIP string
	var x float32
	var y float32 = 1.0
	var z float32 = 1.0

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition((float32)(globeRadius+10), (float32)(globeRadius+10), (float32)(globeRadius+20))
	app.Scene().Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	// Set initial camera position, i.e. viewing point
	app.CameraPersp().SetPosition(0.0, 0.0, (float32)(globeRadius*2.0))

	// Create a sphere representing the globe
	globe3D := geometry.NewSphere(globeRadius, 16, 16, 0, math.Pi*2, 0, math.Pi)
	globeMat := material.NewPhong(&math32.Color{R: 0.0, G: 0.5, B: 1.0}) // Azure blue 0, 128, 255
	globeMat.SetTransparent(true)
	globeMat.SetOpacity(0.25)
	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetPosition(0, 0, 0)
	app.Scene().Add(globeMesh)

	if debugFlag {
		fmt.Println("Beginning routerRows.Next loop; adding routers to 3D scene.")
	}
	/*
	* Add the routers to the 3D scene
	 */
	var routers []Router
	var router Router
	routerArrayIndex := 0
	for routerRows.Next() {
		routerRows.Scan(&RouterID, &Name, &Description, &UpTime, &Contact, &Location, &Services, &GpsLat, &GpsLong, &GpsAlt)

		// Load router struct from DB fields
		router.System.RouterID = RouterID
		router.System.UpTime = UpTime
		router.System.Name = Name
		router.System.Contact = Contact
		router.System.Location = Location
		router.System.Services = Services
		router.System.GPS.Latitude = GpsLat
		router.System.GPS.Longitude = GpsLong
		router.System.GPS.Altitude = GpsAlt

		routers = append(routers, router)

		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		cylinderMesh := graphic.NewMesh(rtr3D, mat)
		/*
		 * Set coordinates and altitude
		 */
		x, y, z = calcCoordinates(GpsLat, GpsLong, GpsAlt)

		if debugFlag {
			fmt.Println("x =", x, "y =", y, "z", z)
			fmt.Println("router =", routers[routerArrayIndex])
			fmt.Println("router.System.GPS =", routers[routerArrayIndex].System.GPS)
			fmt.Println("RouterID=", RouterID, "Name=", Name)
		}

		// Add Router object to 3D scene.
		cylinderMesh.SetPosition(x, y, z)
		app.Scene().Add(cylinderMesh)

		// Add router name to scene
		// Creates Font
		fontfile := os.Getenv("GOPATH") + "/src/govisn/data/fonts/FreeSans.ttf"
		font, err := text.NewFont(fontfile)
		if err != nil {
			app.Log().Fatal(err.Error())
		}

		font.SetLineSpacing(1.0)
		font.SetPointSize(28)
		font.SetDPI(72)
		font.SetFgColor(&math32.Color4{R: 0, G: 0, B: 1, A: 1})
		font.SetBgColor(&math32.Color4{R: 1, G: 1, B: 0, A: 0.8})
		canvas := text.NewCanvas(300, 200, &math32.Color4{R: 0, G: 1, B: 0, A: 0.8})
		rtext := "RouterID: " + strconv.Itoa(routers[routerArrayIndex].System.RouterID) + "\nHostname: " + routers[routerArrayIndex].System.Name
		swidth, sheight := font.MeasureText(rtext)
		canvas = text.NewCanvas(swidth, sheight, &math32.Color4{R: 0, G: 1, B: 1, A: 1})
		canvas.DrawText(0, 0, rtext, font)
		tex3 := texture.NewTexture2DFromRGBA(canvas.RGBA)
		mat3 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat3.AddTexture(tex3)
		aspect := float32(swidth) / float32(sheight)
		mesh3 := graphic.NewSprite(aspect, 1, mat3)
		mesh3.SetPosition(x, y, z+1.0)
		app.Scene().Add(mesh3)

		queryErr = routerRows.Err()
		if queryErr != nil {
			log.Fatal(queryErr)
		}

		routerArrayIndex++
	}
	defer routerRows.Close()

	if debugFlag {
		fmt.Println("\nBeginning linkRows.Next loop; adding links to the 3D scene")
		fmt.Println()
	}
	/*
	* Add the links to the 3D scene
	 */
	var FromRouterX, FromRouterY, FromRouterZ, ToRouterX, ToRouterY, ToRouterZ string
	for linkRows.Next() {
		err := linkRows.Scan(&LinkID, &FromRouterName, &FromRouterIP, &ToRouterName, &ToRouterIP)
		if err != nil {
			log.Fatal(err)
		}

		// Exclude false routes
		if FromRouterIP == "127.0.0.0" || FromRouterIP == "127.0.0.1" || FromRouterIP == "224.0.0.0" || FromRouterIP == "0.0.0.0" || ToRouterIP == "127.0.0.0" || ToRouterIP == "127.0.0.1" || ToRouterIP == "224.0.0.0" || ToRouterIP == "0.0.0.0" {
			continue
		}

		// Load link struct from DB fields
		link.LinkID = LinkID
		link.FromRouterName = FromRouterName
		link.FromRouterIP = FromRouterIP
		link.ToRouterName = ToRouterName
		link.ToRouterIP = ToRouterIP

		// retrieve FromRouter coordinates
		if debugFlag {
			fmt.Println("link =", link)
			fmt.Println("FromRouterName=", link.FromRouterName)
			fmt.Println("FromRouterIP=", link.FromRouterIP)
		}

		//  Query database for FromRouter GPS coordinates

		routerGpsRows, err := databaseForRead.Query("SELECT Name, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name = $1", FromRouterName)
		if err != nil {
			log.Fatalln("databaseForRead Query error", err.Error())
		}
		if debugFlag {
			fmt.Println("Successful Query for FromRouter GPS Coordinates")
		}
		defer routerGpsRows.Close()
		var linksFromRouterName string
		var linksFromRouterGpsLat, linksFromRouterGpsLong, linksFromRouterGpsAlt string
		for routerGpsRows.Next() {
			routerGpsRows.Scan(&linksFromRouterName, &linksFromRouterGpsLat, &linksFromRouterGpsLong, &linksFromRouterGpsAlt)
		}

		FromRouterX = linksFromRouterGpsLat
		FromRouterY = linksFromRouterGpsLong
		FromRouterZ = linksFromRouterGpsAlt
		if debugFlag {
			fmt.Println("returned from getRouterCoordinatesName func: FromRouterX=", FromRouterX, "FromRouterY=", FromRouterY, "FromRouterZ=", FromRouterZ)
		}

		//  Query database for FromRouter GPS coordinates
		if debugFlag {
			fmt.Println("ToRouterName=", link.ToRouterName)
			fmt.Println("ToRouterIP=", link.ToRouterIP)
		}
		routerGpsRows, err = databaseForRead.Query("SELECT Name, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name = $1", ToRouterName)
		if err != nil {
			log.Fatalln("databaseForRead Query error", err.Error())
		}
		if debugFlag {
			fmt.Println("Successful Query for ToRouter GPS Coordinates")
		}
		var linksToRouterName string
		var linksToRouterGpsLat, linksToRouterGpsLong, linksToRouterGpsAlt string
		for routerGpsRows.Next() {
			routerGpsRows.Scan(&linksToRouterName, &linksToRouterGpsLat, &linksToRouterGpsLong, &linksToRouterGpsAlt)
		}
		ToRouterX = linksToRouterGpsLat
		ToRouterY = linksToRouterGpsLong
		ToRouterZ = linksToRouterGpsAlt

		if debugFlag {
			fmt.Println("router", Name, "GPS coordinates =", GpsLat, GpsLong, GpsAlt)
			fmt.Println("returned from getRouterCoordinatesIP func: ToRouterX=", ToRouterX, "ToRouterY=", ToRouterY, "ToRouterZ=", ToRouterZ)
		}

		// Add link object to the 3D scene
		fromX, fromY, fromZ := calcCoordinates(FromRouterX, FromRouterY, FromRouterZ)
		toX, toY, toZ := calcCoordinates(ToRouterX, ToRouterY, ToRouterZ)

		linkGeom := geometry.NewGeometry()
		vertices := math32.NewArrayF32(0, 0)
		vertices.Append(
			fromX, fromY, fromZ,
			toX, toY, toZ,
		)
		if debugFlag {
			fmt.Println("link vertices=", vertices)
			fmt.Println()
		}
		colors := math32.NewArrayF32(0, 0)
		colors.Append(
			0.0, 0.0, 1.0,
			0.0, 0.0, 1.0,
		)
		linkGeom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

		linkGeom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

		// Creates basic material
		mat := material.NewBasic()

		// Creates lines with the specified geometry and material
		link3D := graphic.NewLines(linkGeom, mat)

		app.Scene().Add(link3D)

		err = linkRows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Run()

	if debugFlag {
		fmt.Println("visualizeNetwork", VISUALIZENETWORKVERSION, "func ending")
	}
	return databaseForRead
}

func calcCoordinates(GpsLat string, GpsLong string, GpsAlt string) (float32, float32, float32) {
	var x, y, z float32

	var GpsLatFloat64 = 0.0
	var parseErr error
	if GpsLat != "" {
		GpsLatFloat64, parseErr = strconv.ParseFloat(GpsLat, 64)
	}
	if parseErr != nil {
		fmt.Println("Error parsing GpsLat =", GpsLat)
		log.Fatal(parseErr)
	}
	xRadianLat := Rad(GpsLatFloat64)

	var GpsLongFloat64 = 0.0
	if GpsLong != "" {
		GpsLongFloat64, parseErr = strconv.ParseFloat(GpsLong, 64)
		if parseErr != nil {
			fmt.Println("Error parsing GpsLong", GpsLong)
			log.Fatal(parseErr)
		}
	}

	xRadianLong := Rad(GpsLongFloat64)
	x = (float32)(globeRadius * math.Sin(xRadianLat) * math.Cos(xRadianLong))

	yRadianLat := Rad(GpsLatFloat64)
	yRadianLong := Rad(GpsLongFloat64)
	y = (float32)(globeRadius * math.Sin(yRadianLat) * math.Sin(yRadianLong))

	GpsAltFloat64, parseErr := strconv.ParseFloat(GpsAlt, 64)
	GpsAltFloat64 = GpsAltFloat64 / 100000.0
	z = (float32)(globeRadius*(math.Cos(yRadianLat)) + GpsAltFloat64)

	return x, y, z
}
