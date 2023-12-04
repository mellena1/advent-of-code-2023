package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	games, err := getGamesFromText(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get games: %s", err)
		os.Exit(1)
	}

	partOneSum := 0

	for _, game := range games {
		partOneSum += game.getScore()
	}
	fmt.Printf("Part one sum: %d\n", partOneSum)

	partTwoSum := calcNumberOfCardsWithCopies(games)
	fmt.Printf("Part two sum: %d\n", partTwoSum)
}

func getGamesFromText(f io.Reader) ([]Game, error) {
	scanner := bufio.NewScanner(f)

	games := []Game{}

	for scanner.Scan() {
		line := scanner.Text()

		game, err := getGameFromLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to get game from line %q: %s", line, err)
		}

		games = append(games, game)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return games, nil
}

func getGameFromLine(line string) (Game, error) {
	splitOnColon := strings.Split(line, ":")
	gameID, err := strconv.Atoi(strings.Fields(splitOnColon[0])[1])
	if err != nil {
		return Game{}, fmt.Errorf("invalid game id: %w", err)
	}

	splitOnBar := strings.Split(splitOnColon[1], "|")
	winningNumsStrs := strings.Fields(splitOnBar[0])
	playerNumsStrs := strings.Fields(splitOnBar[1])

	winningNums, err := strSliceToIntSlice(winningNumsStrs)
	if err != nil {
		return Game{}, fmt.Errorf("invalid winning nums: %w", err)
	}

	playerNums, err := strSliceToIntSlice(playerNumsStrs)
	if err != nil {
		return Game{}, fmt.Errorf("invalid player nums: %w", err)
	}

	return Game{
		ID:             gameID,
		WinningNumbers: winningNums,
		Numbers:        playerNums,
		score:          -1,
		matches:        -1,
	}, nil
}

func calcNumberOfCardsWithCopies(games []Game) int {
	// memoize how many cards each card gives you in total
	cardsGiven := make([]int, len(games))

	// recursively determine how many cards you end up with
	//	starting from card idx
	var calc func(idx int) int
	calc = func(idx int) int {
		if cardsGiven[idx] > 0 {
			return cardsGiven[idx]
		}

		matches := games[idx].getMatchingNums()

		numCards := 1

		for i := idx + 1; i < len(games) && i <= idx+matches; i++ {
			numCards += calc(i)
		}

		cardsGiven[idx] = numCards

		return numCards
	}

	numCards := 0
	for i := range games {
		numCards += calc(i)
	}

	return numCards
}

type Game struct {
	ID             int
	WinningNumbers []int
	Numbers        []int
	score          int
	matches        int
}

func (g *Game) getScore() int {
	// memoize score so we don't recalc every call
	if g.score >= 0 {
		return g.score
	}

	score := 0

	matches := g.getMatchingNums()
	if matches > 0 {
		// 2 to the power of (matches - 1)
		// if one match, 1, if two matches, 2, etc
		score = 1 << (matches - 1)
	}

	g.score = score
	return score
}

func (g *Game) getMatchingNums() int {
	// memoize matches so we don't recalc every call
	if g.matches >= 0 {
		return g.matches
	}

	matches := 0

	for _, n := range g.WinningNumbers {
		if slices.Contains(g.Numbers, n) {
			matches += 1
		}
	}

	g.matches = matches
	return matches
}

func strSliceToIntSlice(strs []string) ([]int, error) {
	nums := make([]int, len(strs))
	for i, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		nums[i] = n
	}
	return nums, nil
}
