package gogen

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"dgen/config"
	"dgen/parser"
)

type Gogen struct {
	Name        string
	Output      string
	EncodeType  string
	StructStats []*structStats
	StructMap   map[string]struct{}

	parser *parser.Parser
}

type structStats struct {
	Name    string
	Members []*structMember
}

type structMember struct {
	Seq      uint8
	Optional bool
	Type     string
	Name     string
}

func Gen(config *config.CodegenConfig) error {
	gogen, err := NewGogen(config.Filename, config.OutputDir, config.EncodeType)
	if err != nil {
		return err
	}
	return gogen.Gen()
}

func NewGogen(filepath string, outputDir string, encode string) (*Gogen, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	_, filename := path.Split(filepath)
	filename = strings.Split(filename, ".")[0]
	return &Gogen{
		Name:       filename,
		Output:     path.Join(outputDir, filename),
		parser:     parser.NewParser(f),
		EncodeType: encode,
		StructMap:  make(map[string]struct{}),
	}, nil
}

func (g *Gogen) Gen() error {
	if err := g.parser.Parse(); err != nil {
		log.Println("dgen:parse error:", err)
	}

	g.convertType()
	if err := g.gen1(); err != nil {
		return err
	}

	if err := g.gen2(); err != nil {
		return err
	}

	return nil
}

// enum and struct are defined here
func (g *Gogen) gen1() error {
	if err := os.MkdirAll(g.Output, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path.Join(g.Output, fmt.Sprintf("%s.go", g.Name)))
	if err != nil {
		return err
	}

	if err := g.genHeader1(f); err != nil {
		return err
	}

	if err := g.genEnum(f); err != nil {
		return err
	}

	if err := g.genStruct(f); err != nil {
		return err
	}

	if err := g.genSerializerFunction(f); err != nil {
		return err
	}

	return nil
}

// the content about drpc is defined here
func (g *Gogen) gen2() error {
	if err := os.MkdirAll(g.Output, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path.Join(g.Output, fmt.Sprintf("%s.drpc.go", g.Name)))
	if err != nil {
		return err
	}

	if err := g.genHeader2(f); err != nil {
		return err
	}

	if err := g.genService(f); err != nil {
		return err
	}

	if err := g.genRegisterFunc(f); err != nil {
		return err
	}

	return nil
}

func (g *Gogen) genHeader1(w io.Writer) error {
	if err := header1Tmpl.Execute(w, g); err != nil {
		return err
	}
	return nil
}

func (g *Gogen) genHeader2(w io.Writer) error {
	if err := header2Tmpl.Execute(w, g); err != nil {
		return err
	}
	return nil
}

func (g *Gogen) genEnum(w io.Writer) error {
	if err := enumTmpl.Execute(w, g.parser); err != nil {
		return err
	}
	if g.EncodeType != "json" {
		if err := enumSerializationTmpl.Execute(w, g.parser); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gogen) genStruct(w io.Writer) error {
	if err := structTmpl.Execute(w, g); err != nil {
		return err
	}
	return nil
}

func (g *Gogen) genService(w io.Writer) error {
	return serviceTmpl.Execute(w, g.parser)
}

func (g *Gogen) genRegisterFunc(w io.Writer) error {
	if err := registerTmpl.Execute(w, g.parser); err != nil {
		return err
	}
	return nil
}

// convert the message into struct
func (g *Gogen) convertType() {
	for _, message := range g.parser.MessageStats {
		g.StructMap[message.Name] = struct{}{}
	}

	for _, message := range g.parser.MessageStats {
		ss := &structStats{
			Name: message.Name,
		}
		for _, m := range message.Members {
			sm := &structMember{
				Seq:      m.Seq,
				Optional: m.Optional,
				Type:     g.getType(m.Type),
				Name:     m.Name,
			}
			// if type is a struct, we use the pointer of the struct
			if _, ok := g.StructMap[sm.Type]; ok {
				sm.Type = fmt.Sprintf("*%s", sm.Type)
			}
			ss.Members = append(ss.Members, sm)
		}
		g.StructStats = append(g.StructStats, ss)
	}
}

func (g *Gogen) getType(v interface{}) string {
	mapVal, ok := v.(parser.MapType)
	if ok {
		return g.getMap(mapVal)
	}

	listVal, ok := v.(parser.ListType)
	if ok {
		return g.getList(listVal)
	}

	stringVal, ok := v.(string)
	if ok {
		if _, ok := g.StructMap[stringVal]; ok {
			stringVal = "*" + stringVal
		}
		return stringVal
	}

	return ""
}

func (g *Gogen) getMap(val parser.MapType) string {
	s := ""
	s += fmt.Sprintf("map[%s]", val.Key)
	s += g.getType(val.Val)

	return s
}

func (g *Gogen) getList(val parser.ListType) string {
	s := "[]"
	s += g.getType(val.Ele)

	return s
}
