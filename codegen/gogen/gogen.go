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
	"dgen/utils"
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
	Seq      int
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

func NewGogen(filename string, outputDir string, encode string) (*Gogen, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	name := strings.Split(filename, ".")[0]
	return &Gogen{
		Name:       name,
		Output:     path.Join(outputDir, name),
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
	return nil
}

func (g *Gogen) genStruct(w io.Writer) error {
	if err := structTmpl.Execute(w, g); err != nil {
		return err
	}
	return nil
}

func (g *Gogen) genSerializerFunction(w io.Writer) error {
	if g.EncodeType == "json" {
		return g.genJsonSerializerFunction(w)
	}
	return nil
}

func (g *Gogen) genJsonSerializerFunction(w io.Writer) error {
	return jsonSerializerTmpl.Execute(w, g.parser)
}

func (g *Gogen) genService(w io.Writer) error {
	// s := ""
	// for _, service := range g.parser.ServiceStats {
	// 	s += fmt.Sprintf("type %s interface{\n", service.Name)
	// 	for _, m := range service.Members {
	// 		s += fmt.Sprintf("\t%s(*%s, *%s) error\n", m.Name, m.Req, m.Resp)
	// 	}
	// 	s += "}\n\n"
	// 	s += g.genServiceHandler(service)
	// }
	// return s
	return serviceTmpl.Execute(w, g.parser)
}

func (g *Gogen) genServiceHandler(service *parser.ServiceStat) string {
	// s := fmt.Sprintf("type %sHandler interface{\n", service.Name)
	// for _, m := range service.Members {
	// 	s += fmt.Sprintf("\t%s(req []byte) (data []byte, err error)\n", m.Name)
	// }
	// s += "}\n\n"

	s := fmt.Sprintf("type %sHandler struct {\n", utils.FirstLower(service.Name))
	s += fmt.Sprintf("\t%s %s\n", utils.FirstLower(service.Name), service.Name)
	s += "}\n\n"

	for _, m := range service.Members {
		s += fmt.Sprintf("func (h *%sHandler) %sHandler(req []byte) (data []byte, err error) {\n", utils.FirstLower(service.Name), m.Name)
		s += fmt.Sprintf("\targs := new(%s)\n", m.Req)
		s += "\t if err := args.Unmarshal(req); err != nil {\n"
		s += "\t\treturn nil, err"
		s += "}\n\n"
		s += fmt.Sprintf("\treply := new(%s)\n\n", m.Resp)
		s += fmt.Sprintf("\tif err := h.%s.%s(args, reply); err != nil {\n", utils.FirstLower(service.Name), m.Name)
		s += "\t\treturn nil, err\n"
		s += "\t}\n"
		s += "\treturn reply.Marshal()\n"
		s += "}\n\n"
	}
	return s
}

func (g *Gogen) genRegisterFunc(w io.Writer) error {
	// s := ""
	// for _, service := range g.parser.ServiceStats {
	// 	s += fmt.Sprintf("func Register%sService(s *drpc.Server, serviceName string, %s %s) {\n", service.Name, utils.FirstLower(service.Name), service.Name)
	// 	s += fmt.Sprintf("\t%sHandler := &%sHandler{\n", utils.FirstLower(service.Name), utils.FirstLower(service.Name))
	// 	s += fmt.Sprintf("\t\t%s:%s,\n", utils.FirstLower(service.Name), utils.FirstLower(service.Name))
	// 	s += "\t}\n"
	// 	for _, m := range service.Members {
	// 		s += fmt.Sprintf("\tdrpc.RegisterService(s, serviceName+\".%s\", %sHandler.%sHandler)\n", m.Name, utils.FirstLower(service.Name), m.Name)
	// 	}
	// 	s += "}\n\n"
	// }
	// return s
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
				Type:     getType(m.Type),
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

func getType(v interface{}) string {
	mapVal, ok := v.(parser.MapType)
	if ok {
		return getMap(mapVal)
	}

	listVal, ok := v.(parser.ListType)
	if ok {
		return getList(listVal)
	}

	stringVal, ok := v.(string)
	if ok {
		return stringVal
	}

	return ""
}

func getMap(val parser.MapType) string {
	s := ""
	s += fmt.Sprintf("map[%s]", val.Key)
	s += getType(val.Val)

	return s
}

func getList(val parser.ListType) string {
	s := "[]"
	s += getType(val.Ele)

	return s
}
