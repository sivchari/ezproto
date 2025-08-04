package ezproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type Context struct {
	plugin     *Plugin
	gen        *protogen.Plugin
	file       *protogen.File
	output     GeneratedFile
	parameters map[string]string
}

// GeneratedFile interface for abstraction
type GeneratedFile interface {
	P(v ...any)
	QualifiedGoIdent(ident protogen.GoIdent) string
}

func (c *Context) Code() *CodeBuilder {
	return c.NewCodeBuilder()
}

func (c *Context) NewOutputFile(filename string) GeneratedFile {
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}
	c.output = c.gen.NewGeneratedFile(filename, c.file.GoImportPath)
	return c.output
}

func (c *Context) createOutputFile() {
	if c.output != nil {
		return
	}

	base := filepath.Base(c.file.Desc.Path())
	name := strings.TrimSuffix(base, ".proto") + ".pb.go"
	c.output = c.gen.NewGeneratedFile(name, c.file.GoImportPath)
}

func (c *Context) Import(importPath string) string {
	if c.output == nil {
		c.createOutputFile()
	}
	return c.output.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: protogen.GoImportPath(importPath),
	})
}

func (c *Context) Files() []*File {
	var files []*File
	for _, f := range c.gen.Files {
		if f.Generate {
			files = append(files, &File{
				proto: f,
				Name:  f.Desc.Path(),
			})
		}
	}
	return files
}

func (c *Context) Debug(format string, args ...interface{}) {
	if c.plugin.options.Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// Parameters returns the plugin parameters passed from protoc
func (c *Context) Parameters() map[string]string {
	return c.parameters
}

// GetParameter returns a specific parameter value
func (c *Context) GetParameter(key string) (string, bool) {
	value, exists := c.parameters[key]
	return value, exists
}

// GetParameterWithDefault returns a parameter value or default if not found
func (c *Context) GetParameterWithDefault(key, defaultValue string) string {
	if value, exists := c.parameters[key]; exists {
		return value
	}
	return defaultValue
}
