package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	logger "github.com/alouca/gologger"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/util/stats"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/window"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"

	//	"github.com/g3n/engine/util/application"

	"github.com/g3n/engine/experimental/collision"

	_ "github.com/mattn/go-sqlite3"
)

// App contains the application state
type App struct {
	*app.Application                // Embedded standard application object
	log              *logger.Logger // Application logger
	//currentDemo      IDemo            // Current demo
	dirData    string           // Full path of the data directory
	scene      *core.Node       // Scene rendered
	demoScene  *core.Node       // Scene populated by individual demos
	ambLight   *light.Ambient   // Scene ambient light
	frameRater *util.FrameRater // Render loop frame rater

	// GUI
	mainPanel  *gui.Panel
	demoPanel  *gui.Panel
	labelFPS   *gui.Label         // header FPS label
	treeTests  *gui.Tree          // tree with test names
	stats      *stats.Stats       // statistics object
	statsTable *stats.StatsTable  // statistics table panel
	control    *gui.ControlFolder // Pointer to gui control panel

	// Camera and orbit control
	camera *camera.Camera       // Camera
	orbit  *camera.OrbitControl // Orbit control
}

// GuiMenu is the structure containing the Menus for the Gui
type GuiMenu struct {
}

// Raycast is the structure containing the Raycaster
type Raycast struct {
	//	rayCast *Raycaster
	rayCast *collision.Raycaster
}

func visualizeNetwork(debugFlag bool, databaseForRead *sql.DB) *sql.DB {
	const VISUALIZENETWORKVERSION = "0.2.1"
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

	// Create application and scene
	gv := new(gvapp)
	a := app.App()
	gv.Application = a
	//	scene := core.NewNode()
	gv.scene = core.NewNode()

	// Create perspective camera
	//	cam := camera.New(1)
	//	cam.SetPosition(0, 0, (float32)(globeRadius*2.0))
	gv.cam = camera.New(1)
	gv.cam.SetPosition(0, 0, (float32)(globeRadius*2.0))

	// Setup orbit control for the camera
	//	camera.NewOrbitControl(cam)
	camera.NewOrbitControl(gv.cam)

	//	scene.Add(cam)
	//	gv.scene.Add(cam)
	gv.scene.Add(gv.cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		//		cam.SetAspect(float32(width) / float32(height))
		gv.cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create and Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.0)
	//	scene.Add(ambientLight)
	gv.scene.Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition((float32)(globeRadius+10), (float32)(globeRadius+10), (float32)(globeRadius+20))
	//	scene.Add(pointLight)
	gv.scene.Add(pointLight)

	// Add an axis helper to the scene
	//	axis := graphic.NewAxisHelper(0.5)
	//	app.Scene().Add(axis)
	axes := helper.NewAxes(1)
	//	scene.Add(axes)
	gv.scene.Add(axes)

	// Set background color to black
	a.Gls().ClearColor(0.0, 0.0, 0.0, 0.0)

	// Set the scene to be managed by the gui manager
	//gui.Manager().Set(scene)

	// Build Menus
	//	buildMenus(debugFlag, app)
	buildMenus(debugFlag, gv, a)

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

	// Setup Mouse clicking of objects within the 3D scene
	var t Raycast
	//	t.Initialize(debugFlag, scene, cam, a)
	t.Initialize(debugFlag, gv.scene, gv.cam, a)

	// Create Globe texture
	gobinDir := os.Getenv("GOBIN")
	texfile := gobinDir + "/data/images/earth_clouds_big.jpg"
	globeTex, err := texture.NewTexture2DFromImage(texfile)
	if err != nil {
		//		app.Log().Fatal("Error loading texture: %s", err, "\n Insure govisn /data/images is copied to GOBIN")
		log.Fatalln("Error loading texture:", err, "\n Insure govisn /data/images is copied to GOBIN")
	}
	globeTex.SetFlipY(false)

	// Create a sphere representing the globe
	//	globe3D := geometry.NewSphere(globeRadius, 16, 16, 0, math.Pi*2, 0, math.Pi)
	globe3D := geometry.NewSphere(globeRadius, 16, 16)
	//	globeMat := material.NewPhong(&math32.Color{R: 0.0, G: 0.5, B: 1.0}) // Azure blue 0, 128, 255
	globeMat := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0}) // White 255, 255, 255
	globeMat.AddTexture(globeTex)
	globeMat.SetTransparent(true)
	globeMat.SetOpacity(.50)

	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetPosition(0, 0, 0)
	//	app.Scene().Add(globeMesh)
	//	scene.Add(globeMesh)
	gv.scene.Add(globeMesh)

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

		//		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 16, 2, true, true)
		//		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
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
		cylinderMesh.SetName(string(router.System.RouterID))
		cylinderMesh.SetUserData(string(router.System.RouterID))
		//		app.Scene().Add(cylinderMesh)
		//		scene.Add(cylinderMesh)
		gv.scene.Add(cylinderMesh)

		// Add router name to scene
		// Creates Font
		//		fontfile := os.Getenv("GOPATH") + "/src/govisn/data/fonts/FreeSans.ttf"
		fontfile := os.Getenv("GOBIN") + "/data/fonts/FreeSans.ttf"
		font, err := text.NewFont(fontfile)
		if err != nil {
			//			app.Log().Fatal(err.Error())
			//			app.Log().Fatal("Error loading font: %s", err, "\n Insure govisn /data/fonts is copied to GOBIN")
			log.Fatalln("Error loading font:", err, "\n Insure govisn /data/fonts is copied to GOBIN")
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
		//		app.Scene().Add(mesh3)
		//		scene.Add(mesh3)
		gv.scene.Add(mesh3)

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
		mat.SetLineWidth(3.0) // Set line width. Default is 1.0

		// Creates lines with the specified geometry and material
		link3D := graphic.NewLines(linkGeom, mat)
		link3D.SetName(string(link.LinkID))

		//		app.Scene().Add(link3D)
		//		scene.Add(link3D)
		gv.scene.Add(link3D)

		err = linkRows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run the application
	//	app.Run()
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		//		renderer.Render(scene, cam)
		renderer.Render(gv.scene, gv.cam)
	})

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

