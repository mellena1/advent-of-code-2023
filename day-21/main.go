package main

import (
	"fmt"
	"io"

	"github.com/mellena1/advent-of-code-2023/utils"
)

const (
	START utils.Char = 'S'
	OPEN  utils.Char = '.'
	ROCK  utils.Char = '#'
)

func main() {
	f := utils.ReadFile("sample.txt")
	defer f.Close()

	grid := parseGrid(f)

	// fmt.Printf("Part one solution: %d\n", grid.AvailableSpotsFromSteps(64))
	fmt.Printf("Part two solution: %d\n", grid.AvailableSpotsFromSteps(10))
}

type Grid [][]utils.Char

func (g Grid) AvailableSpotsFromSteps(maxSteps int) int {
	cMap := g.toConnectionMap()
	distances, _ := cMap.Dijkstra(g.findStart())

	s := 0
	for _, d := range distances {
		if d >= 0 && d <= maxSteps && d%2 == 0 {
			s++
		}
	}
	return s
}

func (g Grid) findStart() utils.Coordinate {
	for i, row := range g {
		for j, v := range row {
			if v == START {
				return utils.NewCoordinate(j, i)
			}
		}
	}
	panic("no start")
}

func (g Grid) toConnectionMap() utils.ConnectionMap[utils.Coordinate] {
	cMap := utils.ConnectionMap[utils.Coordinate]{}

	for i, row := range g {
		for j, v := range row {
			if v == ROCK {
				continue
			}

			coor := utils.NewCoordinate(j, i)
			cMap[coor] = map[utils.Coordinate]int{}

			for k := -1; k <= 1; k++ {
				if i == 0 {
					continue
				}

				yCoor := coor.Y + k
				if yCoor >= 0 && yCoor < len(g) {
					if g[yCoor][coor.X] != ROCK {
						cMap[coor][utils.NewCoordinate(coor.X, yCoor)] = 1
					}
				}

				xCoor := coor.X + k
				if xCoor >= 0 && xCoor < len(g[coor.Y]) {
					if g[coor.Y][xCoor] != ROCK {
						cMap[coor][utils.NewCoordinate(xCoor, coor.Y)] = 1
					}
				}
			}
		}
	}

	return cMap
}

func parseGrid(r io.Reader) Grid {
	grid := Grid{}

	utils.ExecutePerLine(r, func(line string) error {
		grid = append(grid, []utils.Char(line))
		return nil
	})

	return grid
}
