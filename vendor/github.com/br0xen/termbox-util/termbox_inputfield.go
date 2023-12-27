package termboxUtil

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

// InputField is a field for inputting text
type InputField struct {
	id                  string
	title               string
	value               string
	x, y, width, height int
	cursor              int
	fg, bg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	cursorFg, cursorBg  termbox.Attribute
	bordered            bool
	wrap                bool
	multiline           bool
	tabSkip             bool
	active              bool
	justified           bool

	filter func(*InputField, string, string) string
}

// CreateInputField creates an input field at x, y that is w by h
func CreateInputField(x, y, w, h int, fg, bg termbox.Attribute) *InputField {
	c := InputField{x: x, y: y, width: w, height: h,
		fg: fg, bg: bg, cursorFg: bg, cursorBg: fg, activeFg: fg, activeBg: bg,
	}
	c.filter = func(fld *InputField, o, n string) string { return n }
	return &c
}

func (c *InputField) SetTitle(title string)                 { c.title = title }
func (c *InputField) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *InputField) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *InputField) SetActive(a bool)                      { c.active = a }
func (c *InputField) IsActive() bool                        { return c.active }

// GetID returns this control's ID
func (c *InputField) GetID() string { return c.id }

// SetID sets this control's ID
func (c *InputField) SetID(newID string) {
	c.id = newID
}

// GetValue gets the current text that is in the InputField
func (c *InputField) GetValue() string { return c.value }

// SetValue sets the current text in the InputField to s
func (c *InputField) SetValue(s string) {
	c.value = s
}

// GetX returns the x position of the input field
func (c *InputField) GetX() int { return c.x }

// SetX sets the x position of the input field
func (c *InputField) SetX(x int) { c.x = x }

// GetY returns the y position of the input field
func (c *InputField) GetY() int { return c.y }

// SetY sets the y position of the input field
func (c *InputField) SetY(y int) { c.y = y }

// GetWidth returns the current width of the input field
func (c *InputField) GetWidth() int { return c.width }

// SetWidth sets the current width of the input field
func (c *InputField) SetWidth(w int) { c.width = w }

// GetHeight returns the current height of the input field
func (c *InputField) GetHeight() int { return c.height }

// SetHeight sets the current height of the input field
func (c *InputField) SetHeight(h int) { c.height = h }

// GetFgColor returns the foreground color
func (c *InputField) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *InputField) SetFgColor(fg termbox.Attribute) { c.fg = fg }

// GetBgColor returns the background color
func (c *InputField) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *InputField) SetBgColor(bg termbox.Attribute) { c.bg = bg }

func (c *InputField) SetCursorFg(fg termbox.Attribute) { c.cursorFg = fg }

func (c *InputField) GetCursorFg() termbox.Attribute { return c.cursorFg }

func (c *InputField) SetCursorBg(bg termbox.Attribute) { c.cursorBg = bg }

func (c *InputField) GetCursorBg() termbox.Attribute { return c.cursorBg }

// IsBordered returns true or false if this input field has a border
func (c *InputField) IsBordered() bool { return c.bordered }

