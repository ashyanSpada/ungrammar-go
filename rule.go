package ungrammar

type RuleToken int

type RuleNode int

type NodeData struct {
	Name string
	Rule Rule
}

type TokenData struct {
	Name string
}

type Grammar struct {
	nodes  []NodeData
	tokens []TokenData
}

func (g Grammar) Node(node RuleNode) NodeData {
	return g.nodes[node]
}

func (g Grammar) Token(token RuleToken) TokenData {
	return g.tokens[token]
}

type LabelRule struct {
	Label string
	Rule  Rule
}

type Rule struct {
	Labeled *LabelRule
	Node    *RuleNode
	Token   *RuleToken
	Seq     []Rule
	Alt     []Rule
	Opt     *Rule
	Rep     *Rule
}

func (r Rule) IsDummy() bool {
	return r.Labeled == nil && r.Node == nil && r.Token == nil && r.Seq == nil && r.Alt == nil && r.Opt == nil && r.Rep == nil
}

var DUMMY_RULE = Rule{}
