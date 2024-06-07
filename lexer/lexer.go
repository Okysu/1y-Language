package lexer

import (
	"1ylang/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	inComment    bool // flag to indicate if inside a multi-line comment
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character in the input and advances the position in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken returns the next token in the input
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	if l.inComment {
		l.skipMultiLineComment()
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.SLASH_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.inComment = true
			l.readChar() // consume '*'
			l.readChar() // move to next character
			l.skipMultiLineComment()
			return l.NextToken()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.ASTERISK_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.POW, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MODULUS_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MODULUS, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND_AND, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.AND, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR_OR, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.OR, l.ch)
		}
	case '^':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.XOR_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.XOR, l.ch)
		}
	case '~':
		tok = newToken(token.TILDE, l.ch)
	case '>':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.Token{Type: token.SHR_ASSIGN, Literal: string(ch) + string(l.ch)}
			} else {
				tok = token.Token{Type: token.SHR, Literal: string(ch) + string(l.ch)}
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '<':
		if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.Token{Type: token.SHL_ASSIGN, Literal: string(ch) + string(l.ch)}
			} else {
				tok = token.Token{Type: token.SHL, Literal: string(ch) + string(l.ch)}
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '.':
		if isDigit(l.peekChar()) {
			tok.Type = token.FLOAT
			tok.Literal = l.readFloatFromPosition(l.position)
			return tok
		} else {
			tok = newToken(token.DOT, l.ch)
		}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Type = token.LookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token with the given type and literal
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// readIdentifier reads an identifier from the input
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// skipWhitespace skips whitespace characters in the input
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || l.inComment {
		if l.inComment {
			l.skipMultiLineComment()
		} else {
			l.readChar()
		}
	}
}

// readNumber reads a number from the input
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	// Check for fractional part or scientific notation
	if l.ch == '.' || l.ch == 'e' || l.ch == 'E' {
		return l.readFloatFromPosition(position)
	}
	return l.input[position:l.position]
}

// readFloatFromPosition reads a float from the input starting at a given position
func (l *Lexer) readFloatFromPosition(position int) string {
	if l.ch == '.' {
		l.readChar()
	}
	for isDigit(l.ch) || l.ch == 'e' || l.ch == 'E' {
		if l.ch == 'e' || l.ch == 'E' {
			l.readChar()
			// After 'e' or 'E', we expect an optional '+' or '-' followed by digits
			if l.ch == '+' || l.ch == '-' {
				l.readChar()
			}
			if !isDigit(l.ch) {
				// Invalid scientific notation, return error
				return ""
			}
		} else {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// peekChar returns the next character in the input without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readString reads a string from the input
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	for l.inComment {
		l.readChar()
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // consume '*'
			l.readChar() // consume '/'
			l.inComment = false
		}
		if l.ch == 0 {
			// Reached EOF without closing comment
			l.inComment = false
			break
		}
	}
}
