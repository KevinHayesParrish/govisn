// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
)

func showRouterWindow(log *logger.Logger, gv *gvapp, router Router) {
	log.Debug("Showing router window for: %s", router.System.Name)

	// Create a panel as a window (use Panel instead of gui.NewWindow to avoid type conflicts)
	routerPanel := gui.NewPanel(400, 700)
	routerPanel.SetPosition(15, 50)
	routerPanel.SetBordersColor(math32.NewColor("darkblue"))
	routerPanel.SetBorders(2, 2, 2, 2)
	routerPanel.SetPaddings(10, 10, 10, 10)

	routerPanel.SetColor(math32.NewColor("grey")) // Set background color to grey

	// Create a title label at the top
	titleLabel := gui.NewLabel("Router Details: " + router.System.Name)
	titleLabel.SetFontSize(16)
	titleLabel.SetColor(math32.NewColor("white"))

	// Create a vertical layout container for the content
	vbox := gui.NewVBoxLayout()
	vbox.SetSpacing(3)
	routerPanel.SetLayout(vbox)

	// Add title as first element
	routerPanel.Add(titleLabel)

	// System Information Section
	systemTitle := gui.NewLabel("=== SYSTEM INFORMATION ===")
	systemTitle.SetFontSize(14)
	systemTitle.SetColor(math32.NewColor("lightblue"))
	routerPanel.Add(systemTitle)

	// Router ID
	routerIDLabel := gui.NewLabel("Router ID: " + strconv.Itoa(router.System.RouterID))
	routerPanel.Add(routerIDLabel)

	// Router Name/Hostname
	nameLabel := gui.NewLabel("Hostname: " + router.System.Name)
	routerPanel.Add(nameLabel)

	// Description
	if router.System.Description != "" {
		descLabel := gui.NewLabel("Description: " + router.System.Description)
		routerPanel.Add(descLabel)
	}

	// Location
	if router.System.Location != "" {
		locLabel := gui.NewLabel("Location: " + router.System.Location)
		routerPanel.Add(locLabel)
	}

	// Contact
	if router.System.Contact != "" {
		contactLabel := gui.NewLabel("Contact: " + router.System.Contact)
		routerPanel.Add(contactLabel)
	}

	// Uptime
	uptimeLabel := gui.NewLabel("Uptime: " + formatUptime(router.System.UpTime))
	routerPanel.Add(uptimeLabel)

	// Services
	servicesLabel := gui.NewLabel("Services: " + strconv.Itoa(router.System.Services))
	routerPanel.Add(servicesLabel)

	// GPS/Location Information Section
	gpsTitle := gui.NewLabel("=== GPS COORDINATES ===")
	gpsTitle.SetFontSize(14)
	gpsTitle.SetColor(math32.NewColor("lightgreen"))
	routerPanel.Add(gpsTitle)

	// Latitude
	latLabel := gui.NewLabel("Latitude: " + router.System.GPS.Latitude)
	routerPanel.Add(latLabel)

	// Longitude
	longLabel := gui.NewLabel("Longitude: " + router.System.GPS.Longitude)
	routerPanel.Add(longLabel)

	// Altitude
	altLabel := gui.NewLabel("Altitude: " + router.System.GPS.Altitude + " meters")
	routerPanel.Add(altLabel)

	// IP Addresses Section
	if len(router.Addresses.NetworkAddresses.IPAddress) > 0 {
		ipTitle := gui.NewLabel("=== IP ADDRESSES ===")
		ipTitle.SetFontSize(14)
		ipTitle.SetColor(math32.NewColor("lightyellow"))
		routerPanel.Add(ipTitle)

		for i, ip := range router.Addresses.NetworkAddresses.IPAddress {
			ipLabel := gui.NewLabel(fmt.Sprintf("[%d] %s", i+1, ip))
			routerPanel.Add(ipLabel)
		}
	}

	// MAC Addresses Section
	if len(router.Addresses.MediaAddresses.MediaAddress) > 0 {
		macTitle := gui.NewLabel("=== MAC ADDRESSES ===")
		macTitle.SetFontSize(14)
		macTitle.SetColor(math32.NewColor("lightcoral"))
		routerPanel.Add(macTitle)

		for i, mac := range router.Addresses.MediaAddresses.MediaAddress {
			macLabel := gui.NewLabel(fmt.Sprintf("[%d] %s", i+1, mac))
			routerPanel.Add(macLabel)
		}
	}

	// Neighbors Section
	if len(router.Neighbors.Neighbor) > 0 {
		neighborTitle := gui.NewLabel("=== NEIGHBORS ===")
		neighborTitle.SetFontSize(14)
		neighborTitle.SetColor(math32.NewColor("lightgrey"))
		routerPanel.Add(neighborTitle)

		for i, neighbor := range router.Neighbors.Neighbor {
			destLabel := gui.NewLabel(fmt.Sprintf("[%d] Dest: %s -> Via: %s",
				i+1, neighbor.DestinationAddress, neighbor.NextHop))
			routerPanel.Add(destLabel)
		}
	}

	// Add spacing before close button
	spacer := gui.NewLabel("")
	routerPanel.Add(spacer)

	// Close button
	closeBtn := gui.NewButton("Close")
	closeBtn.SetWidth(100)
	closeBtn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		routerPanel.SetVisible(false)
		gv.scene.Remove(routerPanel)
		log.Debug("Router panel closed for: %s", router.System.Name)
	})
	routerPanel.Add(closeBtn)

	// Make the panel draggable with mouse events
	var dragStartX, dragStartY float32
	var panelDragging bool

	routerPanel.Subscribe(gui.OnMouseDown, func(name string, ev interface{}) {
		mev := ev.(*window.MouseEvent)
		dragStartX = mev.Xpos
		dragStartY = mev.Ypos
		panelDragging = true
	})

	routerPanel.Subscribe(gui.OnMouseUp, func(name string, ev interface{}) {
		panelDragging = false
	})

	routerPanel.Subscribe(gui.OnCursor, func(name string, ev interface{}) {
		if panelDragging {
			cev := ev.(*window.CursorEvent)
			deltaX := cev.Xpos - dragStartX
			deltaY := cev.Ypos - dragStartY
			pos := routerPanel.Position()
			routerPanel.SetPosition(pos.X+deltaX, pos.Y+deltaY)
			dragStartX = cev.Xpos
			dragStartY = cev.Ypos
		}
	})

	// Add the panel to the scene
	gv.scene.Add(routerPanel)
	log.Debug("Router panel added to scene for: %s", router.System.Name)
}

// formatUptime converts uptime in centiseconds to a human-readable format
func formatUptime(uptimeCs uint32) string {
	// Convert centiseconds to total seconds
	totalSeconds := uptimeCs / 100

	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds",
		days, hours, minutes, seconds)
}
