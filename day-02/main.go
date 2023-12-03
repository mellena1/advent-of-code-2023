package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type cubeAmounts map[string]int

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	partOneSum := 0
	partTwoSum := 0

	allowedAmts := cubeAmounts{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		gameNum, err := gameIsPossible(allowedAmts, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to check game %q: %s", line, err)
			os.Exit(1)
		}
		partOneSum += gameNum

		gamePower, err := calcGamePower(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to calc game power %q: %s", line, err)
			os.Exit(1)
		}
		partTwoSum += gamePower
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1 answer: %d\n", partOneSum)
	fmt.Printf("Part 2 answer: %d\n", partTwoSum)
}

func parseGame(line string) (int, []cubeAmounts, error) {
	splitOnColon := strings.Split(line, ":")
	gameId, err := strconv.Atoi(strings.Split(splitOnColon[0], " ")[1])
	if err != nil {
		return -1, nil, fmt.Errorf("invalid game id %q: %w", gameId, err)
	}

	parsedPulls := []cubeAmounts{}

	gameDetails := strings.TrimSpace(splitOnColon[1])
	pulls := strings.Split(gameDetails, ";")
	for i := range pulls {
		pulls[i] = strings.TrimSpace(pulls[i])

		parsedPull, err := parsePull(pulls[i])
		if err != nil {
			return -1, nil, err
		}

		parsedPulls = append(parsedPulls, parsedPull)
	}

	return gameId, parsedPulls, nil
}

func gameIsPossible(allowed cubeAmounts, line string) (int, error) {
	gameId, pulls, err := parseGame(line)
	if err != nil {
		return -1, err
	}

	for _, pull := range pulls {
		if !pullIsPossible(allowed, pull) {
			return 0, nil
		}
	}

	return gameId, nil
}

func calcGamePower(line string) (int, error) {
	_, pulls, err := parseGame(line)
	if err != nil {
		return -1, err
	}

	minsNeeded := cubeAmounts{}

	for _, pull := range pulls {
		for color, num := range pull {
			if curNeeded, ok := minsNeeded[color]; ok {
				if num > curNeeded {
					minsNeeded[color] = num
				}
			} else {
				minsNeeded[color] = num
			}
		}
	}

	power := 1
	for _, needed := range minsNeeded {
		power *= needed
	}

	return power, nil
}

func pullIsPossible(allowed cubeAmounts, pull cubeAmounts) bool {
	for color, num := range pull {
		if num > allowed[color] {
			return false
		}
	}

	return true
}

var pullRegex = regexp.MustCompile(`(\d+) (blue|red|green)`)

func parsePull(pull string) (cubeAmounts, error) {
	cubes := strings.Split(pull, ",")

	cubeAmts := cubeAmounts{}

	for _, cubePull := range cubes {
		matches := pullRegex.FindStringSubmatch(cubePull)

		num, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("not a number %q: %w", matches[1], err)
		}
		color := matches[2]

		cubeAmts[color] = num
	}

	return cubeAmts, nil
}
