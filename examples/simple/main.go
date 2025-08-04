package main

import (
	"github.com/sivchari/ezproto"
)

// HelperGenerator generates helper structs for proto messages
func HelperGenerator(ctx *ezproto.Context, file *ezproto.File) error {
	ctx.Debug("Processing file: %s", file.Name)

	code := ctx.Code().
		Comment("Generated from " + file.Name).
		Package(file.Package()).
		EmptyLine()

	// Generate enum information
	for _, enum := range file.Enums() {
		code.Comment("Enum: " + enum.GoName() + " (" + enum.FullName() + ")")
	}
	if len(file.Enums()) > 0 {
		code.EmptyLine()
	}

	for _, msg := range file.Messages() {
		code.Comment("Message: "+msg.GoName()).
			Struct(msg.GoName()+"Helper", func(sb *ezproto.StructBuilder) {
				sb.Field("msg", "*"+msg.GoName())
			}).
			EmptyLine().
			Function("New"+msg.GoName()+"Helper(msg *"+msg.GoName()+") *"+msg.GoName()+"Helper", func(cb *ezproto.CodeBuilder) {
				cb.Return("&" + msg.GoName() + "Helper{msg: msg}")
			}).
			EmptyLine()
		
		// Debug: show field information including enums, maps, oneofs
		for _, field := range msg.Fields() {
			var fieldInfo string
			if field.IsMap() {
				fieldInfo = "map field"
			} else if field.IsEnum() {
				fieldInfo = "enum field: " + field.Type()
			} else if field.IsMessage() {
				fieldInfo = "message field: " + field.Type()
			} else {
				fieldInfo = "scalar field: " + field.Type()
			}
			code.Comment("Field " + field.GoName() + ": " + fieldInfo)
		}
		
		for _, oneof := range msg.Oneofs() {
			code.Comment("Oneof: " + oneof.GoName())
		}
		code.EmptyLine()
	}

	for _, svc := range file.Services() {
		code.Comment("Service: " + svc.GoName())
		for _, method := range svc.Methods() {
			streaming := ""
			if method.IsClientStreaming() || method.IsServerStreaming() {
				streaming = " (streaming)"
			}
			code.Comment("Method: " + method.InputType() + " -> " + method.OutputType() + streaming)
		}
		code.EmptyLine()
	}

	code.Generate()
	return nil
}

// NewHelperPlugin creates the helper plugin
func NewHelperPlugin() *ezproto.Plugin {
	return ezproto.NewPlugin().
		WithOptions(ezproto.Options{
			Debug: true,
		}).
		GenerateFor("*.proto", HelperGenerator)
}

func main() {
	NewHelperPlugin().Run()
}
