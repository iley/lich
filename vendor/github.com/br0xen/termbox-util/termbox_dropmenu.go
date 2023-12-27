package termboxUtil

import "github.com/nsf/termbox-go"

// DropMenu is a title that, when active drops a menu down
type DropMenu struct {
	id                  string
	title               string
	x, y, width, height int
	bg, fg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	cursorBg, cursorFg  termbox.Attribute
	menu                *Menu
	menuSelected        bool
	showMenu            bool
	bordered            bool
	tabSkip             bool
	active              bool
}

// CreateDropMenu Creates a menu with the specified attributes
func CreateDropMenu(title string, options []string, x, y, width, height int, fg, bg, cursorFg, cursorBg termbox.Attribute) *DropMenu {
	c := DropMenu{
		title: title,
		x:     x, y: y, width: width, height: height,
		fg: fg, bg: bg, activeFg: fg, activeBg: bg,
		cursorFg: fg, cursorBg: bg,
	}
	c.menu = CreateMenu("", options, x, y+2, width, height, fg, bg)
	return &c
}

// GetID returns this control's ID
func (c *DropMenu) GetID() string { return c.id }

// SetID sets this control's ID
func (c *DropMenu) SetID(newID string) {
	c.id = newID
}

func (c *DropMenu) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *DropMenu) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *DropMenu) SetActive(a bool)                      { c.active = a }
func (c *DropMenu) IsActive() bool                        { return c.active }

// GetTitle returns the current title of the menu
func (c *DropMenu) GetTitle() string { return c.title }

// SetTitle sets the current title of the menu to s
func (c *DropMenu) SetTitle(s string) {
	c.title = s
}

// GetMenu returns the menu for this dropmenu
func (c *DropMenu) GetMenu() *Menu {
	return c.menu
}

// GetX returns the current x coordinate of the menu
func (c *DropMenu) GetX() int { return c.x }

// SetX sets the current x coordinate of the menu to x
func (c *DropMenu) SetX(x int) {
	c.x = x
}

// GetY returns the current y coordinate of the menu
func (c *DropMenu) GetY() int { return c.y }

// SetY sets the current y coordinate of the menu to y
func (c *DropMenu) SetY(y int) {
	c.y = y
}

// GetWidth returns the current width of the menu
func (c *DropMenu) GetWidth() int { return c.width }

// SetWidth sets the current menu width to width
func (c *DropMenu) SetWidth(width int) {
	c.width = width
}

// GetHeight returns the current height of the menu
func (c *DropMenu) GetHeight() int { return c.height }

// SetHeight set the height of the menu to height
func (c *DropMenu) SetHeight(height int) {
	c.height = height
}

// GetFgColor returns the foreground color
func (c *DropMenu) GetFgColor() termbox.Attribute { return c.fg }

// SetFgColor sets the foreground color
func (c *DropMenu) SetFgColor(fg termbox.Attribute) {
	c.fg = fg
}

// GetBgColor returns the background color
func (c *DropMenu) GetBgColor() termbox.Attribute { return c.bg }

// SetBgColor sets the current background color
func (c *DropMenu) SetBgColor(bg termbox.Attribute) {
	c.bg = bg
}

// IsBordered returns the bordered flag
func (c *DropMenu) IsBordered() bool { return c.bordered }

// SetBordered sets the bordered flag
func (c *DropMenu) SetBordered(b bool) {
	c.bordered = b
	c.menu.SetBordered(b)
}

// IsDone returns whether the user has answered the modal
func (c *DropMenu) IsDone() bool { return c.menu.isDone }

// SetDone sets whether the modal has completed it's purpose
func (c *DropMenu) SetDone(b bool) {
	c.menu.isDone = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (c *DropMenu) IsTabSkipped() bool {
	return c.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (c *DropMenu) SetTabSkip(b bool) {
	c.tabSkip = b
}

// ShowMenu tells the menu to draw the options
func (c *DropMenu) ShowMenu() {
	c.showMenu = true
	c.menuSelected = true
}

// HideMenu tells the menu to hide the options
func (c *DropMenu) HideMenu() {
	c.showMenu = false
	c.menuSelected = false
}

// HandleEvent handles the termbox event and returns whether it was consumed
func (c *DropMenu) HandleEvent(event termbox.Event) bool {
	moveUp := (event.Key == termbox.KeyArrowUp || (c.menu.vimMode && event.Ch == 'k'))
	moveDown := (event.Key == termbox.KeyArrowDown || (c.menu.vimMode && event.Ch == 'j'))
	if c.menuSelected {
		selIdx := c.menu.GetSelectedIndex()
		if (moveUp && selIdx == 0) || (moveDown && selIdx == (len(c.menu.options)-1)) {
			c.menuSelected = false
		} else {
			if c.menu.HandleEvent(event) {
				if c.menu.IsDone() {
					c.HideMenu()
				}
				return true
			}
		}
	} else {
		c.ShowMenu()
		return true
	}
	return false
}

// Draw draws the menu
func (c *DropMenu) Draw() {
	// The title
	ttlFg, ttlBg := c.fg, c.bg
	if !c.menuSelected {
		ttlFg, ttlBg = c.cursorFg, c.cursorBg
	}
	ttlTxt := c.title
	if c.showMenu {
		ttlTxt = ttlTxt + "-Showing Menu"
	}
	DrawStringAtPoint(AlignText(c.title, c.width, AlignLeft), c.x, c.y, ttlFg, ttlBg)
	if c.showMenu {
		c.menu.Draw()
	}
}
