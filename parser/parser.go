package parser

import (
	"fmt"
	"gocc/ast"
	"gocc/lexer"
	"gocc/token"
	"strconv"
)

type Parser struct {
	lexer *lexer.Lexer
	token *token.Token
	stack *Stack
}

func NewParser(source []byte) *Parser {
	p := &Parser{lexer: lexer.NewLexer(source), token: token.NewToken(), stack: NewStack()}
	p.next()
	return p
}

func (p *Parser) match(t token.TokenKind) bool {
	return p.token.Kind == t
}

func (p *Parser) matchs(ts []token.TokenKind) bool {
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

func (p *Parser) assert(t token.TokenKind) {
	if !p.match(t) {
		str := fmt.Sprintf("expected token is '"+t.String()+"', but got '"+p.token.String()+"' at line %d column %d", p.lexer.Pos().Line, p.lexer.Pos().Column)
		panic(str)
	}
}

func (p *Parser) next() {
	p.token = p.lexer.Next()
}

func (p *Parser) IsEnd() bool {
	return p.match(token.EOF)
}

func (p *Parser) Parse() ast.Node {
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

func (p *Parser) readVarDef() ast.Node {
	t := p.readType()

	p.assert(token.IDENT)
	tok := p.token
	p.next()

	var n ast.Node
	if p.match(token.LBRACK) {

		s := p.readSubscriptInit()
		t.Array = true
		arr := ast.ArrayDef{Type: t, Token: tok, Subscript: s}

		if s == nil && !p.match(token.ASSIGN) {
			panic(fmt.Errorf("definition of variable with array type needs an explicit size or an initializer"))
		}

		if p.match(token.ASSIGN) {
			p.next()
			init := p.readArrayInit()
			arr.Init = &init
			// obj.IsInit = true
		}
		n = arr
	} else {
		v := ast.VarDef{Type: t, Token: tok}

		if p.match(token.ASSIGN) {
			p.next()
			e := p.assignExpr()
			v.Init = &e
		}
		n = v
	}

	p.assert(token.SEMICOLON)
	p.next()

	return n
}

// [0] []
func (p *Parser) readSubscriptInit() *ast.Expr {
	p.assert(token.LBRACK)
	p.next()

	if p.match(token.RBRACK) {
		p.next()
		return nil
	}

	e := p.conditionalExpr()

	p.assert(token.RBRACK)
	p.next()

	return &e
}

func (p *Parser) readArrayInit() ast.ArrayInit {
	p.assert(token.LBRACE)
	p.next()
	n := ast.ArrayInit{}
	for {
		e := p.assignExpr()
		n.List = append(n.List, e)
		if p.match(token.RBRACE) {
			break
		} else if p.match(token.COMMA) {
			p.next()
		} else {
			panic("expected } or ,")
		}
	}
	p.next()
	return n
}

func (p *Parser) isType() bool {
	return p.matchs([]token.TokenKind{token.INT, token.CHAR, token.VOID, token.FLOAT, token.LONG, token.SHORT, token.DOUBLE})
}

func (p *Parser) readType() ast.CType {
	var t ast.CType
	for {
		if p.isType() {
			switch p.token.Kind {
			case token.INT:
				t.Primitive = ast.C_int
			case token.CHAR:
				t.Primitive = ast.C_char
			case token.VOID:
				t.Primitive = ast.C_void
			case token.FLOAT:
				t.Primitive = ast.C_float
			case token.LONG:
				t.Primitive = ast.C_long
			case token.SHORT:
				t.Primitive = ast.C_short
			case token.DOUBLE:
				t.Primitive = ast.C_double
			default:
				panic("readType")
			}
		} else if p.match(token.MUL) { // * as pointer
			t.Ptr = true
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
		case p.match(token.RPAREN):
			p.next()
			if p.match(token.LBRACE) {
				isFunc = true
				break loop
			}
		case p.match(token.SEMICOLON):
			break loop
		default:
			p.next()
		}

		if p.match(token.EOF) {
			break loop
		}
	}
	return isFunc
}

func (p *Parser) readFuncDef() ast.FuncDef {
	t := p.readType()

	p.assert(token.IDENT)
	name := string(p.token.Str)
	p.next()

	p.assert(token.LPAREN)
	p.next()

	args := p.readFuncArgs()

	p.assert(token.RPAREN)
	p.next()

	block := p.blockStmt()

	return ast.FuncDef{Type: t, Name: name, Args: args, Block: block}
}

func (p *Parser) readFuncArgs() []ast.FuncArg {
	var res []ast.FuncArg
	for p.isType() {
		res = append(res, p.readFuncArg())
		if !p.match(token.COMMA) {
			break
		}
		p.next()
	}
	return res
}

func (p *Parser) readFuncArg() ast.FuncArg {
	var n ast.FuncArg
	n.Type = p.readType()

	p.assert(token.IDENT)
	n.Name = p.token
	p.next()

	return n
}

/**
expression
*/

func (p *Parser) expr() ast.Expr {
	return p.assignExpr()
}

func (p *Parser) assignExpr() ast.Expr {
	p.push()

	var hasAssign bool
	for !p.match(token.SEMICOLON) {
		if p.match(token.ASSIGN) {
			hasAssign = true
			break
		}
		p.next()

		if p.match(token.EOF) {
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
		n := ast.AssignExpr{L: L, Op: op, R: R}
		return n
	} else {
		return p.conditionalExpr()
	}
}

func (p *Parser) isAssignOp() bool {
	return p.matchs(assignOps)
}

func (p *Parser) conditionalExpr() ast.Expr {
	e := p.logOrExpr()
	if p.match(token.QUE) {
		p.next()
		L := p.expr()
		p.assert(token.COLON)
		p.next()
		n := ast.CondExpr{Cond: e, L: L, R: p.conditionalExpr()}
		return n
	}
	return e
}

func (p *Parser) logOrExpr() ast.Expr {
	e := p.logAndExpr()
	return p.logOrExpr2(e)
}

func (p *Parser) logOrExpr2(e ast.Expr) ast.Expr {
	if p.match(token.LOR) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.logAndExpr()}
		return p.logOrExpr2(n)
	}
	return e
}

func (p *Parser) logAndExpr() ast.Expr {
	e := p.incOrExpr()
	return p.logAndExpr2(e)
}

func (p *Parser) logAndExpr2(e ast.Expr) ast.Expr {
	if p.match(token.LAND) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.incOrExpr()}
		return p.logAndExpr2(n)
	}
	return e
}

