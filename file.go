package ezproto

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type File struct {
	proto *protogen.File
	Name  string
}

func (f *File) Package() string {
	return string(f.proto.GoPackageName)
}

func (f *File) GoImportPath() string {
	return string(f.proto.GoImportPath)
}

func (f *File) Messages() []*Message {
	var messages []*Message
	for _, msg := range f.proto.Messages {
		messages = append(messages, &Message{
			proto: msg,
			Name:  string(msg.Desc.Name()),
		})
	}
	return messages
}

func (f *File) Services() []*Service {
	var services []*Service
	for _, svc := range f.proto.Services {
		services = append(services, &Service{
			proto: svc,
			Name:  string(svc.Desc.Name()),
		})
	}
	return services
}

func (f *File) Enums() []*Enum {
	var enums []*Enum
	for _, enum := range f.proto.Enums {
		enums = append(enums, &Enum{
			proto: enum,
			Name:  string(enum.Desc.Name()),
		})
	}
	return enums
}

type Message struct {
	proto *protogen.Message
	Name  string
}

func (m *Message) Fields() []*Field {
	var fields []*Field
	for _, field := range m.proto.Fields {
		fields = append(fields, &Field{
			proto: field,
			Name:  string(field.Desc.Name()),
		})
	}
	return fields
}

func (m *Message) GoName() string {
	return m.proto.GoIdent.GoName
}

func (m *Message) Oneofs() []*Oneof {
	var oneofs []*Oneof
	for _, oneof := range m.proto.Oneofs {
		oneofs = append(oneofs, &Oneof{
			proto: oneof,
			Name:  string(oneof.Desc.Name()),
		})
	}
	return oneofs
}

type Field struct {
	proto *protogen.Field
	Name  string
}

func (f *Field) GoName() string {
	return f.proto.GoName
}

func (f *Field) GoType() string {
	return f.proto.GoIdent.GoName
}

func (f *Field) IsRepeated() bool {
	return f.proto.Desc.Cardinality() == protoreflect.Repeated
}

func (f *Field) IsOptional() bool {
	return f.proto.Desc.HasOptionalKeyword()
}

func (f *Field) IsMap() bool {
	return f.proto.Desc.IsMap()
}

func (f *Field) IsEnum() bool {
	return f.proto.Desc.Kind() == protoreflect.EnumKind
}

func (f *Field) IsMessage() bool {
	return f.proto.Desc.Kind() == protoreflect.MessageKind
}

func (f *Field) Type() string {
	if f.IsEnum() {
		return string(f.proto.Desc.Enum().FullName())
	}
	if f.IsMessage() {
		return string(f.proto.Desc.Message().FullName())
	}
	return f.proto.Desc.Kind().String()
}

type Service struct {
	proto *protogen.Service
	Name  string
}

func (s *Service) Methods() []*Method {
	var methods []*Method
	for _, method := range s.proto.Methods {
		methods = append(methods, &Method{
			proto: method,
			Name:  string(method.Desc.Name()),
		})
	}
	return methods
}

func (s *Service) GoName() string {
	return s.proto.GoName
}

type Method struct {
	proto *protogen.Method
	Name  string
}

func (m *Method) GoName() string {
	return m.proto.GoName
}

func (m *Method) InputType() string {
	return m.proto.Input.GoIdent.GoName
}

func (m *Method) OutputType() string {
	return m.proto.Output.GoIdent.GoName
}

func (m *Method) IsClientStreaming() bool {
	return m.proto.Desc.IsStreamingClient()
}

func (m *Method) IsServerStreaming() bool {
	return m.proto.Desc.IsStreamingServer()
}

type Enum struct {
	proto *protogen.Enum
	Name  string
}

func (e *Enum) Values() []*EnumValue {
	var values []*EnumValue
	for _, value := range e.proto.Values {
		values = append(values, &EnumValue{
			proto: value,
			Name:  string(value.Desc.Name()),
		})
	}
	return values
}

func (e *Enum) GoName() string {
	return e.proto.GoIdent.GoName
}

func (e *Enum) FullName() string {
	return string(e.proto.Desc.FullName())
}

type EnumValue struct {
	proto *protogen.EnumValue
	Name  string
}

func (ev *EnumValue) GoName() string {
	return ev.proto.GoIdent.GoName
}

func (ev *EnumValue) Number() int32 {
	return int32(ev.proto.Desc.Number())
}

type Oneof struct {
	proto *protogen.Oneof
	Name  string
}

func (o *Oneof) GoName() string {
	return o.proto.GoName
}

func (o *Oneof) Fields() []*Field {
	var fields []*Field
	for _, field := range o.proto.Fields {
		fields = append(fields, &Field{
			proto: field,
			Name:  string(field.Desc.Name()),
		})
	}
	return fields
}