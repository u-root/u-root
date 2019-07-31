// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package less

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/nsf/termbox-go"
	"github.com/u-root/u-root/pkg/lineio"
	"github.com/u-root/u-root/pkg/sortedmap"
)

type size struct {
	x int
	y int
}

type Event int

const (
	// EventQuit requests an application exit.
	EventQuit Event = iota

	// EventRefresh requests a display refresh.
	EventRefresh
)

type Mode int

const (
	// ModeNormal is the standard mode, allowing file navigation.
	ModeNormal Mode = iota

	// ModeSearchEntry is search entry mode. Key presses are added
	// to the search string.
	ModeSearchEntry
)

type Less struct {
	// src is the source file being displayed.
	src *lineio.LineReader

	// tabStop is the number of spaces per tab.
	tabStop int

	// events is used to notify the main goroutine of events.
	events chan Event

	// mu locks the fields below.
	mu sync.Mutex

	// size is the size of the file display.
	// There is a statusbar beneath the display.
	size size

	// line is the line number of the first line of the display.
	line int64

	// mode is the viewer mode.
	mode Mode

	// regexp is the search regexp specified by the user.
	// Must only be modified by the event goroutine.
	regexp string

	// searchResults are the results for the current search.
	// They should be highlighted.
	searchResults *searchResults
}

// lastLine returns the last line on the display.  It may be beyond the end
// of the file, if the file is short enough.
// mu must be held on call.
func (l *Less) lastLine() int64 {
	return l.line + int64(l.size.y) - 1
}

// Scroll describes a scroll action.
type Scroll int

const (
	// ScrollTop goes to the first line.
	ScrollTop Scroll = iota
	// ScrollBottom goes to the last line.
	ScrollBottom
	// ScrollUp goes up one line.
	ScrollUp
	// ScrollDown goes down one line.
	ScrollDown
	// ScrollUpPage goes up one page full.
	ScrollUpPage
	// ScrollDownPage goes down one page full.
	ScrollDownPage
	// ScrollUpHalfPage goes up one half page full.
	ScrollUpHalfPage
	// ScrollDownHalfPage goes down one half page full.
	ScrollDownHalfPage
)

// scrollLine tries to scroll the display to the given line,
// but will not scroll beyond the first or last lines in the file.
// l.mu must be held when calling scrollLine.
func (l *Less) scrollLine(dest int64) {
	var delta int64
	if dest > l.line {
		delta = 1
	} else {
		delta = -1
	}

	for l.line != dest && l.line+delta > 0 && l.src.LineExists(l.lastLine()+delta) {
		l.line += delta
	}
}

// scroll moves the display based on the passed scroll action, without
// going past the beginning or end of the file.
func (l *Less) scroll(s Scroll) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var dest int64
	switch s {
	case ScrollTop:
		dest = 1
	case ScrollBottom:
		// Just try to go to int64 max.
		dest = 0x7fffffffffffffff
	case ScrollUp:
		dest = l.line - 1
	case ScrollDown:
		dest = l.line + 1
	case ScrollUpPage:
		dest = l.line - int64(l.size.y)
	case ScrollDownPage:
		dest = l.line + int64(l.size.y)
	case ScrollUpHalfPage:
		dest = l.line - int64(l.size.y)/2
	case ScrollDownHalfPage:
		dest = l.line + int64(l.size.y)/2
	}

	l.scrollLine(dest)
}

