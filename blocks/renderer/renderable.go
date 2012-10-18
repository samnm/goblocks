package renderer

import (
	"github.com/go-gl/gl"
	"unsafe"
)

type Vertex struct {
	position [4]float32
	texcoord [2]float32
}

type Renderable struct {
	vertexBuffer  gl.Buffer
	elementBuffer gl.Buffer
	numIndicies   int
}

func NewRenderable() *Renderable {
	vertexPositions := []float32{
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
	}
	indicies := []gl.GLushort{
		4, 5, 6, // front (-Z)
		6, 5, 7,
		1, 0, 2, // back (+Z)
		2, 3, 1,
		0, 4, 2, // right (+X)
		2, 4, 6,
		7, 5, 1, // left (-X)
		1, 3, 7,
		0, 1, 5, // top (+Y)
		5, 4, 0,
		3, 2, 6, // bottom (-Y)
		6, 7, 3,
	}

	ro := new(Renderable)

	var size int
	size = len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	ro.vertexBuffer = MakeBuffer(gl.ARRAY_BUFFER, size, vertexPositions)

	size = len(indicies) * int(unsafe.Sizeof(indicies[0]))
	ro.elementBuffer = MakeBuffer(gl.ELEMENT_ARRAY_BUFFER, size, indicies)

	ro.numIndicies = len(indicies)

	return ro
}
