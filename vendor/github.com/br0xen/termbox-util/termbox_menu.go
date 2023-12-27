package termboxUtil

import "github.com/nsf/termbox-go"

// Menu is a menu with a list of options
type Menu struct {
	id      string
	title   string
	options []MenuOption
	// If height is -1, then it is adaptive to the menu
	x, y, width, height    int
	showHelp               bool
	cursor                 int
	bg, fg                 termbox.Attribute
	selectedBg, selectedFg termbox.Attribute
	disabledBg, disabledFg termbox.Attribute
	selectedDisabledBg     termbox.Attribute
	selectedDisabledFg     termbox.Attribute
	activeFg, activeBg     termbox.Attribute
	isDone                 bool
	bordered               bool
	vimMode                bool
	tabSkip                bool
	active                 bool
	canSelectDisabled      bool
}

// CreateMenu Creates a menu with the specified attributes
func CreateMenu(title string, options []string, x, y, width, height int, fg, bg termbox.Attribute) *Menu {
	c := Menu{
		title: title,
		x:     x, y: y, width: width, height: height,
		fg: fg, bg: bg, selectedFg: bg, selectedBg: fg,
		disabledFg: bg, disabledBg: bg,
		activeFg: fg, activeBg: bg,
		bordered: true,
		tabSkip:  false,
	}
	for _, line := range options {
		c.options = append(c.options, MenuOption{text: line})
	}
	if len(c.options) > 0 {
		c.SetSelectedOption(&c.options[0])
	}
	return &c
}

func (c *Menu) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }

func (c *Menu) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }

func (c *Menu) SetActive(a bool) { c.active = a }

func (c *Menu) IsActive() bool { return c.active }

// GetID returns this control's ID
func (c *Menu) GetID() string { return c.id }

// SetID sets this control's ID
func (c *Menu) SetID(newID string) { c.id = newID }

func (c *Menu) IsTabSkipped() bool { return c.tabSkip }
func (c *Menu) SetTabSkip(b bool)  { c.tabSkip = b }

// GetTitle returns the current title of the menu
func (c *Menu) GetTitle() string { return c.title }

// SetTitle sets the current title of the menu to s
func (c *Menu) SetTitle(s string) {
	c.title = s
}

// GetOptions returns the current options of the menu
func (c *Menu) GetOptions() []MenuOption {
	return c.options
}

// SetOptions set the menu's options to opts
func (c *Menu) SetOptions(opts []MenuOption) {
	c.options = opts
}

// SetOptionsFromStrings sets the options of this menu from a slice of strings
func (c *Menu) SetOptionsFromStrings(opts []string) {
	var newOpts []MenuOption
	for _, v := range opts {
		newOpts = append(newOpts, *CreateOptionFromText(v))
	}
	c.SetOptions(newOpts)
	c.SetSelectedOption(c.GetOptionFromIndex(0))
}

// GetX returns the current x coordinate of the menu
func (c *Menu) GetX() int { return c.x }

// SetX sets the current x coordinate of the menu to x
func (c *Menu) SetX(x int) {
	c.x = x
}

// GetY returns the current y coordinate of the menu
func (c *Menu) GetY() int { return c.y }

// SetY sets the current y coordinate of the menu to y
func (c *Menu) SetY(y int) {
	c.y = y
}

// GetWidth returns the current width of the menu
func (c *Menu) GetWidth() int { return c.width }

// SetWidth sets the current menu width to width
func (c *Menu) SetWidth(width int) {
	c.width = width
}

// GetHeight returns the current height of the menu
func (c *Menu) GetHeight() int { return c.height }

// SetHeight set the height of the menu to height
func (c *Menu) SetHeight(height int) {
	c.height = height
}

// GetSelectedOption returns the current selected option
func (c *Menu) GetSelectedOption() *MenuOption {
	idx := c.GetSelectedIndex()
	if idx != -1 {
		return &c.options[idx]
	}
	return nil
}

// GetOptionFromIndex Returns the
func (c *Menu) GetOptionFromIndex(idx int) *MenuOption {
	if idx >= 0 && idx < len(c.options) {
		return &c.options[idx]
	}
	return nil
}

