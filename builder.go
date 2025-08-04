// Package ezproto provides tools for generating code from Protocol Buffer definitions.
package ezproto

import (
	"fmt"
	"strings"
)

// CodeBuilder provides a fluent API for generating Go code.
type CodeBuilder struct {
	ctx    *Context
	lines  []string
	indent int
}

// NewCodeBuilder creates a new CodeBuilder instance for building code.
func (c *Context) NewCodeBuilder() *CodeBuilder {
	return &CodeBuilder{
		ctx:   c,
		lines: make([]string, 0),
	}
}

// Line adds a formatted line of code with proper indentation.
func (cb *CodeBuilder) Line(format string, args ...interface{}) *CodeBuilder {
	indentStr := strings.Repeat("\t", cb.indent)
	line := fmt.Sprintf(format, args...)
	cb.lines = append(cb.lines, indentStr+line)

	return cb
}

// EmptyLine adds an empty line to the code.
func (cb *CodeBuilder) EmptyLine() *CodeBuilder {
	cb.lines = append(cb.lines, "")

	return cb
}

// Block creates a code block with curly braces and proper indentation.
func (cb *CodeBuilder) Block(header string, fn func(*CodeBuilder)) *CodeBuilder {
	cb.Line("%s {", header)

	cb.indent++
	fn(cb)

	cb.indent--
	cb.Line("}")

	return cb
}

// Comment adds a single-line comment to the code.
func (cb *CodeBuilder) Comment(text string) *CodeBuilder {
	return cb.Line("// %s", text)
}

// Package adds a package declaration.
func (cb *CodeBuilder) Package(name string) *CodeBuilder {
	return cb.Line("package %s", name)
}

// Import adds a single import statement.
func (cb *CodeBuilder) Import(path string) *CodeBuilder {
	return cb.Line("import \"%s\"", path)
}

// ImportBlock adds an import block with multiple imports.
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

// Struct creates a struct type definition.
func (cb *CodeBuilder) Struct(name string, fn func(*StructBuilder)) *CodeBuilder {
	sb := &StructBuilder{cb: cb}
	cb.Line("type %s struct {", name)

	cb.indent++

	fn(sb)

	cb.indent--
	cb.Line("}")

	return cb
}

// Function creates a function definition.
func (cb *CodeBuilder) Function(signature string, fn func(*CodeBuilder)) *CodeBuilder {
	return cb.Block("func "+signature, fn)
}

// Method creates a method definition with receiver.
func (cb *CodeBuilder) Method(receiver, name, params, returns string, fn func(*CodeBuilder)) *CodeBuilder {
	signature := fmt.Sprintf("(%s) %s(%s)", receiver, name, params)
	if returns != "" {
		signature += " " + returns
	}

	return cb.Function(signature, fn)
}

// Return adds a return statement with optional values.
func (cb *CodeBuilder) Return(values ...string) *CodeBuilder {
	if len(values) == 0 {
		return cb.Line("return")
	}

	return cb.Line("return %s", strings.Join(values, ", "))
}

// Assign creates an assignment statement.
func (cb *CodeBuilder) Assign(left, right string) *CodeBuilder {
	return cb.Line("%s = %s", left, right)
}

// DeclareAssign creates a short variable declaration with assignment.
func (cb *CodeBuilder) DeclareAssign(left, right string) *CodeBuilder {
	return cb.Line("%s := %s", left, right)
}

// Interface creates an interface type definition.
func (cb *CodeBuilder) Interface(name string, fn func(*InterfaceBuilder)) *CodeBuilder {
	ib := &InterfaceBuilder{cb: cb}
	cb.Line("type %s interface {", name)

	cb.indent++

	fn(ib)

	cb.indent--
	cb.Line("}")

	return cb
}

// Const creates a single constant declaration.
func (cb *CodeBuilder) Const(name, value string) *CodeBuilder {
	return cb.Line("const %s = %s", name, value)
}

// ConstBlock creates a const block with multiple constants.
func (cb *CodeBuilder) ConstBlock(fn func(*ConstBuilder)) *CodeBuilder {
	constB := &ConstBuilder{cb: cb}
	cb.Line("const (")

	cb.indent++

	fn(constB)

	cb.indent--
	cb.Line(")")

	return cb
}

// Var creates a variable declaration with optional initialization.
func (cb *CodeBuilder) Var(name, typ string, value ...string) *CodeBuilder {
	if len(value) > 0 {
		return cb.Line("var %s %s = %s", name, typ, value[0])
	}

	return cb.Line("var %s %s", name, typ)
}

// VarBlock creates a var block with multiple variable declarations.
func (cb *CodeBuilder) VarBlock(fn func(*VarBuilder)) *CodeBuilder {
	varB := &VarBuilder{cb: cb}
	cb.Line("var (")

	cb.indent++

	fn(varB)

	cb.indent--
	cb.Line(")")

	return cb
}

