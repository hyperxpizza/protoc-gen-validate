package module

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/fatih/structtag"
)

type applier struct {
	err  error
	tags map[string]string
}

// ApplyTags updates the existing tags with the generated parquet tag
func ApplyTags(n ast.Node, tags map[string]map[string]string) error {
	r := applier{}
	f := func(n ast.Node) ast.Visitor {
		if r.err != nil {
			return nil
		}

		if tp, ok := n.(*ast.TypeSpec); ok {
			r.tags = tags[tp.Name.String()]
			return &r
		}

		return nil
	}

	ast.Walk(structVisitor{f}, n)

	return r.err
}

type structVisitor struct {
	visitor func(n ast.Node) ast.Visitor
}

func (v structVisitor) Visit(n ast.Node) ast.Visitor {
	if tp, ok := n.(*ast.TypeSpec); ok {
		if _, ok := tp.Type.(*ast.StructType); ok {
			ast.Walk(v.visitor(n), n)
			return nil
		}
	}
	return v
}

func (a *applier) Visit(n ast.Node) ast.Visitor {
	if a.err != nil {
		return nil
	}

	if f, ok := n.(*ast.Field); ok {
		if len(f.Names) == 0 {
			return nil
		}
		parquetTag := a.tags[f.Names[0].String()]
		if parquetTag == "" {
			return nil
		}

		parsedParquetTag, err := structtag.Parse(parquetTag)
		if err != nil {
			a.err = err
			return nil
		}

		if f.Tag == nil {
			f.Tag = &ast.BasicLit{
				Kind: token.STRING,
			}
		}

		oldTags, err := structtag.Parse(strings.Trim(f.Tag.Value, "`"))
		if err != nil {
			a.err = err
			return nil
		}

		sort.Stable(parsedParquetTag) // sort tags according to keys
		for _, t := range parsedParquetTag.Tags() {
			oldTags.Set(t)
		}

		f.Tag.Value = "`" + oldTags.String() + "`"

		return nil
	}

	return a
}
