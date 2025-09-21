package module

import (
	"fmt"

	"github.com/hyperxpizza/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type tagGenerator struct {
	pgs.Visitor
	pgs.DebuggerCommon
	pgsgo.Context
	tags map[string]map[string]string
}

func newTagGenerator(d pgs.DebuggerCommon, ctx pgsgo.Context) *tagGenerator {
	tagGenerator := &tagGenerator{
		Context:        ctx,
		DebuggerCommon: d,
		tags:           map[string]map[string]string{},
	}
	tagGenerator.Visitor = pgs.PassThroughVisitor(tagGenerator)

	return tagGenerator
}

func (e *tagGenerator) GenerateValidateTags(f pgs.File) map[string]map[string]string {
	e.tags = map[string]map[string]string{}

	e.CheckErr(pgs.Walk(e, f))

	return e.tags
}

func (e *tagGenerator) VisitFile(f pgs.File) (pgs.Visitor, error) {
	e.Debug("Visiting file:", f.Name().String())
	return e, nil
}

func (e *tagGenerator) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
	e.Debug("Visiting message:", m.Name().String())
	return e, nil
}

func (e *tagGenerator) VisitField(f pgs.Field) (pgs.Visitor, error) {

	e.Debugf("visiting field: %s", f.Name().String())

	msgName := e.Context.Name(f.Message()).String()
	fieldName := e.Context.Name(f).String()

	var tag string
	ok, err := f.Extension(validate.E_Tag, &tag)
	if err != nil {
		return e, err
	}

	if !ok {
		return e, nil
	}

	if e.tags[msgName] == nil {
		e.tags[msgName] = map[string]string{}
	}

	e.tags[msgName][fieldName] = fmt.Sprintf(`validate:"%s"`, tag)

	return e, nil
}
