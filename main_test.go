package main

import (
	"fmt"
	"os"
	"testing"
)

type testCase struct {
	op           rune
	x, y, result string
	err          string
}

func (t *testCase) Name() string {
	return fmt.Sprintf("input: %s %c %s, output: %#v, error: %s", t.x, t.op, t.y, t.result, t.err)
}

func TestCalculator(t *testing.T) {
	var tests = []testCase{
		// op, x, y, result, err
		// basic
		{'+', "2", "2", "4", ""},
		{'+', "II", "II", "IV", ""},
		{'-', "2", "1", "1", ""},
		{'-', "II", "I", "I", ""},
		{'*', "2", "2", "4", ""},
		{'*', "II", "II", "IV", ""},
		{'/', "2", "2", "1", ""},
		{'/', "II", "II", "I", ""},
		// fail on unsupported operator
		{'%', "2", "2", "", "bad expression: % is not a supported operator"},
		// fail on input not in range of 1 to 10
		{'+', "0", "2", "", "bad operand: 0 not in range of possible values from 1 to 10"},
		{'+', "11", "2", "", "bad operand: 11 not in range of possible values from 1 to 10"},
		{'+', "XI", "II", "", "bad operand: XI not in range of possible values from I to X"},
		// big output
		{'*', "10", "10", "100", ""},
		{'*', "X", "X", "C", ""},
		// fail on non integer value
		{'+', "2.5", "2", "", "bad operand: 2.5 is not an arabic or roman integer number"},
		// fail on mixed numeric systems
		{'+', "2", "II", "", "bad expression: operands from different numeric systems"},
		// bad input format or expression
		{'+', "2 + 2", "2", "", "bad input: newline in format does not match input"},
		// result >= I for roman
		{'-', "I", "I", "", "bad evaluation: result of operation 0 cannot be expressed by roman letters"},
	}

	calc := newCalculator(standardOperations)

	for _, tt := range tests {
		testname := tt.Name()
		t.Run(testname, func(t *testing.T) {
			inputReader, inputWriter, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			outputReader, outputWriter, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			inputWriter.WriteString(fmt.Sprintf("%s %c %s\n", tt.x, tt.op, tt.y))
			if err != nil {
				t.Fatal(err)
			}
			inputWriter.Close()
			if err != nil {
				t.Fatal(err)
			}

			var result string
			calcErr := calc.REPL(inputReader, outputWriter)
			if calcErr == nil {
				fmt.Fscanf(outputReader, "input:\noutput:\n%s\n", &result)
			}

			var fail bool
			var errRepr, explainer string
			if tt.result != result {
				fail = true
				explainer = fmt.Sprintf("got result %#v, want %#v; ", result, tt.result)
			}
			if calcErr != nil {
				errRepr = calcErr.Error()
			}
			if tt.err != errRepr {
				if fail {
					explainer = fmt.Sprintf("%s; got error %#v want %#v", explainer, errRepr, tt.err)
				} else {
					fail = true
					explainer = fmt.Sprintf("got error %#v want %#v", errRepr, tt.err)
				}
			}

			if fail {
				t.Error(explainer)
			}
		})
	}
}
