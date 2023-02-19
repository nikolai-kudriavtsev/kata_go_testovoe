package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func fatal(e error) {
	fatal := fmt.Errorf("fatal error: %w", e)
	fmt.Fprintln(os.Stderr, fatal)
	os.Exit(1)
}

func main() {
	c := newCalculator(standardOperations)
	err := c.REPL(os.Stdin, os.Stdout)
	if err != nil {
		fatal(err)
	}
}

type operation func(int, int) int
type operationTable map[rune]operation

var standardOperations = operationTable{
	'+': func(x, y int) int { return x + y },
	'-': func(x, y int) int { return x - y },
	'*': func(x, y int) int { return x * y },
	'/': func(x, y int) int { return x / y },
}

type calculator struct {
	operations operationTable
}

func newCalculator(ot operationTable) *calculator {
	return &calculator{ot}
}

func (c *calculator) newExpression(operator rune, x, y *operand) (*expression, error) {
	op, exists := c.operations[operator]
	if !exists {
		return nil, errors.New("no such operator")
	}

	if x.roman != y.roman {
		return nil, errors.New("operands from different numeric systems")
	}

	return &expression{&op, x, y}, nil
}

func (c *calculator) REPL(input io.Reader, output io.Writer) error {
	var operandXInput, operandYInput string
	var operatorInput rune

	for {
		fmt.Println("input:")
		_, err := fmt.Fscanf(input, "%s %c %s\n", &operandXInput, &operatorInput, &operandYInput)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("bad input: %w", err)
		}

		oX, err := newOperand(operandXInput)
		if err != nil {
			return fmt.Errorf("bad operand: %w", err)
		}

		oY, err := newOperand(operandYInput)
		if err != nil {
			return fmt.Errorf("bad operand: %w", err)
		}

		exp, err := c.newExpression(operatorInput, oX, oY)
		if err != nil {
			return fmt.Errorf("bad expression: %w", err)
		}

		r, err := exp.eval()
		if err != nil {
			return fmt.Errorf("bad evaluation: %w", err)
		}

		var result string
		if exp.isRoman() {
			result = intToRoman(r)
		} else {
			result = strconv.Itoa(r)
		}

		fmt.Fprintf(output, "output:\n%s\n", result)
	}

	fmt.Println("exit")

	return nil
}

type expression struct {
	op   *operation
	x, y *operand
}

func (e *expression) isRoman() bool {
	return e.x.roman
}

func (e *expression) eval() (int, error) {
	result := (*e.op)(e.x.value, e.y.value)

	if e.isRoman() && result < 1 {
		return result, fmt.Errorf("result of operation %d cannot be expressed by roman letters", result)
	}

	return result, nil
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

	return nil, errors.New("not an arabic or roman integer number")
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
	greatest := 0

	for i := len(s) - 1; i >= 0; i-- {
		letter := s[i]

		num, exists := RomanNumerals[rune(letter)]
		if !exists {
			return 0, fmt.Errorf("%c is not a roman number", letter)
		}

		if num < greatest {
			sum = sum - num
			continue
		}

		greatest = num
		sum = sum + num
	}

	return sum, nil
}

func intToRoman(number int) string {
	var roman strings.Builder

	intToRomanTable := []struct {
		value int
		digit string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	for _, row := range intToRomanTable {
		for number >= row.value {
			roman.WriteString(row.digit)
			number -= row.value
		}
	}

	return roman.String()
}
