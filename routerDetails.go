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
)

func showRouterWindow(log *logger.Logger, gv *gvapp, router Router) {
	log.Debug("Showing router window for: %s", router.System.Name)

	// Create a modal panel window
	routerWindow := gui.NewWindow(400, 600)
	routerWindow.SetTitle("Router Details: " + router.System.Name)
	routerWindow.SetPosition(100, 50)

	// Create a vertical layout container for the content
	vbox := gui.NewVBoxLayout()
	vbox.SetSpacing(5)
	contentPanel := gui.NewPanel(580, 650)
	contentPanel.SetLayout(vbox)
	contentPanel.SetBordersColor(math32.NewColor("grey"))
	contentPanel.SetBorders(1, 1, 1, 1)
	contentPanel.SetPaddings(10, 10, 10, 10)

	// System Information Section
	systemTitle := gui.NewLabel("=== SYSTEM INFORMATION ===")
	systemTitle.SetFontSize(14)
	systemTitle.SetColor(math32.NewColor("lightblue"))
	contentPanel.Add(systemTitle)

	// Router ID
	routerIDLabel := gui.NewLabel("Router ID: " + strconv.Itoa(router.System.RouterID))
	contentPanel.Add(routerIDLabel)

	// Router Name/Hostname
	nameLabel := gui.NewLabel("Hostname: " + router.System.Name)
	contentPanel.Add(nameLabel)

	// Description
	if router.System.Description != "" {
		descLabel := gui.NewLabel("Description: " + router.System.Description)
		contentPanel.Add(descLabel)
	}

	// Location
	if router.System.Location != "" {
		locLabel := gui.NewLabel("Location: " + router.System.Location)
		contentPanel.Add(locLabel)
	}

	// Contact
	if router.System.Contact != "" {
		contactLabel := gui.NewLabel("Contact: " + router.System.Contact)
		contentPanel.Add(contactLabel)
	}

	// Uptime
	uptimeLabel := gui.NewLabel("Uptime: " + formatUptime(router.System.UpTime))
	contentPanel.Add(uptimeLabel)

	// Services
	servicesLabel := gui.NewLabel("Services: " + strconv.Itoa(router.System.Services))
	contentPanel.Add(servicesLabel)

	// GPS/Location Information Section
	gpsTitle := gui.NewLabel("=== GPS COORDINATES ===")
	gpsTitle.SetFontSize(14)
	gpsTitle.SetColor(math32.NewColor("lightgreen"))
	contentPanel.Add(gpsTitle)

	// Latitude
	latLabel := gui.NewLabel("Latitude: " + router.System.GPS.Latitude)
	contentPanel.Add(latLabel)

	// Longitude
	longLabel := gui.NewLabel("Longitude: " + router.System.GPS.Longitude)
	contentPanel.Add(longLabel)

	// Altitude
	altLabel := gui.NewLabel("Altitude: " + router.System.GPS.Altitude + " meters")
	contentPanel.Add(altLabel)

	// IP Addresses Section
	if len(router.Addresses.NetworkAddresses.IPAddress) > 0 {
		ipTitle := gui.NewLabel("=== IP ADDRESSES ===")
		ipTitle.SetFontSize(14)
		ipTitle.SetColor(math32.NewColor("lightyellow"))
		contentPanel.Add(ipTitle)

		for i, ip := range router.Addresses.NetworkAddresses.IPAddress {
			ipLabel := gui.NewLabel(fmt.Sprintf("[%d] %s", i+1, ip))
			contentPanel.Add(ipLabel)
		}
	}

	// MAC Addresses Section
	if len(router.Addresses.MediaAddresses.MediaAddress) > 0 {
		macTitle := gui.NewLabel("=== MAC ADDRESSES ===")
		macTitle.SetFontSize(14)
		macTitle.SetColor(math32.NewColor("lightcoral"))
		contentPanel.Add(macTitle)

		for i, mac := range router.Addresses.MediaAddresses.MediaAddress {
			macLabel := gui.NewLabel(fmt.Sprintf("[%d] %s", i+1, mac))
			contentPanel.Add(macLabel)
		}
	}

	// Neighbors Section
	if len(router.Neighbors.Neighbor) > 0 {
		neighborTitle := gui.NewLabel("=== NEIGHBORS ===")
		neighborTitle.SetFontSize(14)
		neighborTitle.SetColor(math32.NewColor("lightgrey"))
		contentPanel.Add(neighborTitle)

		for i, neighbor := range router.Neighbors.Neighbor {
			destLabel := gui.NewLabel(fmt.Sprintf("[%d] Dest: %s -> Via: %s",
				i+1, neighbor.DestinationAddress, neighbor.NextHop))
			contentPanel.Add(destLabel)
		}
	}

	// Add spacing before close button
	spacer := gui.NewLabel("")
	contentPanel.Add(spacer)

	routerWindow.Add(contentPanel)

	// Create a footer panel for the close button
	footerPanel := gui.NewPanel(600, 40)
	footerLayout := gui.NewHBoxLayout()
	footerLayout.SetSpacing(5)
	footerPanel.SetLayout(footerLayout)

	// Close button
	closeBtn := gui.NewButton("Close")
	closeBtn.SetWidth(100)
	closeBtn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		routerWindow.SetVisible(false)
		gv.scene.Remove(routerWindow)
		log.Debug("Router window closed for: %s", router.System.Name)
	})
	footerPanel.Add(closeBtn)

	routerWindow.Add(footerPanel)

	// Add the window to the scene
	gv.scene.Add(routerWindow)
	log.Debug("Router window added to scene for: %s", router.System.Name)
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
