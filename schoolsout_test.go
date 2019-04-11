package schoolsout

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SchoolsOut", func() {

	Describe("Holiday Calculators", func() {

		Context("FixedDay", func() {
			It("Should return a function for calculating a specific day", func() {
				c := FixedDay(1, 1)
				Expect(c).ShouldNot(BeNil())

				d := c(2000)
				Expect(d).Should(Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)))
			})
		})

		Context("NthWeekdayOf", func() {
			It("Should return a function for calculating a specific day", func() {
				c := NthWeekdayOf(1, time.Monday, time.January)
				Expect(c).ShouldNot(BeNil())

				d := c(2000)
				Expect(d).Should(Equal(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)))

				d = c(2001)
				Expect(d).Should(Equal(time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)))
			})

			It("Should support calculating a specific day with a negative offset", func() {
				c := NthWeekdayOf(5, time.Wednesday, time.January)
				Expect(c).ShouldNot(BeNil())

				d := c(2001)
				Expect(d).Should(Equal(time.Date(2001, 1, 31, 0, 0, 0, 0, time.UTC)))
			})
		})

		Context("LastWeekdayOf", func() {
			It("Should return a function for calculating a specific day", func() {
				c := LastWeekdayOf(time.Monday, time.January)
				Expect(c).ShouldNot(BeNil())

				d := c(2000)
				Expect(d).Should(Equal(time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC)))
			})

			It("Should return a function for calculating a specific day, even with Leap Years", func() {
				c := LastWeekdayOf(time.Tuesday, time.February)
				Expect(c).ShouldNot(BeNil())

				d := c(2000)
				Expect(d).Should(Equal(time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)))
			})
		})

	})

	Describe("Calendar", func() {
		var (
			target *Calendar
		)

		BeforeEach(func() {
			target = &Calendar{}
		})

		Context("shiftForWeekend", func() {
			It("should shift saturday to friday", func() {
				sf := target.shiftForWeekend(FixedDay(1, time.January))
				d := sf(2000)

				Expect(d).Should(Equal(time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)))
			})

			It("should shift sunday to monday", func() {
				sf := target.shiftForWeekend(FixedDay(31, time.December))
				d := sf(2000)

				Expect(d).Should(Equal(time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)))
			})

			It("shouldn't shift if shifting is disabled", func() {
				target.DisableShiftSunday = true

				sf := target.shiftForWeekend(FixedDay(31, time.December))
				d := sf(2000)

				Expect(d).Should(Equal(time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC)))
			})
		})

		Context("AddHoliday", func() {
			It("can add a holiday", func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)

				Expect(target.holidays).To(HaveLen(1))
				Expect(target.holidays[0].Name).To(Equal("New Years"))
				Expect(target.holidays[0].calculation).ToNot(BeNil())
				Expect(target.holidays[0].checkForYearShift).To(BeTrue())
			})
		})

		Context("ClearHolidays", func() {
			It("clears the list of holidays", func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)
				Expect(target.holidays).To(HaveLen(1))

				target.ClearHolidays()
				Expect(target.holidays).To(HaveLen(0))
			})
		})

		Context("AllHolidaysForYear", func() {
			It("calculates for specified holidays", func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)
				target.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)
				target.AddHoliday("Thanksgiving", NthWeekdayOf(4, time.Thursday, time.November), false)

				r := target.AllHolidaysForYear(2001)

				Expect(r).To(HaveLen(3))
				Expect(r).To(ContainElement(Equal(Holiday{Name: "New Years", Date: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)})))
				Expect(r).To(ContainElement(Equal(Holiday{Name: "Memorial Day", Date: time.Date(2001, 5, 28, 0, 0, 0, 0, time.UTC)})))
				Expect(r).To(ContainElement(Equal(Holiday{Name: "Thanksgiving", Date: time.Date(2001, 11, 22, 0, 0, 0, 0, time.UTC)})))
			})

			It("calculates for holidays that appear twice in a year", func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)

				r := target.AllHolidaysForYear(1999)

				Expect(r).To(HaveLen(2))
				Expect(r).To(ContainElement(Equal(Holiday{Name: "New Years", Date: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)})))
				Expect(r).To(ContainElement(Equal(Holiday{Name: "New Years", Date: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)})))
			})

			It("excludes holidays that shift to back a year (new years of 2000, shifted to 1999 so we shouldn't have results)", func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)

				r := target.AllHolidaysForYear(2000)

				Expect(r).To(HaveLen(0))
			})

			It("excludes holidays that shift to forward a year (new years eve of 2000, shifted to 2001 so we shouldn't have results)", func() {
				target.AddHoliday("New Years Eve", FixedDay(31, time.December), true)

				r := target.AllHolidaysForYear(2000)

				Expect(r).To(HaveLen(0))
			})
		})

		Context("IsHoliday", func() {

			BeforeEach(func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)
				target.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)
				target.AddHoliday("Thanksgiving", NthWeekdayOf(4, time.Thursday, time.November), false)
			})

			It("should match for new years", func() {
				Expect(target.IsHoliday(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			})

			It("should match for memorial day", func() {
				Expect(target.IsHoliday(time.Date(2019, 5, 27, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			})

			It("should match for thanksgiving", func() {
				Expect(target.IsHoliday(time.Date(2019, 11, 28, 0, 0, 0, 0, time.UTC))).To(BeTrue())
			})

			It("should not match for random day", func() {
				Expect(target.IsHoliday(time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC))).To(BeFalse())
			})

			It("should not match for another random day", func() {
				Expect(target.IsHoliday(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC))).To(BeFalse())
			})
		})

		Context("ListHolidays", func() {
			It("should return a list of registered holidays", func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)
				target.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)

				l := target.ListHolidays()

				Expect(l).To(HaveLen(2))
				Expect(l).To(ConsistOf([]string{"New Years", "Memorial Day"}))
			})

			It("should return empty if there are no registered holidays", func() {
				Expect(target.ListHolidays()).To(HaveLen(0))
			})
		})

		Context("HolidayDateForYears", func() {

			BeforeEach(func() {
				target.AddHoliday("New Years", FixedDay(1, time.January), true)
				target.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)
			})

			It("should return a list of dates applicable for the specified holiday", func() {
				r, err := target.HolidayDateForYears("New Years", []int{1999, 2000})

				Expect(err).ToNot(HaveOccurred())
				Expect(r).To(HaveLen(2))

				Expect(r).To(ConsistOf([]time.Time{
					time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC),
				}))
			})

			It("should return error when holiday is not found", func() {
				_, err := target.HolidayDateForYears("Festivus", []int{1999, 2000})

				Expect(err).To(HaveOccurred())
			})
		})

	})
})
