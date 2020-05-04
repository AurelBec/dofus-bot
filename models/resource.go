package models

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
)

type Resource struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`

	Invert    bool    `json:"inversionFlag"`
	Gray      float64 `json:"gray"`
	Lightness float64 `json:"lightness"`
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
		X:         x,
		Y:         y,
		Invert:    invert,
		Gray:      gray,
		Lightness: lightness,
	}
}

func (r Resource) SquareDistanceTo(other Resource) float64 {
	return float64((r.X-other.X)*(r.X-other.X) + (r.Y-other.Y)*(r.Y-other.Y))
}

func (r Resource) IsActive() bool {
	color := getPixelColor(r.X, r.Y)

	limit := 10.
	grayOk := math.Abs(color.toGray()-r.Gray) < limit
	lightnessOk := math.Abs(color.toLightness()-r.Lightness) < limit

	return (grayOk && lightnessOk) != r.Invert
}

func (r Resource) Collect(react bool) {
	// simulate reaction time
	if react {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(800)+500))
	}

	robotgo.KeyToggle("lshift", "down")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveClick(r.X, r.Y, "left", true)
	time.Sleep(time.Millisecond * 20)
	robotgo.KeyToggle("lshift", "up")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveMouse(400, 0)
}
