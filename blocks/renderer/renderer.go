package renderer

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
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

	timer float32

	timerUniform    gl.UniformLocation
	textureUniforms [2]gl.UniformLocation
	positionAttrib  gl.AttribLocation
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

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.ClearDepth(1.0)

	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	width, height := glfw.WindowSize()
	SetViewport(width, height)

	glfw.SetWindowSizeCallback(SetViewport)
}

func InitProgram() {
	vertexShader = MakeShader(gl.VERTEX_SHADER, "resources/shaders/view-frustum-rotation.v.glsl")
	fragShader = MakeShader(gl.FRAGMENT_SHADER, "resources/shaders/hello-gl.f.glsl")

	program = gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragShader)
	program.Link()
}

func InitBuffers() {
	vertexPositions := []float32{
		-1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0,
		-1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 0.0, 1.0,
	}
	indicies := []gl.GLushort{0, 1, 2, 3}

	var size int
	size = len(vertexPositions) * int(unsafe.Sizeof(vertexPositions[0]))
	vertexBuffer = MakeBuffer(gl.ARRAY_BUFFER, size, vertexPositions)

	size = len(indicies) * int(unsafe.Sizeof(indicies[0]))
	elementBuffer = MakeBuffer(gl.ELEMENT_ARRAY_BUFFER, size, indicies)
}

func InitAttribs() {
	timerUniform = program.GetUniformLocation("timer")
	textureUniforms[0] = program.GetUniformLocation("textures[0]")
	textureUniforms[1] = program.GetUniformLocation("textures[1]")
	positionAttrib = program.GetAttribLocation("position")
}

const sizeOfGLFloat int = int(unsafe.Sizeof(float32(0.0)))

func Tick() {
	timer = float32(glfw.Time())

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	program.Use()

	timerUniform.Uniform1f(timer)

	gl.ActiveTexture(gl.TEXTURE0)
	textures[0].Bind(gl.TEXTURE_2D)
	textureUniforms[0].Uniform1i(0)

	gl.ActiveTexture(gl.TEXTURE1)
	textures[1].Bind(gl.TEXTURE_2D)
	textureUniforms[1].Uniform1i(1)

	vertexBuffer.Bind(gl.ARRAY_BUFFER)
	positionAttrib.AttribPointer(4, gl.FLOAT, false, sizeOfGLFloat*4, nil)
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
