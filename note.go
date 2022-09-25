package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type Category struct {
	Color lipgloss.Color
	Name  textinput.Model
	Code  int
}

type NoteExport struct {
	Title   string
	Content string

	Pinned       bool
	CategoryCode int
}

type Note struct {
	/* represents an individual sticky note*/
	Title   textinput.Model
	Content textarea.Model

	Pinned       bool
	CategoryCode int
	Style        lipgloss.Style
}

// TODO: entirely unfocused view mode, support for home/end keys, pgup, pgdwn

func initialNote(defaultColor []lipgloss.Color, category int) *Note {
	st := lipgloss.NewStyle().Background(defaultColor[category]) // base style

	ti := textinput.New()
	ti.BackgroundStyle = st
	ti.PromptStyle = st

	ts := st.Copy().Italic(true).Faint(true) // title style

	ti.TextStyle = ts
	ti.PlaceholderStyle = ts

	ti.Prompt = "  "
	ti.Placeholder = "(untitled)" + strings.Repeat(" ", (NOTEWIDTH-12-len(ti.Prompt)))
	ti.CharLimit = 40
	ti.Width = NOTEWIDTH - 5

	taf := st.Copy().Faint(true)  // text area focussed
	tab := st.Copy().Faint(false) // text area blurred
	ci := textarea.New()
	ci.BlurredStyle = textarea.Style{
		Base:       st,
		CursorLine: taf,
		Text:       taf,
	}
	ci.FocusedStyle = textarea.Style{
		Base:       st,
		CursorLine: tab,
		Text:       tab,
	}

	ci.Placeholder = ""
	ci.ShowLineNumbers = false
	ci.Focus()
	ci.SetWidth(NOTEWIDTH - 2)
	ci.SetHeight(NOTEHEIGHT - 4)

	return &Note{
		Title:   ti,
		Content: ci,

		CategoryCode: category,

		Pinned: false,
		Style:  st,
	}
}

func (n Note) Init() tea.Cmd {
	return textarea.Blink
}

func (n Note) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return n, tea.Quit
		case tea.KeyCtrlT:
			n.Title.Prompt = "┃ "
			n.Content.Prompt = "  "

			n.Content.Blur()
			n.Title.Focus()
			n.Title.TextStyle.Faint(false)

			n.Title.PlaceholderStyle.Faint(false)

			//return n, nil
		case tea.KeyCtrlG:
			n.Title.Prompt = "  "
			n.Content.Prompt = "┃ "

			n.Content.Focus()
			n.Title.Blur()
			n.Title.TextStyle.Faint(true)
			n.Title.PlaceholderStyle.Faint(true)
			//return n, nil
		case tea.KeyCtrlV:

			if n.Content.Focused() {
				tmp := textarea.Paste()
				n.Content, cmd = n.Content.Update(tmp)
			} else {
				n.Title, cmd = n.Title.Update(textinput.Paste())
			}
			return n, cmd

		}
	case error:
		panic(msg)
	}
	if n.Content.Focused() {
		n.Content, cmd = n.Content.Update(msg)

	} else {
		n.Title, cmd = n.Title.Update(msg)
	}

	return n, cmd
}

func (n Note) View() string {
	noteStyle := n.Style.Copy().
		Height(NOTEHEIGHT).
		Width(NOTEWIDTH).
		Padding(1, 1, 1, 1)
	return noteStyle.Render(lipgloss.JoinVertical(0, n.Title.View()+"\n", n.Content.View()))
}
