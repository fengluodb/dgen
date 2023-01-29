package parser

import (
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	EnumStats    []enumStat
	MessageStats []messageStat
	ServiceStats []serviceStat

	lexer *lexer
	cur   int // 目前解析到的token数
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{
		lexer: NewLexer(rd),
	}
}

func (p *Parser) Parse() error {
	if err := p.lexer.Scan(); err != nil {
		return err
	}

	tokens := p.lexer.tokens
	for p.cur < len(tokens) {
		token := tokens[p.cur]
		p.cur++
		var err error
		switch token.typ {
		case T_Enum:
			err = p.parseEnum()
		case T_Message:
			err = p.parseMessage()
		case T_Service:
			err = p.parseService()
		default:
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseEnum() error {
	tokens := p.lexer.tokens
	es := enumStat{}

	token := tokens[p.cur]
	if token.typ != T_Identifier || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	es.name = token.val
	p.cur++

	token = tokens[p.cur]
	if token := tokens[p.cur]; token.typ != T_LCurlyBracket || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	for p.cur < len(tokens) {
		token = tokens[p.cur]
		if token.typ != T_Identifier {
			break
		}
		es.members = append(es.members, token.val)
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_Comma {
			break
		}
		p.cur++
	}

	token = tokens[p.cur]
	if token := tokens[p.cur]; token.typ != T_RCurlyBracket || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	p.EnumStats = append(p.EnumStats, es)
	return nil
}

func (p *Parser) parseMessage() error {
	tokens := p.lexer.tokens
	ms := messageStat{}

	token := tokens[p.cur]
	if token.typ != T_Identifier || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	ms.name = token.val
	p.cur++

	token = tokens[p.cur]
	if token := tokens[p.cur]; token.typ != T_LCurlyBracket || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	for p.cur < len(tokens) {
		m := messageMember{}

		token := tokens[p.cur]
		if token.typ == T_RCurlyBracket {
			p.cur++
			break
		}
		if token.typ == T_Optional {
			p.cur++
			m.optional = true
			token = tokens[p.cur]
		}
		if token.typ != T_Seq {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_Assign {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_Num {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		m.seq, _ = strconv.Atoi(token.val)
		p.cur++

		token = tokens[p.cur]
		typ, err := p.parseType()
		m.typ = typ
		if err != nil {
			return err
		}

		token = tokens[p.cur]
		if token.typ != T_Identifier {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		m.name = token.val
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_Semicolon {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		p.cur++

		ms.members = append(ms.members, m)
	}

	p.MessageStats = append(p.MessageStats, ms)
	return nil
}

func (p *Parser) parseService() error {
	tokens := p.lexer.tokens
	ss := serviceStat{}

	token := tokens[p.cur]
	if token.typ != T_Identifier || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	ss.name = token.val
	p.cur++

	token = tokens[p.cur]
	if token := tokens[p.cur]; token.typ != T_LCurlyBracket || p.cur >= len(tokens) {
		return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	for p.cur < len(tokens) {
		m := serviceMember{}

		token = tokens[p.cur]
		if token.typ == T_RCurlyBracket {
			p.cur++
			break
		}
		if token.typ != T_Identifier {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		m.name = token.val
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_LSmallBracket {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_Identifier {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		m.req = token.val
		p.cur++

		token = tokens[p.cur]
		if token.typ != T_RSmallBracket {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		p.cur++

		token = tokens[p.cur]
		if token.typ == T_Return {
			p.cur++

			token = tokens[p.cur]
			if token.typ != T_LSmallBracket {
				return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
			}
			p.cur++

			token = tokens[p.cur]
			if token.typ != T_Identifier {
				return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
			}
			m.resp = token.val
			p.cur++

			token = tokens[p.cur]
			if token.typ != T_RSmallBracket {
				return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
			}
			p.cur++

			token = tokens[p.cur]
		}
		if token.typ != T_Semicolon {
			return fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
		}
		p.cur++

		ss.members = append(ss.members, m)
	}

	p.ServiceStats = append(p.ServiceStats, ss)
	return nil

}

func (p *Parser) parseType() (interface{}, error) {
	tokens := p.lexer.tokens

	token := tokens[p.cur]
	p.cur++
	if token.typ == T_Builtin {
		if token.val == "map" {
			return p.parseMap()
		}
		if token.val == "list" {
			return p.parseList()
		}
		return token.val, nil
	} else if token.typ == T_Identifier {
		return token.val, nil
	} else {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
}

func (p *Parser) parseMap() (interface{}, error) {
	tokens := p.lexer.tokens
	m := mapType{}

	token := tokens[p.cur]
	if token.typ != T_LBracket {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	token = tokens[p.cur]
	if token.typ != T_Builtin || token.val == "map" || token.val == "list" {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	m.key = token.val
	p.cur++

	token = tokens[p.cur]
	if token.typ != T_RBracket {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	token = tokens[p.cur]
	if token.typ != T_LBracket {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	token = tokens[p.cur]
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	m.val = typ

	token = tokens[p.cur]
	if token.typ != T_RBracket {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	return m, nil
}

func (p *Parser) parseList() (interface{}, error) {
	tokens := p.lexer.tokens
	l := listType{}

	token := tokens[p.cur]
	if token.typ != T_LBracket {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	token = tokens[p.cur]
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	l.ele = typ

	token = tokens[p.cur]
	if token.typ != T_RBracket {
		return nil, fmt.Errorf("raw:%d, column:%d is invalid grammar", token.row, token.column)
	}
	p.cur++

	return l, nil
}
