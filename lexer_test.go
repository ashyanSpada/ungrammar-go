package ungrammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	input := "Grammar = Node *\n Node = name:'ident' '=' Rule"
	ans := Tokenize(input)
	assert.Equal(t, len(ans), 11)
	assert.Equal(t, ans[0], Token{KIND_NODE, "Grammar"})
	assert.Equal(t, ans[1], Token{KIND_EQ, ""})
	assert.Equal(t, ans[2], Token{KIND_NODE, "Node"})
	assert.Equal(t, ans[3], Token{KIND_STAR, ""})
	assert.Equal(t, ans[4], Token{KIND_NODE, "Node"})
	assert.Equal(t, ans[5], Token{KIND_EQ, ""})
	assert.Equal(t, ans[6], Token{KIND_NODE, "name"})
	assert.Equal(t, ans[7], Token{KIND_COLON, ""})
	assert.Equal(t, ans[8], Token{KIND_TOKEN, "ident"})
	assert.Equal(t, ans[9], Token{KIND_TOKEN, "="})
	assert.Equal(t, ans[10], Token{KIND_NODE, "Rule"})
}

func TestSkipComment(t *testing.T) {
	s := "//hahahaha  "
	skipComment(&s)
	assert.Equal(t, s, "")
	s = "//jajjajjaja\nhaha"
	skipComment(&s)
	assert.Equal(t, s, "haha")
	// multiple comments
	s = "  //jajajaja\n  //dsfadsfdsa\nhha"
	skipComment(&s)
	assert.Equal(t, s, "hha")
}
