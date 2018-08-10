package gocc

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer *Lexer
	token *Token
	stack *Stack
}

func NewParser(source []byte) *Parser {
	p := &Parser{lexer: NewLexer(source), token: NewToken(), stack: &Stack{}}
	p.next()
	return p
}

func (p *Parser) match(t TokenKind) bool {
	return p.token.Kind == t
}

func (p *Parser) matchs(ts []TokenKind) bool {
	for _, v := range ts {
		if p.match(v) {
			return true
		}
	}
	return false
}

// push current token and position to stack
func (p *Parser) push() {
	pos := p.lexer.Pos()
	p.stack.push(*p.token, pos)
}

// pop last pushed token and position
func (p *Parser) pop() {
	t, pos := p.stack.pop()
	p.token = t
	p.lexer.Reset(pos)
}

var (
	typeKeys = []TokenKind{
		INT,
		VOID,
		CHAR,
		FLOAT,
		LONG,
		SHORT,
		DOUBLE,
		STRUCT,
		UNION,
		SIGNED,
		UNSIGNED,
		STATIC,
		AUTO,
		EXTERN,
		REGISTER,
		CONST,
		VOLATILE,
	}

	unaryOps = []TokenKind{
		AND,
		MUL,
		ADD,
		SUB,
		TILDE,
		NOT,
	}

	assignOps = []TokenKind{
		ASSIGN,
		MUL_ASSIGN,
		DIV_ASSIGN,
		REM_ASSIGN,
		ADD_ASSIGN,
		SUB_ASSIGN,
		LEFT_ASSIGN,
		RIGHT_ASSIGN,
		AND_ASSIGN,
		OR_ASSIGN,
		XOR_ASSIGN,
	}

	storageSpecifiers = []TokenKind{
		AUTO,
		REGISTER,
		STATIC,
		EXTERN,
		TYPEDEF,
	}

	typeSpecifiers = []TokenKind{
		VOID,
		CHAR,
		SHORT,
		INT,
		SHORT,
		LONG,
		FLOAT,
		DOUBLE,
		SIGNED,
		UNSIGNED,
	}
)

func (p *Parser) assert(t TokenKind) {
	if !p.match(t) {
		str := fmt.Sprintf("expected token is '"+t.String()+"', but got '"+p.token.String()+"' at line %d column %d", p.lexer.Pos().Line, p.lexer.Pos().Column)
		panic(str)
	}
}

func (p *Parser) next() {
	p.token = p.lexer.Next()
}

func (p *Parser) IsEnd() bool {
	return p.match(EOF)
}

func (p *Parser) Parse() Node {
	if p.isFuncDef() {
		return p.readFuncDef()
	} else if p.isType() {
		return p.readVarDef()
	} else {
		panic("unexpected")
	}
}

/**
read def
*/

func (p *Parser) readVarDef() Node {
	t := p.readType()

	p.assert(IDENT)
	tok := p.token
	p.next()

	var n Node
	if p.match(LBRACK) {

		s := p.readSubscriptInit()
		arr := ArrayDef{Type: t, Token: tok, Subscript: s}

		if s == nil && !p.match(ASSIGN) {
			panic(fmt.Errorf("definition of variable with array type needs an explicit size or an initializer"))
		}

		if p.match(ASSIGN) {
			p.next()
			init := p.readArrayInit()
			arr.Init = &init
			// obj.IsInit = true
		}
		n = arr
	} else {
		v := VarDef{Type: t, Token: tok}

		if p.match(ASSIGN) {
			p.next()
			e := p.assignExpr()
			v.Init = &e
		}
		n = v
	}

	p.assert(SEMICOLON)
	p.next()

	return n
}

// [0] []
func (p *Parser) readSubscriptInit() *Expr {
	p.assert(LBRACK)
	p.next()

	if p.match(RBRACK) {
		p.next()
		return nil
	}

	e := p.conditionalExpr()

	p.assert(RBRACK)
	p.next()

	return &e
}

