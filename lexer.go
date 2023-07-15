package ungrammar

import (
	"fmt"
	"strings"
)

type Kind uint

const (
	KIND_NODE = iota + 1
	KIND_TOKEN
	KIND_EQ
	KIND_STAR
	KIND_PIPE
	KIND_QMARK
	KIND_COLON
	KIND_LPAREN
	KIND_RPAREN
	KIND_EOF
	KIND_INVALID
)

func (k Kind) String() string {
	switch k {
	case KIND_NODE:
		return "Node"
	case KIND_TOKEN:
		return "Token"
	case KIND_EQ:
		return "="
	case KIND_STAR:
		return "*"
	case KIND_PIPE:
		return "|"
	case KIND_QMARK:
		return "?"
	case KIND_COLON:
		return ":"
	case KIND_LPAREN:
		return "("
	case KIND_RPAREN:
		return ")"
	case KIND_EOF:
		return "EOF"
	case KIND_INVALID:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Kind Kind
	// Loc   Location
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("[%s]: %s", t.Kind.String(), t.Value)
}

func (t Token) IsValid() bool {
	return t.Kind != KIND_EOF && t.Kind != KIND_INVALID
}

type Location struct {
	Line   int
	Column int
}

func (l *Location) Advance(input string) {
	index := strings.IndexByte(input, '\n')
	if index == -1 {
		l.Column += len(input)
	}
	l.Line += 1
}

func Tokenize(input string) []Token {
	var ans []Token
	for input != "" {
		skipWhiteSpace(&input)
		skipComment(&input)
		// fmt.Println(input)
		token := Advance(&input)
		fmt.Println(token)
		if token.IsValid() {

			ans = append(ans, token)
		} else {
			break
		}
	}
	return ans
}

func Advance(input *string) Token {
	if input == nil || len(*input) == 0 {
		return Token{
			Kind: KIND_EOF,
		}
	}
	b := (*input)[0]
	*input = (*input)[1:]
	switch b {
	case '=':
		return Token{
			Kind: KIND_EQ,
		}
	case '*':
		return Token{
			Kind: KIND_STAR,
		}
	case '?':
		return Token{
			Kind: KIND_QMARK,
		}
	case '(':
		return Token{
			Kind: KIND_LPAREN,
		}
	case ')':
		return Token{
			Kind: KIND_RPAREN,
		}
	case '|':
		return Token{
			Kind: KIND_PIPE,
		}
	case ':':
		return Token{
			Kind: KIND_COLON,
		}
	case '\'':
		var buf []byte
	loop:
		for {
			nextByte := next(input)
			if nextByte == nil {
				return Token{
					Kind: KIND_INVALID,
				}
			}
			switch *nextByte {
			case '\\':
				nextByte = next(input)
				if nextByte == nil {
					return Token{
						Kind: KIND_INVALID,
					}
				}
				buf = append(buf, *nextByte)
			case '\'':
				break loop
			default:
				buf = append(buf, *nextByte)
			}
		}
		return Token{
			Kind:  KIND_TOKEN,
			Value: string(buf),
		}
	default:
		if isIdent(b) {
			var buf []byte
			buf = append(buf, b)
			for peekByte := peek(input); peekByte != nil && isIdent(*peekByte); peekByte = peek(input) {
				buf = append(buf, *next(input))
			}
			return Token{
				Kind:  KIND_NODE,
				Value: string(buf),
			}
		} else {
			return Token{
				Kind: KIND_INVALID,
			}
		}
	}
}

func next(s *string) *byte {
	if s == nil || len(*s) == 0 {
		return nil
	}
	b := (*s)[0]
	*s = (*s)[1:]
	return &b
}

func peek(s *string) *byte {
	if s == nil || len(*s) == 0 {
		return nil
	}
	b := (*s)[0]
	return &b
}

func isIdent(b byte) bool {
	return (b == '_') || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func skipWhiteSpace(s *string) {
	*s = strings.TrimSpace(*s)
}

func skipComment(s *string) {
	if s == nil {
		return
	}
	// maybe there are multiple lines of comments
	for {
		*s = strings.TrimSpace(*s)
		if !strings.HasPrefix(*s, "//") {
			break
		}
		index := strings.IndexByte(*s, '\n')
		if index == -1 {
			index = len(*s) - 1
		}
		*s = (*s)[index+1:]
	}
}
