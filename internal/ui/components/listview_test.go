package components_test

import (
	"testing"

	"github.com/sorokin-vladimir/tele/internal/ui/components"
	"github.com/stretchr/testify/assert"
)

func TestListView_SetCount_ClampsCursor(t *testing.T) {
	l := components.NewListView(false)
	l.SetCount(10)
	l.SetCursor(7)
	l.SetCount(3) // list shrank
	assert.Equal(t, 2, l.Cursor())
}

func TestListView_MoveDownUp_StopAtEnds_NoWrap(t *testing.T) {
	l := components.NewListView(false)
	l.SetCount(3)
	l.MoveUp()
	assert.Equal(t, 0, l.Cursor())
	l.MoveDown()
	l.MoveDown()
	l.MoveDown() // past end
	assert.Equal(t, 2, l.Cursor())
}

func TestListView_MoveWraps_WhenEnabled(t *testing.T) {
	l := components.NewListView(true)
	l.SetCount(3)
	l.MoveUp() // wrap to last
	assert.Equal(t, 2, l.Cursor())
	l.MoveDown() // wrap to first
	assert.Equal(t, 0, l.Cursor())
}

func TestListView_SkipsNonSelectable(t *testing.T) {
	l := components.NewListView(false)
	l.SetSelectable(func(i int) bool { return i != 1 }) // index 1 is a separator
	l.SetCount(3)
	l.MoveDown() // 0 -> skip 1 -> 2
	assert.Equal(t, 2, l.Cursor())
}

func TestListView_SetCount_LandsOnNearestSelectable(t *testing.T) {
	l := components.NewListView(false)
	l.SetSelectable(func(i int) bool { return i == 2 })
	l.SetCount(3)
	assert.Equal(t, 2, l.Cursor()) // 0 not selectable -> nearest is 2
}

func TestListView_Render_AllRowsWhenFits(t *testing.T) {
	l := components.NewListView(false)
	l.SetCount(3)
	rows := l.Render(8, func(i int, selected bool) string {
		if selected {
			return "*"
		}
		return "."
	})
	assert.Equal(t, []string{"*", ".", "."}, rows)
}

func TestListView_Render_CentersCursor_WhenScrolling(t *testing.T) {
	l := components.NewListView(false)
	l.SetCount(10)
	l.SetCursor(5)
	rows := l.Render(4, func(i int, selected bool) string {
		if selected {
			return "X"
		}
		return "o"
	})
	// maxRows=4, cursor=5 -> offset = 5 - 2 = 3 -> rows for i=3,4,5,6
	assert.Len(t, rows, 4)
	assert.Equal(t, "X", rows[2]) // cursor (i=5) is the 3rd visible row
}

func TestListView_Render_ClampsAtBottom(t *testing.T) {
	l := components.NewListView(false)
	l.SetCount(10)
	l.SetCursor(9)
	rows := l.Render(4, func(i int, selected bool) string {
		if selected {
			return "X"
		}
		return "o"
	})
	// offset clamped to 10-4 = 6 -> i=6,7,8,9; cursor=9 is last row
	assert.Len(t, rows, 4)
	assert.Equal(t, "X", rows[3])
}

func TestListView_Render_EmptyWhenNoItems(t *testing.T) {
	l := components.NewListView(false)
	l.SetCount(0)
	rows := l.Render(8, func(i int, selected bool) string { return "x" })
	assert.Empty(t, rows)
}
