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
		result, err := os.ReadFile(exampleResultFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(yaml2yaml.ProcessFile(exampleFile)).To(Equal(result))
	})
})
