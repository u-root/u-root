package bubbline

import (
	"context"
	"errors"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/bubbline/complete"
	"github.com/knz/bubbline/editline"
	"github.com/knz/bubbline/history"
)

// Editor represents an input line editor.
type Editor struct {
	*editline.Model

	autoSaveHistory bool
	histFile        string
}

// New instantiates an editor.
func New() *Editor {
	return &Editor{
		Model: editline.New(0, 0),
	}
}

var _ tea.Model = (*Editor)(nil)

// Update is part of the tea.Model interface.
func (m *Editor) Update(imsg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := imsg.(type) {
	case tea.WindowSizeMsg:
		m.Model.SetSize(msg.Width, msg.Height)

	case editline.InputCompleteMsg:
		return m, tea.Quit
	}
	_, next := m.Model.Update(imsg)
	return m, next
}

// Close should be called when the editor is not used any more.
func (m *Editor) Close() {}

// ErrInterrupted is returned when the input was interrupted with
// e.g. Ctrl+C, or when SIGINT was received.
var ErrInterrupted = editline.ErrInterrupted

// ErrTerminated is returned when the input was interrupted
// by receiving SIGTERM.
var ErrTerminated = errors.New("terminated")

// Getline runs the editor and returns the line that was read.
func (m *Editor) GetLine() (string, error) {
	// We don't like the default handling of SIGINT/SIGTERM. Provide our own.
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, stopSignals...)
	defer signal.Stop(ch)
	var sig os.Signal
	go func() {
		select {
		case sig = <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()
	// Create a Bubbletea program to handle our input.
	p := tea.NewProgram(m, tea.WithoutSignalHandler(), tea.WithContext(ctx))
	m.Reset()
	if _, err := p.Run(); err != nil {
		// Was a signal received?
		if ctx.Err() != nil {
			// Yes: choose the resulting error depending on which signal was
			// received.
			if sig == os.Interrupt {
				err = ErrInterrupted
			} else {
				err = ErrTerminated
			}
		}
		return "", err
	}
	return m.Value(), m.Err
}

// AddHistory adds a history entry and optionally saves
// the history to file.
func (m *Editor) AddHistory(line string) error {
	m.AddHistoryEntry(line)
	if m.autoSaveHistory && m.histFile != "" {
		return m.SaveHistory()
	}
	return nil
}

// LoadHistory loads the entry history from file.
func (m *Editor) LoadHistory(file string) error {
	h, err := history.LoadHistory(file)
	if err != nil {
		return err
	}
	m.SetHistory(h)
	return nil
}

// SaveHistory saves the current history to the file
// previously configured with SetAutoSaveHistory.
func (m *Editor) SaveHistory() error {
	if m.histFile == "" {
		return errors.New("no savefile configured")
	}
	h := m.GetHistory()
	if h == nil {
		return errors.New("history not configured")
	}
	return history.SaveHistory(h, m.histFile)
}

// SetAutoSaveHistory enables/disables auto-saving of entered lines
// to the history.
func (m *Editor) SetAutoSaveHistory(file string, autoSave bool) {
	m.autoSaveHistory = autoSave
	m.histFile = file
}

// Values is the interface to the values displayed by the completion
// bubble.
type Values = complete.Values

// Entry is the interface to one completion candidate in the menu
// visualizer.
type Entry = complete.Entry

// AutoCompleteFn is called upon the user pressing the
// autocomplete key. The callback is provided the text of the input
// and the position of the cursor in the input.
// The returned msg is printed above the input box.
type AutoCompleteFn = editline.AutoCompleteFn

// Completions is the return value of AutoCompleteFn. It is a
// combination of Values and a Candidate function that converts a
// display Entry into a replacement Candidate.
type Completions = editline.Completions

// Candidate is the type of one completion replacement candidate.
type Candidate = editline.Candidate
