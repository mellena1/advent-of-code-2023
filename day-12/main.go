package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/mellena1/advent-of-code-2023/utils"
)

type SpringState utils.Char

func (s SpringState) String() string {
	return utils.Char(s).String()
}

const (
	WORKING  SpringState = '.'
	BROKEN   SpringState = '#'
	QUESTION SpringState = '?'
)

func main() {
	f := utils.ReadFile("sample.txt")
	defer f.Close()

	lines := parseLines(f)
	partOne := 0
	for _, l := range lines {
		partOne += bruteForceCombos(l.springs, l.groups)
	}
	fmt.Printf("Part one solution: %d\n", partOne)

	unfoldedLines := utils.SliceMap(lines, func(v LineOfSprings) LineOfSprings {
		return v.Unfold()
	})
	partTwo := 0
	for _, l := range unfoldedLines {
		partTwo += bruteForceCombos(l.springs, l.groups)
	}
	fmt.Printf("Part two solution: %d\n", partTwo)
}

type LineOfSprings struct {
	springs []SpringState
	groups  []int
}

func (l LineOfSprings) PossibleArrangements() int {
	// var splitAndCount func(springs []SpringState, groups []int) int
	// splitAndCount = func(springs []SpringState, groups []int) int {
	// 	totalArrangements := 1

	// 	subGroups := strings.Split(string(springs), ".")
	// 	subGroups = utils.SliceFilter(subGroups, func(v string) bool { return v != "" })

	// 	removedOne := true
	// 	for removedOne {
	// 		if len(subGroups) == 0 || len(groups) == 0 {
	// 			fmt.Println(l.springs, l.groups, subGroups, groups)
	// 		}

	// 		removedOne = false
	// 		if len(subGroups[0]) == groups[0] {
	// 			subGroups = subGroups[1:]
	// 			groups = groups[1:]
	// 			removedOne = true
	// 		}

	// 		if len(subGroups) == 0 {
	// 			break
	// 		}
	// 		if len(groups) == 0 {
	// 			break
	// 			// fmt.Println(l.springs, subGroups, groups)
	// 		}

	// 		if len(subGroups[len(subGroups)-1]) == groups[len(groups)-1] {
	// 			subGroups = subGroups[:len(subGroups)-1]
	// 			groups = groups[:len(groups)-1]
	// 			removedOne = true
	// 		}
	// 	}

	// 	if len(subGroups) == 0 {
	// 		return 1
	// 	}

	// 	if len(subGroups) == len(groups) {
	// 		for i, subGroup := range subGroups {
	// 			totalArrangements *= numberOfCombos(groups[i], []SpringState(subGroup))
	// 		}
	// 		return totalArrangements
	// 	}

	// 	return bruteForceCombos([]SpringState(strings.Join(subGroups, ".")), groups)

	// 	// for gIdx, group := range groups {
	// 	// 	newGroup := []SpringState{}

	// 	// 	if gIdx == len(groups)-1 {
	// 	// 		newGroup = []SpringState(subGroups[0])
	// 	// 	} else {
	// 	// 		for _, r := range subGroups[0] {
	// 	// 			spring := SpringState(r)
	// 	// 			newGroup = append(newGroup, spring)
	// 	// 			if spring == BROKEN {
	// 	// 				continue
	// 	// 			}
	// 	// 			if len(newGroup) >= group+1 {
	// 	// 				break
	// 	// 			}
	// 	// 		}
	// 	// 		// need an extra spot for the dot to separate
	// 	// 		newGroup = newGroup[:len(newGroup)-1]
	// 	// 	}

	// 	// 	totalArrangements *= splitAndCount(newGroup, []int{group})
	// 	// 	if len(newGroup) == len(subGroups[0]) {
	// 	// 		subGroups = subGroups[1:]
	// 	// 	} else {
	// 	// 		subGroups[0] = subGroups[0][len(newGroup)+1:]
	// 	// 	}
	// 	// }

	// 	// if totalArrangements == 0 {
	// 	// 	return 1
	// 	// }

	// 	return totalArrangements
	// }

	return 0
}

func (l LineOfSprings) Unfold() LineOfSprings {
	unfolded := LineOfSprings{
		springs: []SpringState{},
		groups:  []int{},
	}
	for i := 0; i < 5; i++ {
		unfolded.springs = append(unfolded.springs, l.springs...)
		if i < 4 {
			unfolded.springs = append(unfolded.springs, QUESTION)
		}
		unfolded.groups = append(unfolded.groups, l.groups...)
	}
	return unfolded
}

func (l *LineOfSprings) DedupDots() {
	l.springs = dedupDots(l.springs)
}

func (l LineOfSprings) String() string {
	s := ""
	for _, spring := range l.springs {
		s += string(spring)
	}
	s += " ("
	for i, group := range l.groups {
		s += strconv.Itoa(group)
		if i < len(l.groups)-1 {
			s += ", "
		}
	}
	return s + ")"
}

var dotsRegex = regexp.MustCompile(`\.{2,}`)

func dedupDots(s []SpringState) []SpringState {
	newSprings := dotsRegex.ReplaceAllString(string(s), ".")
	return []SpringState(newSprings)
}

func numberOfCombos(needed int, springs []SpringState) int {
	start := -1
	end := len(springs)

	for i, spring := range springs {
		if spring == BROKEN {
			if i > start {
				start = i
			}
			if i < end {
				end = i
			}
		}
	}

	// no BROKEN springs found
	if start == -1 {
		return len(springs) - needed + 1
	}

	numOfBROKEN := end - start + 1

	// shrink springs to only possible positions
	highestPossibleEnd := end + (needed - numOfBROKEN)
	if highestPossibleEnd < len(springs) {
		springs = springs[:highestPossibleEnd+1]
	}
	lowestPossibleStart := start - (needed - numOfBROKEN)
	if lowestPossibleStart >= 0 {
		springs = springs[lowestPossibleStart:]
	}

	return len(springs) - needed + 1
}

func springCopyWithChange(springs []SpringState, i int, newVal SpringState) []SpringState {
	newSprings := make([]SpringState, len(springs))
	copy(newSprings, springs)
	newSprings[i] = newVal
	return newSprings
}

func bruteForceCombos(springs []SpringState, groups []int) int {
	total := 0

	questionIdx := strings.Index(string(springs), "?")
	if questionIdx > -1 {
		total += bruteForceCombos(springCopyWithChange(springs, questionIdx, WORKING), groups)
		total += bruteForceCombos(springCopyWithChange(springs, questionIdx, BROKEN), groups)
	} else {
		springsSplit := strings.Split(string(springs), ".")
		springsSplit = utils.SliceFilter(springsSplit, func(v string) bool { return v != "" })
		if len(groups) != len(springsSplit) {
			return 0
		}
		for i, group := range groups {
			if len(springsSplit[i]) != group {
				return 0
			}
		}
		return 1
	}

	return total
}

func parseLines(r io.Reader) []LineOfSprings {
	lines := []LineOfSprings{}
	utils.ExecutePerLine(r, func(line string) error {
		springs, groups, _ := strings.Cut(line, " ")

		groupsInts, err := utils.StrSliceToIntSlice(strings.FieldsFunc(groups, func(r rune) bool {
			return r == ','
		}))
		if err != nil {
			return fmt.Errorf("failed to parse groups %q: %w", groups, err)
		}

		newLineOfSprings := LineOfSprings{
			springs: []SpringState(springs),
			groups:  groupsInts,
		}
		newLineOfSprings.DedupDots()

		lines = append(lines, newLineOfSprings)

		return nil
	})
	return lines
}
