package directives

import (
	"strings"

	"github.com/willoma/yaml2yaml/errors"
	"github.com/willoma/yaml2yaml/internal/expressions"
	"github.com/willoma/yaml2yaml/internal/node"
)

const (
	NameFmt    = "!fmt"
	fmtOpener  = "$("
	fmtCloser  = ")"
	fmtEscaper = `\`
)

type dirFmt struct {
	src     node.Source
	sources []node.Source
}

func Fmt(src node.Source) (*dirFmt, error) {
	strSrc, ok := src.(node.String)
	if !ok {
		return nil, errors.ErrNotAString
	}

	sources, err := parseFmt(strSrc)
	if err != nil {
		return nil, err
	}

	return &dirFmt{src: src, sources: sources}, nil
}

func parseFmt(srcNode node.String) ([]node.Source, error) {
	var (
		errs   errors.JoinedError
		result []node.Source
	)

	pos := srcNode.Pos()
	src := srcNode.String()

	for src != "" {
		beforeOpener, afterOpener, openerFound := strings.Cut(src, fmtOpener)
		if !openerFound {
			result = append(result, node.NewString(pos, beforeOpener))

			break
		}

		if beforeEscaper, escaperFound := strings.CutSuffix(beforeOpener, fmtEscaper); escaperFound {
			// Opener found, but it is escaped, ignoring and going on.
			result = append(result, node.NewString(pos, beforeEscaper+fmtOpener))
			src = afterOpener

			continue
		}

		result = append(result, node.NewString(pos, beforeOpener))

		expr, remaining, err := findFmtCloser(pos, afterOpener)

		src = remaining

		if err != nil {
			errs = errs.Append(err)

			continue
		}

		result = append(result, expr)
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return result, nil
}

func findFmtCloser(pos [2]int, src string) (result node.Source, remaining string, err error) {
	var expr strings.Builder

	for {
		beforeCloser, afterCloser, closerFound := strings.Cut(src, fmtCloser)
		if !closerFound {
			expr.WriteString(beforeCloser)

			return node.NewString(pos, expr.String()), "", nil
		}

		if beforeEscaper, escaperFound := strings.CutSuffix(beforeCloser, fmtEscaper); escaperFound {
			// Closer found, but it is escaped, ignoring and searching again.
			expr.WriteString(beforeEscaper)
			expr.WriteString(fmtCloser)

			continue
		}

		expr.WriteString(beforeCloser)

		exprSource, err := expressions.Parse(pos, expr.String())

		return exprSource, afterCloser, err
	}
}

func (f *dirFmt) Pos() [2]int {
	return f.src.Pos()
}

func (f *dirFmt) String() string {
	return NameApply + " " + f.src.String()
}

func (f *dirFmt) EqualsSource(other node.Source) bool {
	if other, ok := other.(*dirFmt); ok {
		return f.src.EqualsSource(other.src)
	}

	return false
}

func (f *dirFmt) Process(defs node.Definitions) (node.Node, error) {
	var (
		buf  strings.Builder
		errs errors.JoinedError
	)

	for _, src := range f.sources {
		child, err := src.Process(defs)
		if err != nil {
			errs = errs.Append(err)

			continue
		}

		buf.WriteString(child.String())
	}

	if err := errs.Compact(); err != nil {
		return nil, err
	}

	return node.NewString(f.Pos(), buf.String()), nil
}

func (f *dirFmt) PopulateDefs(node.Definitions) error {
	return nil
}
