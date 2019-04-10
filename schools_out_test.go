package schools_out

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedDay(t *testing.T) {
	c := FixedDay(1, 1)
	assert.NotNil(t, c)

	d := c(2000)
	assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), d)
}

func TestNthWeekdayOf(t *testing.T) {
	c := NthWeekdayOf(1, time.Monday, time.January)
	assert.NotNil(t, c)

	d := c(2000)
	assert.Equal(t, time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), d)

	d = c(2001)
	assert.Equal(t, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), d)
}

func TestNthWeekdayOf_NegativeOffset(t *testing.T) {
	c := NthWeekdayOf(5, time.Wednesday, time.January)
	assert.NotNil(t, c)

	d := c(2001)
	assert.Equal(t, time.Date(2001, 1, 31, 0, 0, 0, 0, time.UTC), d)
}


func TestLastWeekdayOf(t *testing.T) {
	c := LastWeekdayOf(time.Monday, time.January)
	assert.NotNil(t, c)

	d := c(2000)
	assert.Equal(t, time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), d)
}

func TestLastWeekdayOf_LeapYear(t *testing.T) {
	c := LastWeekdayOf(time.Tuesday, time.February)
	assert.NotNil(t, c)

	d := c(2000)
	assert.Equal(t, time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC), d)
}

func TestSchoolsOut_shiftForWeekend_Saturday(t *testing.T) {
	so := SchoolsOut{}
	sf := so.shiftForWeekend(FixedDay(1, time.January))

	d := sf(2000)
	assert.Equal(t, time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), d)
}

func TestSchoolsOut_shiftForWeekend_Sunday(t *testing.T) {
	so := SchoolsOut{}
	sf := so.shiftForWeekend(FixedDay(31, time.December))

	d := sf(2000)
	assert.Equal(t, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), d)
}

func TestSchoolsOut_shiftForWeekend_Disabled(t *testing.T) {
	so := SchoolsOut{
		DisableShiftSunday: true,
	}
	sf := so.shiftForWeekend(FixedDay(31, time.December))

	d := sf(2000)
	assert.Equal(t, time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC), d)
}

func TestSchoolsOut_AddHoliday(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)

	assert.Len(t, so.holidays, 1)
	assert.Equal(t, "New Years", so.holidays[0].Name)
	assert.NotNil(t, so.holidays[0].calculation)
	assert.True(t, so.holidays[0].checkForYearShift)
}

func TestSchoolsOut_ClearHolidays(t *testing.T) {
	so := SchoolsOut{}
	assert.Len(t, so.holidays, 0)

	so.AddHoliday("New Years", FixedDay(1, time.January), true)
	assert.Len(t, so.holidays, 1)

	so.ClearHolidays()
	assert.Len(t, so.holidays, 0)
}

func TestSchoolsOut_AllHolidaysForYear(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)
	so.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)
	so.AddHoliday("Thanksgiving", NthWeekdayOf(4, time.Thursday, time.November), false)

	r :=  so.AllHolidaysForYear(2001)

	assert.Len(t, r, 3)
	assert.Equal(t, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), r[0].Date)
	assert.Equal(t, time.Date(2001, 5, 28, 0, 0, 0, 0, time.UTC), r[1].Date)
	assert.Equal(t, time.Date(2001, 11, 22, 0, 0, 0, 0, time.UTC), r[2].Date)
}

func TestSchoolsOut_AllHolidaysForYear_Double(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)

	r := so.AllHolidaysForYear(1999)

	assert.Len(t, r, 2, "there should be two new years holidays in 1999")
	assert.Equal(t, time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC), r[0].Date)
	assert.Equal(t, time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), r[1].Date)
}

func TestSchoolsOut_AllHolidaysForYear_ShiftYearBack(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)

	r := so.AllHolidaysForYear(2000)
	assert.Len(t, r, 0, "new years of 2000, shifted to 1999 so we shouldn't have results")
}

func TestSchoolsOut_AllHolidaysForYear_ShiftYearForward(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years Eve", FixedDay(31, time.December), true)

	r := so.AllHolidaysForYear(2000)
	assert.Len(t, r, 0, "new years eve of 2000, shifted to 2001 so we shouldn't have results")
}

func TestSchoolsOut_IsHoliday(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)
	so.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)
	so.AddHoliday("Thanksgiving", NthWeekdayOf(4, time.Thursday, time.November), false)

	assert.True(t, so.IsHoliday(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 5, 27, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 11, 28, 0, 0, 0, 0, time.UTC)))

	assert.False(t, so.IsHoliday(time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC)))
	assert.False(t, so.IsHoliday(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)))
}

func TestSchoolsOut_ListHolidays(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)
	so.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)

	l := so.ListHolidays()
	assert.Len(t, l, 2)

	assert.Equal(t, "New Years", l[0])
	assert.Equal(t, "Memorial Day", l[1])
}

func TestSchoolsOut_ListHolidays_NoHolidays(t *testing.T) {
	so := SchoolsOut{}
	l := so.ListHolidays()
	assert.Len(t, l, 0)
}

func TestSchoolsOut_HolidayDateForYears(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)
	so.AddHoliday("Memorial Day", LastWeekdayOf(time.Monday, time.May), false)

	r, err := so.HolidayDateForYears("New Years", []int{1999, 2000})
	assert.NoError(t, err)
	assert.Len(t, r, 2, "there should be two new years holidays between 1999 & 2000")
	assert.Equal(t, time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC), r[0])
	assert.Equal(t, time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), r[1])
}

func TestSchoolsOut_HolidayDateForYears_NotFound(t *testing.T) {
	so := SchoolsOut{}
	so.AddHoliday("New Years", FixedDay(1, time.January), true)

	_, err := so.HolidayDateForYears("Festivus", []int{1999, 2000})

	assert.Error(t, err)
}

