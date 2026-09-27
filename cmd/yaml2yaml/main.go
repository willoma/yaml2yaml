// Command-line tool for processing YAML data.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/willoma/yaml2yaml"
)

func main() {
	result, err := read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}

	err = write(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
}

func read() ([]byte, error) {
	if len(os.Args) > 1 {
		result, err := yaml2yaml.ProcessFile(os.Args[1])
		if err != nil {
			return nil, fmt.Errorf("converting %s: %w", os.Args[1], err)
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

func write(result []byte) error {
	if len(os.Args) > 2 {
		//nolint:gosec // Using the same permissions as os.Create.
		err := os.WriteFile(os.Args[2], result, 0o666)
		if err != nil {
			return fmt.Errorf("writing to %s: %w", os.Args[2], err)
		}

		return nil
	}

	_, err := os.Stdout.Write(result)
	if err != nil {
		return fmt.Errorf("writing to stdout: %w", err)
	}

	return nil
}
