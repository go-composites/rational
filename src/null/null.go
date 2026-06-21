package NullRational

import (
	MethodNotImplementedError "github.com/go-composites/error/src/method_not_implemented"
	Rational "github.com/go-composites/rational/src"
	Result "github.com/go-composites/result/src"
)

/*
NullRational is the Null-Object variant of Rational.

It satisfies Rational.Interface so callers never have to test for a bare nil:
its value is zero, its arithmetic yields a Result carrying a
"method not implemented" Error, its comparisons are false (except Equal against
another null), and IsNull() returns true.
*/
type Interface interface {
	Rational.Interface
}

type data struct{}

/*
New returns a NullRational.
*/
func New() Interface {
	return &data{}
}

func (d data) Numerator() int64 {
	return 0
}

func (d data) Denominator() int64 {
	return 0
}

func (d data) ToGoString() string {
	return ``
}

func (d data) ToFloat() float64 {
	return 0
}

func (d data) IsNull() bool {
	return true
}

func notImplemented(methodName string) Result.Interface {
	return Result.New(
		Result.WithError(
			MethodNotImplementedError.New(methodName),
		),
	)
}

func (d data) Add(Rational.Interface) Result.Interface {
	return notImplemented(`Add`)
}

func (d data) Sub(Rational.Interface) Result.Interface {
	return notImplemented(`Sub`)
}

func (d data) Mul(Rational.Interface) Result.Interface {
	return notImplemented(`Mul`)
}

func (d data) Div(Rational.Interface) Result.Interface {
	return notImplemented(`Div`)
}

func (d data) Abs() Result.Interface {
	return notImplemented(`Abs`)
}

func (d data) Neg() Result.Interface {
	return notImplemented(`Neg`)
}

func (d data) Equal(other Rational.Interface) bool {
	return other.IsNull()
}

func (d data) LessThan(Rational.Interface) bool {
	return false
}

func (d data) GreaterThan(Rational.Interface) bool {
	return false
}

func (d data) Inspect() Rational.String {
	return `<NullRational>`
}
