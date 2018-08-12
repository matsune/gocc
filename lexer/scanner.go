package lexer

import "gocc/token"

type Scanner struct {
	source []byte
	pos    token.Position
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{
		source: source,
		pos: token.Position{
			Line:   1,
			Column: 1,
			Offset: 0,
		},
	}
}

func (s *Scanner) Pos() token.Position {
	return s.pos
}

func (s *Scanner) Get() byte {
	return s.source[s.pos.Offset]
}

func (s *Scanner) Step() {
	if s.Get() == '\n' {
		s.pos.Line++
		s.pos.Column = 1
	} else {
		s.pos.Column++
	}
	s.pos.Offset++
}

func (s *Scanner) IsEnd() bool {
	return int(s.pos.Offset) >= len(s.source)
}

func (s *Scanner) Reset(pos token.Position) {
	s.pos = pos
}
