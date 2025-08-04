package main

import (
	"testing"

	"github.com/sivchari/ezproto"
)

func TestHelperGenerator(t *testing.T) {
	test := ezproto.NewTest(t)
	test.TestGenerator("plugin_output", "test.proto", HelperGenerator)
}
