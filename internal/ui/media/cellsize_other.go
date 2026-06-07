//go:build !unix

package media

// detectCellAspect has no portable implementation off unix; callers fall back
// to defaultCellAspect.
func detectCellAspect() float64 { return 0 }
