package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseEnum(t *testing.T) {
	s := `
	enum color {
		red, green, blue, alpha
	}

	enum fruit {
		orange, banana, tomato,
	}

	enum xxx {
	}
	`

	parser := NewParser(strings.NewReader(s))
	if err := parser.Parse(); err != nil {
		t.Fatalf("failed to parse enum: %v", err)
	}

	if len(parser.EnumStats) != 3 {
		t.Fatalf("expected 2 enum stats, got %d", len(parser.EnumStats))
	}

	for _, es := range parser.EnumStats {
		fmt.Printf("enum %s {\n", es.name)
		for _, m := range es.members {
			fmt.Printf("\t%s,\n", m)
		}
		fmt.Printf("}\n")
	}
}

func TestParseMessage(t *testing.T) {
	s := `
	message SearchResponse {
		seq=1 string result;
		optional seq=20 list[map[int32][list[map[int32][string]]]] page_number;
		optional seq=3 map[string][list[map[string][string]]] page_number;
	}

	message SearchRequest {
		seq=1 map[string][map[string][string]] query;
		optional seq=2 list[list[int32]] page_number;
		optional seq=3 int32 result_per_page;
	}

	message xxx {

	}
	`

	parser := NewParser(strings.NewReader(s))
	if err := parser.Parse(); err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}

	if len(parser.MessageStats) != 3 {
		t.Fatalf("expected 3 message stats, got %d", len(parser.EnumStats))
	}

	for _, ms := range parser.MessageStats {
		fmt.Printf("message %s {\n", ms.name)
		for _, m := range ms.members {
			fmt.Printf("\t")
			if m.optional {
				fmt.Printf("optional ")
			}
			fmt.Printf("seq=%d ", m.seq)
			PrintTyep(m.typ)
			fmt.Printf(" %s;\n", m.name)
		}
		fmt.Printf("}\n")
	}
}

func TestParseService(t *testing.T) {
	s := `
	service SearchService {
    	Search(SearchRequest) return (SearchResponse);
    	Query(Queryquest);
	}

	service OrderService {
		Order(OrderRequeset);
	}

	service Clear {

	}
	`

	parser := NewParser(strings.NewReader(s))
	if err := parser.Parse(); err != nil {
		t.Fatalf("failed to parse service: %v", err)
	}

	if len(parser.ServiceStats) != 3 {
		t.Fatalf("expected 2 enum stats, got %d", len(parser.EnumStats))
	}

	for _, ss := range parser.ServiceStats {
		fmt.Printf("service %s {\n", ss.name)
		for _, m := range ss.members {
			fmt.Printf("\t%s(%s)", m.name, m.req)
			if m.resp != "" {
				fmt.Printf(" return (%s)", m.resp)
			}
			fmt.Printf(";\n")
		}
		fmt.Printf("}\n")
	}
}

func PrintTyep(v interface{}) {
	mapVal, ok := v.(mapType)
	if ok {
		PrintMap(mapVal)
	}

	listVal, ok := v.(listType)
	if ok {
		PrintList(listVal)
	}

	stringVal, ok := v.(string)
	if ok {
		fmt.Printf("%s", stringVal)
	}
}

func PrintMap(val mapType) {
	fmt.Printf("map[%s][", val.key)
	PrintTyep(val.val)
	fmt.Printf("]")
}

func PrintList(val listType) {
	fmt.Printf("list[")
	PrintTyep(val.ele)
	fmt.Printf("]")
}
