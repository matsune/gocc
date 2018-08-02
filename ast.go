package main

type Kind int

const (
	BINARY_EXPR Kind = iota
	COND_EXPR
	UNARY_EXPR
	ASSIGN_EXPR

	INT_VAL
)

type (
	Node interface {
		Kind() Kind
	}

	Ident struct {
		Token *Token
	}

	Expr interface {
		Node
		Expr()
	}

	BinaryExpr struct {
		X  Expr
		Op *Token
		Y  Expr
	}

	CondExpr struct {
		Cond Expr
		L    Expr
		R    Expr
	}

	UnaryExpr struct {
		Op *Token
		E  Expr
	}

	AssignExpr struct {
		L  Expr
		Op *Token
		R  Expr
	}

	IntVal struct {
		Token *Token
	}
)

// func (Ident) Kind() Kind { return IDENT }

func (BinaryExpr) Kind() Kind { return BINARY_EXPR }
func (CondExpr) Kind() Kind   { return COND_EXPR }
func (UnaryExpr) Kind() Kind  { return UNARY_EXPR }
func (AssignExpr) Kind() Kind { return ASSIGN_EXPR }

func (BinaryExpr) Expr() {}
func (CondExpr) Expr()   {}
func (UnaryExpr) Expr()  {}
func (AssignExpr) Expr() {}

func (IntVal) Kind() Kind { return INT_VAL }
func (IntVal) Expr()      {}
