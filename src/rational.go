package Rational

import (
	"fmt"
	"math/big"

	Error "github.com/go-composites/error/src"
	MethodNotImplementedError "github.com/go-composites/error/src/method_not_implemented"
	Result "github.com/go-composites/result/src"
)

/*
Rational is an exact rational-number composite over a math/big.Rat.

It mirrors Ruby's Rational: a fraction is always stored in lowest terms and is
exact — 1/3 + 1/6 is exactly 1/2, with no floating-point error. Its fallible
operations (notably the constructors and Div) return a Result.Interface so that
failures — such as a zero denominator or a division by zero — are values rather
than panics, and they never return a bare nil.
*/
type Interface interface {
	Numerator() int64
	Denominator() int64
	ToGoString() string
	ToFloat() float64
	IsNull() bool
	Add(Interface) Result.Interface
	Sub(Interface) Result.Interface
	Mul(Interface) Result.Interface
	Div(Interface) Result.Interface
	Abs() Result.Interface
	Neg() Result.Interface
	Equal(Interface) bool
	LessThan(Interface) bool
	GreaterThan(Interface) bool
	Inspect() String
}

// String is the lightweight inspection representation of a Rational.
type String = string

type data struct {
	value *big.Rat
}

/*
FromInts is the Rational constructor from two Go int64 values.

It returns a Result whose payload is the reduced fraction num/den. When the
denominator is zero the Result carries an Error instead of a payload — the
construction never panics and never returns nil.

	r := Rational.FromInts(3, 4) // 3/4
*/
func FromInts(num, den int64) Result.Interface {
	if den == 0 {
		return Result.New(
			Result.WithError(
				Error.New("zero denominator"),
			),
		)
	}
	return payload(
		new(big.Rat).SetFrac(
			big.NewInt(num),
			big.NewInt(den),
		),
	)
}

/*
FromString parses a fraction string such as "3/4" into a Rational.

It returns a Result whose payload is the parsed, reduced Rational. When the
input is not a valid rational the Result carries an Error instead of a payload —
the parse never panics and never returns nil.

	r := Rational.FromString("3/4")
	if !r.HasError() {
	    q := r.Payload().(Rational.Interface)
	}
*/
func FromString(s string) Result.Interface {
	value, ok := new(big.Rat).SetString(s)
	if !ok {
		return Result.New(
			Result.WithError(
				Error.New("invalid rational: " + s),
			),
		)
	}
	return payload(value)
}

/*
Null returns the Null-Object variant of Rational.

It is defined in src/null; this thin re-export keeps a Null next to the
concrete constructors. The returned value satisfies Interface and reports
IsNull() == true.
*/
func Null() Interface {
	return newNull()
}

/*
Numerator returns the numerator of the reduced fraction as a Go int64.
*/
func (d data) Numerator() int64 {
	return d.value.Num().Int64()
}

/*
Denominator returns the denominator of the reduced fraction as a Go int64.
*/
func (d data) Denominator() int64 {
	return d.value.Denom().Int64()
}

/*
ToGoString returns the "num/den" representation of the reduced fraction.
*/
func (d data) ToGoString() string {
	return d.value.RatString()
}

/*
ToFloat returns the value as a Go float64.

This conversion is lossy; the exact value is preserved by ToGoString and the
arithmetic methods.
*/
func (d data) ToFloat() float64 {
	f, _ := d.value.Float64()
	return f
}

/*
IsNull reports whether the Rational is the Null-Object variant.

A concrete Rational is never null.
*/
func (d data) IsNull() bool {
	return false
}

/*
Add returns a Result whose payload is the sum of the receiver and other.

A fresh big.Rat backs the payload; the operands are never mutated.
*/
func (d data) Add(other Interface) Result.Interface {
	return payload(
		new(big.Rat).Add(d.value, fromInterface(other)),
	)
}

/*
Sub returns a Result whose payload is the difference of the receiver and other.
*/
func (d data) Sub(other Interface) Result.Interface {
	return payload(
		new(big.Rat).Sub(d.value, fromInterface(other)),
	)
}

