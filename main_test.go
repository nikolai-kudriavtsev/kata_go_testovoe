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
		// basic
		{'+', "2", "2", "4", ""},
		{'+', "II", "II", "IV", ""},
		{'-', "2", "1", "1", ""},
		{'-', "II", "I", "I", ""},
		{'*', "2", "2", "4", ""},
		{'*', "II", "II", "IV", ""},
		{'/', "2", "2", "1", ""},
		{'/', "II", "II", "I", ""},
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
				explainer = fmt.Sprintf("got result %s, want %s; ", result, tt.result)
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
