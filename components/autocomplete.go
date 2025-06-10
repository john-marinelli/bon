package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct{}

type AutoComplete struct {
	textInput textinput.Model
}

func NewAutoComplete(suggestions []string, width int) AutoComplete {
	ti := textinput.New()
	ti.Placeholder = "throw it in the bag..."
	ti.Prompt = "bon/"
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#dbc89a"))
	ti.ShowSuggestions = true
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#dbc89a"))
	ti.Focus()
	ti.CharLimit = 50
	ti.ShowSuggestions = true
	ti.Width = width
	ti.SetSuggestions(suggestions)

	return AutoComplete{
		textInput: ti,
	}
}

func (a AutoComplete) Init() tea.Cmd {
	return textinput.Blink
}

func (a AutoComplete) Update(msg tea.Msg) (NoteComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.textInput.Width = msg.Width
	}
	var cmd tea.Cmd
	a.textInput, cmd = a.textInput.Update(msg)
	return a, cmd
}

func (a AutoComplete) View() string {
	return a.textInput.View()
}

func (a AutoComplete) Blur() NoteComponent {
	a.textInput.Blur()
	return a
}

func (a AutoComplete) Focus() NoteComponent {
	a.textInput.Focus()
	return a
}

func (a AutoComplete) Text() string {
	return a.textInput.Value()
}

func (a AutoComplete) Clear() NoteComponent {
	a.textInput.SetValue("")

	return a
}
