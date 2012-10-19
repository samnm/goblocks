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
	fmt.Println(shader.GetInfoLog())

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
	verticies := []Vertex{
		Vertex{[3]float32{-0.5, 0.5, 0.5}, [2]float32{0.0, 1.0}}, //+Z Face
		Vertex{[3]float32{-0.5, -0.5, 0.5}, [2]float32{0.0, 0.0}},
		Vertex{[3]float32{0.5, -0.5, 0.5}, [2]float32{1.0, 0.0}},
		Vertex{[3]float32{0.5, 0.5, 0.5}, [2]float32{1.0, 1.0}},
		Vertex{[3]float32{-0.5, 0.5, 0.5}, [2]float32{0.0, 1.0}},
		Vertex{[3]float32{0.5, -0.5, 0.5}, [2]float32{1.0, 0.0}},
		Vertex{[3]float32{-0.5, 0.5, -0.5}, [2]float32{1.0, 1.0}}, //-Z Face
		Vertex{[3]float32{0.5, -0.5, -0.5}, [2]float32{0.0, 0.0}},
		Vertex{[3]float32{-0.5, -0.5, -0.5}, [2]float32{1.0, 0.0}},
		Vertex{[3]float32{0.5, 0.5, -0.5}, [2]float32{0.0, 1.0}},
		Vertex{[3]float32{0.5, -0.5, -0.5}, [2]float32{0.0, 0.0}},
		Vertex{[3]float32{-0.5, 0.5, -0.5}, [2]float32{1.0, 1.0}},
	}

	ro := new(RenderObject)

	size := len(verticies) * int(unsafe.Sizeof(verticies[0]))
	ro.vertexBuffer = MakeBuffer(gl.ARRAY_BUFFER, size, verticies)
	ro.numVerticies = len(verticies)

	return ro
}
