package models

import (
	"encoding/hex"

	"github.com/gerow/go-color"
	"github.com/go-vgo/robotgo"
)

type colorData struct {
	color.RGB
}

func getPixelColor(x, y int) colorData {
	str := robotgo.GetPixelColor(x, y)
	rgb, _ := hex.DecodeString(str)

	return colorData{
		color.RGB{
			R: float64(rgb[0]) / 255,
			G: float64(rgb[1]) / 255,
			B: float64(rgb[2]) / 255,
		},
	}
}

func (c colorData) toGray() float64 {
	return (c.R + c.G + c.B) * 255. / 3.
}

func (c colorData) toLightness() float64 {
	return c.ToHSL().L * 100.
}
