package schoolsout

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Holiday_Def_US", func() {

	var (
		target *Calendar
	)

	BeforeEach(func() {
		target = &Calendar{}
	})

	Context("AddUSHoliday", func() {
		It("should add US holidays", func() {
			AddUSHolidays(target)

			Expect(target.holidays).To(HaveLen(10))

			Expect(target.IsHoliday(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 1, 21, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 2, 18, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 5, 27, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 7, 4, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 9, 2, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 10, 14, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 11, 28, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			Expect(target.IsHoliday(time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC))).To(BeTrue())
		})
	})
})
