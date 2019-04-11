package schoolsout

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddUSHolidays(t *testing.T) {
	so := &Calendar{}
	AddUSHolidays(so)

	assert.Len(t, so.holidays, 10)
	assert.True(t, so.IsHoliday(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 1, 21, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 2, 18, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 5, 27, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 7, 4, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 9, 2, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 10, 14, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 11, 28, 0, 0, 0, 0, time.UTC)))
	assert.True(t, so.IsHoliday(time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC)))
}