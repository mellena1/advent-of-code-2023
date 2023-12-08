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

	directions, maps := parseMaps(f)

	fmt.Printf("Part one solution: %d\n", maps.stepsToZZZ(directions))
	fmt.Printf("Part two solution: %d\n", maps.stepsToAllZs(directions))
}

type Node struct {
	Left  string
	Right string
}

type Maps map[string]Node

func (m Maps) stepsToZZZ(directions string) int {
	steps := 0

	curNode := "AAA"
	dirIdx := 0

	for {
		if dirIdx >= len(directions) {
			dirIdx = 0
		}
		dir := rune(directions[dirIdx])

		curNode = m.getNextNode(curNode, dir)

		steps++
		dirIdx++

		if curNode == "ZZZ" {
			break
		}
	}

	return steps
}

func (m Maps) stepsToAllZs(directions string) int {
	steps := 0

	// start from all nodes that end in A
	curNodes := []string{}
	for nodeName := range m {
		if nodeName[2] == 'A' {
			curNodes = append(curNodes, nodeName)
		}
	}
	eachNodeSteps := make([]int, len(curNodes))

	dirIdx := 0
	nodesDone := 0

	for nodesDone < len(eachNodeSteps) {
		if dirIdx >= len(directions) {
			dirIdx = 0
		}
		dir := rune(directions[dirIdx])

		steps++
		dirIdx++

		for i := range curNodes {
			// don't care about nodes that are already done
			if eachNodeSteps[i] > 0 {
				continue
			}

			curNodes[i] = m.getNextNode(curNodes[i], dir)
			if curNodes[i][2] == 'Z' {
				eachNodeSteps[i] = steps
				nodesDone++
			}
		}
	}

	return leastCommonMultiple(eachNodeSteps)
}

func (m Maps) getNextNode(curNode string, dir rune) string {
	switch dir {
	case 'L':
		return m[curNode].Left
	case 'R':
		return m[curNode].Right
	}
	panic("unknown direction: " + string(dir))
}

func leastCommonMultiple(nums []int) int {
	dedupFactors := map[int]any{}

	for _, n := range nums {
		factors := primeFactorization(n)
		for _, f := range factors {
			dedupFactors[f] = nil
		}
	}

	lcm := 1
	for f := range dedupFactors {
		lcm *= f
	}

	return lcm
}

func primeFactorization(num int) []int {
	primeFactors := []int{}

	// get all 2s
	for num%2 == 0 {
		primeFactors = append(primeFactors, 2)
		num /= 2
	}

	// go through odd nums
	for i := 3; i*i <= num; i += 2 {
		for num%i == 0 {
			primeFactors = append(primeFactors, i)
			num /= i
		}
	}

	// whatever is left must still be prime
	if num > 2 {
		primeFactors = append(primeFactors, num)
	}

	return primeFactors
}

func parseMaps(r io.Reader) (string, Maps) {
	directions := ""
	maps := Maps{}

	utils.ExecutePerLine(r, func(line string) error {
		if len(line) == 0 {
			return nil
		}

		if strings.Contains(line, "=") {
			nodeName, lr, _ := strings.Cut(line, " = ")

			maps[nodeName] = Node{
				Left:  lr[1:4],
				Right: lr[6:9],
			}

			return nil
		}

		directions = line

		return nil
	})

	return directions, maps
}
