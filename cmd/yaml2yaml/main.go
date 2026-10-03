// Command-line tool for processing YAML data.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/willoma/yaml2yaml"
)

func main() {
	input := flag.String("input", "", "Input file path (defaults to stdin)")
	output := flag.String("output", "", "Output file path (defaults to stdin)")

	flag.StringVar(input, "i", "", "Input file path (shorthand)")
	flag.StringVar(output, "o", "", "Output file path (shorthand)")
	flag.Parse()

	result, err := read(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}

	err = write(*output, result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
}

func read(input string) ([]byte, error) {
	if input != "" {
		result, err := yaml2yaml.ProcessFile(input)
		if err != nil {
			return nil, fmt.Errorf("converting %s:\n%w", input, err)
		}

		return result, nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("reading stdin: %w", err)
	}

	result, err := yaml2yaml.ProcessBytes(data)
	if err != nil {
		return nil, fmt.Errorf("converting stdin: %w", err)
	}

	return result, nil
}

func write(output string, result []byte) error {
	if output != "" {
		//nolint:gosec // Using the same permissions as os.Create.
		err := os.WriteFile(output, result, 0o666)
		if err != nil {
			return fmt.Errorf("writing to %s: %w", output, err)
		}

		return nil
	}

	_, err := os.Stdout.Write(result)
	if err != nil {
		return fmt.Errorf("writing to stdout: %w", err)
	}

	return nil
}
