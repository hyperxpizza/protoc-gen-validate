package module

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"path/filepath"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type module struct {
	*pgs.ModuleBase
	pgsgo.Context
	tpl *template.Template
}

func New() pgs.Module {
	return &module{ModuleBase: &pgs.ModuleBase{}}
}

func (module) Name() string {
	return "protoc-gen-validate"
}

func (g *module) InitContext(c pgs.BuildContext) {
	g.ModuleBase.InitContext(c)
	g.Context = pgsgo.InitContext(c.Parameters())
	tpl := template.New("jsonify").Funcs(map[string]interface{}{
		"package": g.PackageName,
		"name":    g.Context.Name,
	})

	g.tpl = template.Must(tpl.Parse(validatorTmpl))
}

func (g module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {

	tagGenerator := newTagGenerator(g, g.Context)
	module := g.Parameters().Str("module")
	outdir := g.Parameters().Str("outdir")
	g.Debugf("output path: %s", g.Context.Params().OutputPath())
	for protoFilePath, file := range targets {

		g.Debugf("proto file path: %s", protoFilePath)

		generatedTags := tagGenerator.GenerateValidateTags(file)
		g.Debugf("%v", generatedTags)

		filename := g.Context.OutputPath(file).SetExt(".go").String()

		if module != "" {
			filename = strings.ReplaceAll(filename, string(filepath.Separator), "/")
			trim := module + "/"
			if !strings.HasPrefix(filename, trim) {
				g.Debug(fmt.Sprintf("%v: generated file does not match prefix %q", filename, module))
				g.Exit(1)
			}
			filename = strings.TrimPrefix(filename, trim)
		}

		if outdir != "" {
			filename = filepath.Join(outdir, filename)
		}

		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
		g.CheckErr(err)

		g.CheckErr(ApplyTags(fn, generatedTags))
		var buf strings.Builder
		g.CheckErr(printer.Fprint(&buf, fs, fn))
		g.Debugf("filename: %s", filename)
		g.OverwriteGeneratorFile(filename, buf.String())

		filename = strings.TrimSuffix(filename, ".go") + ".validate.go"
		g.AddGeneratorTemplateFile(filename, g.tpl, file)
	}

	return g.Artifacts()
}

const validatorTmpl = `package {{ package . }}

import (
	validator "github.com/go-playground/validator/v10"
)

{{ range .AllMessages }}

func (m *{{ name . }}) Validate() error {
	return validator.New().Struct(m)
}

{{ end }}
`
