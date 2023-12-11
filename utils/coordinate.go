package utils

import (
	"fmt"
)

type Coordinate struct {
	X int
	Y int
}

func NewCoordinate(x int, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func (c Coordinate) StepsToCoordinate(c2 Coordinate) int {
	xDiff := c.X - c2.X
	if xDiff < 0 {
		xDiff *= -1
	}

	yDiff := c.Y - c2.Y
	if yDiff < 0 {
		yDiff *= -1
	}

	return xDiff + yDiff
}