func (p *Parser) incOrExpr() ast.Expr {
	e := p.excOrExpr()
	return p.incOrExpr2(e)
}

func (p *Parser) incOrExpr2(e ast.Expr) ast.Expr {
	if p.match(token.OR) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.excOrExpr()}
		return p.incOrExpr2(n)
	}
	return e
}

func (p *Parser) excOrExpr() ast.Expr {
	e := p.andExpr()
	return p.excOrExpr2(e)
}

func (p *Parser) excOrExpr2(e ast.Expr) ast.Expr {
	if p.match(token.XOR) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.andExpr()}
		return p.excOrExpr2(n)
	}
	return e
}

func (p *Parser) andExpr() ast.Expr {
	e := p.eqExpr()
	return p.andExpr2(e)
}

func (p *Parser) andExpr2(e ast.Expr) ast.Expr {
	if p.match(token.AND) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.eqExpr()}
		return p.andExpr2(n)
	}
	return e
}

func (p *Parser) eqExpr() ast.Expr {
	e := p.relExpr()
	return p.eqExpr2(e)
}

func (p *Parser) eqExpr2(e ast.Expr) ast.Expr {
	if p.match(token.EQ) || p.match(token.NE) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.relExpr()}
		return p.eqExpr2(n)
	}
	return e
}

func (p *Parser) relExpr() ast.Expr {
	e := p.shiftExpr()
	return p.relExpr2(e)
}

func (p *Parser) relExpr2(e ast.Expr) ast.Expr {
	if p.match(token.LT) || p.match(token.GT) || p.match(token.LE) || p.match(token.GE) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.shiftExpr()}
		return p.relExpr2(n)
	}
	return e
}

