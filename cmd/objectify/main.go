package main

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"go/importer"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/imports"
)

var (
	schema   = flag.String("schema", "", "schema for struct")
	pkgdir   = flag.String("dir", ".", "input package")
	output   = flag.String("out", "-", "output file (default is stdout)")
	typename = flag.String("type", "", "type to generate methods for")
)

func init() {
	flag.Parse()
}

func main() {
	gen := Generator{
		Dir:  *pkgdir,
		Type: *typename,
	}

	code, err := gen.process()
	if err != nil {
		log.Fatal(err)
	}

	if *output == "-" {
		os.Stdout.Write(code)
	} else if err := ioutil.WriteFile(*output, code, 0644); err != nil {
		log.Fatal(err)
	}
}

type Generator struct {
	Dir      string
	Type     string
	Importer types.Importer
	FileSet  *token.FileSet
}

type Values struct {
	Package      string
	StructName   string
	StructFields []*Field
	Schema       string
	Imports      map[string]bool
}

type Field struct {
	Skip     bool
	Name     string
	Tag      string
	TypePtr  string
	Type     string
	Hint     string
	IsObject bool
	IsSlice  bool
	CanBeNil bool
}

func (gen *Generator) process() (code []byte, err error) {
	if gen.FileSet == nil {
		gen.FileSet = token.NewFileSet()
	}
	if gen.Importer == nil {
		gen.Importer = importer.Default()
	}
	pkg, err := gen.loadPackage()
	if err != nil {
		return nil, err
	}

	typ, err := lookupStructType(pkg.Scope(), gen.Type)
	if err != nil {
		return nil, fmt.Errorf("can't find %s in %q: %v", gen.Type, pkg.Path(), err)
	}

	values := &Values{
		Package:      pkg.Name(),
		StructName:   gen.Type,
		StructFields: []*Field{},
		Schema:       *schema,
		Imports: map[string]bool{
			"nimona.io/go/encoding": true,
		},
	}

	fmt.Printf("Objectifying %s.%s\n", values.Package, gen.Type)

	styp := typ.Underlying().(*types.Struct)

	// obj := pkg.Scope().Lookup(gen.Type)
	// ptr:=types.NewPointer(obj.Type())
	// imp:=types.Implements(obj.Name(), ifff)

	for i := 0; i < styp.NumFields(); i++ {
		f := styp.Field(i)
		if !f.Exported() {
			continue
		}
		if f.Anonymous() {
			fmt.Fprintf(os.Stderr, "Warning: ignoring embedded field %s\n", f.Name())
			continue
		}

		tag := reflect.StructTag(styp.Tag(i))
		ftag := tag.Get("fluffy")
		if ftag == "" {
			// TODO should we fallback to json?
			ftag = tag.Get("json")
		}
		vf := getMetaFromTag(ftag)
		if vf == nil {
			// no tags found
			vf = &Field{
				Tag: toLowerFirst(f.Name()),
			}
		}
		if vf.Skip {
			continue
		}
		vf.Name = f.Name()
		// vf.TypePtr = removePackageFromTypePtr(f.Type().String(), pkg.Path(), f.Pkg().Name())
		// vf.Type = removePackageFromType(f.Type().String())

		tp, tpkg := getPackageAndType(f.Type().String(), pkg.Path(), false)
		vf.TypePtr = tp

		if tpkg != "" {
			values.Imports[tpkg] = true
		}

		tp, _ = getPackageAndType(f.Type().String(), pkg.Path(), true)
		vf.Type = tp

		// vf.Type = removePackageFromType(f.Type().String())

		hint := getHint(f.Type())
		if vf.Hint == "" {
			vf.Hint = hint
		} else if vf.Hint != hint {
			panic(fmt.Errorf("existing hint of %s for field %s does not match infered %s", vf.Hint, vf.Name, hint))
		}

		if strings.Contains(vf.Hint, "O") {
			vf.IsObject = true
		}

		if vf.TypePtr[0] == '*' {
			vf.CanBeNil = true
		}

		if _, ok := f.Type().(*types.Map); ok {
			vf.CanBeNil = true
		}

		if fi, ok := f.Type().(*types.Slice); ok {
			vf.IsSlice = true
			vf.CanBeNil = true
			if _, ok := fi.Elem().(*types.Basic); ok && vf.IsSlice {
				vf.Type = "[]" + vf.Type
				vf.TypePtr = "[]" + vf.TypePtr
			}

		}

		values.StructFields = append(values.StructFields, vf)

		fmt.Printf("  - field=%s; tag=%s, type=%s, hint=%s, skipping=%t\n", vf.Name, vf.Tag, vf.Type, vf.Hint, vf.Skip)
	}

	tpl := `// Code generated by nimona.io/go/cmd/objectify. DO NOT EDIT.

// +build !generate

package {{ .Package }}

import (
	"fmt"

	{{- range $pkg, $ok := .Imports }}
	"{{ $pkg }}"
	{{- end }}
)

// ToMap returns a map compatible with f12n
func (s {{ .StructName }}) ToMap() map[string]interface{} {
	{{- range .StructFields }}
	{{- if and .IsObject .IsSlice }}
	s{{ .Name }} := []map[string]interface{}{}
	for _, v := range s.{{ .Name }} {
		s{{ .Name }} = append(s{{ .Name }}, v.ToMap())
	}
	{{- end }}
	{{- end }}
	m := map[string]interface{}{
		"@ctx:s": "{{ .Schema }}",
		{{- range .StructFields }}
		{{- if eq .Tag "@" }}
		{{- else if .CanBeNil }}
		{{- else if and .IsObject .IsSlice }}
		"{{ .Tag }}:{{ .Hint }}": s{{ .Name }},
		{{- else if .IsObject }}
		"{{ .Tag }}:{{ .Hint }}": s.{{ .Name }}.ToMap(),
		{{- else }}
		"{{ .Tag }}:{{ .Hint }}": s.{{ .Name }},
		{{- end }}
		{{- end }}
	}
	{{- range .StructFields }}
	{{- if eq .Tag "@" }}
	{{- else if .CanBeNil }}
	if s.{{ .Name }} != nil {
		{{- if .IsObject }}
		m["{{ .Tag }}:{{ .Hint }}"] = s.{{ .Name }}.ToMap()
		{{- else }}
		m["{{ .Tag }}:{{ .Hint }}"] = s.{{ .Name }}
		{{- end }}
	}
	{{- end }}
	{{- end }}
	return m
}

// ToObject returns a f12n object
func (s {{ .StructName }}) ToObject() *encoding.Object {
	return encoding.NewObjectFromMap(s.ToMap())
}

// FromMap populates the struct from a f12n compatible map
func (s *{{ .StructName }}) FromMap(m map[string]interface{}) error {
	{{- range .StructFields }}
	{{- if eq .Tag "@" }}
	s.{{ .Name }} = encoding.NewObjectFromMap(m)
	{{- else if and .IsObject .IsSlice }}
	s.{{ .Name }} = []{{ .TypePtr }}{}
	if ss, ok := m["{{ .Tag }}:{{ .Hint }}"].([]interface{}); ok {
		for _, si := range ss {
			if v, ok := si.(map[string]interface{}); ok {
				s{{ .Name }} := {{ .Type }}{}
				if err := s{{ .Name }}.FromMap(v); err != nil {
					return err
				}
				s.{{ .Name }} = append(s.{{ .Name }}, s{{ .Name }})
			} else if v, ok := m["{{ .Tag }}:{{ .Hint }}"].({{ .TypePtr }}); ok {
				s.{{ .Name }} = append(s.{{ .Name }}, v)
			}
		}
	}
	{{- else if .IsObject }}
	if v, ok := m["{{ .Tag }}:{{ .Hint }}"].(map[string]interface{}); ok {
		s.{{ .Name }} = {{ .Type }}{}
		if err := s.{{ .Name }}.FromMap(v); err != nil {
			return err
		}
	} else if v, ok := m["{{ .Tag }}:{{ .Hint }}"].({{ .TypePtr }}); ok {
		s.{{ .Name }} = v
	}
	{{- else }}
	if v, ok := m["{{ .Tag }}:{{ .Hint }}"].({{ .TypePtr }}); ok {
		s.{{ .Name }} = v
	}
	{{- end }}
	{{- end }}
	return nil
}

// FromObject populates the struct from a f12n object
func (s *{{ .StructName }}) FromObject(o *encoding.Object) error {
	return s.FromMap(o.Map())
}

// GetType returns the object's type
func (s {{ .StructName }}) GetType() string {
	return "{{ .Schema }}"
}`

	f, err := os.Create(*output)
	if err != nil {
		panic(err)
	}

	t, err := template.New("t").Parse(tpl)
	if err != nil {
		panic(err)
	}

	err = t.Execute(f, values)
	if err != nil {
		panic(err)
	}

	opt := &imports.Options{
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}
	code, err = imports.Process(*output, code, opt)
	if err != nil {
		panic(fmt.Errorf("BUG: can't gofmt generated code: %v", err))
	}
	return code, nil
}

