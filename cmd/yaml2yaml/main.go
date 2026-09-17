// Command-line tool for processing YAML data.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/willoma/yaml2yaml"
)

//nolint:funlen // It's not THAT long
func main() {
	var err error

	closeInput := false
	inputName := "input"
	input := os.Stdin

	if len(os.Args) > 1 {
		inputName = os.Args[1]

		//nolint:gosec // Input name is provided by the user, do as they say.
		input, err = os.Open(inputName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Opening %s: %v\n", inputName, err)
			os.Exit(1)
		}

		closeInput = true
	}

	data, err := io.ReadAll(input)

	if closeInput {
		if err := input.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Closing %s: %v\n", inputName, err)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Reading %s: %v\n", inputName, err)
		os.Exit(1)
	}

	result, err := yaml2yaml.Bytes(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Converting %s:\n%v\n", inputName, err)
	}

	closeOutput := false
	outputName := "output"
	output := os.Stdout

	if len(os.Args) > 2 {
		outputName = os.Args[2]

		//nolint:gosec // Input name is provided by the user, do as they say.
		output, err = os.Create(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Creating %s: %v\n", outputName, err)
			os.Exit(1)
		}

		closeOutput = true
	}

	_, err = output.Write(result)

	if closeOutput {
		if err := output.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Closing %s: %v\n", outputName, err)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Writing to %s: %v\n", outputName, err)
	}
}
