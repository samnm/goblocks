package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"os"
	"sam.vg/goblocks/blocks"
	"sam.vg/util/fps"
)

const (
	title     = "Goblocks!"
	appWidth  = 800
	appHeight = 640
)

func main() {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.Terminate()

	if err := glfw.OpenWindow(appWidth, appHeight, 8, 8, 8, 8, 24, 8, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	if err := gl.Init(); err != 0 {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	glfw.SetWindowTitle(title)

	fps := fps.NewFPS(glfw.Time())

	blocks.Init()
	for glfw.WindowParam(glfw.Opened) == 1 {
		blocks.Tick()
		fps.Tick(glfw.Time())
	}
}
