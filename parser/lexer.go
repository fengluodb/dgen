package parser

import (
	"bufio"
	"io"
)

type lexer struct {
	row    int // 当前行
	column int // 当前列

	buf    *bufio.Reader
	tokens []token
}

func NewLexer(rd io.Reader) *lexer {
	return &lexer{
		row:    0,
		column: 0,
		buf:    bufio.NewReader(rd),
		tokens: []token{},
	}
}

func (lex *lexer) Scan() error {
	for {
		if err := lex.skipWhitespace(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := lex.nextToken(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (lex *lexer) skipComment() error {
	_, _, err := lex.buf.ReadLine()
	lex.row++
	lex.column = 0

	return err
}

func (lex *lexer) skipWhitespace() error {
	for {
		b, err := lex.buf.Peek(1)
		if err != nil {
			return err
		}

		if b[0] == ' ' || b[0] == '\t' {
			lex.buf.ReadByte()
			lex.column++
		} else if b[0] == '\n' {
			lex.buf.ReadByte()
			lex.column = 0
			lex.row++
		} else {
			return nil
		}
	}
}

func (lex *lexer) nextToken() error {
	s := []byte{}
	for {
		b, err := lex.buf.ReadByte()
		if err != nil {
			return err
		}

		lex.column++

		if b == '=' || b == ',' || b == ';' ||
			b == '(' || b == ')' || b == '{' || b == '}' {
			if len(s) != 0 {
				lex.addToken(string(s))
				s = []byte{}
			}
			lex.addToken(string(b))
			continue
		} else if b == ' ' || b == '\t' {
			if len(s) != 0 {
				lex.addToken(string(s))
			}
			return nil
		} else if b == '\n' {
			if len(s) != 0 {
				lex.addToken(string(s))
			}
			lex.column = 0
			lex.row++
			return nil
		} else if b == '#' {
			if len(s) != 0 {
				lex.addToken(string(s))
			}
			lex.skipComment()
		} else {
			s = append(s, b)
		}
	}
}

func (lex *lexer) addToken(s string) error {
	token := token{
		val:    s,
		row:    lex.row,
		column: lex.column,
	}

	if typ, ok := tokenTypeMap[s]; ok {
		token.typ = typ
	} else if isNum(s) {
		token.typ = T_Num
	} else {
		token.typ = T_Identifier
	}

	lex.tokens = append(lex.tokens, token)
	return nil
}
