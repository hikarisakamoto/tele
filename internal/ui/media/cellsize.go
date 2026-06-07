package media

import "sync"

// defaultCellAspect is the cell height:width ratio assumed when the terminal
// does not report its pixel size. Half-block art is built around this 2:1
// approximation; it is also a safe fallback for the Kitty path.
const defaultCellAspect = 2.0

var (
	cellAspectOnce sync.Once
	cellAspectVal  float64
)

// CellAspect returns the terminal's real cell aspect ratio (height/width),
// detected once from the controlling TTY. Falls back to defaultCellAspect when
// the terminal does not report pixel dimensions.
func CellAspect() float64 {
	cellAspectOnce.Do(func() {
		a := detectCellAspect()
		if a <= 0 {
			a = defaultCellAspect
		}
		cellAspectVal = a
	})
	return cellAspectVal
}
