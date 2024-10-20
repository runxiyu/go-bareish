package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"git.sr.ht/~sircmpwn/getopt"

	"git.sr.ht/~runxiyu/go-bareish/schema"
)

const templateString = `
package {{.package}}

// Code generated by go-bare/cmd/gen, DO NOT EDIT.

import (
{{- if .schema.NeedErrors }}
	"errors"
{{- end }}
	"git.sr.ht/~runxiyu/go-bareish"
)

{{ define "type" }}
	{{- if eq (typeKind .) "PrimitiveType"  -}}
		{{ primitiveType .Kind }}
	{{- else if eq (typeKind .) "DataType" -}}
		[{{if gt .Length 0 }}{{.Length}}{{end}}]byte
	{{- else if eq (typeKind .) "ArrayType" -}}
		[{{if gt .Length 0 }}{{.Length}}{{end}}]{{template "type" .Member}}
	{{- else if eq (typeKind .) "StructType" -}}
		struct {
			{{- range .Fields }}
				{{ capitalize .Name }} {{ template "type" .Type }} {{ structTag .Name }}
			{{- end -}}
		}
	{{- else if eq (typeKind .) "NamedUserType" -}}
		{{.Name}}
	{{- else if eq (typeKind .) "MapType" -}}
		map[{{template "type" .Key}}]{{template "type" .Value}}
	{{- else if eq (typeKind .) "OptionalType" -}}
		*{{template "type" .Subtype}}
	{{- end -}}
{{ end }}

{{with .schema}}

{{range .UserTypes}}
	type {{ .Name }} {{ template "type" .Type }}

	func (t *{{ .Name }}) Decode(data []byte) error {
		return bare.Unmarshal(data, t)
	}

	func (t *{{ .Name }}) Encode() ([]byte, error) {
		return bare.Marshal(t)
	}
{{end}}

{{range .Enums}}
type {{ .Name }} {{ primitiveType .Kind }}

{{ $name := .Name }}

const (
		{{- range $i, $el := .Values }}
			{{ .Name }} {{ $name }} = {{ .Value }}
		{{- end -}}
	)

	func (t {{ .Name }}) String() string {
		switch (t) {
		{{- range .Values }}
		case {{ .Name }}:
			return "{{ .Name }}"
		{{- end -}}
		}
		panic(errors.New("Invalid {{.Name}} value"))
	}
{{end}}

{{ if gt (len .Unions) 0 }}
	{{range .Unions}}
		type {{ .Name }} interface {
			bare.Union
		}

		{{range .Type.Types}}
			func (_ {{.Type.Name}}) IsUnion() {}
		{{end}}
	{{end}}

	func init() {
		{{- range .Unions}}
		bare.RegisterUnion((*{{.Name}})(nil)).
			{{ $len := len .Type.Types }}
			{{range $i, $el := .Type.Types}}
				Member(*new({{ template "type" $el.Type}}), {{$el.Tag}}){{- if not (last $len $i) -}}.{{end}}
			{{end}}
		{{ end }}
	}
{{ end}}

{{end}}
`

var funcs = template.FuncMap{
	"typeKind": func(ty interface{}) string {
		switch ty := ty.(type) {
		case *schema.PrimitiveType:
			return "PrimitiveType"
		case *schema.DataType:
			return "DataType"
		case *schema.StructType:
			return "StructType"
		case *schema.NamedUserType:
			return "NamedUserType"
		case *schema.MapType:
			return "MapType"
		case *schema.ArrayType:
			return "ArrayType"
		case *schema.OptionalType:
			return "OptionalType"
		default:
			panic(fmt.Sprintf("Unimplemented schema type: %T", ty))
		}
	},
	"primitiveType": func(t schema.TypeKind) string {
		switch t {
		case schema.U8:
			return "uint8"
		case schema.U16:
			return "uint16"
		case schema.U32:
			return "uint32"
		case schema.U64:
			return "uint64"
		case schema.UINT:
			return "uint"
		case schema.I8:
			return "int8"
		case schema.I16:
			return "int16"
		case schema.I32:
			return "int32"
		case schema.I64:
			return "int64"
		case schema.INT:
			return "int"
		case schema.F32:
			return "float32"
		case schema.F64:
			return "float64"
		case schema.Bool:
			return "bool"
		case schema.String:
			return "string"
		case schema.Void:
			return "struct{}"
		}
		panic(fmt.Errorf("Invalid primitive type %d", t))
	},
	"structTag": func(name string) string {
		return fmt.Sprintf("`bare:\"%s\"`", name)
	},
	"capitalize": func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"last": func(len, i int) bool {
		return i+1 == len
	},
}

func main() {
	cfg := parseArgs()
	out := &bytes.Buffer{}

	tmpl, err := template.New("").Funcs(funcs).Parse(templateString)
	if err != nil {
		log.Fatalf("error parsing template: %v", err)
	}

	types := parseSchema(cfg.In, cfg.Skip)

	data := make(map[string]interface{})

	data["package"] = cfg.PackageName
	data["schema"] = types

	err = tmpl.Execute(out, data)
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}

	// Format generated code
	formatted, err := format.Source(out.Bytes())
	if err != nil {
		log.Println(out.String())
		log.Fatalf("--- error formatting source code: %v", err)
	}

	err = ioutil.WriteFile(cfg.Out, formatted, 0644)
	if err != nil {
		log.Fatalf("error writing output to %s: %e", cfg.Out, err)
	}
}

type Config struct {
	PackageName string
	In          string
	Out         string
	Skip        map[string]bool
}

func parseArgs() *Config {
	cfg := &Config{}

	log.SetFlags(0)
	opts, optind, err := getopt.Getopts(os.Args, "hs:p:")
	if err != nil {
		log.Fatalf("error: %e", err)
	}

	cfg.PackageName = "gen"
	cfg.Skip = make(map[string]bool)

	for _, opt := range opts {
		switch opt.Option {
		case 'p':
			cfg.PackageName = opt.Value
		case 's':
			cfg.Skip[opt.Value] = true
		case 'h':
			log.Println("Usage: gen [-p <package>] [-s <skip type>] <input.bare> <output.go>")
			os.Exit(0)
		}
	}

	args := os.Args[optind:]
	if len(args) != 2 {
		log.Fatal("Usage: gen [-p <package>] <input.bare> <output.go>")
	}

	cfg.In = args[0]
	cfg.Out = args[1]

	return cfg
}

type Types struct {
	UserTypes  []*schema.UserDefinedType
	Enums      []*schema.UserDefinedEnum
	Unions     []*schema.UserDefinedType
	NeedErrors bool
}

func parseSchema(path string, skip map[string]bool) Types {
	inf, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening %s: %e", path, err)
	}
	defer inf.Close()

	schemaTypes, err := schema.Parse(inf)
	if err != nil {
		log.Fatalf("error parsing %s: %e", path, err)
	}

	types := Types{}

	for _, ty := range schemaTypes {
		if skip[ty.Name()] {
			continue
		}

		switch ty := ty.(type) {
		case *schema.UserDefinedType:
			if ty.Type().Kind() == schema.Union {
				types.Unions = append(types.Unions, ty)
				continue
			}
			types.UserTypes = append(types.UserTypes, ty)

		case *schema.UserDefinedEnum:
			types.Enums = append(types.Enums, ty)

		}
	}

	if len(types.Enums) > 0 {
		types.NeedErrors = true
	}

	return types
}