func (gen *Generator) loadPackage() (*types.Package, error) {
	// Find the import path of the package in the given directory.
	cwd, _ := os.Getwd()
	dir := filepath.Join(gen.Dir, "*.go") // glob is stripped by ContainingPackage
	pkg, err := buildutil.ContainingPackage(&build.Default, cwd, dir)
	if err != nil {
		return nil, err
	}

	// Load the actual package.
	nocheck := func(path string) bool {
		return false
	}

	lcfg := loader.Config{
		Fset:                gen.FileSet,
		TypeCheckFuncBodies: nocheck,
	}
	lcfg.ImportWithTests(pkg.ImportPath)
	prog, err := lcfg.Load()
	if err != nil {
		return nil, err
	}
	return prog.Package(pkg.ImportPath).Pkg, nil
}

func lookupStructType(scope *types.Scope, name string) (*types.Named, error) {
	typ, err := lookupType(scope, name)
	if err != nil {
		return nil, err
	}
	_, ok := typ.Underlying().(*types.Struct)
	if !ok {
		return nil, errors.New("not a struct type")
	}
	return typ, nil
}

func lookupType(scope *types.Scope, name string) (*types.Named, error) {
	obj := scope.Lookup(name)
	if obj == nil {
		return nil, errors.New("no such identifier")
	}
	typ, ok := obj.(*types.TypeName)
	if !ok {
		return nil, errors.New("not a type")
	}
	return typ.Type().(*types.Named), nil
}

