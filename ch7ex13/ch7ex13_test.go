package ch7ex13

import (
	"fmt"
	"github.com/fenegroni/TGPL-exercise-solutions/ch7ex13/expr"
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		exp  string
		env  expr.Env
		want string
	}{
		{"sqrt(A / pi)", expr.Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", expr.Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", expr.Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", expr.Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", expr.Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", expr.Env{"F": 212}, "100"},
	}
	var prevExpr string
	for _, test := range tests {
		if test.exp != prevExpr {
			fmt.Printf("\n%s\n", test.exp)
			prevExpr = test.exp
		}
		exp, err := expr.Parse(test.exp)
		if err != nil {
			t.Errorf("Parse: %s", err)
			continue
		}
		got := fmt.Sprintf("%.6g", exp.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q",
				test.exp, test.env, got, test.want)
		}
	}
}
