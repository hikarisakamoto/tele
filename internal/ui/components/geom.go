package components

// Rect is a rectangle in terminal cells. It describes the on-screen position of
// a UI element (e.g. the selected message bubble) so overlays can be anchored to
// it. Top/Left are the upper-left corner; Height/Width are the extents.
type Rect struct {
	Top, Left, Height, Width int
}