//func buildMenus(debugFlag bool, app *application.Application) *application.Application {
func buildMenus(debugFlag bool, gv *gvapp, a *app.Application) *app.Application {
	//func buildMenus(debugFlag bool, a *App) *App {
	if debugFlag {
		fmt.Println("Starting func buildMenus")
	}

	// Event handler for menu clicks
	onClick := func(evname string, ev interface{}) {
		switch ev.(*gui.MenuItem).Id() {
		case "Exit":
			{
				fmt.Println("GoVisn terminating. File/Exit selected.")
				//				os.Exit(0)
				gv.Exit()
			}
		}
	}

	// Show and enable demo panel

	//dl := gui.NewDockLayout()
	//width, height := a.GetSize()
	//a.mainPanel = gui.NewPanel(float32(width), float32(height))
	//a.mainPanel.SetRenderable(true)
	//a.mainPanel.SetEnabled(true)
	//a.mainPanel.SetLayout(dl)
	//a.scene.Add(a.mainPanel)
	//gui.Manager().Set(a.mainPanel)

	gui.Manager().Set(gv.scene)

	// Create menu bar
	mb := gui.NewMenuBar()
	mb.Subscribe(gui.OnClick, onClick)
	mb.SetPosition(10, 10)

	// Create fileMenu and adds it to the menu bar
	m1 := gui.NewMenu()
	m1.AddOption("File/Exit").
		SetId("Exit")
	mb.AddMenu("File", m1).
		SetId("File").
		SetShortcut(window.ModAlt, window.Key1)

	//a.mainPanel.Add(mb)
	mb.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		material.NewStandard(math32.NewColor("DarkRed"))
	})
	gv.scene.Add(mb)
	//app.Gui().Add(mb)
	//app.Gui().Root().SetKeyFocus(mb)

	gui.Manager().SetKeyFocus(mb)

	if debugFlag {
		fmt.Println("func buildMenus ended")
	}
	return (a)
}