// GetOptionFromText Returns the first option with the text v
func (c *Menu) GetOptionFromText(v string) *MenuOption {
	for idx := range c.options {
		testOption := &c.options[idx]
		if testOption.GetText() == v {
			return testOption
		}
	}
	return nil
}

// GetSelectedIndex returns the index of the selected option
// Returns -1 if nothing is selected
func (c *Menu) GetSelectedIndex() int {
	for idx := range c.options {
		if c.options[idx].IsSelected() {
			return idx
		}
	}
	return -1
}

// SetSelectedIndex sets the selection to setIdx
func (c *Menu) SetSelectedIndex(idx int) {
	if len(c.options) > 0 {
		if idx < 0 {
			idx = 0
		} else if idx >= len(c.options) {
			idx = len(c.options) - 1
		}
		c.SetSelectedOption(&c.options[idx])
	}
}

// SetSelectedOption sets the current selected option to v (if it's valid)
func (c *Menu) SetSelectedOption(v *MenuOption) {
	for idx := range c.options {
		if &c.options[idx] == v {
			c.options[idx].Select()
		} else {
			c.options[idx].Unselect()
		}
	}
}

// SelectPrevOption Decrements the selected option (if it can)
func (c *Menu) SelectPrevOption() {
	idx := c.GetSelectedIndex()
	for idx >= 0 {
		idx--
		testOption := c.GetOptionFromIndex(idx)
		if testOption != nil {
			if c.canSelectDisabled || !testOption.IsDisabled() {
				c.SetSelectedOption(testOption)
				return
			}
		}
	}
}

// SelectNextOption Increments the selected option (if it can)
func (c *Menu) SelectNextOption() {
	idx := c.GetSelectedIndex()
	for idx < len(c.options) {
		idx++
		testOption := c.GetOptionFromIndex(idx)
		if testOption != nil {
			if c.canSelectDisabled || !testOption.IsDisabled() {
				c.SetSelectedOption(testOption)
				return
			}
		}
	}
}

// SelectPageUpOption Goes up 'menu height' options
func (c *Menu) SelectPageUpOption() {
	idx := c.GetSelectedIndex()
	idx -= c.height
	if idx < 0 {
		idx = 0
	}
	c.SetSelectedIndex(idx)
	return
}

// SelectPageDownOption Goes down 'menu height' options
func (c *Menu) SelectPageDownOption() {
	idx := c.GetSelectedIndex()
	idx += c.height
	if idx >= len(c.options) {
		idx = len(c.options) - 1
	}
	c.SetSelectedIndex(idx)
	return
}

// SelectFirstOption Goes to the top
func (c *Menu) SelectFirstOption() {
	c.SetSelectedIndex(0)
	return
}

// SelectLastOption Goes to the bottom
func (c *Menu) SelectLastOption() {
	c.SetSelectedIndex(len(c.options) - 1)
	return
}

// SetOptionDisabled Disables the specified option
func (c *Menu) SetOptionDisabled(idx int) {
	if len(c.options) > idx {
		c.GetOptionFromIndex(idx).Disable()
	}
}

// SetOptionEnabled Enables the specified option
func (c *Menu) SetOptionEnabled(idx int) {
	if len(c.options) > idx {
		c.GetOptionFromIndex(idx).Enable()
	}
}

// HelpIsShown returns true or false if the help is displayed
func (c *Menu) HelpIsShown() bool { return c.showHelp }

// ShowHelp sets whether or not to display the help text
func (c *Menu) ShowHelp(b bool) {
	c.showHelp = b
}

func (c *Menu) GetFgColor() termbox.Attribute   { return c.fg }
func (c *Menu) SetFgColor(fg termbox.Attribute) { c.fg = fg }
func (c *Menu) GetBgColor() termbox.Attribute   { return c.bg }
func (c *Menu) SetBgColor(bg termbox.Attribute) { c.bg = bg }

func (c *Menu) GetSelectedFgColor() termbox.Attribute   { return c.selectedFg }
func (c *Menu) SetSelectedFgColor(fg termbox.Attribute) { c.selectedFg = fg }
func (c *Menu) GetSelectedBgColor() termbox.Attribute   { return c.selectedBg }
func (c *Menu) SetSelectedBgColor(bg termbox.Attribute) { c.selectedBg = bg }

