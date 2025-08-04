package ezproto

import (
	"fmt"
	"strings"
)

type CodeBuilder struct {
	ctx    *Context
	lines  []string
	indent int
}

func (c *Context) NewCodeBuilder() *CodeBuilder {
	return &CodeBuilder{
		ctx:   c,
		lines: make([]string, 0),
	}
}

func (cb *CodeBuilder) Line(format string, args ...interface{}) *CodeBuilder {
	indentStr := strings.Repeat("\t", cb.indent)
	line := fmt.Sprintf(format, args...)
	cb.lines = append(cb.lines, indentStr+line)
	return cb
}

func (cb *CodeBuilder) EmptyLine() *CodeBuilder {
	cb.lines = append(cb.lines, "")
	return cb
}

func (cb *CodeBuilder) Block(header string, fn func(*CodeBuilder)) *CodeBuilder {
	cb.Line("%s {", header)
	cb.indent++
	fn(cb)
	cb.indent--
	cb.Line("}")
	return cb
}

func (cb *CodeBuilder) Comment(text string) *CodeBuilder {
	return cb.Line("// %s", text)
}

func (cb *CodeBuilder) Package(name string) *CodeBuilder {
	return cb.Line("package %s", name)
}

func (cb *CodeBuilder) Import(path string) *CodeBuilder {
	return cb.Line("import \"%s\"", path)
}

func (cb *CodeBuilder) ImportBlock(imports []string) *CodeBuilder {
	if len(imports) == 0 {
		return cb
	}
	
	if len(imports) == 1 {
		return cb.Import(imports[0])
	}
	
	cb.Line("import (")
	cb.indent++
	for _, imp := range imports {
		cb.Line("\"%s\"", imp)
	}
	cb.indent--
	cb.Line(")")
	return cb
}

func (cb *CodeBuilder) Struct(name string, fn func(*StructBuilder)) *CodeBuilder {
	sb := &StructBuilder{cb: cb}
	cb.Line("type %s struct {", name)
	cb.indent++
	fn(sb)
	cb.indent--
	cb.Line("}")
	return cb
}

func (cb *CodeBuilder) Function(signature string, fn func(*CodeBuilder)) *CodeBuilder {
	return cb.Block("func "+signature, fn)
}

func (cb *CodeBuilder) Method(receiver, name, params, returns string, fn func(*CodeBuilder)) *CodeBuilder {
	signature := fmt.Sprintf("(%s) %s(%s)", receiver, name, params)
	if returns != "" {
		signature += " " + returns
	}
	return cb.Function(signature, fn)
}

func (cb *CodeBuilder) Return(values ...string) *CodeBuilder {
	if len(values) == 0 {
		return cb.Line("return")
	}
	return cb.Line("return %s", strings.Join(values, ", "))
}

func (cb *CodeBuilder) Assign(left, right string) *CodeBuilder {
	return cb.Line("%s = %s", left, right)
}

func (cb *CodeBuilder) DeclareAssign(left, right string) *CodeBuilder {
	return cb.Line("%s := %s", left, right)
}

func (cb *CodeBuilder) Interface(name string, fn func(*InterfaceBuilder)) *CodeBuilder {
	ib := &InterfaceBuilder{cb: cb}
	cb.Line("type %s interface {", name)
	cb.indent++
	fn(ib)
	cb.indent--
	cb.Line("}")
	return cb
}

func (cb *CodeBuilder) Const(name, value string) *CodeBuilder {
	return cb.Line("const %s = %s", name, value)
}

func (cb *CodeBuilder) ConstBlock(fn func(*ConstBuilder)) *CodeBuilder {
	constB := &ConstBuilder{cb: cb}
	cb.Line("const (")
	cb.indent++
	fn(constB)
	cb.indent--
	cb.Line(")")
	return cb
}

func (cb *CodeBuilder) Var(name, typ string, value ...string) *CodeBuilder {
	if len(value) > 0 {
		return cb.Line("var %s %s = %s", name, typ, value[0])
	}
	return cb.Line("var %s %s", name, typ)
}

func (cb *CodeBuilder) VarBlock(fn func(*VarBuilder)) *CodeBuilder {
	varB := &VarBuilder{cb: cb}
	cb.Line("var (")
	cb.indent++
	fn(varB)
	cb.indent--
	cb.Line(")")
	return cb
}

