package yaml2yaml

import (
	"strconv"

	"go.yaml.in/yaml/v4"

	"github.com/willoma/yaml2yaml/internal/errors"
)

func locateMappingEntryError(i int, basePath string, src *yaml.Node, err error) error {
	var keyStr string

	key, kErr := parseNode(basePath, src)
	if kErr != nil {
		keyStr = "[" + strconv.Itoa(i/2) + "]"
	} else {
		keyStr = key.String()
	}

	return errors.Locate(keyStr, pos(src), err)
}

func pos(src *yaml.Node) [2]int {
	return [2]int{src.Line, src.Column}
}
