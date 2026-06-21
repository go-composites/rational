<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/rational" width="720"></p>

# rational

[![ci](https://github.com/go-composites/rational/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/rational/actions/workflows/ci.yml)

An **exact rational-number** composite for **Composition-Oriented Programming**.
A `Rational` wraps Go's `math/big.Rat` and mirrors Ruby's `Rational`: every
fraction is stored in lowest terms and is **exact** — `1/3 + 1/6` is exactly
`1/2`, with no floating-point error. Its arithmetic is exposed as **fallible
operations that return a `Result`** — so failures (a zero denominator, a
division by zero) are *values*, never panics and never `nil`.

```golang
quotient := numerator.Div(denominator)
if quotient.HasError() {
    fmt.Println(quotient.Error().Message()) // "division by zero"
} else {
    fmt.Println(quotient.Payload().(Rational.Interface).ToGoString())
}
```

`Rational` follows the org's Null-Object / never-nil invariant (enforced by the
`nonnil` CI analyzer): the `NullRational` variant in `src/null` satisfies the
same `Interface` and reports `IsNull() == true`.

## Install

```bash
export GOPRIVATE=github.com/go-composites GOPROXY=direct GOSUMDB=off
go get github.com/go-composites/rational@main
```

## Usage

> [!NOTE] main.go

```golang
package main

import (
    "fmt"

    Rational "github.com/go-composites/rational/src"
    Result "github.com/go-composites/result/src"
)

func payloadOf(r Result.Interface) Rational.Interface {
    return r.Payload().(Rational.Interface)
}

func main() {
    oneThird := payloadOf(Rational.FromInts(1, 3))
    oneSixth := payloadOf(Rational.FromInts(1, 6))
    oneHalf := payloadOf(Rational.FromInts(1, 2))

    // Exactness: 1/3 + 1/6 is exactly 1/2 — no floating-point error.
    sum := oneThird.Add(oneSixth)
    fmt.Println(payloadOf(sum).ToGoString())   // 1/2
    fmt.Println(payloadOf(sum).Equal(oneHalf)) // true

    // A zero denominator is a value, not a panic.
    bad := Rational.FromInts(1, 0)
    fmt.Println("has error:", bad.HasError()) // true

    // Division by a zero rational is a value, not a panic.
    div := oneHalf.Div(payloadOf(Rational.FromInts(0, 1)))
    fmt.Println(div.Error().Message()) // division by zero

    fmt.Println(oneHalf.LessThan(oneThird)) // false
    fmt.Println(payloadOf(Rational.FromInts(3, 4)).Inspect())
}
```

```bash
$ go run .
```

## API

Constructors

- `FromInts(num, den int64) Result.Interface` — build the reduced fraction
  `num/den`; a `Result` carrying `Error.New("zero denominator")` when `den` is
  zero.
- `FromString(s string) Result.Interface` — parse a fraction such as `"3/4"`; a
  `Result` carrying `Error.New(...)` when the input is not a valid rational.
- `Null() Interface` — the `NullRational` Null-Object (`IsNull() == true`).
- `null.New() Interface` — the importable `NullRational` Null-Object.

Conversions

- `Numerator() int64` / `Denominator() int64` — the reduced fraction's terms.
- `ToGoString() string` (`"num/den"`), `ToFloat() float64`, `IsNull() bool`.

Arithmetic (each returns `Result.Interface`)

- `Add(other)` / `Sub(other)` / `Mul(other)` — sum, difference, product.
- `Div(other)` — quotient; a `Result` carrying `Error.New("division by zero")`
  when `other` is a zero rational.
- `Abs()` / `Neg()` — absolute value and negation.

Every operation works on a fresh `big.Rat`, so operands are never mutated.

Comparisons (each returns `bool`)

- `Equal(other)` / `LessThan(other)` / `GreaterThan(other)`.

Inspection

- `Inspect() string` — `<Rational:0x... value=...>`.

## License

BSD-3-Clause — see [LICENSE](./LICENSE).
