package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"os"
	"sam.vg/goblocks/blocks"
	"sam.vg/goblocks/util/fps"
)

const (
	title     = "Goblocks!"
	appWidth  = 1024
	appHeight = 768
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

	blocks.Init(appWidth, appHeight)
	for glfw.WindowParam(glfw.Opened) == 1 {
		blocks.Tick()
		fps.Tick(glfw.Time())
		if glfw.WindowParam(glfw.Active) == 1 {
			glfw.Sleep(0.001)
		} else {
			glfw.Sleep(0.05)
		}
	}
}
