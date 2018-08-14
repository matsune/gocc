package parser

import (
	"gocc/ast"
	"gocc/token"
	"reflect"
	"testing"
)

func intValExpect(t *testing.T, v ast.IntVal, n int) {
	if v.Num != n {
		t.Errorf("expected num is %d, but got %d", n, v.Num)
	}
}

func TestIntVal(t *testing.T) {
	p := NewParser([]byte("1"))

	e := p.expr()
	v, ok := e.(ast.IntVal)
	if !ok {
		t.Errorf("expected type is ast.IntVal")
		return
	}
	intValExpect(t, v, 1)
}

func TestBinaryExpr(t *testing.T) {
	/**
	   Binary {
	     X:  Binary{
	      X: 1,
	      Op: +,
	      Y:  Binary{ X: 2, Op: *, Y: 3 }
	    },
	    Op: -,
	    Y:  Binary{
	      X:  Binary{ X: 4, Op: /, Y: 5 },
	      Op: +,
	      Y: 6
	    }
	  }
	*/
	p := NewParser([]byte("1 + 2 * 3 - (4 / 5 + 6)"))

	e := p.expr()
	v, ok := e.(ast.BinaryExpr)
	if !ok {
		t.Errorf("expected type isast.BinaryExpr")
		return
	}

	// 1 + 2 * 3
	x, ok := v.X.(ast.BinaryExpr)
	if !ok {
		t.Errorf("expected x type isast.BinaryExpr")
		return
	}

	// 1
	xx, ok := x.X.(ast.IntVal)
	if !ok {
		t.Errorf("expected xx type is ast.IntVal")
		return
	}
	intValExpect(t, xx, 1)

	// +
	if x.Op.Kind != token.ADD {
		t.Errorf("expected op is %s", token.ADD)
		return
	}

	// 2 * 3
	xy, ok := x.Y.(ast.BinaryExpr)
	if !ok {
		t.Errorf("expected xy type isast.BinaryExpr")
		return
	}

	xyx, ok := xy.X.(ast.IntVal)
	if !ok {
		t.Errorf("expected xyx type is ast.IntVal")
		return
	}
	intValExpect(t, xyx, 2)

	xyy, ok := xy.Y.(ast.IntVal)
	if !ok {
		t.Errorf("expected xyy type is ast.IntVal")
		return
	}
	intValExpect(t, xyy, 3)

	if xy.Op.Kind != token.MUL {
		t.Errorf("expected op is %s", token.MUL)
		return
	}

	// 4 / 5 + 6
	y, ok := v.Y.(ast.BinaryExpr)
	if !ok {
		t.Errorf("expected y type isast.BinaryExpr")
		return
	}

	// 4 / 5
	yx, ok := y.X.(ast.BinaryExpr)
	if !ok {
		t.Errorf("expected yx type isast.BinaryExpr")
		return
	}

	yxx, ok := yx.X.(ast.IntVal)
	if !ok {
		t.Errorf("expected yxx type is ast.IntVal")
		return
	}
	intValExpect(t, yxx, 4)

	yxy, ok := yx.Y.(ast.IntVal)
	if !ok {
		t.Errorf("expected yxy type is ast.IntVal")
		return
	}
	intValExpect(t, yxy, 5)

	if yx.Op.Kind != token.DIV {
		t.Errorf("expected op is %s", token.DIV)
		return
	}

	yy, ok := y.Y.(ast.IntVal)
	if !ok {
		t.Errorf("expected yy type is ast.IntVal")
		return
	}
	intValExpect(t, yy, 6)

	if y.Op.Kind != token.ADD {
		t.Errorf("expected op is %s", token.ADD)
		return
	}
}

func varTypeExpect(t *testing.T, v ast.VarDef, ty ast.CType) {
	if v.Type != ty {
		t.Errorf("expected type is %s, but got %s", ty, v.Type)
	}
}

