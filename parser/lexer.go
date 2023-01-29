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
		buf: bufio.NewReader(rd),
	}
}

func (lex *lexer) Scan() error {
	err := lex.scanToken()
	if err == io.EOF {
		return nil
	}
	return err
}

func (lex *lexer) skipComment() error {
	_, _, err := lex.buf.ReadLine()
	lex.row++
	lex.column = 0

	return err
}

func (lex *lexer) scanToken() error {
	s := []byte{}
	for {
		b, err := lex.buf.ReadByte()
		if err != nil {
			return err
		}

		if _, ok := tokenTypeMap[string(b)]; ok {
			if len(s) != 0 {
				lex.addToken(string(s))
				s = []byte{}
			}
			lex.addToken(string(b))
		} else if _, ok := ignoreCharMap[b]; ok {
			if len(s) != 0 {
				lex.addToken(string(s))
				s = []byte{}
			}
			switch b {
			case '\n':
				lex.column = 0
				lex.row++
			case '#':
				lex.skipComment()
			default:
				lex.column++
			}
		} else {
			s = append(s, b)
		}
	}
}

func (lex *lexer) addToken(s string) {
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

	lex.column += len(s)
	lex.tokens = append(lex.tokens, token)
}
