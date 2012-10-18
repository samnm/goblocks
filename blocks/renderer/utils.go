package renderer

import (
	"fmt"
	"github.com/go-gl/gl"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"unsafe"
)

func MakeBuffer(target gl.GLenum, size int, data interface{}) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(target)
	gl.BufferData(target, size, data, gl.STATIC_DRAW)
	gl.BufferUnbind(target)
	return buffer
}

func MakeShader(shaderType gl.GLenum, filename string) gl.Shader {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
	}
	defer file.Close()

	data := make([]byte, 2048)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	source := string(data[:count])

	shader := gl.CreateShader(shaderType)
	shader.Source(source)
	shader.Compile()
	// fmt.Println(shader.GetInfoLog())

	return shader
}

func LoadTexture(filename string) gl.Texture {
	r, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return 0
	}
	defer r.Close()

	img, _, err := image.Decode(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return 0
	}

	// lazy way to ensure we have the correct image format for opengl
	rgbaImg := image.NewNRGBA(img.Bounds())
	draw.Draw(rgbaImg, img.Bounds(), img, image.ZP, draw.Src)

	tex := gl.GenTexture()
	tex.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// flip image: first pixel is lower left corner
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, w*h*4)
	lineLen := w * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgbaImg.Pix); src += rgbaImg.Stride {
		copy(data[dest:dest+lineLen], rgbaImg.Pix[src:src+rgbaImg.Stride])
		dest -= lineLen
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)

	return tex
}

func NewUnitCubeRenderObject() *RenderObject {
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

	ro := new(RenderObject)

	var size int
	size = len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	ro.vertexBuffer = MakeBuffer(gl.ARRAY_BUFFER, size, vertexPositions)

	size = len(indicies) * int(unsafe.Sizeof(indicies[0]))
	ro.elementBuffer = MakeBuffer(gl.ELEMENT_ARRAY_BUFFER, size, indicies)

	ro.numIndicies = len(indicies)

	return ro
}
