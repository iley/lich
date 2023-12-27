package termboxUtil

import "github.com/nsf/termbox-go"

// ProgressBar Just contains the data needed to display a progress bar
type ProgressBar struct {
	id             string
	total          int
	progress       int
	allowOverflow  bool
	allowUnderflow bool
	fullChar       rune
	emptyChar      rune
	bordered       bool
	alignment      TextAlignment
	colorized      bool

	x, y               int
	width, height      int
	fg, bg             termbox.Attribute
	activeFg, activeBg termbox.Attribute
	active             bool
}

// CreateProgressBar Create a progress bar object
func CreateProgressBar(tot, x, y int, fg, bg termbox.Attribute) *ProgressBar {
	c := ProgressBar{total: tot,
		fullChar: '#', emptyChar: ' ',
		x: x, y: y, height: 1, width: 10,
		bordered: true, fg: fg, bg: bg,
		activeFg: fg, activeBg: bg,
		alignment: AlignLeft,
	}
	return &c
}

func (c *ProgressBar) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *ProgressBar) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *ProgressBar) SetActive(a bool)                      { c.active = a }
func (c *ProgressBar) IsActive() bool                        { return c.active }

// GetID returns this control's ID
func (c *ProgressBar) GetID() string { return c.id }

// SetID sets this control's ID
func (c *ProgressBar) SetID(newID string) {
	c.id = newID
}

// GetProgress returns the curret progress value
func (c *ProgressBar) GetProgress() int {
	return c.progress
}

// SetProgress sets the current progress of the bar
func (c *ProgressBar) SetProgress(p int) {
	if (p <= c.total || c.allowOverflow) || (p >= 0 || c.allowUnderflow) {
		c.progress = p
	}
}

// IncrProgress increments the current progress of the bar
func (c *ProgressBar) IncrProgress() {
	if c.progress < c.total || c.allowOverflow {
		c.progress++
	}
}

// DecrProgress decrements the current progress of the bar
func (c *ProgressBar) DecrProgress() {
	if c.progress > 0 || c.allowUnderflow {
		c.progress--
	}
}

// GetPercent returns the percent full of the bar
func (c *ProgressBar) GetPercent() int {
	return int(float64(c.progress) / float64(c.total) * 100)
}

// EnableOverflow Tells the progress bar that it can go over the total
func (c *ProgressBar) EnableOverflow() {
	c.allowOverflow = true
}

// DisableOverflow Tells the progress bar that it can NOT go over the total
func (c *ProgressBar) DisableOverflow() {
	c.allowOverflow = false
}

// EnableUnderflow Tells the progress bar that it can go below zero
func (c *ProgressBar) EnableUnderflow() {
	c.allowUnderflow = true
}

// DisableUnderflow Tells the progress bar that it can NOT go below zero
func (c *ProgressBar) DisableUnderflow() {
	c.allowUnderflow = false
}

// GetFullChar returns the rune used for 'full'
func (c *ProgressBar) GetFullChar() rune {
	return c.fullChar
}

// SetFullChar sets the rune used for 'full'
func (c *ProgressBar) SetFullChar(f rune) {
	c.fullChar = f
}

// GetEmptyChar gets the rune used for 'empty'
func (c *ProgressBar) GetEmptyChar() rune {
	return c.emptyChar
}

// SetEmptyChar sets the rune used for 'empty'
func (c *ProgressBar) SetEmptyChar(f rune) {
	c.emptyChar = f
}

// GetX Return the x position of the Progress Bar
func (c *ProgressBar) GetX() int { return c.x }

// SetX set the x position of the ProgressBar to x
func (c *ProgressBar) SetX(x int) {
	c.x = x
}

// GetY Return the y position of the ProgressBar
func (c *ProgressBar) GetY() int { return c.y }

// SetY Set the y position of the ProgressBar to y
func (c *ProgressBar) SetY(y int) {
	c.y = y
}

// GetHeight returns the height of the progress bar
// Defaults to 1 (3 if bordered)
func (c *ProgressBar) GetHeight() int {
	return c.height
}

// SetHeight Sets the height of the progress bar
func (c *ProgressBar) SetHeight(h int) {
	c.height = h
}

// GetWidth returns the width of the progress bar
func (c *ProgressBar) GetWidth() int {
	return c.width
}

// SetWidth Sets the width of the progress bar
func (c *ProgressBar) SetWidth(w int) {
	c.width = w
}

// GetFgColor returns the foreground color
func (c *ProgressBar) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *ProgressBar) SetFgColor(fg termbox.Attribute) {
	c.fg = fg
}

// GetBgColor returns the background color
func (c *ProgressBar) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *ProgressBar) SetBgColor(bg termbox.Attribute) {
	c.bg = bg
}

// Align Tells which direction the progress bar empties
func (c *ProgressBar) Align(a TextAlignment) {
	c.alignment = a
}

// SetColorized sets whether the progress bar should be colored
// depending on how full it is:
//  10% - Red
//	50% - Yellow
//	80% - Green
func (c *ProgressBar) SetColorized(color bool) {
	c.colorized = color
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (c *ProgressBar) HandleEvent(event termbox.Event) bool {
	return false
}

// Draw outputs the input field on the screen
func (c *ProgressBar) Draw() {
	// For now, just draw a [####  ] bar
	// TODO: make this more advanced
	useFg := c.fg
	if c.colorized {
		if c.GetPercent() < 10 {
			useFg = termbox.ColorRed
		} else if c.GetPercent() < 50 {
			useFg = termbox.ColorYellow
		} else {
			useFg = termbox.ColorGreen
		}
	}
	drawX, drawY := c.x, c.y
	fillWidth, fillHeight := c.width-2, c.height
	DrawStringAtPoint("[", drawX, drawY, c.fg, c.bg)
	numFull := int(float64(fillWidth) * float64(c.progress) / float64(c.total))
	FillWithChar(c.fullChar, drawX+1, drawY, drawX+1+numFull, drawY+(fillHeight-1), useFg, c.bg)
	DrawStringAtPoint("]", drawX+c.width-1, drawY, c.fg, c.bg)

	/*
		drawX, drawY := c.x, c.y
		drawWidth, drawHeight := c.width, c.height
		if c.bordered {
			if c.height == 1 && c.width > 2 {
				// Just using [ & ] for the border
				DrawStringAtPoint("[", drawX, drawY, c.fg, c.bg)
				DrawStringAtPoint("]", drawX+c.width-1, drawY, c.fg, c.bg)
				drawX++
				drawWidth -= 2
			} else if c.height >= 3 {
				DrawBorder(drawX, drawY, drawX+c.width, drawY+c.height, c.fg, c.bg)
				drawX++
				drawY++
				drawWidth -= 2
				drawHeight -= 2
			}
		}

		// Figure out how many chars are full
		numFull := drawWidth * (c.progress / c.total)
		switch c.alignment {
		case AlignRight: // TODO: Fill from right to left
		case AlignCenter: // TODO: Fill from middle out
		default: // Fill from left to right
			FillWithChar(c.fullChar, drawX, drawY, drawX+numFull, drawY+(drawHeight-1), c.fg, c.bg)
			if numFull < drawWidth {
				FillWithChar(c.emptyChar, drawX+numFull, drawY, drawX+drawWidth-1, drawY+(drawHeight-1), c.fg, c.bg)
			}
		}
	*/
}
