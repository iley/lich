package termboxUtil

import termbox "github.com/nsf/termbox-go"

type Button struct {
	id                  string
	x, y, width, height int
	label               string
	fg, bg              termbox.Attribute
	activeFg, activeBg  termbox.Attribute
	bordered            bool
	tabSkip             bool
	active              bool
}

func CreateButton(x, y, w, h int, fg, bg termbox.Attribute) *Button {
	c := Button{
		x: x, y: y, width: w, height: h,
		fg: fg, bg: bg, activeFg: bg, activeBg: fg,
		bordered: true,
		tabSkip:  true,
	}
	return &c
}

func (c *Button) SetActiveFgColor(fg termbox.Attribute) { c.activeFg = fg }
func (c *Button) SetActiveBgColor(bg termbox.Attribute) { c.activeBg = bg }
func (c *Button) SetActive(a bool)                      { c.active = a }
func (c *Button) IsActive() bool                        { return c.active }
func (c *Button) GetID() string                         { return c.id }
func (c *Button) SetID(newID string)                    { c.id = newID }
func (c *Button) GetX() int                             { return c.x }
func (c *Button) SetX(x int)                            { c.x = x }
func (c *Button) GetY() int                             { return c.y }
func (c *Button) SetY(y int)                            { c.y = y }
func (c *Button) GetWidth() int                         { return c.width }
func (c *Button) SetWidth(w int)                        { c.width = w }
func (c *Button) GetHeight() int                        { return c.height }
func (c *Button) SetHeight(h int)                       { c.height = h }
func (c *Button) GetFgColor() termbox.Attribute         { return c.fg }
func (c *Button) SetFgColor(fg termbox.Attribute)       { c.fg = fg }
func (c *Button) GetBgColor() termbox.Attribute         { return c.bg }
func (c *Button) SetBgColor(bg termbox.Attribute)       { c.bg = bg }
func (c *Button) IsBordered() bool                      { return c.bordered }
func (c *Button) SetBordered(bordered bool)             { c.bordered = bordered }
func (c *Button) SetTabSkip(skip bool)                  { c.tabSkip = skip }
func (c *Button) IsTabSkipped() bool                    { return c.tabSkip }
func (c *Button) HandleEvent(e termbox.Event) bool {
	return false
}
func (c *Button) Draw() {
	stX, stY := c.x, c.y
	if c.bordered {
		DrawBorder(c.x, c.y, c.x+c.width, c.y+c.height, c.fg, c.bg)
		stX++
		stY++
	}
	DrawStringAtPoint(c.label, stX, stY, c.fg, c.bg)
}
