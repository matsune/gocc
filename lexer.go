package main

type Lexer struct {
	scanner *Scanner
}

func NewLexer(source []byte) *Lexer {
	return &Lexer{scanner: NewScanner(source)}
}

func (l *Lexer) Pos() Position {
	return l.scanner.Pos()
}

func (l *Lexer) Next() *Token {
	t := NewToken()
	if l.scanner.IsEnd() {
		t.Kind = EOF
		t.Pos = l.scanner.Pos()
		return t
	}

	c := l.skipSpace()
	pos := l.scanner.Pos()

	if isAlpha(c) || c == '_' {
		l.parseAlpha(t)
	} else if isDigit(c) {
		l.parseNumber(t)
	} else if k, ok := checkSingleToken(c); ok {
		t.Kind = k
		t.Str = append(t.Str, c)
		l.consume()
	} else if isDoubleQuote(c) {
		l.parseString(t)
	} else if isSingleQuote(c) {
		l.parseChar(t)
	} else if isOperator(c) {
		l.parseOperator(t)
	} else if isPeriod(c) {
		l.parsePeriod(t)
	} else {
		t.Kind = EOF
	}
	if t.Kind == COMMENT {
		return l.Next()
	}
	t.Pos = pos
	return t
}

func (l *Lexer) Reset(pos Position) {
	l.scanner.Reset(pos)
}

func (l *Lexer) consume() (byte, bool) {
	l.scanner.Step()
	if l.scanner.IsEnd() {
		return 0, false
	}
	return l.scanner.Get(), true
}

func (l *Lexer) parseAlpha(t *Token) {
	var s []byte
	c := l.scanner.Get()
	ok := false
	for isAlpha(c) || isDigit(c) || c == '_' {
		s = append(s, c)

		if c, ok = l.consume(); !ok {
			break
		}
	}
	t.Str = s
	t.Kind = checkKeyword(string(t.Str))
}

func checkKeyword(s string) TokenKind {
	if v, ok := TypeKeys[s]; ok {
		return v
	}
	if v, ok := Keywords[s]; ok {
		return v
	}
	return IDENT
}

func (l *Lexer) parseNumber(t *Token) {
	var s []byte
	c := l.scanner.Get()
	var ok bool
	for isDigit(c) {
		s = append(s, c)
		if c, ok = l.consume(); !ok {
			break
		}
	}
	t.Str = s
	t.Kind = INT_CONST
}

func (l *Lexer) skipSpace() byte {
	c := l.scanner.Get()
	ok := false
	for isWhitespace(c) || isReturn(c) {
		if c, ok = l.consume(); !ok {
			break
		}
	}
	return c
}

func isWhitespace(c byte) bool {
	return (c == ' ' || c == '\t' || c == '\f' || c == '\r')
}

func isReturn(c byte) bool {
	return c == '\n'
}

func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isDigit(c byte) bool {
	return ('0' <= c && c <= '9')
}

func checkSingleToken(c byte) (TokenKind, bool) {
	if v, ok := SingleTokens[c]; ok {
		return v, ok
	}
	return EOF, false
}

func isSingleQuote(c byte) bool {
	return c == '\''
}

func (l *Lexer) parseChar(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}
	if !isSingleQuote(c) {
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	}
	if !isSingleQuote(c) {
		panic("parseChar")
	}
	t.Kind = CHAR_CONST
	if c, ok = l.consume(); !ok {
		return
	}
}

func isDoubleQuote(c byte) bool {
	return c == '"'
}

func (l *Lexer) parseString(t *Token) {
	t.Str = l.readString()
	t.Kind = STRING_CONST
}

func (l *Lexer) readString() []byte {
	ok := false
	c := l.scanner.Get()

	if c != '"' {
		panic("string must start with double quote")
	}

	var s []byte

	if c, ok = l.consume(); !ok {
		return s
	}

	for c != '"' {
		s = append(s, l.scanner.Get())
		if c, ok = l.consume(); !ok {
			break
		}
	}

	l.consume()

	return s
}

func isOperator(c byte) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%' || c == '=' || c == '<' || c == '>' || c == '|' || c == '&' || c == '!' || c == '^'
}

func (l *Lexer) readAND(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}

	switch c {
	case '&': // &&
		t.Kind = LAND
		t.Str = append(t.Str, c)
		l.scanner.Step()
	case '=': // &=
		t.Kind = AND_ASSIGN
		t.Str = append(t.Str, c)
		l.scanner.Step()
	default: // &
		t.Kind = AND
	}
}

func (l *Lexer) readOR(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}

	switch c {
	case '|': // ||
		t.Kind = LOR
		t.Str = append(t.Str, c)
		l.scanner.Step()
	case '=': // |=
		t.Kind = OR_ASSIGN
		t.Str = append(t.Str, c)
		l.scanner.Step()
	default: // |
		t.Kind = OR
	}
}

