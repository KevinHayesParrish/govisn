package main

import (
	"fmt"
	"os"

	"github.com/g3n/engine/camera"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"

	"github.com/g3n/engine/core"

	_ "github.com/mattn/go-sqlite3"
)

func dumpSceneToFile(scene core.INode, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "Scene Graph Dump\n")
	fmt.Fprintf(f, "================\n\n")

	// Recursively dump the scene
	dumpNodeToFile(f, scene, 0)

	return nil
}

// dumpNodeToFile recursively writes node information to a file
func dumpNodeToFile(file *os.File, node core.INode, depth int) {
	// Create indentation
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	// Write node basic info
	fmt.Fprintf(file, "%sNode: %s (Type: %T)\n", indent, node.Name(), node)

	// Write node-specific details based on type
	switch n := node.(type) {
	case *core.Node:
		pos := n.Position()
		rot := n.Rotation()
		scl := n.Scale()
		fmt.Fprintf(file, "%s  Position: (%.3f, %.3f, %.3f)\n", indent, pos.X, pos.Y, pos.Z)
		fmt.Fprintf(file, "%s  Rotation: (%.3f, %.3f, %.3f)\n", indent, rot.X, rot.Y, rot.Z)
		fmt.Fprintf(file, "%s  Scale: (%.3f, %.3f, %.3f)\n", indent, scl.X, scl.Y, scl.Z)

	case *graphic.Mesh:
		fmt.Fprintf(file, "%s  Geometry: %v\n", indent, n.GetGeometry())
		fmt.Fprintf(file, "%s  Material Count: %d\n", indent, n.MaterialCount())

	case *camera.Camera:
		fmt.Fprintf(file, "%s  Aspect Ratio: %.3f\n", indent, n.Aspect())
		fmt.Fprintf(file, "%s  Near Plane: %.3f\n", indent, n.Near())
		fmt.Fprintf(file, "%s  Far Plane: %.3f\n", indent, n.Far())

	case *light.Point:
		pos := n.Position()
		fmt.Fprintf(file, "%s  Point Light Position: (%.3f, %.3f, %.3f)\n", indent, pos.X, pos.Y, pos.Z)

	case *light.Directional:
		fmt.Fprintf(file, "%s  Directional Light\n", indent)

	case *light.Ambient:
		fmt.Fprintf(file, "%s  Ambient Light\n", indent)
	}

	// Recursively dump child nodes
	children := node.Children()
	if len(children) > 0 {
		fmt.Fprintf(file, "%s  Children: %d\n", indent, len(children))
		for _, child := range children {
			dumpNodeToFile(file, child, depth+1)
		}
	}

	fmt.Fprintf(file, "\n")
}