func (c *Menu) GetSelectedDisabledFgColor() termbox.Attribute   { return c.selectedDisabledFg }
func (c *Menu) SetSelectedDisabledFgColor(fg termbox.Attribute) { c.selectedDisabledFg = fg }
func (c *Menu) GetSelectedDisabledBgColor() termbox.Attribute   { return c.selectedDisabledBg }
func (c *Menu) SetSelectedDisabledBgColor(bg termbox.Attribute) { c.selectedDisabledBg = bg }

func (c *Menu) GetDisabledFgColor() termbox.Attribute   { return c.disabledFg }
func (c *Menu) SetDisabledFgColor(fg termbox.Attribute) { c.disabledFg = fg }
func (c *Menu) GetDisabledBgColor() termbox.Attribute   { return c.disabledBg }
func (c *Menu) SetDisabledBgColor(bg termbox.Attribute) { c.disabledBg = bg }

// IsDone returns whether the user has answered the modal
func (c *Menu) IsDone() bool { return c.isDone }

// SetDone sets whether the modal has completed it's purpose
func (c *Menu) SetDone(b bool) {
	c.isDone = b
}

// IsBordered returns true or false if this menu has a border
func (c *Menu) IsBordered() bool { return c.bordered }

// SetBordered sets whether we render a border around the menu
func (c *Menu) SetBordered(b bool) {
	c.bordered = b
}

// EnableVimMode Enables h,j,k,l navigation
func (c *Menu) EnableVimMode() {
	c.vimMode = true
}

// DisableVimMode Disables h,j,k,l navigation
func (c *Menu) DisableVimMode() {
	c.vimMode = false
}

func (c *Menu) SetCanSelectDisabled(b bool) {
	c.canSelectDisabled = b
}

// HandleEvent handles the termbox event and returns whether it was consumed
func (c *Menu) HandleEvent(event termbox.Event) bool {
	if event.Key == termbox.KeyEnter || event.Key == termbox.KeySpace {
		c.isDone = true
		return true
	}
	currentIdx := c.GetSelectedIndex()
	switch event.Key {
	case termbox.KeyArrowUp:
		c.SelectPrevOption()
	case termbox.KeyArrowDown:
		c.SelectNextOption()
	case termbox.KeyArrowLeft:
		c.SelectPageUpOption()
	case termbox.KeyArrowRight:
		c.SelectPageDownOption()
	}
	if c.vimMode {
		switch event.Ch {
		case 'j':
			c.SelectNextOption()
		case 'k':
			c.SelectPrevOption()
		}
		if event.Key == termbox.KeyCtrlF {
			c.SelectPageDownOption()
		} else if event.Key == termbox.KeyCtrlB {
			c.SelectPageUpOption()
		}
	}
	if c.GetSelectedIndex() != currentIdx {
		return true
	}
	return false
}