// TypeAlias creates a type alias declaration.
func (cb *CodeBuilder) TypeAlias(name, typ string) *CodeBuilder {
	return cb.Line("type %s = %s", name, typ)
}

// If creates an if statement block.
func (cb *CodeBuilder) If(condition string, fn func(*CodeBuilder)) *CodeBuilder {
	return cb.Block("if "+condition, fn)
}

// IfErr creates an if err != nil block.
func (cb *CodeBuilder) IfErr(fn func(*CodeBuilder)) *CodeBuilder {
	return cb.If("err != nil", fn)
}

// For creates a for loop with init, condition, and post statements.
func (cb *CodeBuilder) For(init, condition, post string, fn func(*CodeBuilder)) *CodeBuilder {
	forStmt := "for"
	if init != "" || condition != "" || post != "" {
		forStmt += " " + init + "; " + condition + "; " + post
	}

	return cb.Block(forStmt, fn)
}

// ForRange creates a for range loop.
func (cb *CodeBuilder) ForRange(variable, iterable string, fn func(*CodeBuilder)) *CodeBuilder {
	return cb.Block("for "+variable+" := range "+iterable, fn)
}

// Switch creates a switch statement.
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

// RawString adds a raw string literal using backticks.
func (cb *CodeBuilder) RawString(content string) *CodeBuilder {
	return cb.Line("`%s`", content)
}

// BuildTag adds a //go:build constraint comment.
func (cb *CodeBuilder) BuildTag(tag string) *CodeBuilder {
	return cb.Line("//go:build %s", tag)
}

// GoGenerate adds a //go:generate comment.
func (cb *CodeBuilder) GoGenerate(tag string) *CodeBuilder {
	return cb.Line("//go:generate %s", tag)
}

// Generate outputs all accumulated code to the context's output file.
func (cb *CodeBuilder) Generate() {
	if cb.ctx.output == nil {
		cb.ctx.createOutputFile()
	}

	for _, line := range cb.lines {
		cb.ctx.output.P(line)
	}
}

// StructBuilder provides methods for building struct definitions.
type StructBuilder struct {
	cb *CodeBuilder
}

// Field adds a field to the struct with optional tags.
func (sb *StructBuilder) Field(name, typ string, tags ...string) *StructBuilder {
	line := fmt.Sprintf("%s %s", name, typ)
	if len(tags) > 0 {
		line += " `" + strings.Join(tags, " ") + "`"
	}

	sb.cb.Line("%s", line)

	return sb
}

// EmbeddedField adds an embedded field to the struct.
func (sb *StructBuilder) EmbeddedField(typ string) *StructBuilder {
	sb.cb.Line("%s", typ)

	return sb
}

// InterfaceBuilder provides methods for building interface definitions.
type InterfaceBuilder struct {
	cb *CodeBuilder
}

// Method adds a method signature to the interface.
func (ib *InterfaceBuilder) Method(name, params, returns string) *InterfaceBuilder {
	signature := name + "(" + params + ")"
	if returns != "" {
		signature += " " + returns
	}

	ib.cb.Line("%s", signature)

	return ib
}

// EmbeddedInterface adds an embedded interface to the interface.
func (ib *InterfaceBuilder) EmbeddedInterface(typ string) *InterfaceBuilder {
	ib.cb.Line("%s", typ)

	return ib
}

// ConstBuilder provides methods for building const blocks.
type ConstBuilder struct {
	cb *CodeBuilder
}

// Const adds a constant declaration to the const block.
func (constB *ConstBuilder) Const(name, value string) *ConstBuilder {
	constB.cb.Line("%s = %s", name, value)

	return constB
}

// ConstWithType adds a typed constant declaration to the const block.
func (constB *ConstBuilder) ConstWithType(name, typ, value string) *ConstBuilder {
	constB.cb.Line("%s %s = %s", name, typ, value)

	return constB
}

// VarBuilder provides methods for building var blocks.
type VarBuilder struct {
	cb *CodeBuilder
}

// Var adds a variable declaration to the var block.
func (vb *VarBuilder) Var(name, typ string, value ...string) *VarBuilder {
	if len(value) > 0 {
		vb.cb.Line("%s %s = %s", name, typ, value[0])
	} else {
		vb.cb.Line("%s %s", name, typ)
	}

	return vb
}

// SwitchBuilder provides methods for building switch statements.
type SwitchBuilder struct {
	cb *CodeBuilder
}

// Case adds a case clause to the switch statement.
func (sb *SwitchBuilder) Case(value string, fn func(*CodeBuilder)) *SwitchBuilder {
	sb.cb.Line("case %s:", value)

	sb.cb.indent++
	fn(sb.cb)

	sb.cb.indent--

	return sb
}

// Default adds a default clause to the switch statement.
func (sb *SwitchBuilder) Default(fn func(*CodeBuilder)) *SwitchBuilder {
	sb.cb.Line("default:")

	sb.cb.indent++
	fn(sb.cb)

	sb.cb.indent--

	return sb
}
