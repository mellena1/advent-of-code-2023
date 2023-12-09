package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/mellena1/advent-of-code-2023/utils"
)

func main() {
	f := utils.ReadFile("input.txt")
	defer f.Close()

	reading := parseOasis(f)

	partOneSum := 0
	for _, r := range reading {
		partOneSum += r.GetNextNumber()
	}
	fmt.Printf("Part one solution: %d\n", partOneSum)

	partTwoSum := 0
	for _, r := range reading {
		partTwoSum += r.GetPrevNumber()
	}
	fmt.Printf("Part two solution: %d\n", partTwoSum)
}

type OasisReading []int

func (o OasisReading) GetNextNumber() int {
	return o.nevilleInterpolation(len(o))
}

func (o OasisReading) GetPrevNumber() int {
	return o.nevilleInterpolation(-1)
}

func (o OasisReading) nevilleInterpolation(x int) int {
	n := len(o)
	p := make([]float64, n)

	for k := 0; k < n; k++ {
		for i := 0; i < n-k; i++ {
			if k == 0 {
				p[i] = float64(o[i])
			} else {
				p[i] = (float64(x-i-k)*p[i] + float64(i-x)*p[i+1]) / float64(-k)
			}
		}
	}

	return int(p[0])
}

func parseOasis(r io.Reader) []OasisReading {
	readings := []OasisReading{}

	utils.ExecutePerLine(r, func(line string) error {
		reading, err := utils.StrSliceToIntSlice(strings.Fields(line))
		if err != nil {
			return fmt.Errorf("failed to parse line %q: %w", line, err)
		}

		readings = append(readings, reading)

		return nil
	})

	return readings
}
