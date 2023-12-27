// +build windows
package boltbrowser

// Windows doesn't support process backgrounding like *nix.
// So we have a separate loop for it.

import (
	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
)

func Browse(db *bolt.DB, readOnly bool) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)

	memBolt := newModel(db, readOnly)
	style := defaultStyle()

	screens := defaultScreensForData(memBolt)
	displayScreen := screens[browserScreenIndex]
	layoutAndDrawScreen(displayScreen, style)
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			newScreenIndex := displayScreen.handleKeyEvent(event)
			if newScreenIndex < len(screens) {
				displayScreen = screens[newScreenIndex]
				layoutAndDrawScreen(displayScreen, style)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize {
			layoutAndDrawScreen(displayScreen, style)
		}
	}
}
