package termboxUtil

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Frame is a frame for holding other elements
// It manages it's own x/y, tab index
type Frame struct {
	id                  string
	x, y, width, height int
	tabIdx              int
	fg, bg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	bordered            bool
	controls            []termboxControl
	tabSkip             bool
	active              bool
	title               string
	status              string
	rightStatus         string
}

// CreateFrame creates a Frame at x, y that is w by h
func CreateFrame(x, y, w, h int, fg, bg termbox.Attribute) *Frame {
	c := Frame{x: x, y: y, width: w, height: h,
		fg: fg, bg: bg, activeFg: fg, activeBg: bg,
		bordered: true,
	}
	return &c
}

func (c *Frame) SetTitle(title string) { c.title = title }

func (c *Frame) SetStatus(status string) { c.status = status }

// Setting color attributes on a frame trickles down to its controls
func (c *Frame) SetActiveFgColor(fg termbox.Attribute) {
	c.activeFg = fg
	for _, v := range c.controls {
		v.SetActiveFgColor(fg)
	}
}
func (c *Frame) SetActiveBgColor(bg termbox.Attribute) {
	c.activeBg = bg
	for _, v := range c.controls {
		v.SetActiveBgColor(bg)
	}
}
func (c *Frame) SetActive(a bool) {
	c.active = a
	for idx := range c.controls {
		if idx == c.tabIdx && a {
			c.controls[idx].SetActive(true)
		} else {
			c.controls[idx].SetActive(false)
		}
	}
}
func (c *Frame) IsActive() bool { return c.active }

// GetID returns this control's ID
func (c *Frame) GetID() string { return c.id }

// SetID sets this control's ID
func (c *Frame) SetID(newID string) { c.id = newID }

// GetX returns the x position of the frame
func (c *Frame) GetX() int { return c.x }

// SetX sets the x position of the frame
func (c *Frame) SetX(x int) { c.x = x }

// GetY returns the y position of the frame
func (c *Frame) GetY() int { return c.y }

// SetY sets the y position of the frame
func (c *Frame) SetY(y int) { c.y = y }

// GetWidth returns the current width of the frame
func (c *Frame) GetWidth() int { return c.width }

// SetWidth sets the current width of the frame
func (c *Frame) SetWidth(w int) { c.width = w }

// GetHeight returns the current height of the frame
func (c *Frame) GetHeight() int { return c.height }

// SetHeight sets the current height of the frame
func (c *Frame) SetHeight(h int) { c.height = h }

// GetFgColor returns the foreground color
func (c *Frame) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *Frame) SetFgColor(fg termbox.Attribute) { c.fg = fg }

// GetBgColor returns the background color
func (c *Frame) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *Frame) SetBgColor(bg termbox.Attribute) { c.bg = bg }

// IsBordered returns true or false if this frame has a border
func (c *Frame) IsBordered() bool { return c.bordered }

// SetBordered sets whether we render a border around the frame
func (c *Frame) SetBordered(b bool) { c.bordered = b }

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (c *Frame) IsTabSkipped() bool { return c.tabSkip }

// SetTabSkip sets the tabskip flag for this control
func (c *Frame) SetTabSkip(b bool) { c.tabSkip = b }

// AddControl adds a control to the frame
func (c *Frame) AddControl(t termboxControl) { c.controls = append(c.controls, t) }

func (c *Frame) ResetTabIndex() {
	for k, v := range c.controls {
		if !v.IsTabSkipped() {
			c.tabIdx = k
			break
		}
	}
}

// GetActiveControl returns the control at tabIdx
func (c *Frame) GetActiveControl() termboxControl {
	if len(c.controls) >= c.tabIdx {
		if c.controls[c.tabIdx].IsTabSkipped() {
			c.FindNextTabStop()
		}
		return c.controls[c.tabIdx]
	}
	return nil
}

// GetControls returns a slice of all controls
func (c *Frame) GetControls() []termboxControl {
	return c.controls
}

// GetControl returns the control at index i
func (c *Frame) GetControl(idx int) termboxControl {
	if len(c.controls) >= idx {
		return c.controls[idx]
	}
	return nil
}

