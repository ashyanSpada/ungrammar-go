package ungrammar

import (
	"errors"
	"fmt"
)

func Parse(tokens []Token) (*Grammar, error) {
	parser := NewParser(tokens)
	for !parser.isEOF() {
		node(parser)
	}
	return parser.finish()
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:     tokens,
		nodeTable:  make(map[string]RuleNode),
		tokenTable: make(map[string]RuleToken),
	}
}

type Parser struct {
	grammar    Grammar
	tokens     []Token
	nodeTable  map[string]RuleNode
	tokenTable map[string]RuleToken
}

func (p *Parser) peek() *Token {
	return p.peekN(0)
}

func (p *Parser) peekN(n int) *Token {
	if n < len(p.tokens) {
		return &p.tokens[n]
	}
	return nil
}

func (p *Parser) bump() (*Token, error) {
	if len(p.tokens) > 0 {
		tmp := p.tokens[0]
		p.tokens = p.tokens[1:]
		return &tmp, nil
	}
	return nil, errors.New("unexpected EOF")
}

func (p *Parser) expect(kind Kind) error {
	token, err := p.bump()
	if err != nil {
		return err
	}
	if token.Kind != kind {
		return errors.New("unexpected token")
	}
	return nil
}

func (p *Parser) isEOF() bool {
	return len(p.tokens) == 0
}

func (p *Parser) finish() (*Grammar, error) {
	for _, nodeData := range p.grammar.Nodes {
		if nodeData.Rule.IsDummy() {
			// return nil, fmt.Errorf("unexpected node: %s", nodeData.Name)
		}
	}
	return &p.grammar, nil
}

func (p *Parser) internNode(name string) RuleNode {
	if node, ok := p.nodeTable[name]; ok {
		return node
	}
	node := RuleNode(len(p.nodeTable))
	p.grammar.Nodes = append(p.grammar.Nodes, NodeData{
		Name: name,
		Rule: DUMMY_RULE,
	})
	p.nodeTable[name] = node
	return node
}

func (p *Parser) internToken(name string) RuleToken {
	if token, ok := p.tokenTable[name]; ok {
		return token
	}
	token := RuleToken(len(p.tokenTable))
	p.grammar.Tokens = append(p.grammar.Tokens, TokenData{name})
	p.tokenTable[name] = token
	return token
}

func node(p *Parser) error {
	token, err := p.bump()
	if err != nil {
		return err
	}

	if token.Kind != KIND_NODE {
		return errors.New("expect ident")
	}

	node := p.internNode(token.Value)
	if err := p.expect(KIND_EQ); err != nil {
		return err
	}

	if !p.grammar.Node(node).Rule.IsDummy() {
		return fmt.Errorf("duplicate rule: %s", p.grammar.Node(node).Name)
	}

	rule, err := rule(p)
	if err != nil {
		return err
	}
	p.grammar.Nodes[node].Rule = *rule
	return nil
}

func rule(p *Parser) (*Rule, error) {
	token := p.peek()
	if token == nil {
		return nil, errors.New("empty rule")
	}
	lhs, err := seqRule(p)
	if err != nil {
		return nil, err
	}
	alt := []Rule{*lhs}
	for token := p.peek(); token != nil; token = p.peek() {
		if token.Kind != KIND_PIPE {
			break
		}
		p.bump()
		rule, err := seqRule(p)
		if err != nil {
			return nil, err
		}
		alt = append(alt, *rule)
	}
	if len(alt) == 1 {
		return &alt[0], nil
	}
	return &Rule{
		Alt: alt,
	}, nil
}

func seqRule(p *Parser) (*Rule, error) {
	var rule *Rule
	var err error
	rule, err = atomRule(p)
	if err != nil {
		return nil, err
	}

	seq := []Rule{*rule}
	for {
		rule, err = optAtomRule(p)
		if err != nil {
			return nil, err
		}
		if rule != nil {
			seq = append(seq, *rule)
		} else {
			break
		}
	}
	if len(seq) == 1 {
		return &seq[0], nil
	}
	return &Rule{
		Seq: seq,
	}, nil
}

func atomRule(p *Parser) (*Rule, error) {
	rule, err := optAtomRule(p)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("unexpected token")
	}
	return rule, nil
}

func optAtomRule(p *Parser) (*Rule, error) {
	token := p.peek()
	if token == nil {
		return nil, nil
	}
	var tmp *Rule
	var err error
	switch token.Kind {
	case KIND_NODE:
		lookAhead := p.peekN(1)
		if lookAhead != nil {
			switch lookAhead.Kind {
			case KIND_EQ:
				return nil, nil
			case KIND_COLON:
				label := token.Value
				p.bump()
				p.bump()
				rule, err := atomRule(p)
				if err != nil {
					return nil, err
				}
				return &Rule{
					Labeled: &LabelRule{
						label, *rule,
					},
				}, nil
			}
		}
		p.bump()
		node := p.internNode(token.Value)
		tmp = &Rule{
			Node: &node,
		}
	case KIND_TOKEN:
		p.bump()
		t := p.internToken(token.Value)
		tmp = &Rule{
			Token: &t,
		}
	case KIND_LPAREN:
		p.bump()
		tmp, err = rule(p)
		if err != nil {
			return nil, err
		}
		err = p.expect(KIND_RPAREN)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	token = p.peek()
	if token != nil {
		switch token.Kind {
		case KIND_QMARK:
			p.bump()
			tmp = &Rule{
				Opt: tmp,
			}
		case KIND_STAR:
			p.bump()
			tmp = &Rule{
				Rep: tmp,
			}
		}
	}
	return tmp, nil
}
