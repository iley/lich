package termboxUtil

import termbox "github.com/nsf/termbox-go"

type Checkbox struct {
	id                  string
	title               string
	isChecked           bool
	x, y, width, height int
	fg, bg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	bordered            bool
	tabSkip             bool
	active              bool
}

func CreateCheckbox(lbl string, x, y, w, h int, fg, bg termbox.Attribute) *Checkbox {
	c := Checkbox{
		x: x, y: y, width: w, height: h,
		fg: fg, bg: bg, activeFg: fg, activeBg: bg,
	}
	return &c
}

func (c *Checkbox) SetTitle(title string)                 { c.title = title }
func (c *Checkbox) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *Checkbox) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *Checkbox) SetActive(a bool)                      { c.active = a }
func (c *Checkbox) IsActive() bool                        { return c.active }

// GetID returns this control's ID
func (c *Checkbox) GetID() string { return c.id }

// SetID sets this control's ID
func (c *Checkbox) SetID(newID string) {
	c.id = newID
}

// GetX returns the x position of the input field
func (c *Checkbox) GetX() int { return c.x }

// SetX sets the x position of the input field
func (c *Checkbox) SetX(x int) {
	c.x = x
}

// GetY returns the y position of the input field
func (c *Checkbox) GetY() int { return c.y }

// SetY sets the y position of the input field
func (c *Checkbox) SetY(y int) {
	c.y = y
}

// GetWidth returns the current width of the input field
func (c *Checkbox) GetWidth() int { return c.width }

// SetWidth sets the current width of the input field
func (c *Checkbox) SetWidth(w int) {
	c.width = w
}

// GetHeight returns the current height of the input field
func (c *Checkbox) GetHeight() int { return c.height }

// SetHeight sets the current height of the input field
func (c *Checkbox) SetHeight(h int) {
	c.height = h
}

// GetFgColor returns the foreground color
func (c *Checkbox) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *Checkbox) SetFgColor(fg termbox.Attribute) {
	c.fg = fg
}

// GetBgColor returns the background color
func (c *Checkbox) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *Checkbox) SetBgColor(bg termbox.Attribute) {
	c.bg = bg
}

// IsBordered returns true or false if this input field has a border
func (c *Checkbox) IsBordered() bool { return c.bordered }

// SetBordered sets whether we render a border around the input field
func (c *Checkbox) SetBordered(b bool) {
	c.bordered = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (c *Checkbox) IsTabSkipped() bool {
	return c.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (c *Checkbox) SetTabSkip(b bool) {
	c.tabSkip = b
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (c *Checkbox) HandleEvent(event termbox.Event) bool {
	if event.Ch == 0 {
		switch event.Key {
		case termbox.KeySpace, termbox.KeyEnter:
			c.isChecked = !c.isChecked
			return true
		}
	}
	return false
}

func (c *Checkbox) Draw() {
	x, y, _, w := c.x, c.y, c.height, c.width
	useFg, useBg := c.fg, c.bg
	if c.active {
		useFg, useBg = c.activeFg, c.activeBg
	}
	if c.bordered {
		DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, useFg, useBg)
		x++
		w = w - 2
	}

	if c.isChecked {
		DrawStringAtPoint("[X]", x, y, useFg, useBg)
	} else {
		DrawStringAtPoint("[ ]", x, y, useFg, useBg)
	}
	x = x + 3
	w = w - 3

	if c.title != "" {
		wrk := c.title
		if len(wrk) > w {
			wrk = wrk[:w]
		}
		DrawStringAtPoint(wrk, x, y, useFg, useBg)
	}
}
