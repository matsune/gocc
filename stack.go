package main

type Stack struct {
	ts    []Token
	poss  []Position
	count int
}

func (s *Stack) push(t Token, pos Position) {
	s.ts = append(s.ts, t)
	s.poss = append(s.poss, pos)
	s.count++
}

func (s *Stack) pop() (*Token, Position) {
	s.count--
	t, p := &s.ts[s.count], s.poss[s.count]
	s.remove(s.count)
	return t, p
}

func (s *Stack) remove(at int) {
	resT := []Token{}

	for i, v := range s.ts {
		if i != at {
			resT = append(resT, v)
		}
	}
	resP := []Position{}
	for i, v := range s.poss {
		if i != at {
			resP = append(resP, v)
		}
	}
	s.ts = resT
	s.poss = resP
}
