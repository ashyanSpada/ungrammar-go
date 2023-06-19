package ungrammar

import "strings"

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

type Token struct {
	Kind  Kind
	Loc   Location
	Value string
}

type Location struct {
	Line   int
	Column int
}

func (l Location) Advance(input string) {
	index := strings.IndexByte(input, '\n')
	if index == -1 {
		l.Column += len(input)
	}
	l.Line += 1
}

func Advance(input *string) Token {
	if input == nil || len(*input) == 0 {
		return Token{
			Kind: KIND_EOF,
		}
	}
	*input = (*input)[1:]
	b := (*input)[0]
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
			for nextByte := next(input); nextByte != nil && isIdent(*nextByte); nextByte = next(input) {
				buf = append(buf, *nextByte)
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

func isEscapable(b byte) bool {
	return b == '\\' || b == '\''
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n'
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
	if strings.HasPrefix(*s, "//") {
		index := strings.IndexByte(*s, '\n')
		if index == -1 {
			index = len(*s) - 1
		}
		*s = (*s)[index+1:]
	}
}