func (p *Parser) shiftExpr() ast.Expr {
	e := p.additiveExpr()
	return p.shiftExpr2(e)
}

func (p *Parser) shiftExpr2(e ast.Expr) ast.Expr {
	if p.match(token.LSHIFT) || p.match(token.RSHIFT) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.additiveExpr()}
		return p.shiftExpr2(n)
	}
	return e
}

func (p *Parser) additiveExpr() ast.Expr {
	e := p.multiExpr()
	return p.additiveExpr2(e)
}

func (p *Parser) additiveExpr2(e ast.Expr) ast.Expr {
	if p.match(token.ADD) || p.match(token.SUB) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.multiExpr()}
		return p.additiveExpr2(n)
	}
	return e
}

func (p *Parser) multiExpr() ast.Expr {
	e := p.castExpr()
	return p.multiExpr2(e)
}

func (p *Parser) multiExpr2(e ast.Expr) ast.Expr {
	if p.match(token.MUL) || p.match(token.DIV) || p.match(token.REM) {
		op := p.token
		p.next()
		n := ast.BinaryExpr{X: e, Op: op, Y: p.castExpr()}
		return p.multiExpr2(n)
	}
	return e
}

func (p *Parser) castExpr() ast.Expr {
	return p.unaryExpr()
}

func (p *Parser) unaryExpr() ast.Expr {
	if p.match(token.INC) {
		p.next()
		i, ok := p.unaryExpr().(ast.Ident)
		if !ok {
			panic("unimplemented not ident increment")
		}
		return ast.IncExpr{Ident: i}
	} else if p.match(token.DEC) {
		p.next()
		i, ok := p.unaryExpr().(ast.Ident)
		if !ok {
			panic("unimplemented not ident decrement")
		}
		return ast.DecExpr{Ident: i}
	} else if p.isUnaryOp() {
		op := p.token
		p.next()

		switch op.Kind {
		case token.MUL:
			pv := ast.PtrVal{Token: p.token}
			p.next()
			return pv
		case token.AND:
			av := ast.AddressVal{Token: p.token}
			p.next()
			return av
		default:
			return ast.UnaryExpr{Op: op, Expr: p.castExpr()}
		}
	} else {
		return p.postfixExpr()
	}
}

func (p *Parser) isUnaryOp() bool {
	return p.matchs(unaryOps)
}

func (p *Parser) postfixExpr() ast.Expr {
	n := p.primaryExpr()
	return p.postfixExpr2(n)
}

func (p *Parser) postfixExpr2(e ast.Expr) ast.Expr {
	if p.match(token.INC) {
		p.next()
		i, ok := e.(ast.Ident)
		if !ok {
			panic("unimplemented not ident increment")
		}
		return ast.IncExpr{Ident: i}
	} else if p.match(token.DEC) {
		p.next()
		i, ok := e.(ast.Ident)
		if !ok {
			panic("unimplemented not ident decrement")
		}
		return ast.DecExpr{Ident: i}
	} else if p.match(token.LPAREN) {
		switch e.(type) {
		case ast.Ident:
			return p.readFuncCall(e)
		default:
			panic("unimplemented postfixExpr2")
		}
	} else if p.match(token.PERIOD) {
		panic("postfix .")
	} else if p.match(token.ARROW) {
		panic("postfix ->")
	} else {
		return e
	}
}

// [0] [1]
func (p *Parser) readSubscriptExpr(t *token.Token) ast.SubscriptExpr {
	p.assert(token.LBRACK)
	p.next()

	e := p.conditionalExpr()

	p.assert(token.RBRACK)
	p.next()

	se := ast.SubscriptExpr{Token: t, Expr: e}
	return se
}

