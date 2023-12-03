package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

const (
	nothing = -1
	// anything -2 or less must be a symbol
	symbol = -2
	gear   = -3
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	board, err := getBoard(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create board: %s\n", err)
		os.Exit(1)
	}

	partNums, err := findPartNumbers(board)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get part numbers: %s\n", err)
		os.Exit(1)
	}

	partOneSum := 0
	for _, partNum := range partNums {
		partOneSum += partNum
	}

	fmt.Printf("Part one solution: %d\n", partOneSum)

	gearRatios, err := findGearRatios(board)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get gear ratios: %s\n", err)
		os.Exit(1)
	}

	partTwoSum := 0
	for _, gearRatio := range gearRatios {
		partTwoSum += gearRatio
	}

	fmt.Printf("Part two solution: %d\n", partTwoSum)
}

func getBoard(f io.Reader) ([][]int, error) {
	board := [][]int{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		newLine := []int{}
		numBuilder := []rune{}

		addNumToLine := func() error {
			val, err := strconv.Atoi(string(numBuilder))
			if err != nil {
				return fmt.Errorf("built invalid num %q: %w", string(numBuilder), err)
			}

			for i := 0; i < len(numBuilder); i++ {
				newLine = append(newLine, val)
			}

			numBuilder = []rune{}

			return nil
		}

		for _, r := range line {
			if unicode.IsDigit(r) {
				numBuilder = append(numBuilder, r)
				continue
			}

			if len(numBuilder) > 0 {
				if err := addNumToLine(); err != nil {
					return nil, err
				}
			}

			switch r {
			case '.':
				newLine = append(newLine, nothing)
			case '*':
				newLine = append(newLine, gear)
			default:
				newLine = append(newLine, symbol)
			}
		}

		// empty out num if it exists at EOL
		if len(numBuilder) > 0 {
			if err := addNumToLine(); err != nil {
				return nil, err
			}
		}

		board = append(board, newLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return board, nil
}

func findPartNumbers(board [][]int) ([]int, error) {
	partNums := []int{}

	// since I store the value of each number in each of its positions, keep track of
	//		when I have found and added a number as a part number to avoid double counting
	found := foundMap{}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] > symbol {
				continue
			}

			for y := i - 1; y <= i+1; y++ {
				for x := j - 1; x <= j+1; x++ {
					// avoid going out of bounds
					if x < 0 || y < 0 || y >= len(board) || x >= len(board[y]) {
						continue
					}

					num := board[y][x]
					if num >= 0 && !found.isFound(y, x) {
						partNums = append(partNums, num)
						found.markFound(board, y, x)
					}
				}
			}
		}
	}

	return partNums, nil
}

func findGearRatios(board [][]int) ([]int, error) {
	gearRatios := []int{}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] != gear {
				continue
			}

			// keep found map per gear, since we could have 2 gears touching the same part num,
			//		but we don't want to recount the same part num per gear
			found := foundMap{}
			partNums := []int{}

			for y := i - 1; y <= i+1; y++ {
				for x := j - 1; x <= j+1; x++ {
					// avoid going out of bounds
					if x < 0 || y < 0 || y >= len(board) || x >= len(board[y]) {
						continue
					}

					num := board[y][x]
					if num >= 0 && !found.isFound(y, x) {
						partNums = append(partNums, num)

						found.markFound(board, y, x)
					}
				}
			}

			if len(partNums) == 2 {
				gearRatios = append(gearRatios, partNums[0]*partNums[1])
			}
		}
	}

	return gearRatios, nil
}

type foundMap map[int]map[int]any

func (found foundMap) isFound(y, x int) bool {
	if _, ok := found[y]; ok {
		if _, ok := found[y][x]; ok {
			return true
		}
	}
	return false
}

func (found foundMap) markFound(board [][]int, y, x int) {
	num := board[y][x]

	if _, ok := found[y]; !ok {
		found[y] = map[int]any{}
	}
	found[y][x] = nil

	for j := x - 1; j >= 0; j-- {
		if board[y][j] == num {
			found[y][j] = nil
		} else {
			break
		}
	}
	for j := x + 1; j < len(board[y]); j++ {
		if board[y][j] == num {
			found[y][j] = nil
		} else {
			break
		}
	}
}