func (l *Less) handleEvent(e termbox.Event) {
	l.mu.Lock()
	mode := l.mode
	l.mu.Unlock()

	if e.Type != termbox.EventKey {
		return
	}

	c := e.Ch
	k := e.Key
	// Key is only valid is Ch is 0
	if c != 0 {
		k = 0
	}

	switch mode {
	case ModeNormal:
		switch {
		case c == 'q':
			l.events <- EventQuit
		case c == 'j':
			l.scroll(ScrollDown)
			l.events <- EventRefresh
		case c == 'k':
			l.scroll(ScrollUp)
			l.events <- EventRefresh
		case c == 'g':
			l.scroll(ScrollTop)
			l.events <- EventRefresh
		case c == 'G':
			l.scroll(ScrollBottom)
			l.events <- EventRefresh
		case k == termbox.KeyPgup:
			l.scroll(ScrollUpPage)
			l.events <- EventRefresh
		case k == termbox.KeyPgdn:
			l.scroll(ScrollDownPage)
			l.events <- EventRefresh
		case k == termbox.KeyCtrlU:
			l.scroll(ScrollUpHalfPage)
			l.events <- EventRefresh
		case k == termbox.KeyCtrlD:
			l.scroll(ScrollDownHalfPage)
			l.events <- EventRefresh
		case c == '/':
			l.mu.Lock()
			l.mode = ModeSearchEntry
			l.mu.Unlock()
			l.events <- EventRefresh
		case c == 'n':
			l.mu.Lock()
			if r, ok := l.searchResults.Next(l.line); ok {
				l.scrollLine(r.line)
				l.events <- EventRefresh
			}
			l.mu.Unlock()
		case c == 'N':
			l.mu.Lock()
			if r, ok := l.searchResults.Prev(l.line); ok {
				l.scrollLine(r.line)
				l.events <- EventRefresh
			}
			l.mu.Unlock()
		}
	case ModeSearchEntry:
		switch {
		case k == termbox.KeyEnter:
			r := l.search(l.regexp)
			l.mu.Lock()
			l.mode = ModeNormal
			l.regexp = ""
			l.searchResults = r
			// Jump to nearest result
			if r, ok := l.searchResults.Next(l.line); ok {
				l.scrollLine(r.line)
			}
			l.mu.Unlock()
			l.events <- EventRefresh
		default:
			l.mu.Lock()
			l.regexp += string(c)
			l.mu.Unlock()
			l.events <- EventRefresh
		}
	}
}

func (l *Less) listenEvents() {
	for {
		e := termbox.PollEvent()
		l.handleEvent(e)
	}
}

// searchResult describes search matches on a single line.
type searchResult struct {
	line    int64
	matches [][]int
	err     error
}

// matchesChar returns true if the search result contains a match for
// character index c.
func (s searchResult) matchesChar(c int) bool {
	for _, match := range s.matches {
		if len(match) < 2 {
			continue
		}

		if c >= match[0] && c < match[1] {
			return true
		}
	}
	return false
}

type searchResults struct {
	// mu locks the fields below.
	mu sync.Mutex

	// lines maps search results for a specific line to an index in results.
	lines sortedmap.Map

	// results contains the actual search results, in no particular order.
	results []searchResult
}

func NewSearchResults() *searchResults {
	return &searchResults{
		lines: sortedmap.NewMap(),
	}
}

// Add adds a result.
func (s *searchResults) Add(r searchResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := int64(len(s.results))
	s.results = append(s.results, r)
	s.lines.Insert(r.line, i)
}

// Get finds the result for a specific line, returning ok if found
func (s *searchResults) Get(line int64) (searchResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if i, ok := s.lines.Get(line); ok {
		return s.results[i], true
	}

	return searchResult{}, false
}

// Next returns the search result for the nearest line after line,
// noninclusive, if one exists.
func (s *searchResults) Next(line int64) (searchResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, i, err := s.lines.NearestGreater(line)
	if err != nil {
		// Probably ErrNoSuchKey, aka none found.
		return searchResult{}, false
	}

	return s.results[i], true
}

// Prev returns the search result for the nearest line before line,
// noninclusive, if one exists.
func (s *searchResults) Prev(line int64) (searchResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Search for line - 1, since it may be equal.
	_, i, err := s.lines.NearestLessEqual(line - 1)
	if err != nil {
		// Probably ErrNoSuchKey, aka none found.
		return searchResult{}, false
	}

	return s.results[i], true
}

