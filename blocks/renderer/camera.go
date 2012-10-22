package renderer

import (
	"github.com/go-gl/glfw"
	"sam.vg/goblocks/blocks/input"
	"sam.vg/util/matrix"
)

type Camera struct {
	modelViewMatrix  *matrix.Matrix4
	projectionMatrix *matrix.Matrix4
	eyePosition      [3]float32
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

	m := matrix.MakeRotationMatrix(-input.MouseX()/100.0, 0, 1, 0)
	movement = m.MultiplyPoint(movement)

	c.eyePosition[0] += movement[0]
	c.eyePosition[1] += movement[1]
	c.eyePosition[2] += movement[2]

	c.UpdateModelViewMatrix()
}

func (c *Camera) UpdateProjectionMatrix(width, height int) {
	c.projectionMatrix = matrix.MakePerspectiveMatrix(width, height, 60.0, 0.0625, 256.0)
	c.projectionMatrix.Transpose()
}

func (c *Camera) UpdateModelViewMatrix() {
	c.modelViewMatrix = matrix.MakeRotationMatrix(-input.MouseY()/100.0, 1, 0, 0)
	c.modelViewMatrix = c.modelViewMatrix.Multiply(matrix.MakeRotationMatrix(-input.MouseX()/100.0, 0, 1, 0))
	c.modelViewMatrix = c.modelViewMatrix.Multiply(matrix.MakeTranslationMatrix(-c.eyePosition[0], -c.eyePosition[1], -c.eyePosition[2]))
	c.modelViewMatrix.Transpose()
}
