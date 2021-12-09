package graphproc

import (
	"bufio"
	"bytes"
	"io"
)

// Graph processor description language
// grap1hname:a|b;b|d|k;b|c|f|h|k;b|e|g;g|i|l;g|j|l;k|l|m{"k1":"v1", "k2":"v2"};graph2name:x|y|z

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS
	// Literals
	IDENT // vertex names
	// Misc characters
	SEMICOLON    // ;
	COLON        //	:
	PERIOD       // .
	UNDERSCORE   // _
	COMMA        //	,
	PIPE         //	|
	DOUBLEQUOTE  //	"
	LEFTBRACE    // {
	RIGHTBRACE   //	}
	LEFTBRACKET  // [
	RIGHTBRACKET //	]
	ESCAPE       // \
)

var eof = rune(0)

// Scanner : lexical scanner
type Scanner struct {
	r *bufio.Reader
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

// NewScanner : create a new instance of Scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read : Reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned)
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread : places the previously read rune back on the reader
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// Scan : returns the next token and literal value
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace
	// If we see a letter then consume as an ident or reserved word
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character
	switch ch {
	case eof:
		return EOF, ""
	case ';':
		return SEMICOLON, string(ch)
	case '.':
		return PERIOD, string(ch)
	case '_':
		return UNDERSCORE, string(ch)
	case ':':
		return COLON, string(ch)
	case '|':
		return PIPE, string(ch)
	case '{':
		return LEFTBRACE, string(ch)
	case '}':
		return RIGHTBRACE, string(ch)
	case '[':
		return LEFTBRACKET, string(ch)
	case ']':
		return RIGHTBRACKET, string(ch)
	case ',':
		return COMMA, string(ch)
	case '\\':
		return ESCAPE, string(ch)
	case '"':
		return DOUBLEQUOTE, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace : consumes the current rune and all contiguous whitespace
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer
	// Non-whitespace characters and EOF will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent : consumes the current rune and all contiguous ident runes
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer
	// Non-ident characters and EOF will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}
