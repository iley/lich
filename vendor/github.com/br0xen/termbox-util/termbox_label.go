package termboxUtil

import "github.com/nsf/termbox-go"

// Label is a field for inputting text
type Label struct {
	id                  string
	value               string
	x, y, width, height int
	cursor              int
	fg, bg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	bordered            bool
	wrap                bool
	multiline           bool
	active              bool
}

// CreateLabel creates an input field at x, y that is w by h
func CreateLabel(lbl string, x, y, w, h int, fg, bg termbox.Attribute) *Label {
	c := Label{
		value: lbl, x: x, y: y, width: w, height: h,
		fg: fg, bg: bg, activeFg: fg, activeBg: bg,
	}
	return &c
}

func (c *Label) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *Label) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *Label) SetActive(a bool)                      { c.active = a }
func (c *Label) IsActive() bool                        { return c.active }

// IsTabSkipped is always true for a label
func (c *Label) IsTabSkipped() bool { return true }

// This doesn't do anything for a label
func (c *Label) SetTabSkip(b bool) {}

// GetID returns this control's ID
func (c *Label) GetID() string { return c.id }

// SetID sets this control's ID
func (c *Label) SetID(newID string) { c.id = newID }

// GetValue gets the current text that is in the Label
func (c *Label) GetValue() string { return c.value }

// SetValue sets the current text in the Label to s
func (c *Label) SetValue(s string) { c.value = s }

// GetX returns the x position of the input field
func (c *Label) GetX() int { return c.x }

// SetX sets the x position of the input field
func (c *Label) SetX(x int) { c.x = x }

// GetY returns the y position of the input field
func (c *Label) GetY() int { return c.y }

// SetY sets the y position of the input field
func (c *Label) SetY(y int) { c.y = y }

// GetWidth returns the current width of the input field
func (c *Label) GetWidth() int {
	if c.width == -1 {
		if c.bordered {
			return len(c.value) + 2
		}
		return len(c.value)
	}
	return c.width
}

// SetWidth sets the current width of the input field
func (c *Label) SetWidth(w int) {
	c.width = w
}

// GetHeight returns the current height of the input field
func (c *Label) GetHeight() int { return c.height }

// SetHeight sets the current height of the input field
func (c *Label) SetHeight(h int) {
	c.height = h
}

// GetFgColor returns the foreground color
func (c *Label) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *Label) SetFgColor(fg termbox.Attribute) {
	c.fg = fg
}

// GetBgColor returns the background color
func (c *Label) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *Label) SetBgColor(bg termbox.Attribute) {
	c.bg = bg
}

// IsBordered returns true or false if this input field has a border
func (c *Label) IsBordered() bool { return c.bordered }

// SetBordered sets whether we render a border around the input field
func (c *Label) SetBordered(b bool) {
	c.bordered = b
}

// DoesWrap returns true or false if this input field wraps text
func (c *Label) DoesWrap() bool { return c.wrap }

// SetWrap sets whether we wrap the text at width.
func (c *Label) SetWrap(b bool) {
	c.wrap = b
}

// IsMultiline returns true or false if this field can have multiple lines
func (c *Label) IsMultiline() bool { return c.multiline }

// SetMultiline sets whether the field can have multiple lines
func (c *Label) SetMultiline(b bool) {
	c.multiline = b
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (c *Label) HandleEvent(event termbox.Event) bool { return false }

// Draw outputs the input field on the screen
func (c *Label) Draw() {
	maxWidth := c.width
	maxHeight := c.height
	x, y := c.x, c.y
	startX := c.x
	startY := c.y
	if c.bordered {
		DrawBorder(c.x, c.y, c.x+c.GetWidth(), c.y+c.height, c.fg, c.bg)
		maxWidth--
		maxHeight--
		x++
		y++
		startX++
		startY++
	}

	DrawStringAtPoint(c.value, x, y, c.fg, c.bg)
}