// Initialize the raycaster
//func (t *Raycast) Initialize(debugFlag bool, app *application.Application) {
func (t *Raycast) Initialize(debugFlag bool, scene *core.Node, cam *camera.Camera, app *app.Application) {
	fmt.Println("Initializing the raycaster") // TESTING ONLY
	// Creates the raycaster
	//	var t *Raycast
	//	t.rayCast = core.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	t.rayCast = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	t.rayCast.LinePrecision = 0.05
	t.rayCast.PointPrecision = 0.05

	// Subscribe to mouse button down events
	//	app.Window().Subscribe(window.OnMouseDown, func(evname string, ev interface{}) {
	app.SubscribeID(window.OnMouseDown, app, func(evname string, ev interface{}) {
		//		t.onMouse(debugFlag, app, ev)
		t.onMouse(debugFlag, scene, cam, app, ev)
	})
}

// onMouse is executed when an object in the 3D scene is selected with a mouse click
//func (t *Raycast) onMouse(debugFlag bool, app *application.Application, ev interface{}) {
func (t *Raycast) onMouse(debugFlag bool, scene *core.Node, cam *camera.Camera, app *app.Application, ev interface{}) {
	// Convert mouse coordinates to normalized device coordinates
	mev := ev.(*window.MouseEvent)
	//	width, height := app.Window().Size()
	width, height := app.GetSize()
	x := 2*(mev.Xpos/float32(width)) - 1
	y := -2*(mev.Ypos/float32(height)) + 1
	if debugFlag {
		fmt.Println("onMouse x=", x)
		fmt.Println("onMouse y=", y)
	}

	// Set the raycaster from the current camera and mouse coordinates
	//	app.Camera().SetRaycaster(t.rayCast, x, y)
	t.rayCast.SetFromCamera(cam, x, y)
	if debugFlag {
		fmt.Printf("rayCast:%+v\n", t.rayCast.Ray)
	}

	// Checks intersection with all objects in the scene
	//	intersects := t.rayCast.IntersectObjects(app.Scene().Children(), true)
	intersects := t.rayCast.IntersectObjects(scene.Children(), true)
	if debugFlag {
		fmt.Printf("intersects:%+v\n", intersects)
	}
	if len(intersects) == 0 {
		return
	}

	// Get first intersection
	obj := intersects[0].Object
	router3D := obj.GetNode()
	router3DName := router3D.Name()
	if debugFlag {
		fmt.Println("Picked object Name=", router3DName)
		fmt.Println("Picked object UserData=", router3D.UserData())
	}

	// Convert INode to IGraphic
	ig, ok := obj.(graphic.IGraphic)
	if !ok {
		//		app.Log().Debug("Not graphic:%T", obj)
		log.Fatalf("Not graphic:%T", obj)
		return
	}
	// Get graphic object
	gr := ig.GetGraphic()
	imat := gr.GetMaterial(0)

	type matI interface {
		EmissiveColor() math32.Color
		SetEmissiveColor(*math32.Color)
	}

	if v, ok := imat.(matI); ok {
		if em := v.EmissiveColor(); em.R == 1 && em.G == 1 && em.B == 1 {
			v.SetEmissiveColor(&math32.Color{R: 0, B: 0, G: 0})
		} else {
			v.SetEmissiveColor(&math32.Color{R: 1, B: 1, G: 1})
		}
	}

	// Gets the router system information for a pop-up window
	//router := getRouterFromDB(debugFlag, obj)
}

// Render renders the mouse pick action
//func (t *Raycast) Render(app *application.Application) {
func (t *Raycast) Render(a *app.Application) {
}

// Update is called every frame.
func (t *Raycast) Update(a *app.Application, deltaTime time.Duration) {}

// Cleanup is called once at the end of the demo.
func (t *Raycast) Cleanup(a *app.Application) {}
