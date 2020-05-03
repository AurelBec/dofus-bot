package models

import (
	"encoding/hex"

	"github.com/gerow/go-color"
	"github.com/go-vgo/robotgo"
)

type colorData struct {
	str string
	rgb []float64
}

func getPixelColor(x, y int) colorData {
	str := robotgo.GetPixelColor(x, y)
	rgb, _ := hex.DecodeString(str)

	return colorData{
		str: str,
		rgb: []float64{float64(rgb[0]), float64(rgb[1]), float64(rgb[2])},
	}
}

func (c colorData) toGray() float64 {
	return (c.rgb[0] + c.rgb[1] + c.rgb[2]) / 3.0
}

func (c colorData) toLightness() float64 {
	return color.RGB{
		R: c.rgb[0] / 255.,
		G: c.rgb[1] / 255.,
		B: c.rgb[2] / 255.,
	}.ToHSL().L * 100.
}
