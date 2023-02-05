package gogen

import "text/template"

var (
	header1Tmpl        = must(_header1Tmpl)
	header2Tmpl        = must(_header2Tmpl)
	enumTmpl           = must(_enumTmpl)
	structTmpl         = must(_structTmpl)
	serviceTmpl        = must(_serviceTmpl)
	registerTmpl       = must(_registerTmpl)
	jsonSerializerTmpl = must(_jsonSerializerTmpl)
)

func must(s string) *template.Template {
	return template.Must(template.New("").Parse(s))
}

const _header1Tmpl = `package {{.Name}}
{{if eq .EncodeType "json"}}
import "encoding/json"
{{- end}}
`

const _header2Tmpl = `package {{.Name}}

import "github.com/fengluodb/drpc"

`

const _enumTmpl = `
{{- range .EnumStats}}
type {{.Name}} uint32
{{$name := .Name}}
const(
	{{- range $i,$v:=.Members}}
	{{if eq $i 0}}{{$v}} {{$name}} = iota 
	{{- else}}{{$v}} 
	{{- end}}
	{{- end}}
)
{{end -}}
`

const _structTmpl = `
{{- range .StructStats}}
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
	{{.Name}}(*{{.Req}}, *{{.Resp}}) error
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

{{$name := .Name}}
{{- range .Members -}}
func (c *{{$name}}Complement) {{.Name}}Handler(req []byte) (data []byte, err error) {
	args := new({{.Req}})
	if err := args.Unmarshal(req); err != nil {
		return nil, err
	}

	reply := new({{.Resp}})
	if err := c.{{$name}}.{{.Name}}(args, reply); err != nil {
		return nil, err
	}
	return reply.Marshal()
}
{{end}}

{{- end }}
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
