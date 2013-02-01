package fps

import "fmt"

type FPS struct {
	frames int
	time   float64
}

func NewFPS(time float64) *FPS {
	return &FPS{0, time}
}

func (f *FPS) Tick(currentTime float64) {
	f.frames++
	if currentTime-f.time > 1.0 {
		f.time = currentTime
		fmt.Printf("FPS: %d\n", f.frames)
		f.frames = 0
	}
}
