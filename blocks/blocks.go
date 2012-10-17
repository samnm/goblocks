package blocks

import (
	"sam.vg/goblocks/blocks/renderer"
)

func Init(width, height int) {
	renderer.Init(width, height)
}

func Tick() {
	renderer.Tick()
}