func (p *Parser) readArrayInit() ArrayInit {
	p.assert(LBRACE)
	p.next()
	n := ArrayInit{}
	for {
		e := p.assignExpr()
		n.List = append(n.List, e)
		if p.match(RBRACE) {
			break
		} else if p.match(COMMA) {
			p.next()
		} else {
			panic("expected } or ,")
		}
	}
	p.next()
	return n
}

func (p *Parser) isType() bool {
	return p.matchs([]TokenKind{INT, CHAR, VOID, FLOAT, LONG, SHORT, DOUBLE})
}

func (p *Parser) readType() CType {
	var t CType
	for {
		if p.isType() {
			switch p.token.Kind {
			case INT:
				t = C_int
			case CHAR:
				t = C_char
			case VOID:
				t = C_void
			case FLOAT:
				t = C_float
			case LONG:
				t = C_long
			case SHORT:
				t = C_short
			case DOUBLE:
				t = C_double
			default:
				panic("readType")
			}
		} else if p.match(MUL) { // * as pointer
			t = C_pointer
		} else {
			break
		}
		p.next()
	}
	return t
}

func (p *Parser) isFuncDef() bool {
	p.push()
	defer p.pop()

	var isFunc bool
loop:
	for {
		switch {
		case p.match(RPAREN):
			p.next()
			if p.match(LBRACE) {
				isFunc = true
				break loop
			}
		case p.match(SEMICOLON):
			break loop
		default:
			p.next()
		}

		if p.match(EOF) {
			break loop
		}
	}
	return isFunc
}

func (p *Parser) readFuncDef() FuncDef {
	t := p.readType()

	p.assert(IDENT)
	name := string(p.token.Str)
	p.next()

	p.assert(LPAREN)
	p.next()

	args := p.readFuncArgs()

	p.assert(RPAREN)
	p.next()

	block := p.blockStmt()

	return FuncDef{Type: t, Name: name, Args: args, Block: block}
}

func (p *Parser) readFuncArgs() []FuncArg {
	var res []FuncArg
	for p.isType() {
		res = append(res, p.readFuncArg())
		if !p.match(COMMA) {
			break
		}
		p.next()
	}
	return res
}

func (p *Parser) readFuncArg() FuncArg {
	var n FuncArg
	n.Type = p.readType()

	p.assert(IDENT)
	n.Name = p.token
	p.next()

	return n
}

/**
expression
*/

func (p *Parser) expr() Expr {
	return p.assignExpr()
}

func (p *Parser) assignExpr() Expr {
	p.push()

	var hasAssign bool
	for !p.match(SEMICOLON) {
		if p.match(ASSIGN) {
			hasAssign = true
			break
		}
		p.next()

		if p.match(EOF) {
			break
		}
	}

	p.pop()

	if hasAssign {
		L := p.unaryExpr()
		if !p.isAssignOp() {
			panic("expected assign operator")
		}
		op := p.token
		p.next()
		R := p.assignExpr()
		n := AssignExpr{L: L, Op: op, R: R}
		return n
	} else {
		return p.conditionalExpr()
	}
}

func (p *Parser) isAssignOp() bool {
	return p.matchs(assignOps)
}

func (p *Parser) conditionalExpr() Expr {
	e := p.logOrExpr()
	if p.match(QUE) {
		p.next()
		L := p.expr()
		p.assert(COLON)
		p.next()
		n := CondExpr{Cond: e, L: L, R: p.conditionalExpr()}
		return n
	}
	return e
}

func (p *Parser) logOrExpr() Expr {
	e := p.logAndExpr()
	return p.logOrExpr2(e)
}

func (p *Parser) logOrExpr2(e Expr) Expr {
	if p.match(LOR) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.logAndExpr()}
		return p.logOrExpr2(n)
	}
	return e
}

func (p *Parser) logAndExpr() Expr {
	e := p.incOrExpr()
	return p.logAndExpr2(e)
}

func (p *Parser) logAndExpr2(e Expr) Expr {
	if p.match(LAND) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.incOrExpr()}
		return p.logAndExpr2(n)
	}
	return e
}

