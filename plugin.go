package ezproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type GeneratorFunc func(ctx *Context, file *File) error

type Options struct {
	Debug          bool
	PackageMapping map[string]string
}

type Plugin struct {
	options          Options
	generators       map[string]GeneratorFunc
	parameterHandler func(params map[string]string, options *Options)
}

func NewPlugin() *Plugin {
	return &Plugin{
		options: Options{
			Debug:          false,
			PackageMapping: make(map[string]string),
		},
		generators: make(map[string]GeneratorFunc),
	}
}

func (p *Plugin) WithOptions(opts Options) *Plugin {
	p.options = opts
	if p.options.PackageMapping == nil {
		p.options.PackageMapping = make(map[string]string)
	}
	return p
}

func (p *Plugin) GenerateFor(pattern string, generator GeneratorFunc) *Plugin {
	p.generators[pattern] = generator
	return p
}

// WithParameterHandler sets a custom parameter handler
func (p *Plugin) WithParameterHandler(handler func(params map[string]string, options *Options)) *Plugin {
	p.parameterHandler = handler
	return p
}

func (p *Plugin) Run() error {
	opts := protogen.Options{}
	opts.Run(func(gen *protogen.Plugin) error {
		// Parse plugin parameters
		params := parseParameters(gen.Request.GetParameter())
		
		// Update plugin options with parsed parameters
		p.updateOptionsFromParams(params)
		
		// Call custom parameter handler if provided
		if p.parameterHandler != nil {
			p.parameterHandler(params, &p.options)
		}
		
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			file := &File{
				proto: f,
				Name:  f.Desc.Path(),
			}

			ctx := &Context{
				plugin:     p,
				gen:        gen,
				file:       f,
				parameters: params,
			}

			for pattern, generator := range p.generators {
				if p.matchesPattern(f.Desc.Path(), pattern) {
					if p.options.Debug {
						fmt.Fprintf(os.Stderr, "[DEBUG] Generating for %s with pattern %s\n", f.Desc.Path(), pattern)
					}
					
					if err := generator(ctx, file); err != nil {
						return fmt.Errorf("generator failed for %s: %w", f.Desc.Path(), err)
					}
				}
			}
		}
		return nil
	})
	return nil
}

// parseParameters parses plugin parameters from protoc
func parseParameters(parameter string) map[string]string {
	params := make(map[string]string)
	if parameter == "" {
		return params
	}
	
	// Split by comma: "key1=value1,key2=value2"
	pairs := strings.Split(parameter, ",")
	for _, pair := range pairs {
		if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
			params[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else {
			// Boolean flag without value
			params[strings.TrimSpace(pair)] = "true"
		}
	}
	
	return params
}

// updateOptionsFromParams updates plugin options from parsed parameters
func (p *Plugin) updateOptionsFromParams(params map[string]string) {
	for key, value := range params {
		switch key {
		case "debug":
			p.options.Debug = value == "true" || value == "1"
		case "package_mapping":
			// Handle package mapping: package_mapping=proto.package:go.package
			if mapping := strings.SplitN(value, ":", 2); len(mapping) == 2 {
				if p.options.PackageMapping == nil {
					p.options.PackageMapping = make(map[string]string)
				}
				p.options.PackageMapping[mapping[0]] = mapping[1]
			}
		}
	}
}

func (p *Plugin) matchesPattern(path, pattern string) bool {
	if pattern == "*" || pattern == "*.proto" {
		return true
	}
	
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return strings.Contains(path, strings.TrimSuffix(pattern, "*"))
	}
	return matched
}