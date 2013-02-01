package renderer

import (
	"github.com/go-gl/glfw"
	"github.com/samnm/goblocks/blocks/input"
	"github.com/samnm/goblocks/util/matrix"
	"math"
)

type Camera struct {
	modelViewMatrix  *matrix.Matrix4
	projectionMatrix *matrix.Matrix4
	eyePosition      [3]float32
	yaw, pitch       float32
	mouseX, mouseY   float32
}

func NewCamera(viewWidth, viewHeight int) *Camera {
	camera := &Camera{}
	camera.eyePosition = [3]float32{5.0, 10.0, -5.0}

	camera.UpdateProjectionMatrix(viewWidth, viewHeight)
	camera.UpdateModelViewMatrix()

	return camera
}

func (c *Camera) Tick() {
	moveSpeed := float32(0.1)
	movement := [4]float32{}

	if input.IsKeyDown(glfw.KeyUp) || input.IsKeyDown('W') {
		movement[2] += moveSpeed
	}
	if input.IsKeyDown(glfw.KeyDown) || input.IsKeyDown('S') {
		movement[2] -= moveSpeed
	}
	if input.IsKeyDown(glfw.KeyLeft) || input.IsKeyDown('A') {
		movement[0] -= moveSpeed
	}
	if input.IsKeyDown(glfw.KeyRight) || input.IsKeyDown('D') {
		movement[0] += moveSpeed
	}

	c.UpdateRotation()

	m := matrix.MakeRotationMatrix(c.pitch, 1, 0, 0)
	m = m.Multiply(matrix.MakeRotationMatrix(c.yaw, 0, 1, 0))
	movement = m.MultiplyPoint(movement)

	c.eyePosition[0] += movement[0]
	c.eyePosition[1] += movement[1]
	c.eyePosition[2] += movement[2]

	c.UpdateModelViewMatrix()
}

func (c *Camera) UpdateRotation() {
	newMouseX, newMouseY := input.MouseX(), input.MouseY()

	c.yaw += -(newMouseX - c.mouseX) / 100.0
	c.pitch += -(newMouseY - c.mouseY) / 100.0

	c.yaw = float32(math.Mod(float64(c.yaw), math.Pi*2))
	c.pitch = float32(math.Max(-math.Pi/2, math.Min(math.Pi/2, float64(c.pitch))))

	c.mouseX, c.mouseY = newMouseX, newMouseY
}

func (c *Camera) UpdateProjectionMatrix(width, height int) {
	c.projectionMatrix = matrix.MakePerspectiveMatrix(width, height, 60.0, 0.0625, 256.0)
	c.projectionMatrix.Transpose()
}

func (c *Camera) UpdateModelViewMatrix() {
	c.modelViewMatrix = matrix.MakeRotationMatrix(c.pitch, 1, 0, 0)
	c.modelViewMatrix = c.modelViewMatrix.Multiply(matrix.MakeRotationMatrix(c.yaw, 0, 1, 0))
	c.modelViewMatrix = c.modelViewMatrix.Multiply(matrix.MakeTranslationMatrix(-c.eyePosition[0], -c.eyePosition[1], -c.eyePosition[2]))
	c.modelViewMatrix.Transpose()
}
