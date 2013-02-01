package blocks

import (
	"github.com/samnm/goblocks/blocks/input"
	"github.com/samnm/goblocks/blocks/renderer"
)

func Init(width, height int) {
	renderer.Init(width, height)
	input.Init()
}

func Tick() {
	renderer.Tick()
}
