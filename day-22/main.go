package main

import (
	"io"
	"slices"
	"strings"

	"github.com/mellena1/advent-of-code-2023/utils"
)

type Axis utils.Char

const (
	XAxis Axis = 'X'
	YAxis Axis = 'Y'
	ZAxis Axis = 'Z'
)

func main() {
	f := utils.ReadFile("input.txt")
	defer f.Close()

	bricks := parseBricks(f)
	bricks.MoveAllDown()
}

type Bricks []Brick

func (b Bricks) MoveAllDown() {
	slices.SortFunc(b, func(a, b Brick) int {
		return a.Start.Z - b.Start.Z
	})

	for i, brick := range b {
		for {
			brick = brick.Add(ZAxis, -1)

			if brick.Start.Z < 1 {
				break
			}

			noIntersections := true
			brick.ForEachCoor(func(c utils.Coordinate3D) bool {
				for j, otherBrick := range b {
					if i == j {
						continue
					}

					if otherBrick.Intersects(c) {
						noIntersections = false
						return true
					}
				}

				return false
			})

			if !noIntersections {
				break
			}

			b[i] = brick
		}
	}
}

type Brick struct {
	Start utils.Coordinate3D
	End   utils.Coordinate3D
	Axis  Axis
}

func (b Brick) ForEachCoor(f func(c utils.Coordinate3D) bool) {
	switch b.Axis {
	case XAxis:
		for x := b.Start.X; x <= b.End.X; x++ {
			shouldBreak := f(utils.NewCoordinate3D(x, b.Start.Y, b.Start.Z))
			if shouldBreak {
				return
			}
		}
	case YAxis:
		for y := b.Start.Y; y <= b.End.Y; y++ {
			shouldBreak := f(utils.NewCoordinate3D(b.Start.X, y, b.Start.Z))
			if shouldBreak {
				return
			}
		}
	case ZAxis:
		for z := b.Start.Z; z <= b.End.Z; z++ {
			shouldBreak := f(utils.NewCoordinate3D(b.Start.X, b.Start.Y, z))
			if shouldBreak {
				return
			}
		}
	}
}

func (b Brick) Add(axis Axis, num int) Brick {
	newBrick := Brick{
		Axis: b.Axis,
	}

	switch axis {
	case XAxis:
		newBrick.Start = b.Start.Translate(num, 0, 0)
		newBrick.End = b.End.Translate(num, 0, 0)
	case YAxis:
		newBrick.Start = b.Start.Translate(0, num, 0)
		newBrick.End = b.End.Translate(0, num, 0)
	case ZAxis:
		newBrick.Start = b.Start.Translate(0, 0, num)
		newBrick.End = b.End.Translate(0, 0, num)
	}

	return newBrick
}

func (b Brick) Intersects(c utils.Coordinate3D) bool {
	switch b.Axis {
	case XAxis:
		return c.Y == b.Start.Y && c.Z == b.Start.Z && c.X >= b.Start.X && c.X <= b.End.X
	case YAxis:
		return c.X == b.Start.X && c.Z == b.Start.Z && c.Y >= b.Start.Y && c.Y <= b.End.Y
	case ZAxis:
		return c.X == b.Start.X && c.Y == b.Start.Y && c.Z >= b.Start.Z && c.Z <= b.End.Z
	}
	return false
}

func parseBricks(r io.Reader) Bricks {
	bricks := Bricks{}

	utils.ExecutePerLine(r, func(line string) error {
		startCoor, endCoor, _ := strings.Cut(line, "~")

		start, err := strCoorTo3DCoor(startCoor)
		if err != nil {
			return err
		}
		end, err := strCoorTo3DCoor(endCoor)
		if err != nil {
			return err
		}

		newBrick := Brick{
			Start: start,
			End:   end,
		}

		if start.X != end.X {
			newBrick.Axis = XAxis
		}
		if start.Y != end.Y {
			newBrick.Axis = YAxis
		}
		if start.Z != end.Z {
			newBrick.Axis = ZAxis
		}

		bricks = append(bricks, newBrick)

		return nil
	})

	return bricks
}

func strCoorTo3DCoor(s string) (utils.Coordinate3D, error) {
	nums, err := utils.StrSliceToIntSlice(strings.Split(s, ","))
	if err != nil {
		return utils.Coordinate3D{}, err
	}

	return utils.NewCoordinate3D(nums[0], nums[1], nums[2]), nil
}