// GetControlCount returns the number of controls contained
func (c *Frame) GetControlCount() int {
	return len(c.controls)
}

// GetLastControl returns the last control contained
func (c *Frame) GetLastControl() termboxControl {
	return c.controls[len(c.controls)-1]
}

// RemoveAllControls clears the control slice
func (c *Frame) RemoveAllControls() {
	c.controls = []termboxControl{}
}

// DrawControl figures out the relative position of the control,
// sets it, draws it, then resets it.
func (c *Frame) DrawControl(t termboxControl) {
	ctlX, ctlY := t.GetX(), t.GetY()
	t.SetX((c.GetX() + ctlX))
	t.SetY((c.GetY() + ctlY))
	t.Draw()
	t.SetX(ctlX)
	t.SetY(ctlY)
}

// GetBottomY returns the y of the lowest control in the frame
func (c *Frame) GetBottomY() int {
	var ret int
	for idx := range c.controls {
		if c.controls[idx].GetY()+c.controls[idx].GetHeight() > ret {
			ret = c.controls[idx].GetY() + c.controls[idx].GetHeight()
		}
	}
	return ret
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (c *Frame) HandleEvent(event termbox.Event) bool {
	// If the currently active control consumes the event, we don't need to handle it
	if c.controls[c.tabIdx].HandleEvent(event) {
		c.rightStatus = fmt.Sprintf("C (%d/%d)", c.tabIdx, len(c.controls))
		return true
	}
	// All that a frame cares about is tabbing around
	if event.Key == termbox.KeyTab {
		ret := !c.IsOnLastControl()
		c.FindNextTabStop()
		c.rightStatus = fmt.Sprintf("N (%d/%d)", c.tabIdx, len(c.controls))
		return ret
	}
	c.rightStatus = fmt.Sprintf("B (%d/%d)", c.tabIdx, len(c.controls))
	return false
}

// FindNextTabStop finds the next control that can be tabbed to
// A return of true means it found a different one than we started on.
func (c *Frame) FindNextTabStop() bool {
	startTab := c.tabIdx
	c.tabIdx = (c.tabIdx + 1) % len(c.controls)
	for c.controls[c.tabIdx].IsTabSkipped() {
		c.tabIdx = (c.tabIdx + 1) % len(c.controls)
		if c.tabIdx == startTab {
			break
		}
	}
	return c.tabIdx != startTab
}

// IsOnLastControl returns true if the active control
// is the last control that isn't tab skippable.
func (c *Frame) IsOnLastControl() bool {
	for _, v := range c.controls[c.tabIdx+1:] {
		if !v.IsTabSkipped() {
			return false
		}
	}
	return true
}

// Draw outputs the Scoll Frame on the screen
func (c *Frame) Draw() {
	maxWidth := c.width
	maxHeight := c.height
	x, y := c.x, c.y
	startX := c.x
	startY := c.y
	borderFg, borderBg := c.fg, c.bg
	if c.active {
		borderFg, borderBg = c.activeFg, c.activeBg
	}
	if c.bordered {
		// Clear the framed area
		FillWithChar(' ', c.x, c.y, c.x+c.width, c.y+c.height, borderFg, borderBg)
		if c.title == "" {
			DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, borderFg, borderBg)
		} else {
			DrawBorderWithTitle(c.x, c.y, c.x+c.width, c.y+c.height, " "+c.title+" ", borderFg, borderBg)
		}
		maxWidth--
		maxHeight--
		x++
		y++
		startX++
		startY++
	}
	for idx := range c.controls {
		if idx == c.tabIdx {
			c.controls[idx].SetActive(true)
		} else {
			c.controls[idx].SetActive(false)
		}
		c.DrawControl(c.controls[idx])
	}
	if c.status != "" {
		DrawStringAtPoint(" "+c.status+" ", c.x+1, c.y+c.height, borderFg, borderBg)
	}
	if c.rightStatus != "" {
		DrawStringAtPoint(" "+c.rightStatus+" ", (c.x+c.width)-len(c.rightStatus)-2, c.y+c.height, borderFg, borderBg)
	}
}
