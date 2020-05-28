package models

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
	color "github.com/lucasb-eyer/go-colorful"
	"github.com/sirupsen/logrus"
)

type Resource struct {
	ID    string      `json:"id"`
	Pos   Pos         `json:"position"`
	Color color.Color `json:"color"`

	new          bool
	colorUpdated bool

	// old
	Gray      float64 `json:"gray"`
	Lightness float64 `json:"lightness"`
}

func (r Resource) String() string {
	return r.ID
}

func getPixelColor(x, y int) color.Color {
	c, _ := color.Hex("#" + robotgo.GetPixelColor(x, y))
	return c
}

func NewResourceUnderMouse() *Resource {
	x, y := robotgo.GetMousePos()
	defer robotgo.MoveMouse(x, y)

	robotgo.MoveMouse(0, 0)
	time.Sleep(time.Millisecond * 100)

	id := fmt.Sprintf("%vx%v", x, y)
	logrus.Infof("register resource [%s]", id)

	return &Resource{
		ID:    id,
		Pos:   Pos{X: x, Y: y},
		new:   true,
		Color: getPixelColor(x, y),
	}
}

func (r Resource) IsNew() bool {
	return r.new
}

func (r Resource) ColorUpdated() bool {
	return r.colorUpdated
}

func (r *Resource) IsActive() bool {
	colorXY := getPixelColor(r.Pos.X, r.Pos.Y)
	// if d1, d2, d3 := r.Color.DistanceRgb(color), r.Color.DistanceLab(color), r.Color.DistanceLuv(color); d1+d2+d3 != 0 {
	// 	logrus.Infof("%s\tdiffRGB=%.4f, diffLAB=%.4f, diffLUV=%.4f", r, d1, d2, d3)
	// }

	// check old format
	if r.Color.AlmostEqualRgb(color.Color{}) && !r.new {
		limit := 10.
		cr, cg, cb := colorXY.RGB255()
		_, _, cl := colorXY.Hsl()

		// fmt.Printf("%s:\t(%d,%d,%d) %.2f vs %.2f, %.2f vs %.2f\n", r, cr, cg, cb, (float64(cr)+float64(cg)+float64(cb))/3., r.Gray, cl*100, r.Lightness)

		grayOk := math.Abs((float64(cr)+float64(cg)+float64(cb))/3.-r.Gray) < limit
		lightnessOk := math.Abs(cl*100-r.Lightness) < limit

		if grayOk && lightnessOk {
			r.Color = colorXY
			r.colorUpdated = true
			logrus.Infof("convert [%s] to new color format", r)
			return true
		}
		return false
	} else {
		// fmt.Printf("%s  \t%.3f\n", r.ID, r.Color.DistanceLab(colorXY))
		return r.Color.DistanceLab(colorXY) < .05
	}
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
	robotgo.MoveMouse(150, 150)
}
