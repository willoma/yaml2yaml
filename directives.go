package yaml2yaml

const (
	dirDefine = "!define"
	dirApply  = "!apply"
	dirIf     = "!if"
	dirFor    = "!for"
)

//nolint:gochecknoglobals // List of directives tags, never modified...
var directivesTags = []string{dirDefine, dirApply, dirIf, dirFor}
