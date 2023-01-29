package parser

import (
	"io"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	s := `
	enum colors {
		red, green, blue
	 };
	 
	 # test
	 message SearchRequest {#test
		 seq=1 string query;
		 optional seq=2 int32 page_number;
		 optional seq=3 int32 result_per_page;
	 }
	 
	 message SearchResponse {
		 seq=1 string result;
	 }
	 
	 service SearchService {
		 Search(SearchRequest) return (SearchResponse);
	 }

	 service
	`

	lexer := NewLexer(strings.NewReader(s))

	err := lexer.Scan()
	if err != nil && err != io.EOF {
		t.Logf("error:%s", err.Error())
		t.Fail()
	}

	for _, token := range lexer.tokens {
		t.Logf("raw:%d, column%d, tokenType:%d, token.val: %s \n", token.row, token.column, token.typ, token.val)
	}
}
