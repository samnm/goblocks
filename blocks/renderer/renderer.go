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
)

var (
	timerUniform            gl.UniformLocation
	textureUniforms         [2]gl.UniformLocation
	modelViewMatrixUniform  gl.UniformLocation
	projectionMatrixUniform gl.UniformLocation
)

var (
	positionAttrib gl.AttribLocation
)

var (
	modelViewMatrix  []float32
	projectionMatrix []float32
	eyePosition      [3]float32
	timer            float32
)

func Init(width, height int) {
	projectionMatrix = make([]float32, 16)
	modelViewMatrix = make([]float32, 16)
	eyePosition = [3]float32{0.0, 0.0, -5.0}

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
	InitProperties()
}

func InitGL() {
	//enable vertical sync if the card supports it
	glfw.SetSwapInterval(1)

	gl.ShadeModel(gl.SMOOTH)

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.ClearDepth(1.0)

	gl.Enable(gl.TEXTURE_2D)
	// gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	width, height := glfw.WindowSize()
	SetViewport(width, height)
	UpdateModelViewMatrix()

	glfw.SetWindowSizeCallback(SetViewport)
}

func InitProgram() {
	vertexShader = MakeShader(gl.VERTEX_SHADER, "resources/shaders/blocks.vert")
	fragShader = MakeShader(gl.FRAGMENT_SHADER, "resources/shaders/blocks.frag")

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

func InitProperties() {
	timerUniform = program.GetUniformLocation("timer")
	textureUniforms[0] = program.GetUniformLocation("textures[0]")
	textureUniforms[1] = program.GetUniformLocation("textures[1]")
	modelViewMatrixUniform = program.GetUniformLocation("mv_matrix")
	projectionMatrixUniform = program.GetUniformLocation("p_matrix")

	positionAttrib = program.GetAttribLocation("position")
}

const sizeOfGLFloat int = int(unsafe.Sizeof(float32(0.0)))

func Tick() {
	timer = float32(glfw.Time())

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	program.Use()

	modelViewMatrixUniform.UniformMatrix4fv(modelViewMatrix)
	projectionMatrixUniform.UniformMatrix4fv(projectionMatrix)
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

	UpdateProjectionMatrix(width, height)
}

const DegToRad = math.Pi / 180

func UpdateProjectionMatrix(width, height int) {
	fov := 60.0 * DegToRad
	near := 0.0625
	far := 256.0

	w, h := float64(width), float64(height)
	r_xy_factor := math.Min(w, h) * 1.0 / fov
	r_x := r_xy_factor / w
	r_y := r_xy_factor / h
	r_zw_factor := 1.0 / (far - near)
	r_z := (near + far) * r_zw_factor
	r_w := -2.0 * near * far * r_zw_factor

	projectionMatrix[0] = float32(r_x)
	projectionMatrix[1] = 0.0
	projectionMatrix[2] = 0.0
	projectionMatrix[3] = 0.0

	projectionMatrix[4] = 0.0
	projectionMatrix[5] = float32(r_y)
	projectionMatrix[6] = 0.0
	projectionMatrix[7] = 0.0

	projectionMatrix[8] = 0.0
	projectionMatrix[9] = 0.0
	projectionMatrix[10] = float32(r_z)
	projectionMatrix[11] = 1.0

	projectionMatrix[12] = 0.0
	projectionMatrix[13] = 0.0
	projectionMatrix[14] = float32(r_w)
	projectionMatrix[15] = 0.0
}

func UpdateModelViewMatrix() {
	modelViewMatrix[0] = 1.0
	modelViewMatrix[1] = 0.0
	modelViewMatrix[2] = 0.0
	modelViewMatrix[3] = 0.0

	modelViewMatrix[4] = 0.0
	modelViewMatrix[5] = 1.0
	modelViewMatrix[6] = 0.0
	modelViewMatrix[7] = 0.0

	modelViewMatrix[8] = 0.0
	modelViewMatrix[9] = 0.0
	modelViewMatrix[10] = 1.0
	modelViewMatrix[11] = 0.0

	modelViewMatrix[12] = -eyePosition[0]
	modelViewMatrix[13] = -eyePosition[1]
	modelViewMatrix[14] = -eyePosition[2]
	modelViewMatrix[15] = 1.0
}
