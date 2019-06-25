package hermes

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hermes", func() {
	Describe("QueryString", func() {
		It("should get bool values", func() {
			req := newRequest()
			req.Raw().QueryArgs().Set("is_test", "true")

			qs := ParseQuery(req)
			Expect(qs.Bool("is_test")).To(BeTrue())
			Expect(qs.Bool("with_default", true)).To(BeTrue())
			Expect(qs.Bool("without_default")).To(BeZero())
		})

		It("should get string values", func() {
			req := newRequest()
			req.Raw().QueryArgs().Set("username", "gi.joe")

			qs := ParseQuery(req)
			Expect(qs.String("username")).To(Equal("gi.joe"))
			Expect(qs.String("with_default", "max.steel")).To(Equal("max.steel"))
			Expect(qs.String("without_default")).To(BeZero())
		})

		It("should get int values", func() {
			req := newRequest()
			req.Raw().QueryArgs().Set("age", "26")

			qs := ParseQuery(req)
			Expect(qs.Int("age")).To(Equal(26))
			Expect(qs.Int("with_default", 36)).To(Equal(36))
			Expect(qs.Int("without_default")).To(BeZero())
		})

		It("should get int64 values", func() {
			req := newRequest()
			req.Raw().QueryArgs().Set("duration", "86400000")

			qs := ParseQuery(req)
			Expect(qs.Int64("duration")).To(Equal(int64(86400000)))
			Expect(qs.Int64("with_default", 43200000)).To(Equal(int64(43200000)))
			Expect(qs.Int64("without_default")).To(BeZero())
		})

		It("should get float values", func() {
			req := newRequest()
			req.Raw().QueryArgs().Set("price", "3.14")

			qs := ParseQuery(req)
			Expect(qs.Float("price")).To(Equal(3.14))
			Expect(qs.Float("with_default", 4.90)).To(Equal(4.90))
			Expect(qs.Float("without_default")).To(BeZero())
		})
	})
})
