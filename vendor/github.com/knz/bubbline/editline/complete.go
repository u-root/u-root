package editline

import (
	"sort"

	"github.com/knz/bubbline/complete"
	"github.com/knz/bubbline/computil"
	rw "github.com/mattn/go-runewidth"
)

// AutoCompleteFn is called upon the user pressing the
// autocomplete key. The callback is provided the text of the input
// and the position of the cursor in the input.
// The returned msg is printed above the input box.
type AutoCompleteFn func(entireInput [][]rune, line, col int) (msg string, comp Completions)

// Completions is the return value of AutoCompleteFn.
type Completions interface {
	// Values is the set of all completion values.
	complete.Values

	// Candidate converts a complete.Entry to a Candidate.
	Candidate(e complete.Entry) Candidate
}

// Candidate is the type of one completion candidate.
type Candidate interface {
	// Replacement is the string to replace.
	Replacement() string

	// MoveRight returns the number of times the cursor
	// should be moved to the right to arrive at the
	// end of the word being replaced by the completion.
	//
	// For example, if the input is this:
	//
	//       alice
	//        ^
	//
	// where the cursor is on the 2nd character, and
	// the completion is able to replace the entire word,
	// MoveRight should return 4.
	MoveRight() int

	// DeleteLeft returns the total number of characters
	// being replaced by the completion, including the
	// characters to the right of the cursor (as returned by MoveRight).
	//
	// For example, if the input is this:
	//
	//       alice
	//        ^
	//
	// where the cursor is on the 2nd character, and
	// the completion is able to replace the entire word,
	// DeleteLeft should return 5.
	DeleteLeft() int
}

// SingleWordCompletion turns a simple string into a Completions
// interface suitable to return from an AutoCompleteFn.
// The start/end positions refer to the word start and end
// positions on the current line.
func SingleWordCompletion(word string, cursor, start, end int) Completions {
	return &wordsCompletion{[]string{word}, "completion", cursor, start, end}
}

// SimpleWordsCompletion turns an array of simple strings into a
// Completions interface suitable to return from an AutoCompleteFn.
// The start/end positions refer to the word start and end positions
// on the current line.
func SimpleWordsCompletion(words []string, category string, cursor, start, end int) Completions {
	if len(words) == 0 {
		return nil
	}
	return &wordsCompletion{words, category, cursor, start, end}
}

type wordsCompletion struct {
	words              []string
	category           string
	cursor, start, end int
}

func (s *wordsCompletion) NumCategories() int                   { return 1 }
func (s *wordsCompletion) CategoryTitle(_ int) string           { return s.category }
func (s *wordsCompletion) NumEntries(_ int) int                 { return len(s.words) }
func (s *wordsCompletion) Entry(_, i int) complete.Entry        { return wordsEntry{s, i} }
func (s *wordsCompletion) Candidate(e complete.Entry) Candidate { return e.(wordsEntry) }

type wordsEntry struct {
	s *wordsCompletion
	i int
}

func (s wordsEntry) Title() string       { return s.s.words[s.i] }
func (s wordsEntry) Description() string { return "" }
func (s wordsEntry) Replacement() string { return s.Title() }
func (s wordsEntry) MoveRight() int      { return s.s.end - s.s.cursor }
func (s wordsEntry) DeleteLeft() int     { return s.s.end - s.s.start }

// computePrefill computes whether there's a common prefix we can pre-fill.
// For a prefill to exist, the following two conditions must hold:
// - all completions should have the same start and end positions.
// - there's a common prefix.
func computePrefill(
	comp Completions,
) (hasPrefill bool, moveRight, deleteLeft int, prefill string, newCompletions Completions) {
	if comp == nil {
		return false, 0, 0, "", nil
	}
	var candidates []string
	deleteLeft = -1
	moveRight = -1
	numCats := comp.NumCategories()
	for catIdx := 0; catIdx < numCats; catIdx++ {
		numE := comp.NumEntries(catIdx)
		for eIdx := 0; eIdx < numE; eIdx++ {
			e := comp.Entry(catIdx, eIdx)
			c := comp.Candidate(e)
			cdl, cmr := c.DeleteLeft(), c.MoveRight()
			if deleteLeft == -1 {
				deleteLeft = cdl
				moveRight = cmr
			} else {
				if cdl != deleteLeft || cmr != moveRight {
					// Not all candidates start at the same position: there is
					// no common prefix.
					return false, 0, 0, "", comp
				}
			}
			candidates = append(candidates, c.Replacement())
		}
	}
	if len(candidates) == 0 {
		return false, 0, 0, "", nil
	}
	if len(candidates) == 1 {
		return true, moveRight, deleteLeft, candidates[0], nil
	}
	sort.Strings(candidates)
	// TODO(knz): Do we ever need case-insensitive prefix?
	prefix := computil.FindLongestCommonPrefix(candidates[0], candidates[len(candidates)-1], false)
	if len(prefix) == 0 {
		return false, 0, 0, "", comp
	}
	return true, moveRight, deleteLeft, prefix, shiftComp{
		Completions: comp,
		shift:       rw.StringWidth(prefix),
	}
}

type shiftComp struct {
	Completions
	shift int
}

func (s shiftComp) Candidate(e complete.Entry) Candidate {
	return shiftCandidate{c: s.Completions.Candidate(e), shift: s.shift}
}

type shiftCandidate struct {
	c     Candidate
	shift int
}

func (s shiftCandidate) Replacement() string { return s.c.Replacement() }
func (s shiftCandidate) MoveRight() int      { return 0 }
func (s shiftCandidate) DeleteLeft() int     { return s.shift }
