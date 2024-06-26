package ungrammar

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data, err := os.ReadFile("./expression_engine.ungram")
	assert.Nil(t, err)
	tokens := Tokenize(string(data))
	// fmt.Println(tokens)
	grammar, err := Parse(tokens)
	assert.Nil(t, err)
	// fmt.Println(*grammar, err)
	tmp, err := json.MarshalIndent(*grammar, "", "  ")
	assert.Nil(t, err)
	fmt.Println(string(tmp))
}
