package renderer

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"unsafe"
)

type Vertex struct {
	position [3]float32
	texcoord [2]float32
}

type RenderObject struct {
	vertexBuffer gl.Buffer
	numVerticies int
}

var (
	textures     [2]gl.Texture
	program      gl.Program
	vertexShader gl.Shader
	fragShader   gl.Shader

	renderObjects []*RenderObject
	culling       bool = true

	camera *Camera
	timer  float32
)

var (
	timerUniform            gl.UniformLocation
	textureUniforms         [2]gl.UniformLocation
	modelViewMatrixUniform  gl.UniformLocation
	projectionMatrixUniform gl.UniformLocation

	positionAttrib gl.AttribLocation
	texcoordAttrib gl.AttribLocation
)

func Init(width, height int) {
	camera = NewCamera(width, height)

	InitGL(width, height)

	textures[0] = LoadTexture("resources/hello1.png")
	textures[1] = LoadTexture("resources/hello2.png")

	InitProgram()
	InitProperties()

	renderObjects = []*RenderObject{NewUnitCubeRenderObject()}
}

func InitGL(width, height int) {
	//enable vertical sync if the card supports it
	glfw.SetSwapInterval(1)

	gl.ShadeModel(gl.SMOOTH)

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)

	gl.Enable(gl.TEXTURE_2D)
	// gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	SetViewport(width, height)
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

func InitProperties() {
	timerUniform = program.GetUniformLocation("timer")
	textureUniforms[0] = program.GetUniformLocation("textures[0]")
	textureUniforms[1] = program.GetUniformLocation("textures[1]")
	modelViewMatrixUniform = program.GetUniformLocation("mv_matrix")
	projectionMatrixUniform = program.GetUniformLocation("p_matrix")

	positionAttrib = program.GetAttribLocation("position")
	texcoordAttrib = program.GetAttribLocation("texcoord")
}

const sizeOfGLFloat int = int(unsafe.Sizeof(float32(0.0)))

func Tick() {
	timer = float32(glfw.Time())

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	program.Use()

	modelViewMatrixUniform.UniformMatrix4fv(camera.modelViewMatrix)
	projectionMatrixUniform.UniformMatrix4fv(camera.projectionMatrix)
	timerUniform.Uniform1f(timer)

	positionAttrib.EnableArray()
	texcoordAttrib.EnableArray()

	gl.ActiveTexture(gl.TEXTURE0)
	textures[0].Bind(gl.TEXTURE_2D)
	textureUniforms[0].Uniform1i(0)

	gl.ActiveTexture(gl.TEXTURE1)
	textures[1].Bind(gl.TEXTURE_2D)
	textureUniforms[1].Uniform1i(1)

	sizeOfVertex := int(unsafe.Sizeof(Vertex{}))
	posoffset := uintptr(0)
	texoffset := unsafe.Offsetof(Vertex{}.texcoord)
	for _, renderObject := range renderObjects {
		renderObject.vertexBuffer.Bind(gl.ARRAY_BUFFER)
		positionAttrib.AttribPointer(3, gl.FLOAT, false, sizeOfVertex, posoffset)
		texcoordAttrib.AttribPointer(2, gl.FLOAT, false, sizeOfVertex, texoffset)

		gl.DrawArrays(gl.TRIANGLES, 0, renderObject.numVerticies)
	}

	positionAttrib.DisableArray()
	texcoordAttrib.DisableArray()
	gl.ProgramUnuse()

	glfw.SwapBuffers()
}

func OnResizeWindow(width, height int) {
	SetViewport(width, height)
}

func SetViewport(width, height int) {
	gl.Viewport(0, 0, width, height)
	camera.UpdateProjectionMatrix(width, height)
}