func (p *Parser) primaryExpr() ast.Expr {
	switch {
	case p.match(token.IDENT):
		t := p.token
		p.next()

		if p.match(token.LBRACK) {
			return p.readSubscriptExpr(t)
		} else {
			n := ast.Ident{Token: t}
			return n
		}
	case p.match(token.INT_CONST):
		i, err := strconv.Atoi(p.token.String())
		if err != nil {
			panic(err)
		}
		n := ast.IntVal{Num: i}
		p.next()
		return n
	case p.match(token.CHAR_CONST):
		n := ast.CharVal{Token: p.token}
		p.next()
		return n
	case p.match(token.LPAREN):
		p.next()
		e := p.expr()
		p.assert(token.RPAREN)
		p.next()
		return e
	default:
		fmt.Println(p.token.Kind)
		panic("primaryExpr: " + p.token.String())
	}
}

func (p *Parser) readFuncCall(e ast.Expr) ast.FuncCall {
	p.assert(token.LPAREN)
	p.next()

	n := ast.FuncCall{Ident: e.(ast.Ident)}
	for !p.match(token.RPAREN) {
		expr := p.expr()
		n.Args = append(n.Args, expr)
		if p.match(token.COMMA) {
			p.next()
		}
	}
	p.next()

	return n
}

/**
Statement
*/

func (p *Parser) stmt() ast.Stmt {
	switch {
	case p.isSelectionStmt():
		return p.selectionStmt()
	case p.isIterationStmt():
		return p.iterationStmt()
	case p.isJumpStmt():
		return p.jumpStmt()
	case p.match(token.LBRACE):
		return p.blockStmt()
	case p.isLabeledStmt():
		return p.labeledStmt()
	default:
		e := p.expr()
		p.assert(token.SEMICOLON)
		p.next()
		return ast.ExprStmt{Expr: e}
	}
}

func (p *Parser) blockStmt() ast.BlockStmt {
	p.assert(token.LBRACE)
	p.next()
	n := ast.BlockStmt{}

	for !p.match(token.RBRACE) {
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
	return p.match(token.IF) || p.match(token.SWITCH)
}

func (p *Parser) selectionStmt() ast.Stmt {
	if p.match(token.IF) {
		return p.ifStmt()
	} else if p.match(token.SWITCH) {
		panic("selection stmt switch is not implemented.")
	} else {
		panic(fmt.Sprintf("%s not selection Stmt", p.token))
	}
}

func (p *Parser) ifStmt() ast.IfStmt {
	p.assert(token.IF)
	p.next()

	p.assert(token.LPAREN)
	p.next()

	e := p.conditionalExpr()

	p.assert(token.RPAREN)
	p.next()

	b := p.blockStmt()

	return ast.IfStmt{Expr: &e, Block: b, Else: p.elseStmt()}
}

func (p *Parser) elseStmt() *ast.IfStmt {
	if !p.match(token.ELSE) {
		return nil
	}

	p.next()

	if p.match(token.IF) {
		s := p.ifStmt()
		return &s
	} else {
		return &ast.IfStmt{Expr: nil, Block: p.blockStmt(), Else: nil}
	}
}

func (p *Parser) isIterationStmt() bool {
	return p.match(token.WHILE) || p.match(token.DO) || p.match(token.FOR)
}

func (p *Parser) iterationStmt() ast.Stmt {
	panic("iterationStmt")
}

func (p *Parser) isJumpStmt() bool {
	return p.match(token.GOTO) || p.match(token.CONTINUE) || p.match(token.BREAK) || p.match(token.RETURN)
}

func (p *Parser) jumpStmt() ast.Stmt {
	if p.match(token.GOTO) {
		panic("unimplemented goto stmt")
	} else if p.match(token.CONTINUE) {
		panic("unimplemented continue stmt")
	} else if p.match(token.BREAK) {
		panic("unimplemented break stmt")
	} else if p.match(token.RETURN) {
		p.next()

		n := ast.ReturnStmt{Expr: p.expr()}

		p.assert(token.SEMICOLON)
		p.next()
		return n
	} else {
		panic("expected jump statement, but got '" + p.token.String() + "'.")
	}
}

func (p *Parser) isLabeledStmt() bool {
	if p.match(token.CASE) || p.match(token.DEFAULT) {
		return true
	}
	if p.match(token.IDENT) {
		p.push()
		defer p.pop()
		p.next()
		return p.match(token.COLON)
	}
	return false
}

func (p *Parser) labeledStmt() ast.Stmt {
	panic("labeledStmt")
}
