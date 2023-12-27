package termboxUtil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

type termboxControl interface {
	GetID() string
	GetX() int
	SetX(int)
	GetY() int
	SetY(int)
	GetWidth() int
	SetWidth(int)
	GetHeight() int
	SetHeight(int)
	GetFgColor() termbox.Attribute
	SetFgColor(termbox.Attribute)
	GetBgColor() termbox.Attribute
	SetBgColor(termbox.Attribute)
	HandleEvent(termbox.Event) bool
	IsBordered() bool
	SetBordered(bool)
	SetTabSkip(bool)
	IsTabSkipped() bool
	Draw()
	SetActive(bool)
	IsActive() bool
	SetActiveFgColor(termbox.Attribute)
	SetActiveBgColor(termbox.Attribute)
}

// TextAlignment is an int value for how we're aligning text
type TextAlignment int

const (
	// AlignLeft Aligns text to the left
	AlignLeft = iota
	// AlignCenter Aligns text to the center
	AlignCenter
	// AlignRight Aligns text to the right
	AlignRight
)

/* Basic Input Helpers */

// KeyIsAlphaNumeric Returns whether the termbox event is an
// Alpha-Numeric Key Press
func KeyIsAlphaNumeric(event termbox.Event) bool {
	return KeyIsAlpha(event) || KeyIsNumeric(event)
}

// KeyIsAlpha Returns whether the termbox event is a
// alphabetic Key press
func KeyIsAlpha(event termbox.Event) bool {
	k := event.Ch
	if (k >= 'a' && k <= 'z') || (k >= 'A' && k <= 'Z') {
		return true
	}
	return false
}

// KeyIsNumeric Returns whether the termbox event is a
// numeric Key press
func KeyIsNumeric(event termbox.Event) bool {
	k := event.Ch
	if k >= '0' && k <= '9' {
		return true
	}
	return false
}

// KeyIsSymbol Returns whether the termbox event is a
// symbol Key press
func KeyIsSymbol(event termbox.Event) bool {
	symbols := []rune{'!', '@', '#', '$', '%', '^', '&', '*',
		'(', ')', '-', '_', '=', '+', '[', ']', '{', '}', '|',
		';', ':', '"', '\'', ',', '<', '.', '>', '/', '?', '`', '~'}
	k := event.Ch
	for i := range symbols {
		if k == symbols[i] {
			return true
		}
	}
	return false
}

/* Basic Output Helpers */

// DrawStringAtPoint Draw a string of text at x, y with foreground color fg, background color bg
func DrawStringAtPoint(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) (int, int) {
	xPos := x
	for _, runeValue := range str {
		termbox.SetCell(xPos, y, runeValue, fg, bg)
		xPos++
	}
	return xPos, y
}

// FillWithChar Fills from x1,y1 through x2,y2 with the rune r, foreground color fg, background bg
func FillWithChar(r rune, x1, y1, x2, y2 int, fg termbox.Attribute, bg termbox.Attribute) {
	for xx := x1; xx <= x2; xx++ {
		for yx := y1; yx <= y2; yx++ {
			termbox.SetCell(xx, yx, r, fg, bg)
		}
	}
}

// DrawBorder Draw a border around the area inside x1,y1 -> x2, y2
func DrawBorder(x1, y1, x2, y2 int, fg, bg termbox.Attribute) {
	termbox.SetCell(x1, y1, '╔', fg, bg)
	FillWithChar('═', x1+1, y1, x2-1, y1, fg, bg)
	termbox.SetCell(x2, y1, '╗', fg, bg)

	FillWithChar('║', x1, y1+1, x1, y2-1, fg, bg)
	FillWithChar('║', x2, y1+1, x2, y2-1, fg, bg)

	termbox.SetCell(x1, y2, '╚', fg, bg)
	FillWithChar('═', x1+1, y2, x2-1, y2, fg, bg)
	termbox.SetCell(x2, y2, '╝', fg, bg)
}

