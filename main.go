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
