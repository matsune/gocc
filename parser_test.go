package main

import (
	"testing"
)

func intValExpect(t *testing.T, v IntVal, e string) {
	if v.Token.String() != e {
		t.Errorf("expected str is %s, but got %s", e, v.Token.String())
	}
}

func TestIntVal(t *testing.T) {
	p := NewParser([]byte("1"))

	e := p.expr()
	v, ok := e.(IntVal)
	if !ok {
		t.Errorf("expected type is IntVal")
		return
	}
	intValExpect(t, v, "1")
}

func TestBinaryExpr1(t *testing.T) {
	p := NewParser([]byte("1 + 2"))

	e := p.expr()
	v, ok := e.(BinaryExpr)
	if !ok {
		t.Errorf("expected type is BinaryExpr")
		return
	}
	x, ok := v.X.(IntVal)
	if !ok {
		t.Errorf("expected x type is IntVal")
		return
	}
	intValExpect(t, x, "1")

	y, ok := v.Y.(IntVal)
	if !ok {
		t.Errorf("expected y type is IntVal")
		return
	}
	intValExpect(t, y, "2")

	if v.Op.Kind != ADD {
		t.Errorf("expected op is %s", ADD)
		return
	}
}

func TestBinaryExpr2(t *testing.T) {
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
	v, ok := e.(BinaryExpr)
	if !ok {
		t.Errorf("expected type is BinaryExpr")
		return
	}

	// 1 + 2 * 3
	x, ok := v.X.(BinaryExpr)
	if !ok {
		t.Errorf("expected x type is BinaryExpr")
		return
	}

	// 1
	xx, ok := x.X.(IntVal)
	if !ok {
		t.Errorf("expected xx type is IntVal")
		return
	}
	intValExpect(t, xx, "1")

	// +
	if x.Op.Kind != ADD {
		t.Errorf("expected op is %s", ADD)
		return
	}

	// 2 * 3
	xy, ok := x.Y.(BinaryExpr)
	if !ok {
		t.Errorf("expected xy type is BinaryExpr")
		return
	}

	xyx, ok := xy.X.(IntVal)
	if !ok {
		t.Errorf("expected xyx type is IntVal")
		return
	}
	intValExpect(t, xyx, "2")

	xyy, ok := xy.Y.(IntVal)
	if !ok {
		t.Errorf("expected xyy type is IntVal")
		return
	}
	intValExpect(t, xyy, "3")

	if xy.Op.Kind != MUL {
		t.Errorf("expected op is %s", MUL)
		return
	}

	// 4 / 5 + 6
	y, ok := v.Y.(BinaryExpr)
	if !ok {
		t.Errorf("expected y type is BinaryExpr")
		return
	}

	// 4 / 5
	yx, ok := y.X.(BinaryExpr)
	if !ok {
		t.Errorf("expected yx type is BinaryExpr")
		return
	}

	yxx, ok := yx.X.(IntVal)
	if !ok {
		t.Errorf("expected yxx type is IntVal")
		return
	}
	intValExpect(t, yxx, "4")

	yxy, ok := yx.Y.(IntVal)
	if !ok {
		t.Errorf("expected yxy type is IntVal")
		return
	}
	intValExpect(t, yxy, "5")

	if yx.Op.Kind != DIV {
		t.Errorf("expected op is %s", DIV)
		return
	}

	yy, ok := y.Y.(IntVal)
	if !ok {
		t.Errorf("expected yy type is IntVal")
		return
	}
	intValExpect(t, yy, "6")

	if y.Op.Kind != ADD {
		t.Errorf("expected op is %s", ADD)
		return
	}
}
