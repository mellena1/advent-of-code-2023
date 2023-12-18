package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mellena1/advent-of-code-2023/utils"
)

const (
	TRENCH  = '#'
	NOTHING = '.'
)

func main() {
	f := utils.ReadFile("input.txt")
	defer f.Close()

	steps := parseDigInput(f)
	grid := steps.getGrid()
	// fmt.Println(grid.polyCorners)
	// fmt.Println(grid)
	grid.FillIn()
	// fmt.Println(grid)
	fmt.Println(steps.Area())
}

type Polygon []utils.Coordinate

func (pg Polygon) PointIsInside(c utils.Coordinate) bool {
	// https://github.com/soniakeys/raycast/blob/master/raycast.go
	if len(pg) < 3 {
		return false
	}

	a := pg[0]
	in := rayIntersectsSegment(c, pg[len(pg)-1], a)
	for _, b := range pg[1:] {
		if rayIntersectsSegment(c, a, b) {
			in = !in
		}
		a = b
	}

	return in
}

type Grid struct {
	grid        [][]utils.Char
	polyCorners Polygon
}

func NewGrid(initGrid [][]utils.Char) Grid {
	return Grid{
		grid:        initGrid,
		polyCorners: []utils.Coordinate{},
	}
}

func (g *Grid) AddRow(idx int) {
	numCols := len(g.grid[0])
	newRow := make([]utils.Char, numCols)
	for i := range newRow {
		newRow[i] = NOTHING
	}

	if idx == -1 {
		g.grid = append([][]utils.Char{newRow}, g.grid...)

		for i, corner := range g.polyCorners {
			g.polyCorners[i] = corner.MoveDir(utils.DOWN)
		}
	} else {
		g.grid = append(g.grid, newRow)
	}
}

func (g *Grid) AddCol(idx int) {
	if idx == -1 {
		for i, row := range g.grid {
			g.grid[i] = append([]utils.Char{NOTHING}, row...)
		}

		for i, corner := range g.polyCorners {
			g.polyCorners[i] = corner.MoveDir(utils.RIGHT)
		}
	} else {
		for i, row := range g.grid {
			g.grid[i] = append(row, NOTHING)
		}
	}
}

func (g *Grid) Set(coor utils.Coordinate, v utils.Char) {
	g.grid[coor.Y][coor.X] = v
}

func (g *Grid) FillIn() {
	for i, row := range g.grid {
		for j, v := range row {
			if v == TRENCH {
				continue
			}
			coor := utils.NewCoordinate(j, i)
			if g.polyCorners.PointIsInside(coor) {
				g.Set(coor, TRENCH)
			}
		}
	}
}

func (g Grid) String() string {
	s := ""
	for _, row := range g.grid {
		for _, v := range row {
			s += string(v)
		}
		s += "\n"
	}
	return s
}

type DigStep struct {
	Dir      utils.Direction
	NumToDig int
	Color    string
}

type DigSteps []DigStep

func (steps DigSteps) getGrid() Grid {
	grid := NewGrid([][]utils.Char{{TRENCH}})

	coor := utils.NewCoordinate(0, 0)
	grid.polyCorners = append(grid.polyCorners, coor)

	for stepIdx, step := range steps {
		for i := 0; i < step.NumToDig; i++ {
			coor = coor.MoveDir(step.Dir)
			switch step.Dir {
			case utils.RIGHT:
				if coor.X >= len(grid.grid[coor.Y]) {
					grid.AddCol(coor.X)
				}
			case utils.LEFT:
				if coor.X < 0 {
					grid.AddCol(coor.X)
					// -1 is now 0
					coor = coor.MoveDir(utils.RIGHT)
				}
			case utils.DOWN:
				if coor.Y >= len(grid.grid) {
					grid.AddRow(coor.Y)
				}
			case utils.UP:
				if coor.Y < 0 {
					grid.AddRow(coor.Y)
					// -1 is now 0
					coor = coor.MoveDir(utils.DOWN)
				}
			}
			grid.Set(coor, TRENCH)
		}
		if stepIdx < len(steps)-1 {
			grid.polyCorners = append(grid.polyCorners, coor)
		}
	}

	return grid
}

func (s DigSteps) Area() int {
	grid := s.getGrid()
	grid.FillIn()

	sum := 0
	for _, row := range grid.grid {
		for _, v := range row {
			if v == TRENCH {
				sum++
			}
		}
	}

	return sum
}

func parseDigInput(r io.Reader) DigSteps {
	steps := DigSteps{}

	utils.ExecutePerLine(r, func(line string) error {
		lineSplit := strings.Split(line, " ")

		step := DigStep{}

		switch lineSplit[0] {
		case "R":
			step.Dir = utils.RIGHT
		case "L":
			step.Dir = utils.LEFT
		case "U":
			step.Dir = utils.UP
		case "D":
			step.Dir = utils.DOWN
		}

		step.NumToDig, _ = strconv.Atoi(lineSplit[1])
		step.Color = lineSplit[2][2:8]

		steps = append(steps, step)

		return nil
	})

	return steps
}

func rayIntersectsSegment(p, a, b utils.Coordinate) bool {
	return (a.Y > p.Y) != (b.Y > p.Y) &&
		p.X < (b.X-a.X)*(p.Y-a.Y)/(b.Y-a.Y)+a.X
}