// Draw draws the modal
func (c *Menu) Draw() {
	useFg, useBg := c.fg, c.bg
	if c.active {
		useFg, useBg = c.activeFg, c.activeBg
	}
	// First blank out the area we'll be putting the menu
	FillWithChar(' ', c.x, c.y, c.x+c.width, c.y+c.height, useFg, useBg)
	// Now draw the border
	optionStartX := c.x
	optionStartY := c.y
	optionWidth := c.width
	_ = optionWidth
	optionHeight := c.height
	if optionHeight == -1 {
		optionHeight = len(c.options)
	}
	if c.bordered {
		pct := float64(c.GetSelectedIndex()) / float64(len(c.options))
		if c.title == "" {
			if len(c.options) > c.height-2 {
				DrawBorderWithPct(c.x, c.y, c.x+c.width, c.y+c.height, pct, useFg, useBg)
			} else {
				DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, useFg, useBg)
			}
		} else {
			if len(c.options) > c.height-2 {
				DrawBorderWithTitleAndPct(c.x, c.y, c.x+c.width, c.y+c.height, " "+c.title+" ", pct, useFg, useBg)
			} else {
				DrawBorderWithTitle(c.x, c.y, c.x+c.width, c.y+c.height, " "+c.title+" ", useFg, useBg)
			}
		}
		optionStartX = c.x + 1
		optionStartY = c.y + 1
		optionWidth = c.width - 1
		optionHeight -= 2
	}

	if len(c.options) > 0 {
		firstDispIdx := 0
		lastDispIdx := len(c.options) - 1
		if len(c.options) > c.height-2 {
			lastDispIdx = c.height - 2
		}
		if c.GetSelectedIndex() > c.height-2 {
			firstDispIdx = c.GetSelectedIndex() - (c.height - 2)
			lastDispIdx = c.GetSelectedIndex()
		}
		for idx := firstDispIdx; idx < lastDispIdx+1; idx++ {
			currOpt := &c.options[idx]
			outTxt := currOpt.GetText()
			if currOpt.IsDisabled() {
				if c.GetSelectedOption() == currOpt {
					DrawStringAtPoint(outTxt, optionStartX, optionStartY, c.selectedDisabledFg, c.selectedDisabledBg)
				} else {
					DrawStringAtPoint(outTxt, optionStartX, optionStartY, c.disabledFg, c.disabledBg)
				}
			} else if c.GetSelectedOption() == currOpt {
				DrawStringAtPoint(outTxt, optionStartX, optionStartY, c.selectedFg, c.selectedBg)
			} else {
				DrawStringAtPoint(outTxt, optionStartX, optionStartY, useFg, useBg)
			}
			optionStartY++
		}
	}
	/*
			// Print the options
			bldHeight := (optionHeight / 2)
			startIdx := c.GetSelectedIndex()
			endIdx := c.GetSelectedIndex()
			for bldHeight > 0 && startIdx >= 1 {
				startIdx--
				bldHeight--
			}
			bldHeight += (optionHeight / 2)
			for bldHeight > 0 && endIdx < len(c.options) {
				endIdx++
				bldHeight--
			}

			for idx := startIdx; idx < endIdx; idx++ {
				if c.GetSelectedIndex()-idx >= optionHeight-1 {
					// Skip this one
					continue
				}
				currOpt := &c.options[idx]
				outTxt := currOpt.GetText()
				if len(outTxt) >= c.width {
					outTxt = outTxt[:c.width]
				}
				if currOpt.IsDisabled() {
					DrawStringAtPoint(outTxt, optionStartX, optionStartY, c.disabledFg, c.disabledBg)
				} else if c.GetSelectedOption() == currOpt {
					DrawStringAtPoint(outTxt, optionStartX, optionStartY, c.selectedFg, c.selectedBg)
				} else {
					DrawStringAtPoint(outTxt, optionStartX, optionStartY, useFg, useBg)
				}
				optionStartY++
				if optionStartY > c.y+optionHeight-1 {
					break
				}
			}
		}
	*/
}

/* MenuOption Struct & methods */

// MenuOption An option in the menu
type MenuOption struct {
	id       string
	text     string
	selected bool
	disabled bool
	helpText string
	subMenu  []MenuOption
}

// CreateOptionFromText just returns a MenuOption object
// That only has it's text value set.
func CreateOptionFromText(s string) *MenuOption {
	return &MenuOption{text: s}
}

// SetText Sets the text for this option
func (c *MenuOption) SetText(s string) {
	c.text = s
}

// GetText Returns the text for this option
func (c *MenuOption) GetText() string { return c.text }

// Disable Sets this option to disabled
func (c *MenuOption) Disable() {
	c.disabled = true
}

// Enable Sets this option to enabled
func (c *MenuOption) Enable() {
	c.disabled = false
}

// IsDisabled returns whether this option is enabled
func (c *MenuOption) IsDisabled() bool {
	return c.disabled
}

// IsSelected Returns whether this option is selected
func (c *MenuOption) IsSelected() bool {
	return c.selected
}

// Select Sets this option to selected
func (c *MenuOption) Select() {
	c.selected = true
}

// Unselect Sets this option to not selected
func (c *MenuOption) Unselect() {
	c.selected = false
}

// SetHelpText Sets this option's help text to s
func (c *MenuOption) SetHelpText(s string) {
	c.helpText = s
}

// GetHelpText Returns the help text for this option
func (c *MenuOption) GetHelpText() string { return c.helpText }

// AddToSubMenu adds a slice of MenuOptions to this option
func (c *MenuOption) AddToSubMenu(sub *MenuOption) {
	c.subMenu = append(c.subMenu, *sub)
}
