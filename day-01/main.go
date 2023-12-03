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
		fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	partOneSum := 0
	partTwoSum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		num, err := getNumFromLine(line, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unexpected error parsing line p1: %s, err: %s\n", line, err)
			os.Exit(1)
		}
		partOneSum += num

		num2, err := getNumFromLine(line, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unexpected error parsing line p2: %s, err: %s\n", line, err)
			os.Exit(1)
		}
		partTwoSum += num2
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1 answer: %d\n", partOneSum)
	fmt.Printf("Part 2 answer: %d\n", partTwoSum)
}

func getNumFromLine(line string, lookForWords bool) (int, error) {
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

	digits := make([]rune, 2)

	for i, r := range line {
		if unicode.IsDigit(r) {
			if digits[0] == 0 {
				digits[0] = r
			}
			digits[1] = r

			// if this char is a digit, no reason to look for words
			continue
		}

		// need at least three letters for a word number
		if !lookForWords || i < 2 {
			continue
		}

		// starting from the back, look for 3, 4, or 5 letter word numbers
		// if you find one, use it as the value
		for j := i - 2; j >= 0 && j >= i-4; j-- {
			if val, ok := numberMap[line[j:i+1]]; ok {
				if digits[0] == 0 {
					digits[0] = val
				}
				digits[1] = val
				break
			}
		}
	}

	return strconv.Atoi(string(digits))
}
