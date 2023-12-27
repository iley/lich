package termboxUtil

import (
	"github.com/nsf/termbox-go"
)

// InputModal A modal for text input
type InputModal struct {
	id                  string
	title               string
	text                string
	input               *InputField
	x, y, width, height int
	showHelp            bool
	cursor              int
	bg, fg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	isDone              bool
	isAccepted          bool
	isVisible           bool
	bordered            bool
	tabSkip             bool
	inputSelected       bool
	active              bool
}

// CreateInputModal Create an input modal with the given attributes
func CreateInputModal(title string, x, y, width, height int, fg, bg termbox.Attribute) *InputModal {
	c := InputModal{title: title, x: x, y: y, width: width, height: height, fg: fg, bg: bg, bordered: true}
	c.input = CreateInputField(c.x+2, c.y+3, c.width-2, 2, c.fg, c.bg)
	c.showHelp = true
	c.input.bordered = true
	c.isVisible = true
	c.inputSelected = true
	return &c
}

func (c *InputModal) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *InputModal) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *InputModal) SetActive(a bool)                      { c.active = a }
func (c *InputModal) IsActive() bool                        { return c.active }

// GetID returns this control's ID
func (c *InputModal) GetID() string { return c.id }

// SetID sets this control's ID
func (c *InputModal) SetID(newID string) {
	c.id = newID
}

// GetTitle Return the title of the modal
func (c *InputModal) GetTitle() string { return c.title }

// SetTitle Sets the title of the modal to s
func (c *InputModal) SetTitle(s string) {
	c.title = s
}

// GetText Return the text of the modal
func (c *InputModal) GetText() string { return c.text }

// SetText Set the text of the modal to s
func (c *InputModal) SetText(s string) {
	c.text = s
}

// GetX Return the x position of the modal
func (c *InputModal) GetX() int { return c.x }

// SetX set the x position of the modal to x
func (c *InputModal) SetX(x int) {
	c.x = x
}

// GetY Return the y position of the modal
func (c *InputModal) GetY() int { return c.y }

// SetY Set the y position of the modal to y
func (c *InputModal) SetY(y int) {
	c.y = y
}

// GetWidth Return the width of the modal
func (c *InputModal) GetWidth() int { return c.width }

// SetWidth Set the width of the modal to width
func (c *InputModal) SetWidth(width int) {
	c.width = width
}

// GetHeight Return the height of the modal
func (c *InputModal) GetHeight() int { return c.height }

// SetHeight Set the height of the modal to height
func (c *InputModal) SetHeight(height int) {
	c.height = height
}

// SetMultiline returns whether this is a multiline modal
func (c *InputModal) SetMultiline(m bool) {
	c.input.multiline = m
}

// IsMultiline returns whether this is a multiline modal
func (c *InputModal) IsMultiline() bool {
	return c.input.multiline
}

// IsBordered returns whether this control is bordered or not
func (c *InputModal) IsBordered() bool {
	return c.bordered
}

// SetBordered sets whether we render a border around the frame
func (c *InputModal) SetBordered(b bool) {
	c.bordered = b
}

// IsTabSkipped returns whether this control has it's tabskip flag set
func (c *InputModal) IsTabSkipped() bool {
	return c.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (c *InputModal) SetTabSkip(b bool) {
	c.tabSkip = b
}

// HelpIsShown Returns whether the modal is showing it's help text or not
func (c *InputModal) HelpIsShown() bool { return c.showHelp }

// ShowHelp Set the "Show Help" flag
func (c *InputModal) ShowHelp(b bool) {
	c.showHelp = b
}

// GetFgColor returns the foreground color
func (c *InputModal) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *InputModal) SetFgColor(fg termbox.Attribute) {
	c.fg = fg
}

// GetBgColor returns the background color
func (c *InputModal) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *InputModal) SetBgColor(bg termbox.Attribute) {
	c.bg = bg
}

// Show Sets the visibility flag to true
func (c *InputModal) Show() {
	c.isVisible = true
}

// Hide Sets the visibility flag to false
func (c *InputModal) Hide() {
	c.isVisible = false
}

// IsVisible returns the isVisible flag
func (c *InputModal) IsVisible() bool {
	return c.isVisible
}

// SetDone Sets the flag that tells whether this modal has completed it's purpose
func (c *InputModal) SetDone(b bool) {
	c.isDone = b
}

// IsDone Returns the "isDone" flag
func (c *InputModal) IsDone() bool {
	return c.isDone
}

// IsAccepted Returns whether the modal has been accepted
func (c *InputModal) IsAccepted() bool {
	return c.isAccepted
}

// GetValue Return the current value of the input
func (c *InputModal) GetValue() string { return c.input.GetValue() }

// SetValue Sets the value of the input to s
func (c *InputModal) SetValue(s string) {
	c.input.SetValue(s)
}

// SetInputWrap sets whether the input field will wrap long text or not
func (c *InputModal) SetInputWrap(b bool) {
	c.input.SetWrap(b)
}

// Clear Resets all non-positional parameters of the modal
func (c *InputModal) Clear() {
	c.title = ""
	c.text = ""
	c.input.SetValue("")
	c.isDone = false
	c.isVisible = false
}

// HandleEvent Handle the termbox event, return true if it was consumed
func (c *InputModal) HandleEvent(event termbox.Event) bool {
	if event.Key == termbox.KeyEnter {
		if !c.input.IsMultiline() || !c.inputSelected {
			// Done editing
			c.isDone = true
			c.isAccepted = true
		} else {
			c.input.HandleEvent(event)
		}
		return true
	} else if event.Key == termbox.KeyTab {
		if c.input.IsMultiline() {
			c.inputSelected = !c.inputSelected
		}
	} else if event.Key == termbox.KeyEsc {
		// Done editing
		c.isDone = true
		c.isAccepted = false
		return true
	}
	return c.input.HandleEvent(event)
}

// Draw Draw the modal
func (c *InputModal) Draw() {
	if c.isVisible {
		// First blank out the area we'll be putting the modal
		FillWithChar(' ', c.x, c.y, c.x+c.width, c.y+c.height, c.fg, c.bg)
		nextY := c.y + 1
		// The title
		if c.title != "" {
			if len(c.title) > c.width {
				diff := c.width - len(c.title)
				DrawStringAtPoint(c.title[:len(c.title)+diff-1], c.x+1, nextY, c.fg, c.bg)
			} else {
				DrawStringAtPoint(c.title, c.x+1, nextY, c.fg, c.bg)
			}
			nextY++
			FillWithChar('-', c.x+1, nextY, c.x+c.width-1, nextY, c.fg, c.bg)
			nextY++
		}
		if c.text != "" {
			DrawStringAtPoint(c.text, c.x+1, nextY, c.fg, c.bg)
			nextY++
		}
		c.input.SetY(nextY)
		c.input.Draw()
		nextY += 3
		if c.showHelp {
			helpString := " (ENTER) to Accept. (ESC) to Cancel. "
			helpX := (c.x + c.width - len(helpString)) - 1
			DrawStringAtPoint(helpString, helpX, nextY, c.fg, c.bg)
		}
		if c.bordered {
			// Now draw the border
			DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, c.fg, c.bg)
		}
	}
}
