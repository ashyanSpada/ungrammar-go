package ungrammar

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	input := "Grammer = Node *\n Node = name:'ident' '=' Rule"
	ans := Tokenize(input)
	fmt.Println(ans)
}

func TestAdvance(T *testing.T) {
	input := "Grammer"
	ans := Advance(&input)
	fmt.Println(ans)
}
