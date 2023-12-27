package termboxUtil

import "github.com/nsf/termbox-go"

// ScrollFrame is a frame for holding other elements
// It manages it's own x/y, tab index
type ScrollFrame struct {
	id                  string
	x, y, width, height int
	scrollX, scrollY    int
	tabIdx              int
	fg, bg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	bordered            bool
	controls            []termboxControl
	active              bool
}

// CreateScrollFrame creates Scrolling Frame at x, y that is w by h
func CreateScrollFrame(x, y, w, h int, fg, bg termbox.Attribute) *ScrollFrame {
	c := ScrollFrame{
		x: x, y: y, width: w, height: h,
		fg: fg, bg: bg, activeFg: fg, activeBg: bg,
	}
	return &c
}

// GetID returns this control's ID
func (c *ScrollFrame) GetID() string { return c.id }

// SetID sets this control's ID
func (c *ScrollFrame) SetID(newID string) {
	c.id = newID
}

// GetX returns the x position of the scroll frame
func (c *ScrollFrame) GetX() int { return c.x }

// SetX sets the x position of the scroll frame
func (c *ScrollFrame) SetX(x int) {
	c.x = x
}

// GetY returns the y position of the scroll frame
func (c *ScrollFrame) GetY() int { return c.y }

// SetY sets the y position of the scroll frame
func (c *ScrollFrame) SetY(y int) {
	c.y = y
}

// GetWidth returns the current width of the scroll frame
func (c *ScrollFrame) GetWidth() int { return c.width }

// SetWidth sets the current width of the scroll frame
func (c *ScrollFrame) SetWidth(w int) {
	c.width = w
}

// GetHeight returns the current height of the scroll frame
func (c *ScrollFrame) GetHeight() int { return c.height }

// SetHeight sets the current height of the scroll frame
func (c *ScrollFrame) SetHeight(h int) {
	c.height = h
}

// GetFgColor returns the foreground color
func (c *ScrollFrame) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *ScrollFrame) SetFgColor(fg termbox.Attribute) {
	c.fg = fg
}

// GetBgColor returns the background color
func (c *ScrollFrame) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *ScrollFrame) SetBgColor(bg termbox.Attribute) {
	c.bg = bg
}

// IsBordered returns true or false if this scroll frame has a border
func (c *ScrollFrame) IsBordered() bool { return c.bordered }

// SetBordered sets whether we render a border around the scroll frame
func (c *ScrollFrame) SetBordered(b bool) {
	c.bordered = b
}

// GetScrollX returns the x distance scrolled
func (c *ScrollFrame) GetScrollX() int {
	return c.scrollX
}

// GetScrollY returns the y distance scrolled
func (c *ScrollFrame) GetScrollY() int {
	return c.scrollY
}

// ScrollDown scrolls the frame down
func (c *ScrollFrame) ScrollDown() {
	c.scrollY++
}

// ScrollUp scrolls the frame up
func (c *ScrollFrame) ScrollUp() {
	if c.scrollY > 0 {
		c.scrollY--
	}
}

// ScrollLeft scrolls the frame left
func (c *ScrollFrame) ScrollLeft() {
	if c.scrollX > 0 {
		c.scrollX--
	}
}

// ScrollRight scrolls the frame right
func (c *ScrollFrame) ScrollRight() {
	c.scrollX++
}

// AddControl adds a control to the frame
func (c *ScrollFrame) AddControl(t termboxControl) {
	c.controls = append(c.controls, t)
}

// DrawControl figures out the relative position of the control,
// sets it, draws it, then resets it.
func (c *ScrollFrame) DrawControl(t termboxControl) {
	if c.IsVisible(t) {
		ctlX, ctlY := t.GetX(), t.GetY()
		t.SetX((c.GetX() + ctlX))
		t.SetY((c.GetY() + ctlY))
		t.Draw()
		t.SetX(ctlX)
		t.SetY(ctlY)
	}
}

// IsVisible takes a Termbox Control and returns whether
// that control would be visible in the frame
func (c *ScrollFrame) IsVisible(t termboxControl) bool {
	// Check if any part of t should be visible
	cX, cY := t.GetX(), t.GetY()
	if cX+t.GetWidth() >= c.scrollX && cX <= c.scrollX+c.width {
		return cY+t.GetHeight() >= c.scrollY && cY <= c.scrollY+c.height
	}
	return false
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (c *ScrollFrame) HandleEvent(event termbox.Event) bool {
	return false
}

// DrawToStrings generates a slice of strings with what should
// be drawn to the screen
func (c *ScrollFrame) DrawToStrings() []string {
	return []string{}
}

// Draw outputs the Scoll Frame on the screen
func (c *ScrollFrame) Draw() {
	maxWidth := c.width
	maxHeight := c.height
	x, y := c.x, c.y
	startX := c.x
	startY := c.y
	if c.bordered {
		DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, c.fg, c.bg)
		maxWidth--
		maxHeight--
		x++
		y++
		startX++
		startY++
	}
	for idx := range c.controls {
		c.DrawControl(c.controls[idx])
	}
}
