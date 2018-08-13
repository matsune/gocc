package parser

import "gocc/token"

type Stack struct {
	ts    []token.Token
	ps    []token.Position
	count int
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) push(t token.Token, pos token.Position) {
	s.ts = append(s.ts, t)
	s.ps = append(s.ps, pos)
	s.count++
}

func (s *Stack) pop() (*token.Token, token.Position) {
	s.count--
	t, p := &s.ts[s.count], s.ps[s.count]
	s.remove(s.count)
	return t, p
}

func (s *Stack) remove(at int) {
	resT := []token.Token{}

	for i, v := range s.ts {
		if i != at {
			resT = append(resT, v)
		}
	}
	resP := []token.Position{}
	for i, v := range s.ps {
		if i != at {
			resP = append(resP, v)
		}
	}
	s.ts = resT
	s.ps = resP
}
