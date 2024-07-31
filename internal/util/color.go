package util

import "image/color"

func ColorToNRGB(clr color.Color) (r uint8, g uint8, b uint8) {
	r16, g16, b16, a16 := clr.RGBA()
	factor := float32(a16) / float32(65535)
	nr8 := uint32(float32(r16)/factor) >> 8
	ng8 := uint32(float32(g16)/factor) >> 8
	nb8 := uint32(float32(b16)/factor) >> 8
	return uint8(nr8), uint8(ng8), uint8(nb8)
}
