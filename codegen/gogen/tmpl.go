package gogen

import "text/template"

var (
	header1Tmpl           = must(_header1Tmpl)
	header2Tmpl           = must(_header2Tmpl)
	enumTmpl              = must(_enumTmpl)
	enumSerializationTmpl = must(_enumSerializationTmpl)
	structTmpl            = must(_structTmpl)
	serviceTmpl           = must(_serviceTmpl)
	registerTmpl          = must(_registerTmpl)
	jsonSerializerTmpl    = must(_jsonSerializerTmpl)
	defaultSerializerFunc = must(_defaultSerializerFunc)
)

func must(s string) *template.Template {
	return template.Must(template.New("").Parse(s))
}

const _header1Tmpl = `package {{.Name}}
{{if eq .EncodeType "json"}}
import "encoding/json"
{{ else }}
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)
{{- end}}
`

const _header2Tmpl = `package {{.Name}}

import "github.com/fengluodb/drpc"

`

const _enumTmpl = `
{{- range .EnumStats}}
type {{.Name}} uint32
{{end -}}
`

const _enumSerializationTmpl = `
{{- range .EnumStats}}
{{$name := .Name}}
const(
	{{- range $i,$v:=.Members}}
	{{if eq $i 0}}{{$v}} {{$name}} = iota 
	{{- else}}{{$v}} 
	{{- end}}
	{{- end}}
)

func (x *{{.Name}}) Marshal() ([]byte, error) {
	return MarshalUint32(uint32(*x)), nil
}

func (x *{{.Name}}) Unmarshal(data []byte) error {
	v := binary.LittleEndian.Uint32(data)
	*x = {{.Name}}(v)
	return nil
}
{{end -}}
`

const _structTmpl = `
{{- range .StructStats }}
var _ Serializer = (*{{.Name}})(nil)
{{- end }}

type Serializer interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
{{ range .StructStats }}
type {{.Name}} struct {
	{{- range .Members}}
	{{.Name}} {{.Type}}
	{{- end}}
}
{{ end -}}
`

const _serviceTmpl = `
{{- range .ServiceStats}}
type {{.Name}} interface {
	{{- range .Members}}
	{{.Name}}(*{{.Req}} {{- if ne .Resp ""}}, *{{.Resp}} {{- end}}) error
	{{- end}}
}

type {{.Name}}Handler interface {
	{{- range .Members}}
	{{.Name}}Handler(req []byte) (data []byte, err error)
	{{- end}}
}

type {{.Name}}Complement struct {
	{{.Name}} {{.Name}}
}
{{- $name := .Name}}
{{ range .Members }}
func (c *{{$name}}Complement) {{.Name}}Handler(req []byte) (data []byte, err error) {
	args := new({{.Req}})
	if err := args.Unmarshal(req); err != nil {
		return nil, err
	}
	{{if ne .Resp ""}}
	reply := new({{.Resp}}){{ end }}
	if err := c.{{$name}}.{{.Name}}(args{{- if ne .Resp ""}}, reply {{- end}}); err != nil {
		return nil, err
	}
	{{if ne .Resp ""}}return reply.Marshal(){{else}}return nil, nil{{end}}
}
{{end}}
{{- end -}}
`

const _registerTmpl = `
{{- range .ServiceStats}}
func Register{{.Name}}Service(s *drpc.Server, serviceName string, complement {{.Name}}) {
	c := &{{.Name}}Complement{
		{{.Name}}: complement,
	}
	{{$name := .Name}}
	{{- range .Members}}
	drpc.RegisterService(s, serviceName+".{{.Name}}", c.{{.Name}}Handler)
	{{- end}}
}
{{- end}}
`

const _jsonSerializerTmpl = `
{{- range .MessageStats}}
func (x *{{.Name}}) Marshal() ([]byte, error) {
	return json.Marshal(x)
}

func (x *{{.Name}}) Unmarshal(data []byte) error {
	return json.Unmarshal(data, x)
}
{{end -}}
`

const _defaultSerializerFunc = `
{{- range .StructStats}}
func Marshal{{.Name}}(v *{{.Name}}) []byte {
	data, _ := v.Marshal()
	return data
}

func Unmarshal{{.Name}}(r io.Reader) *{{.Name}} {
	br := r.(*bytes.Reader)
	data, _ := io.ReadAll(br)
	v := new({{.Name}})
	v.Unmarshal(data)
	tmp := Marshal{{.Name}}(v)
	*br = *bytes.NewReader(data[len(tmp):])

	return v
}
{{end }}
func MarshalUint8(v uint8) []byte {
	data := []byte{}
	return append(data, byte(v))
}

func UnmarshalUint8(r io.Reader) uint8 {
	data := make([]byte, 1)
	io.ReadFull(r, data)
	return uint8(data[0])
}

func MarshalUint16(v uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, v)
	return data
}

func UnmarshalUint16(r io.Reader) uint16 {
	data := make([]byte, 2)
	io.ReadFull(r, data)
	return binary.LittleEndian.Uint16(data)
}

func MarshalUint32(v uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, v)
	return data
}

func UnmarshalUint32(r io.Reader) uint32 {
	data := make([]byte, 4)
	io.ReadFull(r, data)
	return binary.LittleEndian.Uint32(data)
}

func MarshalUint64(v uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, v)
	return data
}

func UnmarshalUint64(r io.Reader) uint64 {
	data := make([]byte, 8)
	io.ReadFull(r, data)
	return binary.LittleEndian.Uint64(data)
}

func MarshalInt8(v int8) []byte {
	data := []byte{}
	return append(data, byte(v))
}

func UnmarshalInt8(r io.Reader) int8 {
	data := make([]byte, 1)
	io.ReadFull(r, data)
	return int8(data[0])
}

func MarshalInt16(v int16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(v))
	return data
}

func UnmarshalInt16(r io.Reader) int16 {
	data := make([]byte, 2)
	io.ReadFull(r, data)
	return int16(binary.LittleEndian.Uint16(data))
}

func MarshalInt32(v int32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(v))
	return data
}

func UnmarshalInt32(r io.Reader) int32 {
	data := make([]byte, 4)
	io.ReadFull(r, data)
	return int32(binary.LittleEndian.Uint32(data))
}

func MarshalInt64(v int64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(v))
	return data
}

func UnmarshalInt64(r io.Reader) int64 {
	data := make([]byte, 8)
	io.ReadFull(r, data)
	return int64(binary.LittleEndian.Uint64(data))
}

func MarshalString(s string) []byte {
	data := []byte{}
	data = append(data, MarshalInt32(int32(len(s)))...)
	data = append(data, []byte(s)...)
	return data
}

func UnmarshalString(r io.Reader) string {
	size := UnmarshalInt32(r)

	data := make([]byte, size)
	io.ReadFull(r, data)
	return string(data)
}
`
