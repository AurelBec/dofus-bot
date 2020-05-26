package models

import "math"

type Pos struct {
	X int `json:"x"`
	Y int `json:"Y"`
}

func (p Pos) IsNull() bool {
	return p.X == 0 && p.Y == 0
}

func (lhs Pos) DistanceTo(rhs Pos) int {
	dx := float64(lhs.X - rhs.X)
	dy := float64(lhs.Y - rhs.Y)
	return int(math.Sqrt(dx*dx + dy*dy))
}
