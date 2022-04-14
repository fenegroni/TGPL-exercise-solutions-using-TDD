package expr

import (
	"fmt"
	"math"
	"strings"
)

type Expr interface {
	Eval(env Env) float64
	// Check reports errors in this Expr and adds its Vars to the set.
	Check(vars map[Var]bool) error
	String() string
}

type literal float64

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (literal) Check(_ map[Var]bool) error {
	return nil
}

func (l literal) String() string {
	return fmt.Sprintf("%.6g", l)
}

type unary struct {
	op rune // '+', '-'
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (u unary) String() string {
	return fmt.Sprintf("%c%s", u.op, u.x)
}

type binary struct {
	op   rune // '+', '-', '*', '/'
	x, y Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("*/+-", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (b binary) String() string {
	var left, right string
	if precedence(b.op) < 2 {
		left, right = "(", ")"
	}
	return fmt.Sprintf("%s%s %c %s%s", left, b.x, b.op, b.y, right)
}

type Var string

type Env map[Var]float64

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (v Var) String() string {
	return string(v)
}

type callArgs []Expr

func (c callArgs) String() string {
	b := new(strings.Builder)
	for n, a := range c {
		sep := ", "
		if n == 0 {
			sep = ""
		}
		_, _ = fmt.Fprint(b, sep, a.String())
	}
	return b.String()
}

type call struct {
	fn   string // "pow", "sin", "sqrt"
	args callArgs
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %q", c.fn))
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d", c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) String() string {
	return fmt.Sprintf("%s(%s)", c.fn, c.args)
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}
