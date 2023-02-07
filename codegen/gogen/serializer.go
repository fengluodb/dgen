package gogen

import (
	"bufio"
	"dgen/utils"
	"fmt"
	"io"
	"strings"
)

func (g *Gogen) genSerializerFunction(w io.Writer) error {
	buf := bufio.NewWriter(w)
	if g.EncodeType == "json" {
		return g.genJsonSerializerFunction(w)
	}

	if err := defaultSerializerFunc.Execute(w, g); err != nil {
		return err
	}

	for _, v := range g.StructStats {
		for _, m := range v.Members {
			g.genTypeSerialization(w, m.Type)
		}
	}

	for _, v := range g.StructStats {
		buf.WriteString(fmt.Sprintf("func (x *%s) Marshal() ([]byte, error) {\n", v.Name))
		buf.WriteString("\tdata := []byte{}\n\n")

		for _, m := range v.Members {
			if m.Type == "uint8" || m.Type == "uint16" || m.Type == "uint32" || m.Type == "uint64" ||
				m.Type == "int8" || m.Type == "int16" || m.Type == "int32" || m.Type == "int64" {
				buf.WriteString(fmt.Sprintf("\tif x.%s != 0 {\n", m.Name))
			} else if m.Type == "string" {
				buf.WriteString(fmt.Sprintf("\tif x.%s != \"\" {\n", m.Name))
			} else {
				buf.WriteString(fmt.Sprintf("\tif x.%s != nil {\n", m.Name))
			}
			buf.WriteString(fmt.Sprintf("\t\tdata = append(data, MarshalUint8(%d)...)\n", m.Seq))
			buf.WriteString(fmt.Sprintf("\t\tdata = append(data, Marshal%s(x.%s)...)\n", g.genTypeSerialization(w, m.Type), m.Name))
			if !m.Optional {
				buf.WriteString("\t}")
				buf.WriteString(" else {\n")
				buf.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"marshal failed, %s must have value\")\n", m.Name))
				buf.WriteString("\t}\n\n")
			} else {
				buf.WriteString("\t}\n\n")
			}
		}
		buf.WriteString("\treturn data, nil\n")
		buf.WriteString("}\n\n")
	}

	for _, v := range g.StructStats {
		buf.WriteString(fmt.Sprintf("func (x *%s) Unmarshal(data []byte) error {\n", v.Name))
		buf.WriteString("\tr := bytes.NewReader(data)\n\n")
		buf.WriteString("\tseq := UnmarshalUint8(r)\n")

		for i, m := range v.Members {
			buf.WriteString(fmt.Sprintf("\t if seq == %d {\n", m.Seq))
			buf.WriteString(fmt.Sprintf("\t\tx.%s = Unmarshal%s(r)\n", m.Name, g.genTypeSerialization(w, m.Type)))
			if i != len(v.Members)-1 {
				buf.WriteString("\t\tseq = UnmarshalUint8(r)\n")
			}

			if !m.Optional {
				buf.WriteString("\t}")
				buf.WriteString(" else {\n")
				buf.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"unmarshal failed, don't find %s\")\n", m.Name))
				buf.WriteString("\t}\n\n")
			} else {
				buf.WriteString("\t}\n\n")
			}
		}
		buf.WriteString("\treturn nil\n")
		buf.WriteString("}\n\n")
	}

	buf.Flush()

	return nil
}

func (g *Gogen) genJsonSerializerFunction(w io.Writer) error {
	return jsonSerializerTmpl.Execute(w, g.parser)
}

// avoid repeated generation
var serializationMap = map[string]bool{}

func (g *Gogen) genTypeSerialization(w io.Writer, typ string) string {
	if strings.HasPrefix(typ, "[]") {
		return g.genListSerialization(w, typ)
	} else if strings.HasPrefix(typ, "map") {
		return g.genMapSerialization(w, typ)
	}

	if idx := strings.LastIndex(typ, "*"); idx != -1 {
		typ = typ[idx+1:]
	}
	return utils.FirstUpper(typ)
}

func (g *Gogen) genListSerialization(w io.Writer, typ string) string {
	ele := typ[2:]
	s := g.genTypeSerialization(w, ele)
	tmpl1 := `
func MarshalList%s(v %s) []byte {
	data := []byte{}

	data = append(data, MarshalInt32(int32(len(v)))...)
	for _, val := range v {
		data = append(data, Marshal%s(val)...)
	}
	return data
}
`
	tmpl2 := `
func UnmarshalList%s(r io.Reader) %s {
	size := UnmarshalInt32(r)
	v := make(%s, 0)
	var i int32
	for i = 0; i < size; i++ {
		v = append(v, Unmarshal%s(r))
	}
	return v
}

`
	if _, ok := serializationMap[typ]; !ok {
		w.Write([]byte(fmt.Sprintf(tmpl1, s, typ, s)))
		w.Write([]byte(fmt.Sprintf(tmpl2, s, typ, typ, s)))
		serializationMap[typ] = true
	}
	return "List" + s
}

func (g *Gogen) genMapSerialization(w io.Writer, typ string) string {
	idx := strings.Index(typ, "]")
	key := typ[4:idx]
	val := typ[idx+1:]
	s := g.genTypeSerialization(w, val)
	tmpl1 := `
func MarshalMap%s%s(v %s) []byte {
	data := []byte{}

	data = append(data, MarshalInt32(int32(len(v)))...)
	for key, val := range v {
		data = append(data, Marshal%s(key)...)
		data = append(data, Marshal%s(val)...)
	}
	return data
}
`
	tmpl2 := `
func UnmarshalMap%s%s(r io.Reader) %s {
	size := UnmarshalInt32(r)
	v := make(%s)
	var i int32
	for i = 0; i < size; i++ {
		key := Unmarshal%s(r)
		val := Unmarshal%s(r)
		v[key] = val
	}
	return v
}
`
	if _, ok := serializationMap[typ]; !ok {
		w.Write([]byte(fmt.Sprintf(tmpl1, utils.FirstUpper(key), s, typ, utils.FirstUpper(key), s)))
		w.Write([]byte(fmt.Sprintf(tmpl2, utils.FirstUpper(key), s, typ, typ, utils.FirstUpper(key), s)))
		serializationMap[typ] = true
	}

	return "Map" + utils.FirstUpper(key) + s
}
