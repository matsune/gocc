package lexer

import (
	"gocc/token"
	"testing"
)

func getExpect(t *testing.T, s *Scanner, c byte) {
	if s.Get() != c {
		t.Errorf("expected char is %c, but got %c", c, s.Get())
	}
}

func posExpect(t *testing.T, s *Scanner, line, col, offset int) {
	p := s.Pos()
	if p.Line != line {
		t.Errorf("expected line is %d, but got %d", line, p.Line)
	}
	if p.Column != col {
		t.Errorf("expected column is %d, but got %d", col, p.Column)
	}
	if p.Offset != offset {
		t.Errorf("expected offset is %d, but got %d", offset, p.Offset)
	}
}

func isEndExpect(t *testing.T, s *Scanner, b bool) {
	if s.IsEnd() != b {
		t.Errorf("expected IsEnd is wrong")
	}
}

func TestScanner(t *testing.T) {
	str := `ab 2
cde`
	s := NewScanner([]byte(str))

	getExpect(t, s, 'a')
	posExpect(t, s, 1, 1, 0)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, 'b')
	posExpect(t, s, 1, 2, 1)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, ' ')
	posExpect(t, s, 1, 3, 2)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, '2')
	posExpect(t, s, 1, 4, 3)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, '\n')
	posExpect(t, s, 1, 5, 4)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, 'c')
	posExpect(t, s, 2, 1, 5)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, 'd')
	posExpect(t, s, 2, 2, 6)
	isEndExpect(t, s, false)
	s.Step()
	getExpect(t, s, 'e')
	posExpect(t, s, 2, 3, 7)
	isEndExpect(t, s, false)
	s.Step()
	isEndExpect(t, s, true)

	s.Reset(token.Position{Line: 1, Column: 1, Offset: 0})
	getExpect(t, s, 'a')
}