func varNameExpect(t *testing.T, v ast.VarDef, name string) {
	if v.Token.String() != name {
		t.Errorf("expected name is %s, but got %s", name, v.Token.String())
	}
}

func TestReadVarDef(t *testing.T) {
	p := NewParser([]byte("int a;"))
	n := p.readVarDef()
	v, ok := n.(ast.VarDef)
	if !ok {
		t.Errorf("expected type is VarDef, but got %s", reflect.TypeOf(n))
	}
	varTypeExpect(t, v, ast.C_int)
	varNameExpect(t, v, "a")
	if v.Init != nil {
		t.Errorf("expected varDef is not initialized")
	}
}

func TestReadVarDefWithInit(t *testing.T) {
	p := NewParser([]byte("int a = 3 + 4;"))
	n := p.readVarDef()
	v, ok := n.(ast.VarDef)
	if !ok {
		t.Errorf("expected type is VarDef, but got %s", reflect.TypeOf(n))
	}
	varTypeExpect(t, v, ast.C_int)
	varNameExpect(t, v, "a")
	b, ok := (*v.Init).(ast.BinaryExpr)
	if !ok {
		t.Errorf("expected type isast.BinaryExpr, but got %s", reflect.TypeOf(v.Init))
	}
	intValExpect(t, b.X.(ast.IntVal), 3)
	intValExpect(t, b.Y.(ast.IntVal), 4)
}

func TestReadFuncDef(t *testing.T) {
	p := NewParser([]byte("int main(int argc) { int a = 2 + 4; }"))
	f := p.readFuncDef()
	if f.Type != ast.C_int {
		t.Errorf("expected type is %s, but got %s", ast.C_int, f.Type)
	}
	if f.Name != "main" {
		t.Errorf("expected name is %s, but got %s", "main", f.Name)
	}
	if len(f.Args) != 1 {
		t.Errorf("expected args count is %d, but got %d", 1, len(f.Args))
	}
	if f.Args[0].Type != ast.C_int {
		t.Errorf("expected type is %s, but got %s", ast.C_int, f.Args[0].Type)
	}
	if f.Args[0].Name.String() != "argc" {
		t.Errorf("expected type is %s, but got %s", "argc", f.Args[0].Name)
	}
	if len(f.Block.Nodes) != 1 {
		t.Errorf("expected block nodes count is %d, but got %d", 1, len(f.Block.Nodes))
	}

	v, ok := f.Block.Nodes[0].(ast.VarDef)
	if !ok {
		t.Errorf("expected block nodes[0] is VarDef, but got %s", reflect.TypeOf(f.Block.Nodes[0]))
	}
	if v.Type != ast.C_int {
		t.Errorf("expected type is %s, but got %s", ast.C_int, v.Type)
	}
	if v.Token.String() != "a" {
		t.Errorf("expected name is %s, but got %s", "a", v.Token.String())
	}
}

func TestIsFuncDef(t *testing.T) {
	p := NewParser([]byte("int main(int argc) { int a = 2 + 4; }"))
	if !p.isFuncDef() {
		t.Errorf("expected source is funcDef")
	}

	p = NewParser([]byte("int a = 2 + 4;"))
	if p.isFuncDef() {
		t.Errorf("expected source is not funcDef")
	}
}

func TestFuncCall(t *testing.T) {
	p := NewParser([]byte("func(a, b, c);"))
	f, ok := p.expr().(ast.FuncCall)
	if !ok {
		t.Errorf("expected type is FuncCall, but got %s", reflect.TypeOf(p.expr()))
	}
	if f.Ident.Token.String() != "func" {
		t.Errorf("expected func name is %s, but got %s", "func", f.Ident.Token)
	}
	if len(f.Args) != 3 {
		t.Errorf("expected args count is %d, but got %d", 3, len(f.Args))
	}
	idents := []string{"a", "b", "c"}
	for i, v := range f.Args {
		a, ok := v.(ast.Ident)
		if !ok {
			t.Errorf("expected arg[%d] is not ident", i)
		}
		if a.Token.String() != idents[i] {
			t.Errorf("expected arg[%d] is %s, but got %s", i, idents[i], a.Token.String())
		}
	}
}

