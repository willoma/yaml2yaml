package yaml2yaml_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/willoma/yaml2yaml"
)

func Test(t *testing.T) {
	src, err := os.ReadFile("example.yaml")
	if err != nil {
		t.Fatalf("reading example.yaml: %v", err)
	}
	result, err := yaml2yaml.Bytes(src)
	if err != nil {
		t.Fatalf("converting example.yaml: %v", err)
	}

	ref, err := os.ReadFile("example.output.yaml")

	if !bytes.Equal(result, ref) {
		t.Fatalf("result does not match reference")
	}
}
