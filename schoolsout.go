package schoolsout

import (
	"errors"
	"sync"
	"time"
)

// Calendar is used to determine and calculate applicable holidays.
type Calendar struct {
	sync.RWMutex

	DisableShiftSaturday bool
	DisableShiftSunday   bool
	holidays             []HolidayDefinition
}

// HolidayDefinition is the definition of a single Holiday.
type HolidayDefinition struct {
	Name              string
	calculation       DateCalculation
	checkForYearShift bool // Special case for New Year's Day
}

// Holiday is an instance of a holiday that occurs in a single year.
type Holiday struct {
	Name string
	Date time.Time
}

// ClearHolidays removes all previously loaded holidays.
func (so *Calendar) ClearHolidays() {
	so.Lock()
	defer so.Unlock()

	so.holidays = nil
}

// AddHoliday adds a new holiday definition to the Calendar instance.
func (so *Calendar) AddHoliday(name string, def DateCalculation, mayShiftYear bool) {
	so.Lock()
	defer so.Unlock()

	so.holidays = append(so.holidays, HolidayDefinition{
		Name:              name,
		calculation:       so.shiftForWeekend(def),
		checkForYearShift: mayShiftYear,
	})
}

// Unique list of the currently loaded holiday names (order is not guaranteed).
func (so *Calendar) ListHolidays() []string {
	names := make([]string, len(so.holidays))
	for i, h := range so.holidays {
		names[i] = h.Name
	}
	return names
}

// AllHolidaysForYear generates a list of all holidays for the specified year (order is not guaranteed).
func (so *Calendar) AllHolidaysForYear(year int) []Holiday {
	so.Lock()
	defer so.Unlock()

	c := []Holiday{}
	for _, def := range so.holidays {
		f := func(processYear int) {
			date := def.calculation(processYear)
			// Only push to slice if resulting year matches
			if date.Year() == year {
				c = append(c, Holiday{Name: def.Name, Date: date})
			}
		}

		//Execute for specified year
		f(year)

		//If this holiday might shift years, check previous & next year
		if def.checkForYearShift {
			f(year - 1)
			f(year + 1)
		}
	}
	return c
}

// HolidayDateForYears returns the date(s) applicable for the holiday over the specified years
func (so *Calendar) HolidayDateForYears(name string, years []int) ([]time.Time, error) {
	so.Lock()
	defer so.Unlock()

	for _, def := range so.holidays {
		if def.Name == name {
			d := []time.Time{}
			for _, year := range years {
				d = append(d, def.calculation(year))
			}
			return d, nil
		}
	}

	return nil, errors.New("holiday not found")
}

// IsHoliday returns true if the passed time.Time occurs on a holiday for the specified year.
func (so *Calendar) IsHoliday(date time.Time) bool {
	for _, h := range so.AllHolidaysForYear(date.Year()) {
		if h.Date.Month() == date.Month() && h.Date.Day() == date.Day() {
			return true
		}
	}

	return false
}

// shiftForWeekend will adjust a holiday to Friday or Sunday based to meet federal guidelines.
func (so *Calendar) shiftForWeekend(calc DateCalculation) DateCalculation {
	return func(year int) time.Time {
		date := calc(year)

		if !so.DisableShiftSaturday && date.Weekday() == time.Saturday {
			date = date.AddDate(0, 0, -1)
		}
		if !so.DisableShiftSunday && date.Weekday() == time.Sunday {
			date = date.AddDate(0, 0, 1)
		}

		return date
	}
}

type DateCalculation func(year int) time.Time

// FixedDay returns a function for calculating a fixed date in a particular year.
func FixedDay(day int, month time.Month) DateCalculation {
	return func(year int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
}

// NthWeekdayOf returns a function for determining the Nth (idx) weekday of the specified month in a particular year.
func NthWeekdayOf(idx int, dow time.Weekday, month time.Month) DateCalculation {
	return func(year int) time.Time {
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

		offset := int(firstOfMonth.Weekday()) - int(dow)
		if offset > 0 {
			offset = 7 - offset
		} else {
			offset = -offset
		}

		initialWeekdayOf := firstOfMonth.AddDate(0, 0, offset)

		return initialWeekdayOf.AddDate(0, 0, 7*(idx-1))
	}
}

// LastWeekdayOf returns a function for determining the last weekday of the specified month in a particular year.
func LastWeekdayOf(dow time.Weekday, month time.Month) DateCalculation {
	initFunc := NthWeekdayOf(1, dow, month)

	return func(year int) time.Time {
		lwdo := initFunc(year)
		adjLastDayOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -7)

		for lwdo.Before(adjLastDayOfMonth) {
			lwdo = lwdo.AddDate(0, 0, 7)
		}
		return lwdo
	}
}