func TestFuncCall2(t *testing.T) {
	p := NewParser([]byte("int main() { return a(); }"))
	f, ok := p.Parse().(ast.FuncDef)
	if !ok {
		t.Errorf("expected type is FuncDef, but got %s", reflect.TypeOf(p.Parse()))
	}
	if len(f.Block.Nodes) != 1 {
		t.Errorf("expected nodes count is %d, but got %d", 1, len(f.Block.Nodes))
	}
}

func TestReadType(t *testing.T) {
	p := NewParser([]byte("int"))
	ty := p.readType()
	if ty != ast.C_int {
		t.Errorf("expected type is %s, but got %s", ast.C_int, ty)
	}
	p = NewParser([]byte("char"))
	ty = p.readType()
	if ty != ast.C_char {
		t.Errorf("expected type is %s, but got %s", ast.C_char, ty)
	}
	p = NewParser([]byte("int *"))
	ty = p.readType()
	if ty != ast.C_pointer {
		t.Errorf("expected type is %s, but got %s", ast.C_pointer, ty)
	}
}

func TestPointer(t *testing.T) {
	p := NewParser([]byte("{ *a = *a + b; }"))
	e := p.blockStmt()
	a, ok := e.Nodes[0].(ast.ExprStmt).Expr.(ast.AssignExpr)
	if !ok {
		t.Errorf("expected type is AssignExpr, but got %s", reflect.TypeOf(e.Nodes[0].(ast.ExprStmt).Expr))
	}
	_, ok = a.L.(ast.PtrVal)
	if !ok {
		t.Errorf("expected type is PointerVal, but got %s", reflect.TypeOf(a.L))
	}
}

func TestParseArray(t *testing.T) {
	p := NewParser([]byte("int a[4];"))
	e := p.readVarDef()
	v, ok := e.(ast.ArrayDef)
	if !ok {
		t.Errorf("expected type is ArrayDef, but got %s", reflect.TypeOf(e))
	}
	if v.Type != ast.C_int {
		t.Errorf("expected type is %s, but got %s", ast.C_int, v.Type)
	}
	if v.Token.String() != "a" {
		t.Errorf("expected name is %s, but got %s", "a", v.Token.String())
	}
	vv, ok := (*v.Subscript).(ast.IntVal)
	if !ok {
		t.Errorf("expected type is ast.IntVal, but got %s", reflect.TypeOf(v.Subscript))
	}
	if vv.Num != 4 {
		t.Errorf("expected string is 4, but got %d", vv.Num)
	}
	if v.Init != nil {
		t.Errorf("expected init is nil")
	}
}

func TestParseArrayInit(t *testing.T) {
	p := NewParser([]byte("int a[] = {0, 1, 2, 3};"))
	e := p.readVarDef()
	v, ok := e.(ast.ArrayDef)
	if !ok {
		t.Errorf("expected type is ArrayDef, but got %s", reflect.TypeOf(e))
	}
	if v.Subscript != nil {
		t.Errorf("expected subscript is nil")
	}
	i := *v.Init
	if len(i.List) != 4 {
		t.Errorf("expected count of elements is %d, but got %d", 4, len(i.List))
	}
}