func DrawBorderWithPct(x1, y1, x2, y2 int, pct float64, fg, bg termbox.Attribute) {
	termbox.SetCell(x1, y1, '╔', fg, bg)

	FillWithChar('═', x1+1, y1, x2-1, y1, fg, bg)
	termbox.SetCell(x2, y1, '╗', fg, bg)

	FillWithChar('║', x1, y1+1, x1, y2-1, fg, bg)
	FillWithChar('║', x2, y1+1, x2, y2-1, fg, bg)
	// Now the percent indicator
	pctY := int(((float64(y2)-float64(y1)-2)*pct)+float64(y1)) + 1
	termbox.SetCell(x2, pctY, '▒', fg, bg)

	termbox.SetCell(x1, y2, '╚', fg, bg)
	FillWithChar('═', x1+1, y2, x2-1, y2, fg, bg)
	termbox.SetCell(x2, y2, '╝', fg, bg)
}

func DrawBorderWithTitle(x1, y1, x2, y2 int, title string, fg, bg termbox.Attribute) {
	termbox.SetCell(x1, y1, '╔', fg, bg)

	DrawStringAtPoint(title, x1+1, y1, fg, bg)
	FillWithChar('═', x1+len(title)+1, y1, x2-1, y1, fg, bg)
	termbox.SetCell(x2, y1, '╗', fg, bg)

	FillWithChar('║', x1, y1+1, x1, y2-1, fg, bg)
	FillWithChar('║', x2, y1+1, x2, y2-1, fg, bg)

	termbox.SetCell(x1, y2, '╚', fg, bg)
	FillWithChar('═', x1+1, y2, x2-1, y2, fg, bg)
	termbox.SetCell(x2, y2, '╝', fg, bg)
}

func DrawBorderWithTitleAndPct(x1, y1, x2, y2 int, title string, pct float64, fg, bg termbox.Attribute) {
	termbox.SetCell(x1, y1, '╔', fg, bg)

	DrawStringAtPoint(title, x1+1, y1, fg, bg)
	FillWithChar('═', x1+len(title)+1, y1, x2-1, y1, fg, bg)
	termbox.SetCell(x2, y1, '╗', fg, bg)

	FillWithChar('║', x1, y1+1, x1, y2-1, fg, bg)
	FillWithChar('║', x2, y1+1, x2, y2-1, fg, bg)
	// Now the percent indicator
	pctY := int(((float64(y2)-float64(y1)-2)*pct)+float64(y1)) + 1
	termbox.SetCell(x2, pctY, '▒', fg, bg)

	termbox.SetCell(x1, y2, '╚', fg, bg)
	FillWithChar('═', x1+1, y2, x2-1, y2, fg, bg)
	termbox.SetCell(x2, y2, '╝', fg, bg)
}

// AlignText Aligns the text txt within width characters using the specified alignment
func AlignText(txt string, width int, align TextAlignment) string {
	return AlignTextWithFill(txt, width, align, ' ')
}

// AlignTextWithFill Aligns the text txt within width characters using the specified alignment
// filling any spaces with the 'fill' character
func AlignTextWithFill(txt string, width int, align TextAlignment, fill rune) string {
	fillChar := string(fill)
	numSpaces := width - len(txt)
	switch align {
	case AlignCenter:
		if numSpaces/2 > 0 {
			return fmt.Sprintf("%s%s%s",
				strings.Repeat(fillChar, numSpaces/2),
				txt, strings.Repeat(fillChar, numSpaces/2),
			)
		}
		return txt
	case AlignRight:
		return fmt.Sprintf("%s%s", strings.Repeat(fillChar, numSpaces), txt)
	default:
		if numSpaces >= 0 {
			return fmt.Sprintf("%s%s", txt, strings.Repeat(fillChar, numSpaces))
		}
		return txt
	}
}

func ToLabel(c termboxControl) (*Label, error) {
	v, ok := c.(*Label)
	if ok {
		return v, nil
	}
	return nil, errors.New("Control isn't a Label")
}
func ToInputField(c termboxControl) (*InputField, error) {
	v, ok := c.(*InputField)
	if ok {
		return v, nil
	}
	return nil, errors.New("Control isn't an Input Field")
}

/* More advanced things are in their respective files */
