package renderer

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/larspensjo/Go-simplex-noise/simplexnoise"
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

var cubeVerticies = [8][3]float32{
	[3]float32{0.5, 0.5, 0.5},
	[3]float32{-0.5, 0.5, 0.5},
	[3]float32{0.5, -0.5, 0.5},
	[3]float32{-0.5, -0.5, 0.5},

	[3]float32{0.5, 0.5, -0.5},
	[3]float32{-0.5, 0.5, -0.5},
	[3]float32{0.5, -0.5, -0.5},
	[3]float32{-0.5, -0.5, -0.5},
}

var texCoords = [4][2]float32{
	[2]float32{0, 0},
	[2]float32{1, 0},
	[2]float32{0, 1},
	[2]float32{1, 1},
}

var faces = [6][6]Vertex{
	BuildFace([6]int{1, 2, 3, 0, 2, 1}, [6]int{3, 0, 1, 2, 0, 3}), //+Z Face (far)
	BuildFace([6]int{5, 7, 6, 4, 5, 6}, [6]int{2, 0, 1, 3, 2, 1}), //-Z Face (close)
	BuildFace([6]int{4, 6, 0, 6, 2, 0}, [6]int{2, 0, 3, 0, 1, 3}), //+X Face (right)
	BuildFace([6]int{5, 1, 7, 7, 1, 3}, [6]int{3, 2, 1, 1, 2, 0}), //-X Face (left)
	BuildFace([6]int{4, 0, 5, 5, 0, 1}, [6]int{3, 2, 1, 1, 2, 0}), //+Y Face (top)
	BuildFace([6]int{6, 7, 2, 7, 3, 2}, [6]int{3, 1, 2, 1, 0, 2}), //-Y Face (bottom)
}

func NewChunk() *RenderObject {
	chunkSize := 16
	verticies := make([]Vertex, 0, chunkSize*chunkSize*chunkSize*6)

	for ix := 0; ix < chunkSize; ix++ {
		for iz := 0; iz < chunkSize; iz++ {
			height := int(((simplexnoise.Noise2(float64(ix)/20, float64(iz)/20) + 1) / 2) * float64(chunkSize))
			verticies = append(verticies, OffsetFace(0, ix, height, iz)...)
			verticies = append(verticies, OffsetFace(1, ix, height, iz)...)
			verticies = append(verticies, OffsetFace(2, ix, height, iz)...)
			verticies = append(verticies, OffsetFace(3, ix, height, iz)...)
			verticies = append(verticies, OffsetFace(4, ix, height, iz)...)
			verticies = append(verticies, OffsetFace(5, ix, height, iz)...)
		}
	}

	ro := new(RenderObject)

	numVerticies := len(verticies)
	size := numVerticies * int(unsafe.Sizeof(verticies[0]))
	ro.vertexBuffer = MakeBuffer(gl.ARRAY_BUFFER, size, verticies[:numVerticies])
	ro.numVerticies = numVerticies

	return ro
}

func MovePoint(p [3]float32, x, y, z int) [3]float32 {
	return [3]float32{p[0] + float32(x), p[1] + float32(y), p[2] + float32(z)}
}

func BuildFace(vertIndicies [6]int, texIndicies [6]int) [6]Vertex {
	verticies := [6]Vertex{}
	for i := 0; i < 6; i++ {
		verticies[i] = Vertex{cubeVerticies[vertIndicies[i]], texCoords[texIndicies[i]]}
	}
	return verticies
}

func OffsetFace(face, x, y, z int) []Vertex {
	faceVerts := faces[face]
	offsetFace := make([]Vertex, 6)
	for i, v := range faceVerts {
		offsetFace[i] = Vertex{MovePoint(v.position, x, y, z), v.texcoord}
	}
	return offsetFace
}