func (l *Less) search(s string) *searchResults {
	reg, err := regexp.Compile(s)
	if err != nil {
		// TODO(prattmic): display a better error
		log.Printf("regexp failed to compile: %v", err)
		return NewSearchResults()
	}

	resultChan := make(chan searchResult, 100)

	searchLine := func(line int64) {
		r, err := l.src.SearchLine(reg, line)
		if err != nil {
			r = nil
		}

		resultChan <- searchResult{
			line:    line,
			matches: r,
			err:     err,
		}
	}

	nextLine := int64(1)
	// Spawn initial search goroutines
	for ; nextLine <= 5; nextLine++ {
		go searchLine(nextLine)
	}

	results := NewSearchResults()

	var count int64

	waitResult := func() searchResult {
		ret := <-resultChan
		count++

		// Only store results with matches.
		if len(ret.matches) > 0 {
			results.Add(ret)
		}

		return ret
	}

	// Collect results, start searching next lines until we start
	// hitting EOF.
	for {
		r := waitResult()

		// We started hitting errors on a previous line,
		// there is no reason to search later lines.
		if r.err != nil {
			break
		}

		go searchLine(nextLine)
		nextLine++
	}

	// Collect the remaining results.
	for count < nextLine-1 {
		waitResult()
	}

	return results
}

// statusBar renders the status bar.
// mu must be held on call.
func (l *Less) statusBar() {
	// The statusbar is just below the display.

	// Clear the statusbar
	for i := 0; i < l.size.x; i++ {
		termbox.SetCell(i, l.size.y, ' ', 0, 0)
	}

	switch l.mode {
	case ModeNormal:
		// Just a colon and a cursor
		termbox.SetCell(0, l.size.y, ':', 0, 0)
		termbox.SetCursor(1, l.size.y)
	case ModeSearchEntry:
		// / and search string
		termbox.SetCell(0, l.size.y, '/', 0, 0)
		for i, c := range l.regexp {
			termbox.SetCell(1+i, l.size.y, c, 0, 0)
		}
		termbox.SetCursor(1+len(l.regexp), l.size.y)
	}
}

// alignUp aligns n up to the next multiple of divisor.
func alignUp(n, divisor int) int {
	return n + (divisor - (n % divisor))
}

func (l *Less) refreshScreen() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for y := 0; y < l.size.y; y++ {
		buf := make([]byte, l.size.x)
		line := l.line + int64(y)

		_, err := l.src.ReadLine(buf, line)
		// EOF just means the line was shorter than the display.
		if err != nil && err != io.EOF {
			return err
		}

		highlight, ok := l.searchResults.Get(line)

		var displayColumn int
		for i, c := range buf {
			if displayColumn >= l.size.x {
				break
			}

			fg := termbox.ColorDefault
			bg := termbox.ColorDefault

			// Highlight matches
			if ok && highlight.matchesChar(i) {
				fg = termbox.ColorBlack
				bg = termbox.ColorWhite
			}

			if c == '\t' {
				// Tabs align the display up to the next
				// multiple of tabstop.
				next := alignUp(displayColumn, l.tabStop)

				// Clear the tab spaces
				for j := displayColumn; j < next; j++ {
					termbox.SetCell(j, y, ' ', 0, 0)
				}

				displayColumn = next
			} else {
				termbox.SetCell(displayColumn, y, rune(c), fg, bg)
				displayColumn++
			}
		}
	}

	l.statusBar()

	termbox.Flush()

	return nil
}

func (l *Less) Run() {
	// Start populating the LineReader cache, to speed things up later.
	go l.src.Populate()

	go l.listenEvents()

	err := l.refreshScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to refresh screen: %v\n", err)
		return
	}

	for {
		e := <-l.events

		switch e {
		case EventQuit:
			return
		case EventRefresh:
			err = l.refreshScreen()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to refresh screen: %v\n", err)
				return
			}
		}
	}
}

func NewLess(r io.ReaderAt, ts int) Less {
	x, y := termbox.Size()

	return Less{
		src:     lineio.NewLineReader(r),
		tabStop: ts,
		// Save one line for statusbar.
		size:          size{x: x, y: y - 1},
		line:          1,
		events:        make(chan Event, 1),
		mode:          ModeNormal,
		searchResults: NewSearchResults(),
	}
}
