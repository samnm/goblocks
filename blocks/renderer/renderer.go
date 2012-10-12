package renderer

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"os"
	"unsafe"
)

var (
	textures      [2]gl.Texture
	program       gl.Program
	vertexShader  gl.Shader
	fragShader    gl.Shader
	vertexBuffer  gl.Buffer
	elementBuffer gl.Buffer

	fadeFactor float32

	fadeFactorUniform gl.UniformLocation
	textureUniforms   [2]gl.UniformLocation
	positionAttrib    gl.AttribLocation
)

func Init() {
	InitGL()

	var err error
	textures[0], err = LoadTexture("resources/hello1.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
	}
	textures[1], err = LoadTexture("resources/hello2.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
	}
	InitProgram()
	InitBuffers()
	InitAttribs()
}

func InitGL() {
	//enable vertical sync if the card supports it
	glfw.SetSwapInterval(1)

	gl.ShadeModel(gl.SMOOTH)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)

	// gl.Enable(gl.TEXTURE_2D)
	// gl.Enable(gl.CULL_FACE)

	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LEQUAL)

	// gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	// width, height := glfw.WindowSize()
	// SetViewport(width, height)

	glfw.SetWindowSizeCallback(SetViewport)
}

func InitProgram() {
	vertexShaderSrc := `#version 110

attribute vec2 position;

varying vec2 texcoord;

void main()
{
    gl_Position = vec4(position, 0.0, 1.0);
    texcoord = position * vec2(0.5) + vec2(0.5);
}`
	fragShaderSrc := `#version 110

uniform float fade_factor;
uniform sampler2D textures[2];

varying vec2 texcoord;

void main()
{
    gl_FragColor = mix(
        texture2D(textures[0], texcoord),
        texture2D(textures[1], texcoord),
        fade_factor
    );
}`
	vertexShader = gl.CreateShader(gl.VERTEX_SHADER)
	vertexShader.Source(vertexShaderSrc)
	vertexShader.Compile()

	fragShader = gl.CreateShader(gl.FRAGMENT_SHADER)
	fragShader.Source(fragShaderSrc)
	fragShader.Compile()

	program = gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragShader)
	program.Link()
}

func InitBuffers() {
	vertexPositions := []float32{
		-1.0, -1.0,
		1.0, -1.0,
		-1.0, 1.0,
		1.0, 1.0,
	}
	indicies := []gl.GLushort{0, 1, 2, 3}

	var size int
	size = len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	vertexBuffer = MakeBuffer(gl.ARRAY_BUFFER, size, vertexPositions)

	size = len(indicies) * int(unsafe.Sizeof(indicies[0]))
	elementBuffer = MakeBuffer(gl.ELEMENT_ARRAY_BUFFER, size, indicies)
}

func InitAttribs() {
	fadeFactorUniform = program.GetUniformLocation("fade_factor")
	textureUniforms[0] = program.GetUniformLocation("textures[0]")
	textureUniforms[1] = program.GetUniformLocation("textures[1]")
	positionAttrib = program.GetAttribLocation("position")
}

func Tick() {
	fadeFactor = float32(math.Sin(glfw.Time())*0.5 + 0.5)

	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()

	fadeFactorUniform.Uniform1f(fadeFactor)

	gl.ActiveTexture(gl.TEXTURE0)
	textures[0].Bind(gl.TEXTURE_2D)
	textureUniforms[0].Uniform1i(0)

	gl.ActiveTexture(gl.TEXTURE1)
	textures[1].Bind(gl.TEXTURE_2D)
	textureUniforms[1].Uniform1i(1)

	vertexBuffer.Bind(gl.ARRAY_BUFFER)
	var aGLfloat gl.GLfloat
	positionAttrib.AttribPointer(2, gl.FLOAT, false, int(unsafe.Sizeof(aGLfloat))*2, nil)
	positionAttrib.EnableArray()

	elementBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.DrawElements(gl.TRIANGLE_STRIP, 4, gl.UNSIGNED_SHORT, nil)

	positionAttrib.DisableArray()
	gl.ProgramUnuse()

	glfw.SwapBuffers()
}

func OnResizeWindow(width, height int) {
	SetViewport(width, height)
}

func SetViewport(width, height int) {
	// Setup our viewport
	gl.Viewport(0, 0, width, height)

	// change to the projection matrix and set our viewing volume.
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	// aspect ratio
	aspect := float64(width) / float64(height)

	// Set our perspective.
	// This code is equivalent to using gluPerspective as in the original tutorial.
	fov := 60.0
	near := 0.1
	far := 100.0
	top := math.Tan(fov*math.Pi/360.0) * near
	bottom := -top
	left := aspect * bottom
	right := aspect * top
	gl.Frustum(left, right, bottom, top, near, far)

	// Make sure we're changing the model view and not the projection
	gl.MatrixMode(gl.MODELVIEW)

	// Reset the view
	gl.LoadIdentity()
}

func MakeBuffer(target gl.GLenum, size int, data interface{}) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(target)
	gl.BufferData(target, size, data, gl.STATIC_DRAW)
	gl.BufferUnbind(target)
	return buffer
}

func LoadTexture(filename string) (gl.Texture, error) {
	r, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
	}
	defer r.Close()

	img, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}

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

	return tex, nil
}
