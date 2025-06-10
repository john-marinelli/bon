package components

import tea "github.com/charmbracelet/bubbletea"

type NoteComponent interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (NoteComponent, tea.Cmd)
	View() string
	Blur() NoteComponent
	Focus() NoteComponent
	Text() string
	Clear() NoteComponent
}
