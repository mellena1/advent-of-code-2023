package main

import (
	"fmt"
	"io"

	"github.com/mellena1/advent-of-code-2023/utils"
)

func main() {
	f := utils.ReadFile("sample.txt")
	defer f.Close()

	grid := parseGrid(f)
	fmt.Println(grid)

	fmt.Printf("Part one solution: %d\n", 0)
}

type Grid [][]int

func parseGrid(r io.Reader) Grid {
	grid := Grid{}

	utils.ExecutePerLine(r, func(line string) error {
		intSlice, err := utils.StrSliceToIntSlice(utils.SliceMap([]rune(line), func(v rune) string { return string(v) }))
		if err != nil {
			return fmt.Errorf("failed to parse line as int slice %q: %w", line, err)
		}

		grid = append(grid, intSlice)
		return nil
	})

	return grid
}
