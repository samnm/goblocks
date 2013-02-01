package matrix

import "math"

type Matrix4 [4 * 4]float32

func MakeRotationMatrix(angle, x, y, z float32) *Matrix4 {
	cosA := float32(math.Cos(float64(angle)))
	sinA := float32(math.Sin(float64(angle)))
	return &Matrix4{
		cosA + ((x * x) * (1 - cosA)), ((x * y) * (1 - cosA)) - (z * sinA), ((x * z) * (1 - cosA)) + (y * sinA), 0,
		((y * x) * (1 - cosA)) + (z * sinA), cosA + ((y * y) * (1 - cosA)), ((y * z) * (1 - cosA)) - (x * sinA), 0,
		((z * x) * (1 - cosA)) - (y * sinA), ((z * y) * (1 - cosA)) + (x * sinA), cosA + ((z * z) * (1 - cosA)), 0,
		0, 0, 0, 1,
	}
}

func MakeTranslationMatrix(x, y, z float32) *Matrix4 {
	return &Matrix4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

func MakeScaleMatrix(x, y, z float32) *Matrix4 {
	return &Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

func MakePerspectiveMatrix(width, height int, fov, near, far float32) *Matrix4 {
	fov = fov * math.Pi / 180

	r_xy_factor := float32(math.Min(float64(width), float64(height))) * 1.0 / fov
	r_x := r_xy_factor / float32(width)
	r_y := r_xy_factor / float32(height)
	r_zw_factor := 1.0 / (far - near)
	r_z := (near + far) * r_zw_factor
	r_w := -2.0 * near * far * r_zw_factor

	matrix := &Matrix4{}
	matrix.Set(0, 0, r_x)
	matrix.Set(1, 1, r_y)
	matrix.Set(2, 2, r_z)
	matrix.Set(2, 3, 1)
	matrix.Set(3, 2, r_w)

	return matrix
}

func (m *Matrix4) Get(x, y int) float32 {
	return m[x+y*4]
}

func (m *Matrix4) Set(x, y int, value float32) {
	m[x+y*4] = value
}

func (a *Matrix4) Multiply(b *Matrix4) *Matrix4 {
	result := &Matrix4{}
	for ix := 0; ix < 4; ix++ {
		for iy := 0; iy < 4; iy++ {
			row := [4]float32{a.Get(0, iy), a.Get(1, iy), a.Get(2, iy), a.Get(3, iy)}
			col := [4]float32{b.Get(ix, 0), b.Get(ix, 1), b.Get(ix, 2), b.Get(ix, 3)}
			result.Set(ix, iy, row[0]*col[0]+row[1]*col[1]+row[2]*col[2]+row[3]*col[3])
		}
	}
	return result
}

func (a *Matrix4) MultiplyPoint(point [4]float32) [4]float32 {
	return [4]float32{
		a.Get(0, 0)*point[0] + a.Get(0, 1)*point[1] + a.Get(0, 2)*point[2] + a.Get(0, 3)*point[3],
		a.Get(1, 0)*point[0] + a.Get(1, 1)*point[1] + a.Get(1, 2)*point[2] + a.Get(1, 3)*point[3],
		a.Get(2, 0)*point[0] + a.Get(2, 1)*point[1] + a.Get(2, 2)*point[2] + a.Get(2, 3)*point[3],
		a.Get(3, 0)*point[0] + a.Get(3, 1)*point[1] + a.Get(3, 2)*point[2] + a.Get(3, 3)*point[3],
	}
}

func (m *Matrix4) Transpose() {
	m[1], m[4] = m[4], m[1]
	m[2], m[8] = m[8], m[2]
	m[3], m[12] = m[12], m[3]
	m[6], m[9] = m[9], m[6]
	m[7], m[13] = m[13], m[7]
	m[11], m[14] = m[14], m[11]
}