func (p *Parser) incOrExpr() Expr {
	e := p.excOrExpr()
	return p.incOrExpr2(e)
}

func (p *Parser) incOrExpr2(e Expr) Expr {
	if p.match(OR) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.excOrExpr()}
		return p.incOrExpr2(n)
	}
	return e
}

func (p *Parser) excOrExpr() Expr {
	e := p.andExpr()
	return p.excOrExpr2(e)
}

func (p *Parser) excOrExpr2(e Expr) Expr {
	if p.match(XOR) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.andExpr()}
		return p.excOrExpr2(n)
	}
	return e
}

func (p *Parser) andExpr() Expr {
	e := p.eqExpr()
	return p.andExpr2(e)
}

func (p *Parser) andExpr2(e Expr) Expr {
	if p.match(AND) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.eqExpr()}
		return p.andExpr2(n)
	}
	return e
}

func (p *Parser) eqExpr() Expr {
	e := p.relExpr()
	return p.eqExpr2(e)
}

func (p *Parser) eqExpr2(e Expr) Expr {
	if p.match(EQ) || p.match(NE) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.relExpr()}
		return p.eqExpr2(n)
	}
	return e
}

func (p *Parser) relExpr() Expr {
	e := p.shiftExpr()
	return p.relExpr2(e)
}

func (p *Parser) relExpr2(e Expr) Expr {
	if p.match(LT) || p.match(GT) || p.match(LE) || p.match(GE) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.shiftExpr()}
		return p.relExpr2(n)
	}
	return e
}

func (p *Parser) shiftExpr() Expr {
	e := p.additiveExpr()
	return p.shiftExpr2(e)
}

func (p *Parser) shiftExpr2(e Expr) Expr {
	if p.match(LSHIFT) || p.match(RSHIFT) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.additiveExpr()}
		return p.shiftExpr2(n)
	}
	return e
}

func (p *Parser) additiveExpr() Expr {
	e := p.multiExpr()
	return p.additiveExpr2(e)
}

func (p *Parser) additiveExpr2(e Expr) Expr {
	if p.match(ADD) || p.match(SUB) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.multiExpr()}
		return p.additiveExpr2(n)
	}
	return e
}

func (p *Parser) multiExpr() Expr {
	e := p.castExpr()
	return p.multiExpr2(e)
}

func (p *Parser) multiExpr2(e Expr) Expr {
	if p.match(MUL) || p.match(DIV) || p.match(REM) {
		op := p.token
		p.next()
		n := BinaryExpr{X: e, Op: op, Y: p.castExpr()}
		return p.multiExpr2(n)
	}
	return e
}

func (p *Parser) castExpr() Expr {
	return p.unaryExpr()
}

func (p *Parser) unaryExpr() Expr {
	if p.match(INC) || p.match(DEC) {
		panic("unimplemented unaryExpr")
	} else if p.isUnaryOp() {
		op := p.token
		p.next()

		switch op.Kind {
		case MUL:
			pv := PointerVal{Token: p.token}
			p.next()
			return pv
		case AND:
			av := AddressVal{Token: p.token}
			p.next()
			return av
		default:
			return UnaryExpr{Op: op, Expr: p.castExpr()}
		}
	} else {
		return p.postfixExpr()
	}
}

func (p *Parser) isUnaryOp() bool {
	return p.matchs(unaryOps)
}

func (p *Parser) postfixExpr() Expr {
	n := p.primaryExpr()
	return p.postfixExpr2(n)
}

func (p *Parser) postfixExpr2(e Expr) Expr {
	if p.match(INC) {
		panic("postfix increment")
	} else if p.match(DEC) {
		panic("postfix decrement")
	} else if p.match(LPAREN) {
		switch e.(type) {
		case Ident:
			return p.readFuncCall(e)
		default:
			panic("unimplemented postfixExpr2")
		}
	} else if p.match(PERIOD) {
		panic("postfix .")
	} else if p.match(ARROW) {
		panic("postfix ->")
	} else {
		return e
	}
}

