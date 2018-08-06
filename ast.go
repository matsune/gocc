package main

type Kind int

const (
	VAR_DEF Kind = iota
	FUNC_DEF
	FUNC_ARG
	AST_IDENT

	// expr
	BINARY_EXPR
	COND_EXPR
	UNARY_EXPR
	ASSIGN_EXPR
	FUNC_CALL
	INT_VAL

	// stmt
	BLOCK_STMT
	RETURN_STMT
	EXPR_STMT
)

type Type int

const (
	Int_t Type = iota
	Void_t
	Char_t
	Float_t
	Long_t
	Short_t
	Double_t
)

func (t Type) Size() int {
	switch t {
	case Int_t:
		return 4
	default:
		panic("unimplemented type size")
	}
}

func (t Type) String() string {
	switch t {
	case Int_t:
		return "int"
	case Void_t:
		return "void"
	case Char_t:
		return "char"
	case Float_t:
		return "float"
	case Long_t:
		return "long"
	case Short_t:
		return "short"
	case Double_t:
		return "double"
	default:
		panic("undefined Type")
	}
}

type (
	Node interface {
		Kind() Kind
	}
)

type (
	Ident struct {
		Token *Token
	}

	VarDef struct {
		Type Type
		Name string
		Init *Expr
	}

	FuncDef struct {
		Type  Type
		Name  string
		Args  []FuncArg
		Block BlockStmt
	}

	FuncArg struct {
		Type Type
		Name *Token
	}
)

type (
	Expr interface {
		Node
		expr()
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
		Op   *Token
		Expr Expr
	}

	AssignExpr struct {
		L  Expr
		Op *Token
		R  Expr
	}

	IntVal struct {
		Token *Token
	}

	FuncCall struct {
		Ident Ident
		Args  []Expr
	}
)

type (
	Stmt interface {
		Node
		stmt()
	}

	BlockStmt struct {
		Nodes []Node
	}

	ReturnStmt struct {
		Expr Expr
	}

	ExprStmt struct {
		Expr Expr
	}
)

func (VarDef) Kind() Kind  { return VAR_DEF }
func (FuncDef) Kind() Kind { return FUNC_DEF }
func (FuncArg) Kind() Kind { return FUNC_ARG }
func (Ident) Kind() Kind   { return AST_IDENT }

func (BinaryExpr) Kind() Kind { return BINARY_EXPR }
func (CondExpr) Kind() Kind   { return COND_EXPR }
func (UnaryExpr) Kind() Kind  { return UNARY_EXPR }
func (AssignExpr) Kind() Kind { return ASSIGN_EXPR }

func (FuncCall) Kind() Kind { return FUNC_CALL }

func (IntVal) Kind() Kind { return INT_VAL }

func (BlockStmt) Kind() Kind  { return BLOCK_STMT }
func (ReturnStmt) Kind() Kind { return RETURN_STMT }
func (ExprStmt) Kind() Kind   { return EXPR_STMT }

func (Ident) expr()      {}
func (BinaryExpr) expr() {}
func (CondExpr) expr()   {}
func (UnaryExpr) expr()  {}
func (AssignExpr) expr() {}
func (FuncCall) expr()   {}
func (IntVal) expr()     {}

func (BlockStmt) stmt()  {}
func (ReturnStmt) stmt() {}
func (ExprStmt) stmt()   {}
