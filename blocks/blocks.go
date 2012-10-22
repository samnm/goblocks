package blocks

import (
	"sam.vg/goblocks/blocks/input"
	"sam.vg/goblocks/blocks/renderer"
)

func Init(width, height int) {
	renderer.Init(width, height)
	input.Init()
}

func Tick() {
	renderer.Tick()
}