func (cb *CodeBuilder) TypeAlias(name, typ string) *CodeBuilder {
	return cb.Line("type %s = %s", name, typ)
}

func (cb *CodeBuilder) If(condition string, fn func(*CodeBuilder)) *CodeBuilder {
	return cb.Block("if "+condition, fn)
}

func (cb *CodeBuilder) IfErr(fn func(*CodeBuilder)) *CodeBuilder {
	return cb.If("err != nil", fn)
}

func (cb *CodeBuilder) For(init, condition, post string, fn func(*CodeBuilder)) *CodeBuilder {
	forStmt := "for"
	if init != "" || condition != "" || post != "" {
		forStmt += " " + init + "; " + condition + "; " + post
	}
	return cb.Block(forStmt, fn)
}

func (cb *CodeBuilder) ForRange(variable, iterable string, fn func(*CodeBuilder)) *CodeBuilder {
	return cb.Block("for "+variable+" := range "+iterable, fn)
}

func (cb *CodeBuilder) Switch(expr string, fn func(*SwitchBuilder)) *CodeBuilder {
	sb := &SwitchBuilder{cb: cb}
	if expr != "" {
		cb.Line("switch %s {", expr)
	} else {
		cb.Line("switch {")
	}
	cb.indent++
	fn(sb)
	cb.indent--
	cb.Line("}")
	return cb
}

func (cb *CodeBuilder) RawString(content string) *CodeBuilder {
	return cb.Line("`%s`", content)
}

func (cb *CodeBuilder) BuildTag(tag string) *CodeBuilder {
	return cb.Line("//go:build %s", tag)
}

func (cb *CodeBuilder) GoGenerate(tag string) *CodeBuilder {
	return cb.Line("//go:generate %s", tag)
}

func (cb *CodeBuilder) Generate() {
	if cb.ctx.output == nil {
		cb.ctx.createOutputFile()
	}
	
	for _, line := range cb.lines {
		cb.ctx.output.P(line)
	}
}

type StructBuilder struct {
	cb *CodeBuilder
}

func (sb *StructBuilder) Field(name, typ string, tags ...string) *StructBuilder {
	line := fmt.Sprintf("%s %s", name, typ)
	if len(tags) > 0 {
		line += " `" + strings.Join(tags, " ") + "`"
	}
	sb.cb.Line("%s", line)
	return sb
}

func (sb *StructBuilder) EmbeddedField(typ string) *StructBuilder {
	sb.cb.Line("%s", typ)
	return sb
}

type InterfaceBuilder struct {
	cb *CodeBuilder
}

func (ib *InterfaceBuilder) Method(name, params, returns string) *InterfaceBuilder {
	signature := name + "(" + params + ")"
	if returns != "" {
		signature += " " + returns
	}
	ib.cb.Line("%s", signature)
	return ib
}

func (ib *InterfaceBuilder) EmbeddedInterface(typ string) *InterfaceBuilder {
	ib.cb.Line("%s", typ)
	return ib
}

type ConstBuilder struct {
	cb *CodeBuilder
}

func (constB *ConstBuilder) Const(name, value string) *ConstBuilder {
	constB.cb.Line("%s = %s", name, value)
	return constB
}

func (constB *ConstBuilder) ConstWithType(name, typ, value string) *ConstBuilder {
	constB.cb.Line("%s %s = %s", name, typ, value)
	return constB
}

type VarBuilder struct {
	cb *CodeBuilder
}

func (vb *VarBuilder) Var(name, typ string, value ...string) *VarBuilder {
	if len(value) > 0 {
		vb.cb.Line("%s %s = %s", name, typ, value[0])
	} else {
		vb.cb.Line("%s %s", name, typ)
	}
	return vb
}

type SwitchBuilder struct {
	cb *CodeBuilder
}

func (sb *SwitchBuilder) Case(value string, fn func(*CodeBuilder)) *SwitchBuilder {
	sb.cb.Line("case %s:", value)
	sb.cb.indent++
	fn(sb.cb)
	sb.cb.indent--
	return sb
}

func (sb *SwitchBuilder) Default(fn func(*CodeBuilder)) *SwitchBuilder {
	sb.cb.Line("default:")
	sb.cb.indent++
	fn(sb.cb)
	sb.cb.indent--
	return sb
}