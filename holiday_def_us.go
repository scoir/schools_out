package schoolsout

import "time"

// US Holidays
func AddUSHolidays(so *Calendar) {
	so.AddHoliday("New Years Day", FixedDay(1, time.January), true)
	so.AddHoliday("Martin Luther King Day", NthWeekdayOf(3, time.Monday, time.January), false)
	so.AddHoliday("President's Day", NthWeekdayOf(3, time.Monday, time.February), false)
	so.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)
	so.AddHoliday("Independence Day", FixedDay(4, time.July), false)
	so.AddHoliday("Labor Day", NthWeekdayOf(1, time.Monday, time.September), false)
	so.AddHoliday("Columbus Day", NthWeekdayOf(2, time.Monday, time.October), false)
	so.AddHoliday("Veteran's Day", FixedDay(11, time.November), false)
	so.AddHoliday("Thanksgiving", NthWeekdayOf(4, time.Thursday, time.November), false)
	so.AddHoliday("Christmas Day", FixedDay(25, time.December), false)
}
