package yaml2yaml_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/willoma/yaml2yaml"
)

const (
	exampleFile       = "example.yaml"
	exampleResultFile = "example.result.yaml"
)

var _ = Describe("Example file", func() {
	It("transforms the example file as expected", func() {
		source, err := os.ReadFile(exampleFile)
		Expect(err).ToNot(HaveOccurred())

		result, err := os.ReadFile(exampleResultFile)
		Expect(err).ToNot(HaveOccurred())

		Expect(yaml2yaml.Bytes(source)).To(Equal(result))
	})
})
