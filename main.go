package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func fatal(v ...any) {
	fmt.Println(v...)
	os.Exit(1)
}

func main() {
	var operandXInput, operandYInput string
	var operatorInput rune

	for {
		fmt.Println("input:")
		_, err := fmt.Scanf("%s %c %s\n", &operandXInput, &operatorInput, &operandYInput)
		if err != nil {
			if err == io.EOF {
				break
			}
			fatal(fmt.Errorf("bad input: %w", err))
		}

		oX, err := newOperand(operandXInput)
		if err != nil {
			fatal(fmt.Errorf("bad operand: %w", err))
		}

		oY, err := newOperand(operandYInput)
		if err != nil {
			fatal(fmt.Errorf("bad operand: %w", err))
		}

		exp, err := newExpression(operatorInput, oX, oY)
		if err != nil {
			fatal(fmt.Errorf("bad expression: %w", err))
		}

		result, err := exp.eval()
		if err != nil {
			fatal(fmt.Errorf("bad evaluation: %w", err))
		}

		fmt.Printf("output:\n%d\n", result)
	}

	fmt.Println("exit")
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
		return &operand{v, false}, nil
	}

	return nil, errors.New("not an arabic or roman number")
}

var operations = map[rune]func(int, int) int{
	'+': func(x, y int) int { return x + y },
	'-': func(x, y int) int { return x - y },
	'*': func(x, y int) int { return x * y },
	'/': func(x, y int) int { return x / y },
}

type expression struct {
	operator rune
	x, y     *operand
}

func newExpression(operator rune, x, y *operand) (*expression, error) {
	_, exists := operations[operator]
	if !exists {
		return nil, errors.New("no such operator")
	}

	if x.roman != y.roman {
		return nil, errors.New("operands from different numeric systems")
	}

	return &expression{operator, x, y}, nil
}

func (e *expression) isRoman() bool {
	return e.x.roman
}

func (e *expression) eval() (int, error) {
	op := operations[e.operator]
	result := op(e.x.value, e.y.value)

	if e.isRoman() && result < 1 {
		return result, fmt.Errorf("result of operation %d cannot be expressed by roman letters", result)
	}

	return result, nil
}