// [0] [1]
func (p *Parser) readSubscriptExpr(t *Token) SubscriptExpr {
	p.assert(LBRACK)
	p.next()

	e := p.conditionalExpr()

	p.assert(RBRACK)
	p.next()

	se := SubscriptExpr{Token: t, Expr: e}
	return se
}

func (p *Parser) primaryExpr() Expr {
	switch {
	case p.match(IDENT):
		t := p.token
		p.next()

		if p.match(LBRACK) {
			return p.readSubscriptExpr(t)
		} else {
			n := Ident{Token: t}
			return n
		}
	case p.match(INT_CONST):
		i, err := strconv.Atoi(p.token.String())
		if err != nil {
			panic(err)
		}
		n := IntVal{Num: i}
		p.next()
		return n
	case p.match(CHAR_CONST):
		n := CharVal{Token: p.token}
		p.next()
		return n
	case p.match(LPAREN):
		p.next()
		e := p.expr()
		p.assert(RPAREN)
		p.next()
		return e
	default:
		fmt.Println(p.token.Kind)
		panic("primaryExpr: " + p.token.String())
	}
}

func (p *Parser) readFuncCall(e Expr) FuncCall {
	p.assert(LPAREN)
	p.next()

	n := FuncCall{Ident: e.(Ident)}
	for !p.match(RPAREN) {
		expr := p.expr()
		n.Args = append(n.Args, expr)
		if p.match(COMMA) {
			p.next()
		}
	}
	p.next()

	return n
}

/**
Statement
*/

func (p *Parser) stmt() Stmt {
	switch {
	case p.isSelectionStmt():
		return p.selectionStmt()
	case p.isIterationStmt():
		return p.iterationStmt()
	case p.isJumpStmt():
		return p.jumpStmt()
	case p.match(LBRACE):
		return p.blockStmt()
	case p.isLabeledStmt():
		return p.labeledStmt()
	default:
		e := p.expr()
		p.assert(SEMICOLON)
		p.next()
		return ExprStmt{Expr: e}
	}
}

func (p *Parser) blockStmt() BlockStmt {
	p.assert(LBRACE)
	p.next()
	n := BlockStmt{}

	for !p.match(RBRACE) {
		if p.isType() {
			d := p.readVarDef()
			n.Nodes = append(n.Nodes, d)
		} else {
			stmt := p.stmt()
			n.Nodes = append(n.Nodes, stmt)
		}
	}
	p.next()

	return n
}

func (p *Parser) isSelectionStmt() bool {
	return p.match(IF) || p.match(SWITCH)
}

func (p *Parser) selectionStmt() Stmt {
	panic("selectionStmt")
}

func (p *Parser) isIterationStmt() bool {
	return p.match(WHILE) || p.match(DO) || p.match(FOR)
}

func (p *Parser) iterationStmt() Stmt {
	panic("iterationStmt")
}

func (p *Parser) isJumpStmt() bool {
	return p.match(GOTO) || p.match(CONTINUE) || p.match(BREAK) || p.match(RETURN)
}

func (p *Parser) jumpStmt() Stmt {
	if p.match(GOTO) {
		panic("unimplemented goto stmt")
	} else if p.match(CONTINUE) {
		panic("unimplemented continue stmt")
	} else if p.match(BREAK) {
		panic("unimplemented break stmt")
	} else if p.match(RETURN) {
		p.next()

		n := ReturnStmt{Expr: p.expr()}

		p.assert(SEMICOLON)
		p.next()
		return n
	} else {
		panic("expected jump statement, but got '" + p.token.String() + "'.")
	}
}

func (p *Parser) isLabeledStmt() bool {
	if p.match(CASE) || p.match(DEFAULT) {
		return true
	}
	if p.match(IDENT) {
		p.push()
		defer p.pop()
		p.next()
		return p.match(COLON)
	}
	return false
}

func (p *Parser) labeledStmt() Stmt {
	panic("labeledStmt")
}
