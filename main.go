package main

import (
	"fmt"
	"io"
	"os"
)

func fatal(v ...any) {
	fmt.Println(v...)
	os.Exit(1)
}

func main() {
	var operandString1, operandString2 string
	var operatorRune int

	for {
		fmt.Println("Input:")
		_, err := fmt.Scanf("%s %c %s\n", &operandString1, &operatorRune, &operandString2)
		if err != nil {
			if err == io.EOF {
				break
			}
			fatal(fmt.Errorf("bad input: %w", err))
		}

		fmt.Printf("Output:\n%#v %#v %#v\n", operandString1, operatorRune, operandString2)
	}

	fmt.Println("Exit")
}

var RomanNumerals = map[rune]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

func romanToInt(s string) (int, error) {
	{
		sum := 0
		greatest := 0 // determens if number needs to be subtracted

		// scanning right to left
		for i := len(s) - 1; i >= 0; i-- {
			letter := s[i]

			num, exists := RomanNumerals[rune(letter)]
			if !exists {
				return 0, fmt.Errorf("%c is not a roman number", letter)
			}

			// case for for I in IV, I in IX, X in XL and so on
			if num < greatest {
				sum = sum - num
				continue
			}

			greatest = num
			sum = sum + num
		}

		return sum, nil
	}
}

type operand struct {
	value int
	roman bool
}

func newOperand(s string) (*operand, error) {
	v, err := romanToInt(s)
	if err == nil {
		return &operand{v, true}, nil
	}

	v, err = strconv.Atoi(s)
	if err == nil {
		return &operand{v, true}, nil
	}

	return nil, fmt.Errorf("not an arabic or roman number")
}
