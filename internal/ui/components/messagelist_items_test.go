package components

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDateLabel(t *testing.T) {
	now := time.Now()
	assert.Equal(t, "Today", FormatDateLabel(now), "same day reads Today")

	thisYear := time.Date(now.Year(), time.January, 2, 9, 0, 0, 0, time.Local)
	assert.Equal(t, "January 2", FormatDateLabel(thisYear), "current year omits the year")

	old := time.Date(2000, time.January, 2, 9, 0, 0, 0, time.Local)
	assert.Equal(t, "January 2, 2000", FormatDateLabel(old), "past years include the year")
}
