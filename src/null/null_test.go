package NullRational_test

import (
	Rational "github.com/go-composites/rational/src"
	NullRational "github.com/go-composites/rational/src/null"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func someRational() Rational.Interface {
	r := Rational.FromInts(0, 1)
	gomega.ExpectWithOffset(1, r.HasError()).To(gomega.BeFalse())
	return r.Payload().(Rational.Interface)
}

var _ = ginkgo.Describe("NullRational", func() {
	var n NullRational.Interface
	ginkgo.BeforeEach(func() {
		n = NullRational.New()
	})

	ginkgo.It("satisfies the Rational interface", func() {
		var _ Rational.Interface = n
	})
	ginkgo.It("reports IsNull() true", func() {
		gomega.Expect(n.IsNull()).To(gomega.BeTrue())
	})
	ginkgo.It("converts to zero values", func() {
		gomega.Expect(n.Numerator()).To(gomega.BeEquivalentTo(0))
		gomega.Expect(n.Denominator()).To(gomega.BeEquivalentTo(0))
		gomega.Expect(n.ToFloat()).To(gomega.BeEquivalentTo(0))
		gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
	})

	ginkgo.It("Add returns an error result", func() {
		r := n.Add(someRational())
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Add"))
	})
	ginkgo.It("Sub returns an error result", func() {
		r := n.Sub(someRational())
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Sub"))
	})
	ginkgo.It("Mul returns an error result", func() {
		r := n.Mul(someRational())
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Mul"))
	})
	ginkgo.It("Div returns an error result", func() {
		r := n.Div(someRational())
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Div"))
	})
	ginkgo.It("Abs returns an error result", func() {
		r := n.Abs()
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Abs"))
	})
	ginkgo.It("Neg returns an error result", func() {
		r := n.Neg()
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Neg"))
	})
	ginkgo.It("Equal is true only against another null", func() {
		gomega.Expect(n.Equal(NullRational.New())).To(gomega.BeTrue())
		gomega.Expect(n.Equal(someRational())).To(gomega.BeFalse())
	})
	ginkgo.It("LessThan is always false", func() {
		gomega.Expect(n.LessThan(someRational())).To(gomega.BeFalse())
	})
	ginkgo.It("GreaterThan is always false", func() {
		gomega.Expect(n.GreaterThan(someRational())).To(gomega.BeFalse())
	})
	ginkgo.It("Inspect renders the null marker", func() {
		gomega.Expect(n.Inspect()).To(gomega.Equal(`<NullRational>`))
	})
})
