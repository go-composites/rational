package main

import (
	"fmt"

	Error "github.com/go-composites/error/src"
	Rational "github.com/go-composites/rational/src"
	Result "github.com/go-composites/result/src"
)

func payloadOf(result Result.Interface) Rational.Interface {
	return result.Payload().(Rational.Interface)
}

func report(label string, result Result.Interface) {
	if result.HasError() {
		fmt.Printf("%s -> error: %s\n", label, result.Error().Message())
		return
	}
	fmt.Printf("%s -> %s\n", label, payloadOf(result).ToGoString())
}

func main() {
	oneThird := payloadOf(Rational.FromInts(1, 3))
	oneSixth := payloadOf(Rational.FromInts(1, 6))
	oneHalf := payloadOf(Rational.FromInts(1, 2))

	// Exactness: 1/3 + 1/6 is exactly 1/2, with no floating-point error.
	sum := oneThird.Add(oneSixth)
	report("1/3 + 1/6", sum)
	fmt.Println("1/3 + 1/6 == 1/2 exactly:", payloadOf(sum).Equal(oneHalf))

	// Contrast with floating point, which is NOT exact: 0.1 + 0.2 != 0.3.
	a, b, c := 0.1, 0.2, 0.3
	fmt.Printf("float: 0.1 + 0.2 == 0.3 : %v\n", (a+b) == c)

	three := payloadOf(Rational.FromInts(3, 4))
	report("3/4 - 1/4", three.Sub(payloadOf(Rational.FromInts(1, 4))))
	report("3/4 * 2/3", three.Mul(payloadOf(Rational.FromInts(2, 3))))
	report("3/4 / 1/2", three.Div(oneHalf))

	// Division by a zero rational is a value, not a panic.
	zero := payloadOf(Rational.FromInts(0, 1))
	divByZero := three.Div(zero)
	fmt.Println("3/4 / 0 has error:", divByZero.HasError())
	report("3/4 / 0", divByZero)

	// A zero denominator is a value, not a panic.
	bad := Rational.FromInts(1, 0)
	fmt.Println("1/0 has error:", bad.HasError())

	// Errors are first-class values.
	var _ Error.Interface = divByZero.Error()

	report("|-2/3|", payloadOf(Rational.FromInts(-2, 3)).Abs())
	report("-(2/3)", payloadOf(Rational.FromInts(2, 3)).Neg())

	fmt.Println("3/4 == 1/2 :", three.Equal(oneHalf))
	fmt.Println("1/2 < 3/4  :", oneHalf.LessThan(three))
	fmt.Println("3/4 > 1/2  :", three.GreaterThan(oneHalf))
	fmt.Println("3/4 as float:", three.ToFloat())
	fmt.Println(three.Inspect())
}