// SetBordered sets whether we render a border around the input field
func (c *InputField) SetBordered(b bool) {
	c.bordered = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (c *InputField) IsTabSkipped() bool {
	return c.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (c *InputField) SetTabSkip(b bool) {
	c.tabSkip = b
}

// DoesWrap returns true or false if this input field wraps text
func (c *InputField) DoesWrap() bool { return c.wrap }

// SetWrap sets whether we wrap the text at width.
func (c *InputField) SetWrap(b bool) {
	c.wrap = b
}

// IsMultiline returns true or false if this field can have multiple lines
func (c *InputField) IsMultiline() bool { return c.multiline }

// SetMultiline sets whether the field can have multiple lines
func (c *InputField) SetMultiline(b bool) {
	c.multiline = b
}

func (c *InputField) SetJustified(b bool) {
	c.justified = b
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (c *InputField) HandleEvent(event termbox.Event) bool {
	prev := c.value
	if event.Key == termbox.KeyTab { // There is no tabbing in here
		return false
	}
	if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
		if c.cursor+len(c.value) > 0 {
			crs := len(c.value)
			if c.cursor < 0 {
				crs = c.cursor + len(c.value)
			}
			c.value = c.value[:crs-1] + c.value[crs:]
			//c.value = c.value[:len(c.value)-1]
		}
	} else if event.Key == termbox.KeyArrowLeft {
		if c.cursor+len(c.value) > 0 {
			c.cursor--
		}
	} else if event.Key == termbox.KeyArrowRight {
		if c.cursor < 0 {
			c.cursor++
		}
	} else if event.Key == termbox.KeyCtrlU {
		// Ctrl+U Clears the Input (before the cursor)
		c.value = c.value[c.cursor+len(c.value):]
	} else {
		// Get the rune to add to our value. Space and Tab are special cases where
		// we can't use the event's rune directly
		var ch string
		switch event.Key {
		case termbox.KeySpace:
			ch = " "
		case termbox.KeyTab:
			ch = "\t"
		case termbox.KeyEnter:
			if c.multiline {
				ch = "\n"
			}
		default:
			if KeyIsAlphaNumeric(event) || KeyIsSymbol(event) {
				ch = string(event.Ch)
			}
		}

		// TODO: Handle newlines
		if c.cursor+len(c.value) == 0 {
			c.value = string(ch) + c.value
		} else if c.cursor == 0 {
			c.value = c.value + string(ch)
		} else {
			strPt1 := c.value[:(len(c.value) + c.cursor)]
			strPt2 := c.value[(len(c.value) + c.cursor):]
			c.value = strPt1 + string(ch) + strPt2
		}
	}
	c.value = c.filter(c, prev, c.value)
	return true
}

// Draw outputs the input field on the screen
func (c *InputField) Draw() {
	maxWidth := c.width
	maxHeight := c.height
	x, y := c.x, c.y
	startX := c.x
	startY := c.y
	useFg, useBg := c.fg, c.bg
	if c.active {
		useFg, useBg = c.activeFg, c.activeBg
	}
	if c.bordered {
		DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, useFg, useBg)
		maxWidth--
		maxHeight--
		x++
		y++
		startX++
		startY++
	}

	var strPt1, strPt2 string
	var cursorRune rune
	if len(c.value) > 0 {
		if c.cursor+len(c.value) == 0 {
			strPt1 = ""
			strPt2 = c.value[1:]
			cursorRune = rune(c.value[0])
		} else if c.cursor == 0 {
			strPt1 = c.value
			strPt2 = ""
			cursorRune = ' '
		} else {
			strPt1 = c.value[:(len(c.value) + c.cursor)]
			strPt2 = c.value[(len(c.value)+c.cursor)+1:]
			cursorRune = rune(c.value[len(c.value)+c.cursor])
		}
	} else {
		strPt1, strPt2, cursorRune = "", "", ' '
	}
	if c.title != "" {
		if c.active {
			DrawStringAtPoint(c.title, x, y, c.activeFg, c.activeBg)
		} else {
			DrawStringAtPoint(c.title, x, y, useFg, useBg)
		}
	}
	if c.wrap {
		// Split the text into maxWidth chunks
		for len(strPt1) > maxWidth {
			breakAt := maxWidth
			DrawStringAtPoint(strPt1[:breakAt], x, y, useFg, useBg)
			x = startX
			y++
			strPt1 = strPt1[breakAt:]
		}
		x, y = DrawStringAtPoint(strPt1, x, y, useFg, useBg)
		if x >= maxWidth {
			y++
			x = startX
		}
		termbox.SetCell(x, y, cursorRune, c.cursorFg, c.cursorBg)
		x++
		if len(strPt2) > 0 {
			lenLeft := maxWidth - len(strPt1) - 1
			if lenLeft > 0 && len(strPt2) > lenLeft {
				DrawStringAtPoint(strPt2[:lenLeft], x+1, y, useFg, useBg)
				strPt2 = strPt2[lenLeft:]
			}
			for len(strPt2) > maxWidth {
				breakAt := maxWidth
				DrawStringAtPoint(strPt2[:breakAt], x, y, useFg, useBg)
				x = startX
				y++
				strPt2 = strPt2[breakAt:]
			}
			x, y = DrawStringAtPoint(strPt2, x, y, useFg, useBg)
		}
	} else {
		for len(strPt1)+len(strPt2)+1 > maxWidth {
			if len(strPt1) >= len(strPt2) {
				if len(strPt1) == 0 {
					break
				}
				strPt1 = strPt1[1:]
			} else {
				strPt2 = strPt2[:len(strPt2)-1]
			}
		}
		stX := c.x + len(c.title)
		if c.justified {
			stX = c.x + c.width - len(strPt1) - len(strPt2) - 1
		}
		x, y = DrawStringAtPoint(strPt1, stX, c.y, useFg, useBg)
		if c.active {
			termbox.SetCell(x, y, cursorRune, c.cursorFg, c.cursorBg)
		} else {
			termbox.SetCell(x, y, cursorRune, useFg, useBg)
		}
		DrawStringAtPoint(strPt2, x+1, y, useFg, useBg)
	}
}

func (c *InputField) SetTextFilter(filter func(*InputField, string, string) string) {
	c.filter = filter
}

// Some handy text filters
func (c *InputField) InputFieldNumberFilter(fld *InputField, o, n string) string {
	_, err := strconv.Atoi(n)
	if err != nil {
		return o
	}
	return n
}
