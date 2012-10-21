package renderer

import (
	"math"
)

type Camera struct {
	modelViewMatrix  [16]float32
	projectionMatrix [16]float32
	eyePosition      [3]float32
}

func NewCamera(viewWidth, viewHeight int) *Camera {
	camera := &Camera{}
	camera.projectionMatrix = [16]float32{}
	camera.modelViewMatrix = [16]float32{}
	camera.eyePosition = [3]float32{5.0, 10.0, -5.0}

	camera.UpdateProjectionMatrix(viewWidth, viewHeight)
	camera.UpdateModelViewMatrix()

	return camera
}

const DegToRad = math.Pi / 180

func (c *Camera) UpdateProjectionMatrix(width, height int) {
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

	c.projectionMatrix[0] = float32(r_x)
	c.projectionMatrix[1] = 0.0
	c.projectionMatrix[2] = 0.0
	c.projectionMatrix[3] = 0.0

	c.projectionMatrix[4] = 0.0
	c.projectionMatrix[5] = float32(r_y)
	c.projectionMatrix[6] = 0.0
	c.projectionMatrix[7] = 0.0

	c.projectionMatrix[8] = 0.0
	c.projectionMatrix[9] = 0.0
	c.projectionMatrix[10] = float32(r_z)
	c.projectionMatrix[11] = 1.0

	c.projectionMatrix[12] = 0.0
	c.projectionMatrix[13] = 0.0
	c.projectionMatrix[14] = float32(r_w)
	c.projectionMatrix[15] = 0.0
}

func (c *Camera) UpdateModelViewMatrix() {
	c.modelViewMatrix[0] = 1.0
	c.modelViewMatrix[1] = 0.0
	c.modelViewMatrix[2] = 0.0
	c.modelViewMatrix[3] = 0.0

	c.modelViewMatrix[4] = 0.0
	c.modelViewMatrix[5] = 1.0
	c.modelViewMatrix[6] = 0.0
	c.modelViewMatrix[7] = 0.0

	c.modelViewMatrix[8] = 0.0
	c.modelViewMatrix[9] = 0.0
	c.modelViewMatrix[10] = 1.0
	c.modelViewMatrix[11] = 0.0

	c.modelViewMatrix[12] = -c.eyePosition[0]
	c.modelViewMatrix[13] = -c.eyePosition[1]
	c.modelViewMatrix[14] = -c.eyePosition[2]
	c.modelViewMatrix[15] = 1.0
}