func (l *Lexer) readDIV(t *Token) {
	var c byte
	var ok bool

	if c, ok = l.consume(); !ok {
		return
	}

	switch c {
	case '/': // //comment
		// parse until end of line
		for c != '\n' && !l.scanner.IsEnd() {
			t.Str = append(t.Str, c)
			if c, ok = l.consume(); !ok {
				break
			}
		}
		t.Kind = COMMENT
	case '*': // /* comment */
		t.Str = append(t.Str, c)
		l.scanner.Step()
		prevC := l.scanner.Get()
		if c, ok = l.consume(); !ok {
			return
		}

		t.Str = append(t.Str, prevC)

		for !(prevC == '*' && c == '/') {
			t.Str = append(t.Str, c)
			prevC = c
			if c, ok = l.consume(); !ok {
				return
			}
		}
		t.Str = append(t.Str, c)

		t.Kind = COMMENT
		l.scanner.Step()
	case '=': // /=
		t.Kind = DIV_ASSIGN
		t.Str = append(t.Str, c)
		l.scanner.Step()
	default: // /
		t.Kind = DIV
	}
}

func (l *Lexer) readLShift(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}
	switch c {
	case '<':
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
		if c == '=' {
			t.Str = append(t.Str, c)
			t.Kind = LEFT_ASSIGN
			if c, ok = l.consume(); !ok {
				return
			}
		} else {
			t.Kind = LSHIFT
		}
	case '=':
		t.Kind = LE
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	default: // <
		t.Kind = LT
	}
}

func (l *Lexer) readRShift(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}
	switch c {
	case '>':
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
		if c == '=' {
			t.Str = append(t.Str, c)
			t.Kind = RIGHT_ASSIGN
			if c, ok = l.consume(); !ok {
				return
			}
		} else {
			t.Kind = RSHIFT
		}
	case '=':
		t.Kind = GE
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	default:
		t.Kind = GT
	}
}

func (l *Lexer) readADD(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}
	switch c {
	case '+':
		t.Kind = INC
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	case '=':
		t.Kind = ADD_ASSIGN
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	default:
		t.Kind = ADD
	}
}

func (l *Lexer) readSUB(t *Token) {
	var c byte
	var ok bool
	if c, ok = l.consume(); !ok {
		return
	}
	switch c {
	case '-':
		t.Kind = DEC
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	case '=':
		t.Kind = SUB_ASSIGN
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	case '>':
		t.Kind = ARROW
		t.Str = append(t.Str, c)
		if c, ok = l.consume(); !ok {
			return
		}
	default:
		t.Kind = SUB
	}
}

func (l *Lexer) parseOperator(t *Token) {
	t.Str = append(t.Str, l.scanner.Get())

	switch l.scanner.Get() {
	case '+': // + += ++
		l.readADD(t)
	case '-': // - -= -- ->
		l.readSUB(t)
	case '&': // && &= &
		l.readAND(t)
	case '|': // || |= |
		l.readOR(t)
	case '/': // /* // /= /
		l.readDIV(t)
	case '<': // <<= <= << <
		l.readLShift(t)
	case '>': // >>= >= >> >
		l.readRShift(t)
	case '!', '=', '*', '%', '^': // != %= == *= ^= ! % = * ^
		prevC := l.scanner.Get()
		c := prevC
		var ok bool

		if c, ok = l.consume(); !ok {
			switch prevC {
			case '!':
				t.Kind = NOT
			case '=':
				t.Kind = ASSIGN
			case '*':
				t.Kind = MUL
			case '%':
				t.Kind = REM
			case '^':
				t.Kind = XOR
			default:
				break
			}
			return
		}

		if c == '=' {
			switch prevC {
			case '!':
				t.Kind = NE
			case '=':
				t.Kind = EQ
			case '*':
				t.Kind = MUL_ASSIGN
			case '%':
				t.Kind = REM_ASSIGN
			case '^':
				t.Kind = XOR_ASSIGN
			default:
				break
			}
			t.Str = append(t.Str, c)
			if c, ok = l.consume(); !ok {
				return
			}
		} else {
			switch prevC {
			case '!':
				t.Kind = NOT
			case '=':
				t.Kind = ASSIGN
			case '*':
				t.Kind = MUL
			case '%':
				t.Kind = REM
			case '^':
				t.Kind = XOR
			default:
				break
			}
		}
	}
}

func isPeriod(c byte) bool {
	return c == '.'
}

func (l *Lexer) parsePeriod(t *Token) {
	c := l.scanner.Get()
	ok := false

	var s []byte
	s = append(s, c)

	if c, ok = l.consume(); !ok {
		return
	}

	if isPeriod(c) {
		s = append(s, c)
		if c, ok = l.consume(); !ok {
			return
		}
		if !isPeriod(c) {
			panic("parsePeriod")
		}
		s = append(s, c)
		l.scanner.Step()
		t.Str = s
		t.Kind = ELLIPSIS
	} else if isDigit(c) {
		//  - TODO:
		panic("unimplemented parsePeriod")
	} else {
		t.Str = s
		t.Kind = PERIOD
	}
}
