package ezproto

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sivchari/golden"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Test provides a testing framework for ezproto code generators.
type Test struct {
	t      *testing.T
	golden *golden.Golden
}

// NewTest creates a new test instance with golden file testing.
func NewTest(t *testing.T) *Test {
	t.Helper()

	return &Test{
		t:      t,
		golden: golden.New(t, golden.WithDir("testdata")),
	}
}

// TestGenerator tests a generator function with a proto file.
func (test *Test) TestGenerator(name, protoFile string, generator GeneratorFunc) {
	// Generate descriptor set using protoc
	req, err := generateCodeGeneratorRequest(protoFile)
	if err != nil {
		test.t.Fatalf("Failed to generate CodeGeneratorRequest: %v", err)
	}

	// Create protogen.Plugin
	gen, err := protogen.Options{}.New(req)
	if err != nil {
		test.t.Fatalf("Failed to create protogen plugin: %v", err)
	}

	// Find the file to generate
	var file *protogen.File

	for _, f := range gen.Files {
		if f.Proto.GetName() == protoFile {
			file = f

			break
		}
	}

	if file == nil {
		test.t.Fatalf("File %s not found in generated files", protoFile)
	}

	// Create output buffer
	var output bytes.Buffer

	// Create ezproto context
	ctx := &Context{
		plugin: &Plugin{},
		gen:    gen,
		file:   file,
		output: &testGeneratedFile{buffer: &output},
	}

	// Create ezproto File wrapper
	ezFile := &File{
		proto: file,
		Name:  file.Proto.GetName(),
	}

	// Execute generator
	err = generator(ctx, ezFile)
	if err != nil {
		test.t.Fatalf("Generator failed: %v", err)
	}

	// Compare with golden file
	test.golden.Assert(name, output.String())
}

// testGeneratedFile implements GeneratedFile interface for testing.
type testGeneratedFile struct {
	buffer *bytes.Buffer
}

// P implements the GeneratedFile interface by writing to a buffer.
func (f *testGeneratedFile) P(v ...any) {
	if len(v) == 0 {
		f.buffer.WriteString("\n")

		return
	}

	for i, arg := range v {
		if i > 0 {
			f.buffer.WriteString(" ")
		}

		fmt.Fprintf(f.buffer, "%v", arg)
	}

	f.buffer.WriteString("\n")
}

// QualifiedGoIdent implements the GeneratedFile interface by returning the Go name.
func (f *testGeneratedFile) QualifiedGoIdent(ident protogen.GoIdent) string {
	return ident.GoName
}

// generateCodeGeneratorRequest uses protoc to generate a proper CodeGeneratorRequest
// but with better error handling and validation.
func generateCodeGeneratorRequest(protoFile string) (*pluginpb.CodeGeneratorRequest, error) {
	// Validate input file path
	if protoFile == "" {
		return nil, fmt.Errorf("proto file path cannot be empty")
	}

	if !strings.HasSuffix(protoFile, ".proto") {
		return nil, fmt.Errorf("invalid proto file extension: %s", protoFile)
	}

	// Validate file path doesn't contain suspicious characters
	if strings.ContainsAny(protoFile, ";&|`$()") {
		return nil, fmt.Errorf("invalid characters in proto file path: %s", protoFile)
	}

	// Check if protoc is available
	if _, err := exec.LookPath("protoc"); err != nil {
		return nil, fmt.Errorf("protoc not found in PATH: %w", err)
	}

	// Create temporary file for descriptor set
	tmpFile, err := os.CreateTemp("", "*.pb")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()
	defer func() {
		_ = tmpFile.Close()
	}()

	// Run protoc to generate descriptor set
	//nolint:gosec // protoFile is validated above for security
	cmd := exec.Command("protoc",
		"--descriptor_set_out="+tmpFile.Name(),
		"--include_imports",
		protoFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("protoc failed with output: %s, error: %w", string(output), err)
	}

	// Read descriptor set
	descriptorData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read descriptor set: %w", err)
	}

	// Parse descriptor set
	var descriptorSet descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(descriptorData, &descriptorSet); err != nil {
		return nil, fmt.Errorf("failed to parse descriptor set: %w", err)
	}

	// Create CodeGeneratorRequest
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{filepath.Base(protoFile)},
		ProtoFile:      descriptorSet.File,
	}

	return req, nil
}
