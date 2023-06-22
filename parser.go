package ungrammar

import (
	"errors"
	"fmt"
)

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
	return nil, errors.New("Unexpected EOF")
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

func (p *Parser) finish() error {
	for _, nodeData := range p.grammar.nodes {
		if nodeData.Rule.IsDummy() {
			return fmt.Errorf("unexpected node: %s", nodeData.Name)
		}
	}
	return nil
}

func (p *Parser) internNode(name string) RuleNode {
	if node, ok := p.nodeTable[name]; ok {
		return node
	}
	node := RuleNode(len(p.nodeTable))
	p.grammar.nodes = append(p.grammar.nodes, NodeData{
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
	p.grammar.tokens = append(p.grammar.tokens, TokenData{name})
	p.tokenTable[name] = token
	return token
}

func (p *Parser) node() error {
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
		return fmt.Errorf("duplicate rule: {}", p.grammar.Node(node).Name)
	}

	rule, err := p.rule()
	if err != nil {
		return err
	}
	p.grammar.nodes[node].Rule = *rule
	return nil
}

func (p *Parser) rule() (*Rule, error) {
	token := p.peek()
	if token == nil {
		return nil, errors.New("empty rule")
	}

}