func getMetaFromTag(tag string) *Field {
	if tag == "" {
		return nil
	}

	args := strings.Split(tag, ",")

	vf := &Field{
		Tag: args[0],
	}

	tp := strings.Split(vf.Tag, ":")
	if len(tp) > 1 {
		vf.Tag = tp[0]
		vf.Hint = tp[1]
	}

	for _, t := range args {
		switch t {
		case "object":
			vf.IsObject = true
		}
	}

	if vf.Tag == "-" {
		vf.Skip = true
	}

	return vf
}

func getPackageAndType(t, pkg string, deref bool) (string, string) {
	t = strings.Replace(t, "[]", "", 1)
	ptr := false
	if t[0] == '*' {
		ptr = true
		t = t[1:]
	}

	ct := strings.Replace(t, pkg, "", 1)
	ts := strings.Split(ct, ".")
	tpkg := strings.Join(ts[:len(ts)-1], ".")
	ts = strings.Split(ct, "/")
	tt := ts[len(ts)-1]

	tt = strings.TrimLeft(tt, ".")

	if ptr {
		if deref {
			tt = "&" + tt
		} else {
			tt = "*" + tt
		}
	}
	return tt, tpkg
}

func toLowerFirst(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	s = string(a)
	return s
}

func getHint(t types.Type) string {
	if t.String() == "[]byte" {
		return "d"
	}
	switch v := t.(type) {
	case *types.Basic:
		switch v.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			return "i"
		case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return "i"
		case types.Float32, types.Float64:
			return "f"
		case types.String:
			return "s"
		}
	case *types.Array:
		st := v.Elem()
		ss := getHint(st)
		if ss != "" {
			return "A<" + ss + ">"
		}
	case *types.Slice:
		st := v.Elem()
		ss := getHint(st)
		if ss != "" {
			return "A<" + ss + ">"
		}
	case *types.Struct:
		return "O"
	case *types.Pointer:
		st := v.Elem()
		return getHint(st)
	case *types.Tuple:
	case *types.Signature:
	case *types.Interface:
	case *types.Map:
		return "O"
	case *types.Chan:
	case *types.Named:
	}
	// TODO(geoah) insane hack/assumption
	return "O"
}