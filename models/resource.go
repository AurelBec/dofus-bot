package models

import (
	"fmt"
	"math"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
)

type Resource struct {
	ID   string
	x, y int

	invert    bool
	gray      float64
	lightness float64
}

func NewResourceUnderMouse(invert bool) Resource {
	x, y := robotgo.GetMousePos()
	robotgo.MoveMouse(0, 0)
	time.Sleep(time.Millisecond * 100)
	return NewResource(x, y, invert)
}

func NewResource(x, y int, invert bool) Resource {
	id := fmt.Sprintf("%vx%v", x, y)

	color := getPixelColor(x, y)
	gray := color.toGray()
	lightness := color.toLightness()

	logrus.Infof("register resource [%s] with params %.3f, %.3f", id, gray, lightness)
	robotgo.MoveMouse(x, y)

	return Resource{
		ID:        id,
		x:         x,
		y:         y,
		invert:    invert,
		gray:      gray,
		lightness: lightness,
	}
}

func (r Resource) SquareDistanceTo(other Resource) float64 {
	return float64((r.x-other.x)*(r.x-other.x) + (r.y-other.y)*(r.y-other.y))
}

func (r Resource) IsActive() bool {
	color := getPixelColor(r.x, r.y)

	limit := 10.
	grayOk := math.Abs(color.toGray()-r.gray) < limit
	lightnessOk := math.Abs(color.toLightness()-r.lightness) < limit

	return (grayOk && lightnessOk) != r.invert
}

func (r Resource) Collect() {
	robotgo.KeyToggle("lshift", "down")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveClick(r.x, r.y, "left", true)
	time.Sleep(time.Millisecond * 20)
	robotgo.KeyToggle("lshift", "up")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveMouse(400, 0)
}
