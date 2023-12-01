package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %s", err))
	}
	defer f.Close()

	partOneSum := 0
	partTwoSum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		num, err := getNumFromLinePartOne(line)
		if err != nil {
			panic(fmt.Sprintf("unexpected error parsing line p1: %s, err: %s", line, err))
		}
		partOneSum += num

		num2, err := getNumFromLinePartTwo(line)
		if err != nil {
			panic(fmt.Sprintf("unexpected error parsing line p2: %s, err: %s", line, err))
		}
		partTwoSum += num2
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error reading file: %s", err))
	}

	fmt.Printf("Part 1 answer: %d\n", partOneSum)
	fmt.Printf("Part 2 answer: %d\n", partTwoSum)
}

func getNumFromLinePartOne(line string) (int, error) {
	var first rune
	var last rune

	for _, r := range line {
		if unicode.IsDigit(r) {
			if first == 0 {
				first = r
			}
			last = r
		}
	}

	return strconv.Atoi(string([]rune{first, last}))
}

func getNumFromLinePartTwo(line string) (int, error) {
	numberMap := map[string]rune{
		"one":   '1',
		"two":   '2',
		"three": '3',
		"four":  '4',
		"five":  '5',
		"six":   '6',
		"seven": '7',
		"eight": '8',
		"nine":  '9',
	}

	var first rune
	var last rune

	for i, r := range line {
		if unicode.IsDigit(r) {
			if first == 0 {
				first = r
			}
			last = r

			// if this char is a digit, no reason to look for words
			continue
		}

		// need at least three letters for a word number
		if i < 2 {
			continue
		}

		curWord := ""
		if i >= 4 {
			curWord = line[i-4 : i+1]
		} else {
			curWord = line[0 : i+1]
		}

		// starting from the back, look for 3, 4, or 5 letter numbers
		// if you find one, use it as the value
		for j := len(curWord) - 1; j >= 0; j-- {
			if val, ok := numberMap[curWord[j:]]; ok {
				if first == 0 {
					first = val
				}
				last = val
				break
			}
		}
	}

	return strconv.Atoi(string([]rune{first, last}))
}
