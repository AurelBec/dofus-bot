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
	ID  string `json:"id"`
	Pos Pos    `json:"position"`

	Invert    bool    `json:"inversionFlag"`
	Gray      float64 `json:"gray"`
	Lightness float64 `json:"lightness"`
}

func (r Resource) String() string {
	return r.ID
}

func NewResourceUnderMouse(invert bool) Resource {
	x, y := robotgo.GetMousePos()
	defer robotgo.MoveMouse(x, y)

	robotgo.MoveMouse(0, 0)
	time.Sleep(time.Millisecond * 100)
	color := getPixelColor(x, y)

	id := fmt.Sprintf("%vx%v", x, y)
	gray := color.toGray()
	lightness := color.toLightness()

	logrus.Infof("register resource [%s] with params %.3f, %.3f", id, gray, lightness)

	return Resource{
		ID:        id,
		Pos:       Pos{X: x, Y: y},
		Invert:    invert,
		Gray:      gray,
		Lightness: lightness,
	}
}

func (r Resource) IsActive() bool {
	color := getPixelColor(r.Pos.X, r.Pos.Y)

	limit := 10.
	grayOk := math.Abs(color.toGray()-r.Gray) < limit
	lightnessOk := math.Abs(color.toLightness()-r.Lightness) < limit

	return (grayOk && lightnessOk) != r.Invert
}

func (r Resource) Collect(react bool) {
	// simulate reaction time
	if react {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(750)+750))
	}

	robotgo.KeyToggle("lshift", "down")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveClick(r.Pos.X, r.Pos.Y, "left", true)
	time.Sleep(time.Millisecond * 20)
	robotgo.KeyToggle("lshift", "up")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveMouse(400, 0)
}