func TestReturnSubscript(t *testing.T) {
	p := NewParser([]byte("return a[0];"))
	e := p.stmt()
	v, ok := e.(ast.ReturnStmt)
	if !ok {
		t.Errorf("expected type is ReturnStmt, but got %s", reflect.TypeOf(e))
	}
	vv, ok := v.Expr.(ast.SubscriptExpr)
	if !ok {
		t.Errorf("expected type is SubscriptExpr, but got %s", reflect.TypeOf(v.Expr))
	}
	if vv.Token.String() != "a" {
		t.Errorf("expected ident is %s, but got %s", "a", vv.Token.String())
	}
	i, ok := vv.Expr.(ast.IntVal)
	if !ok {
		t.Errorf("expected type is ast.IntVal, but got %s", reflect.TypeOf(vv.Expr))
	}
	if i.Num != 0 {
		t.Errorf("expected str is %d, but got %d", 0, i.Num)
	}
}

func TestIfStmt(t *testing.T) {
	p := NewParser([]byte("if (a == 0) { return 0; } else if (a == 1) { return 1; } else { return 2; }"))
	if1 := p.ifStmt()
	e := (*if1.Expr).(ast.BinaryExpr)
	if e.Op.Kind != token.EQ {
		t.Errorf("expected token is %s, but got %s", token.EQ, e.Op.Kind)
	}
	x := e.X.(ast.Ident)
	if x.Token.String() != "a" {
		t.Errorf("expected token is %s, but got %s", "a", x)
	}
	y := e.Y.(ast.IntVal)
	if y.Num != 0 {
		t.Errorf("expected num is %d, but got %d", 0, y.Num)
	}

	n := if1.Block.Nodes[0].(ast.ReturnStmt).Expr.(ast.IntVal)
	if n.Num != 0 {
		t.Errorf("expected return is %d, but got %d", 0, n.Num)
	}

	if2 := if1.Else
	e2 := (*if2.Expr).(ast.BinaryExpr)
	x = e2.X.(ast.Ident)
	if x.Token.String() != "a" {
		t.Errorf("expected token is %s, but got %s", "a", x)
	}
	y = e2.Y.(ast.IntVal)
	if y.Num != 1 {
		t.Errorf("expected num is %d, but got %d", 1, y.Num)
	}

	n = if2.Block.Nodes[0].(ast.ReturnStmt).Expr.(ast.IntVal)
	if n.Num != 1 {
		t.Errorf("expected return is %d, but got %d", 1, n.Num)
	}

	if3 := if2.Else
	if if3.Expr != nil {
		t.Errorf("expected if expr is nil")
	}

	n = if3.Block.Nodes[0].(ast.ReturnStmt).Expr.(ast.IntVal)
	if n.Num != 2 {
		t.Errorf("expected return is %d, but got %d", 2, n.Num)
	}

	if if3.Else != nil {
		t.Errorf("expected else stmt is nil")
	}
}

func TestIncrement(t *testing.T) {
	p := NewParser([]byte("{a++; ++a;}"))
	b := p.blockStmt()
	i1 := b.Nodes[0].(ast.ExprStmt).Expr.(ast.IncExpr)
	a := i1.Postfix.(ast.Ident)
	if a.Token.String() != "a" {
		t.Errorf("expected ident is %s, but got %s", "a", a.Token)
	}

	i2 := b.Nodes[1].(ast.ExprStmt).Expr.(ast.IncExpr)
	a = i2.Postfix.(ast.Ident)
	if a.Token.String() != "a" {
		t.Errorf("expected ident is %s, but got %s", "a", a.Token)
	}
}

func TestDecrement(t *testing.T) {
	p := NewParser([]byte("{a--; --a;}"))
	b := p.blockStmt()
	i1 := b.Nodes[0].(ast.ExprStmt).Expr.(ast.DecExpr)
	a := i1.Postfix.(ast.Ident)
	if a.Token.String() != "a" {
		t.Errorf("expected ident is %s, but got %s", "a", a.Token)
	}

	i2 := b.Nodes[1].(ast.ExprStmt).Expr.(ast.DecExpr)
	a = i2.Postfix.(ast.Ident)
	if a.Token.String() != "a" {
		t.Errorf("expected ident is %s, but got %s", "a", a.Token)
	}
}
