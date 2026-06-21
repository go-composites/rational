package Rational_test

import (
	Rational "github.com/go-composites/rational/src"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// payloadOf unwraps a success Result into a Rational.Interface.
func payloadOf(r interface {
	HasError() bool
	Payload() interface{}
}) Rational.Interface {
	gomega.ExpectWithOffset(1, r.HasError()).To(gomega.BeFalse())
	return r.Payload().(Rational.Interface)
}

// foreign is a Rational.Interface implementation that is NOT the package's own
// concrete type. It is used to exercise the string-bridging path of
// fromInterface with a value that DOES parse as a rational (the success
// branch).
type foreign struct{ s string }

func (f foreign) Numerator() int64                      { return 0 }
func (foreign) Denominator() int64                      { return 0 }
func (f foreign) ToGoString() string                    { return f.s }
func (foreign) ToFloat() float64                        { return 0 }
func (foreign) IsNull() bool                            { return false }
func (foreign) Add(Rational.Interface) Result.Interface { return nil }
func (foreign) Sub(Rational.Interface) Result.Interface { return nil }
func (foreign) Mul(Rational.Interface) Result.Interface { return nil }
func (foreign) Div(Rational.Interface) Result.Interface { return nil }
func (foreign) Abs() Result.Interface                   { return nil }
func (foreign) Neg() Result.Interface                   { return nil }
func (foreign) Equal(Rational.Interface) bool           { return false }
func (foreign) LessThan(Rational.Interface) bool        { return false }
func (foreign) GreaterThan(Rational.Interface) bool     { return false }
func (foreign) Inspect() Rational.String                { return `` }

var _ = ginkgo.Describe("Rational", func() {

	ginkgo.Describe("constructors", func() {
		ginkgo.It("builds a reduced fraction from two int64 values", func() {
			q := payloadOf(Rational.FromInts(2, 4))
			gomega.Expect(q.Numerator()).To(gomega.BeEquivalentTo(1))
			gomega.Expect(q.Denominator()).To(gomega.BeEquivalentTo(2))
			gomega.Expect(q.ToGoString()).To(gomega.Equal("1/2"))
			gomega.Expect(q.IsNull()).To(gomega.BeFalse())
		})
		ginkgo.It("returns an error Result on a zero denominator", func() {
			r := Rational.FromInts(1, 0)
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("zero denominator"))
		})
		ginkgo.It("parses a valid fraction string", func() {
			r := Rational.FromString("3/4")
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(Rational.Interface).ToGoString()).To(gomega.Equal("3/4"))
		})
		ginkgo.It("returns an error Result on a bad string", func() {
			r := Rational.FromString("not-a-number")
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("invalid rational"))
		})
		ginkgo.It("exposes a Null-Object", func() {
			n := Rational.Null()
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("exactness", func() {
		ginkgo.It("computes 1/3 + 1/6 as exactly 1/2", func() {
			a := payloadOf(Rational.FromInts(1, 3))
			b := payloadOf(Rational.FromInts(1, 6))
			sum := payloadOf(a.Add(b))
			gomega.Expect(sum.ToGoString()).To(gomega.Equal("1/2"))
			gomega.Expect(sum.Equal(payloadOf(Rational.FromInts(1, 2)))).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("conversions", func() {
		ginkgo.It("reports numerator and denominator of the reduced fraction", func() {
			q := payloadOf(Rational.FromInts(6, 8))
			gomega.Expect(q.Numerator()).To(gomega.BeEquivalentTo(3))
			gomega.Expect(q.Denominator()).To(gomega.BeEquivalentTo(4))
		})
		ginkgo.It("converts to a float64", func() {
			q := payloadOf(Rational.FromInts(1, 2))
			gomega.Expect(q.ToFloat()).To(gomega.BeNumerically("~", 0.5, 1e-12))
		})
	})

	ginkgo.Describe("arithmetic", func() {
		var threeQuarters = payloadOf(Rational.FromInts(3, 4))
		var oneHalf = payloadOf(Rational.FromInts(1, 2))

		ginkgo.It("adds", func() {
			gomega.Expect(payloadOf(threeQuarters.Add(oneHalf)).ToGoString()).To(gomega.Equal("5/4"))
		})
		ginkgo.It("subtracts", func() {
			gomega.Expect(payloadOf(threeQuarters.Sub(oneHalf)).ToGoString()).To(gomega.Equal("1/4"))
		})
		ginkgo.It("multiplies", func() {
			gomega.Expect(payloadOf(threeQuarters.Mul(oneHalf)).ToGoString()).To(gomega.Equal("3/8"))
		})
		ginkgo.It("divides", func() {
			gomega.Expect(payloadOf(threeQuarters.Div(oneHalf)).ToGoString()).To(gomega.Equal("3/2"))
		})
		ginkgo.It("does not mutate its operands", func() {
			_ = threeQuarters.Add(oneHalf)
			gomega.Expect(threeQuarters.ToGoString()).To(gomega.Equal("3/4"))
			gomega.Expect(oneHalf.ToGoString()).To(gomega.Equal("1/2"))
		})

		ginkgo.Describe("division by a zero rational", func() {
			ginkgo.It("returns a Result carrying an error instead of panicking", func() {
				r := threeQuarters.Div(payloadOf(Rational.FromInts(0, 1)))
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).To(gomega.Equal("division by zero"))
			})
		})

		ginkgo.Describe("absolute value", func() {
			ginkgo.It("makes a negative fraction positive", func() {
				gomega.Expect(payloadOf(payloadOf(Rational.FromInts(-2, 3)).Abs()).ToGoString()).
					To(gomega.Equal("2/3"))
			})
			ginkgo.It("leaves a positive fraction unchanged", func() {
				gomega.Expect(payloadOf(payloadOf(Rational.FromInts(2, 3)).Abs()).ToGoString()).
					To(gomega.Equal("2/3"))
			})
		})

		ginkgo.Describe("negation", func() {
			ginkgo.It("negates a positive fraction", func() {
				gomega.Expect(payloadOf(payloadOf(Rational.FromInts(2, 3)).Neg()).ToGoString()).
					To(gomega.Equal("-2/3"))
			})
			ginkgo.It("negates a negative fraction", func() {
				gomega.Expect(payloadOf(payloadOf(Rational.FromInts(-2, 3)).Neg()).ToGoString()).
					To(gomega.Equal("2/3"))
			})
		})
	})

	ginkgo.Describe("operations against a Null operand", func() {
		var threeQuarters = payloadOf(Rational.FromInts(3, 4))
		var null = Rational.Null()

		ginkgo.It("treats a Null operand as zero in addition", func() {
			gomega.Expect(payloadOf(threeQuarters.Add(null)).ToGoString()).To(gomega.Equal("3/4"))
		})
		ginkgo.It("guards division by a Null operand (zero)", func() {
			gomega.Expect(threeQuarters.Div(null).HasError()).To(gomega.BeTrue())
		})
		ginkgo.It("bridges a foreign Interface through its fraction string", func() {
			gomega.Expect(payloadOf(threeQuarters.Add(foreign{s: "1/4"})).ToGoString()).
				To(gomega.Equal("1"))
		})
		ginkgo.It("treats an unparsable foreign operand as zero", func() {
			gomega.Expect(payloadOf(threeQuarters.Add(foreign{s: "xx"})).ToGoString()).
				To(gomega.Equal("3/4"))
		})
	})

	ginkgo.Describe("comparisons", func() {
		var threeQuarters = payloadOf(Rational.FromInts(3, 4))
		var oneHalf = payloadOf(Rational.FromInts(1, 2))

		ginkgo.It("reports equality both ways", func() {
			gomega.Expect(threeQuarters.Equal(threeQuarters)).To(gomega.BeTrue())
			gomega.Expect(threeQuarters.Equal(oneHalf)).To(gomega.BeFalse())
		})
		ginkgo.It("reports less-than both ways", func() {
			gomega.Expect(oneHalf.LessThan(threeQuarters)).To(gomega.BeTrue())
			gomega.Expect(threeQuarters.LessThan(oneHalf)).To(gomega.BeFalse())
		})
		ginkgo.It("reports greater-than both ways", func() {
			gomega.Expect(threeQuarters.GreaterThan(oneHalf)).To(gomega.BeTrue())
			gomega.Expect(oneHalf.GreaterThan(threeQuarters)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("inspection", func() {
		ginkgo.It("renders a Rational", func() {
			gomega.Expect(payloadOf(Rational.FromInts(3, 4)).Inspect()).
				To(gomega.ContainSubstring("value=3/4"))
		})
	})

	ginkgo.Describe("the package-local Null-Object", func() {
		var n = Rational.Null()

		ginkgo.It("converts to zero values", func() {
			gomega.Expect(n.Numerator()).To(gomega.BeEquivalentTo(0))
			gomega.Expect(n.Denominator()).To(gomega.BeEquivalentTo(0))
			gomega.Expect(n.ToFloat()).To(gomega.BeEquivalentTo(0))
			gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
		})
		ginkgo.It("returns error Results for every arithmetic method", func() {
			gomega.Expect(n.Add(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Sub(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Mul(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Div(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Abs().HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Neg().HasError()).To(gomega.BeTrue())
		})
		ginkgo.It("compares as a Null-Object", func() {
			gomega.Expect(n.Equal(Rational.Null())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(payloadOf(Rational.FromInts(0, 1)))).To(gomega.BeFalse())
			gomega.Expect(n.LessThan(payloadOf(Rational.FromInts(1, 1)))).To(gomega.BeFalse())
			gomega.Expect(n.GreaterThan(payloadOf(Rational.FromInts(-1, 1)))).To(gomega.BeFalse())
		})
		ginkgo.It("inspects as the null marker", func() {
			gomega.Expect(n.Inspect()).To(gomega.Equal(`<NullRational>`))
		})
	})
})