/*
Mul returns a Result whose payload is the product of the receiver and other.
*/
func (d data) Mul(other Interface) Result.Interface {
	return payload(
		new(big.Rat).Mul(d.value, fromInterface(other)),
	)
}

/*
Div returns a Result whose payload is the quotient of the receiver and other.

When other is a zero rational the Result carries an Error ("division by zero")
instead of a payload — the division never panics and never returns nil.
*/
func (d data) Div(other Interface) Result.Interface {
	rhs := fromInterface(other)
	if rhs.Sign() == 0 {
		return Result.New(
			Result.WithError(
				Error.New("division by zero"),
			),
		)
	}
	return payload(
		new(big.Rat).Quo(d.value, rhs),
	)
}

/*
Abs returns a Result whose payload is the absolute value of the receiver.
*/
func (d data) Abs() Result.Interface {
	return payload(
		new(big.Rat).Abs(d.value),
	)
}

/*
Neg returns a Result whose payload is the negation of the receiver.
*/
func (d data) Neg() Result.Interface {
	return payload(
		new(big.Rat).Neg(d.value),
	)
}

/*
Equal reports whether the receiver and other hold the same rational value.
*/
func (d data) Equal(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) == 0
}

/*
LessThan reports whether the receiver is strictly less than other.
*/
func (d data) LessThan(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) < 0
}

/*
GreaterThan reports whether the receiver is strictly greater than other.
*/
func (d data) GreaterThan(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) > 0
}

/*
Inspect returns a one-line representation of the Rational with its address and
value — mirroring the style of the other composites.
*/
func (d data) Inspect() String {
	return fmt.Sprintf(
		"<Rational:%p value=%s>",
		&d, d.value.RatString(),
	)
}

// nullData is the Null-Object variant returned by Null(). The importable
// NullRational package in src/null mirrors it; this copy keeps a Null next to
// the concrete constructors without creating an import cycle.
type nullData struct{}

func newNull() Interface {
	return &nullData{}
}

func nullNotImplemented(methodName string) Result.Interface {
	return Result.New(
		Result.WithError(
			MethodNotImplementedError.New(methodName),
		),
	)
}

func (nullData) Numerator() int64               { return 0 }
func (nullData) Denominator() int64             { return 0 }
func (nullData) ToGoString() string             { return `` }
func (nullData) ToFloat() float64               { return 0 }
func (nullData) IsNull() bool                   { return true }
func (nullData) Add(Interface) Result.Interface { return nullNotImplemented(`Add`) }
func (nullData) Sub(Interface) Result.Interface { return nullNotImplemented(`Sub`) }
func (nullData) Mul(Interface) Result.Interface { return nullNotImplemented(`Mul`) }
func (nullData) Div(Interface) Result.Interface { return nullNotImplemented(`Div`) }
func (nullData) Abs() Result.Interface          { return nullNotImplemented(`Abs`) }
func (nullData) Neg() Result.Interface          { return nullNotImplemented(`Neg`) }
func (nullData) Equal(other Interface) bool     { return other.IsNull() }
func (nullData) LessThan(Interface) bool        { return false }
func (nullData) GreaterThan(Interface) bool     { return false }
func (nullData) Inspect() String                { return `<NullRational>` }

// payload wraps a fresh big.Rat in a success Result.
func payload(value *big.Rat) Result.Interface {
	return Result.New(
		Result.WithPayload(
			&data{value: value},
		),
	)
}

// fromInterface extracts a *big.Rat from any Rational.Interface, parsing its
// "num/den" string when the concrete type is unknown (e.g. the Null-Object).
// The returned big.Rat is always a fresh copy, so operands are never shared.
func fromInterface(other Interface) *big.Rat {
	if d, ok := other.(*data); ok {
		return new(big.Rat).Set(d.value)
	}
	value, ok := new(big.Rat).SetString(other.ToGoString())
	if !ok {
		return new(big.Rat)
	}
	return value
}
