package ezproto

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// File represents a proto file being processed.
type File struct {
	proto *protogen.File
	Name  string
}

// Package returns the Go package name for this file.
func (f *File) Package() string {
	return string(f.proto.GoPackageName)
}

// GoImportPath returns the Go import path for this file.
func (f *File) GoImportPath() string {
	return string(f.proto.GoImportPath)
}

// Messages returns all message types defined in this file.
func (f *File) Messages() []*Message {
	messages := make([]*Message, 0, len(f.proto.Messages))
	for _, msg := range f.proto.Messages {
		messages = append(messages, &Message{
			proto: msg,
			Name:  string(msg.Desc.Name()),
		})
	}

	return messages
}

// Services returns all service types defined in this file.
func (f *File) Services() []*Service {
	services := make([]*Service, 0, len(f.proto.Services))
	for _, svc := range f.proto.Services {
		services = append(services, &Service{
			proto: svc,
			Name:  string(svc.Desc.Name()),
		})
	}

	return services
}

// Enums returns all enum types defined in this file.
func (f *File) Enums() []*Enum {
	enums := make([]*Enum, 0, len(f.proto.Enums))
	for _, enum := range f.proto.Enums {
		enums = append(enums, &Enum{
			proto: enum,
			Name:  string(enum.Desc.Name()),
		})
	}

	return enums
}

// Message represents a protobuf message type.
type Message struct {
	proto *protogen.Message
	Name  string
}

// Fields returns all fields defined in this message.
func (m *Message) Fields() []*Field {
	fields := make([]*Field, 0, len(m.proto.Fields))
	for _, field := range m.proto.Fields {
		fields = append(fields, &Field{
			proto: field,
			Name:  string(field.Desc.Name()),
		})
	}

	return fields
}

// GoName returns the Go type name for this message.
func (m *Message) GoName() string {
	return m.proto.GoIdent.GoName
}

// Oneofs returns all oneof fields defined in this message.
func (m *Message) Oneofs() []*Oneof {
	oneofs := make([]*Oneof, 0, len(m.proto.Oneofs))
	for _, oneof := range m.proto.Oneofs {
		oneofs = append(oneofs, &Oneof{
			proto: oneof,
			Name:  string(oneof.Desc.Name()),
		})
	}

	return oneofs
}

// Field represents a field in a protobuf message.
type Field struct {
	proto *protogen.Field
	Name  string
}

// GoName returns the Go field name.
func (f *Field) GoName() string {
	return f.proto.GoName
}

// GoType returns the Go type name for this field.
func (f *Field) GoType() string {
	return f.proto.GoIdent.GoName
}

// IsRepeated returns true if this field is repeated (array/slice).
func (f *Field) IsRepeated() bool {
	return f.proto.Desc.Cardinality() == protoreflect.Repeated
}

// IsOptional returns true if this field is optional.
func (f *Field) IsOptional() bool {
	return f.proto.Desc.HasOptionalKeyword()
}

// IsMap returns true if this field is a map type.
func (f *Field) IsMap() bool {
	return f.proto.Desc.IsMap()
}

// IsEnum returns true if this field is an enum type.
func (f *Field) IsEnum() bool {
	return f.proto.Desc.Kind() == protoreflect.EnumKind
}

// IsMessage returns true if this field is a message type.
func (f *Field) IsMessage() bool {
	return f.proto.Desc.Kind() == protoreflect.MessageKind
}

// Type returns the protobuf type name for this field.
func (f *Field) Type() string {
	if f.IsEnum() {
		return string(f.proto.Desc.Enum().FullName())
	}

	if f.IsMessage() {
		return string(f.proto.Desc.Message().FullName())
	}

	return f.proto.Desc.Kind().String()
}

// Service represents a protobuf service definition.
type Service struct {
	proto *protogen.Service
	Name  string
}

// Methods returns all methods defined in this service.
func (s *Service) Methods() []*Method {
	methods := make([]*Method, 0, len(s.proto.Methods))
	for _, method := range s.proto.Methods {
		methods = append(methods, &Method{
			proto: method,
			Name:  string(method.Desc.Name()),
		})
	}

	return methods
}

// GoName returns the Go type name for this service.
func (s *Service) GoName() string {
	return s.proto.GoName
}

// Method represents a method in a protobuf service.
type Method struct {
	proto *protogen.Method
	Name  string
}

// GoName returns the Go method name.
func (m *Method) GoName() string {
	return m.proto.GoName
}

// InputType returns the Go type name for the input parameter.
func (m *Method) InputType() string {
	return m.proto.Input.GoIdent.GoName
}

// OutputType returns the Go type name for the output parameter.
func (m *Method) OutputType() string {
	return m.proto.Output.GoIdent.GoName
}

// IsClientStreaming returns true if this method uses client streaming.
func (m *Method) IsClientStreaming() bool {
	return m.proto.Desc.IsStreamingClient()
}

// IsServerStreaming returns true if this method uses server streaming.
func (m *Method) IsServerStreaming() bool {
	return m.proto.Desc.IsStreamingServer()
}

// Enum represents a protobuf enum definition.
type Enum struct {
	proto *protogen.Enum
	Name  string
}

// Values returns all values defined in this enum.
func (e *Enum) Values() []*EnumValue {
	values := make([]*EnumValue, 0, len(e.proto.Values))
	for _, value := range e.proto.Values {
		values = append(values, &EnumValue{
			proto: value,
			Name:  string(value.Desc.Name()),
		})
	}

	return values
}

// GoName returns the Go type name for this enum.
func (e *Enum) GoName() string {
	return e.proto.GoIdent.GoName
}

// FullName returns the fully qualified protobuf name for this enum.
func (e *Enum) FullName() string {
	return string(e.proto.Desc.FullName())
}

// EnumValue represents a value in a protobuf enum.
type EnumValue struct {
	proto *protogen.EnumValue
	Name  string
}

// GoName returns the Go constant name for this enum value.
func (ev *EnumValue) GoName() string {
	return ev.proto.GoIdent.GoName
}

// Number returns the numeric value of this enum value.
func (ev *EnumValue) Number() int32 {
	return int32(ev.proto.Desc.Number())
}

// Oneof represents a protobuf oneof field group.
type Oneof struct {
	proto *protogen.Oneof
	Name  string
}

// GoName returns the Go field name for this oneof.
func (o *Oneof) GoName() string {
	return o.proto.GoName
}

// Fields returns all fields in this oneof group.
func (o *Oneof) Fields() []*Field {
	fields := make([]*Field, 0, len(o.proto.Fields))
	for _, field := range o.proto.Fields {
		fields = append(fields, &Field{
			proto: field,
			Name:  string(field.Desc.Name()),
		})
	}

	return fields
}
