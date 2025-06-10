package components

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type NoteInput struct {
	textArea textarea.Model
	yMargin  int
	xMargin  int
	style    lipgloss.Style
	mdRender *glamour.TermRenderer
}

func NewNoteInput(width int, height int, yMargin int) NoteInput {
	ta := textarea.New()
	ta.Placeholder = "Type your note..."
	ta.CharLimit = -1
	ta.Prompt = ""
	ta.SetHeight(height - yMargin)
	ta.SetWidth(width - 2)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Border(lipgloss.ASCIIBorder())
	ta.ShowLineNumbers = false
	mdR, err := glamour.NewTermRenderer()
	if err != nil {
		fmt.Println("failed to create md renderer")
		os.Exit(1)
	}

	style := lipgloss.NewStyle().Border(lipgloss.ASCIIBorder()).Height(height - yMargin).Width(width - 2)

	return NoteInput{
		textArea: ta,
		style:    style,
		mdRender: mdR,
	}
}

func (ni NoteInput) Init() tea.Cmd {
	return textarea.Blink
}

func (ni NoteInput) Update(msg tea.Msg) (NoteComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ni.style.Height(msg.Height - ni.yMargin)
		ni.style.Width(msg.Width - 2)
	}

	var cmd tea.Cmd
	ni.textArea, cmd = ni.textArea.Update(msg)

	return ni, cmd
}

func (ni NoteInput) View() string {
	return ni.style.Render(ni.textArea.View())
}

func (ni NoteInput) Blur() NoteComponent {
	ni.textArea.Blur()
	return ni
}

func (ni NoteInput) Focus() NoteComponent {
	ni.textArea.Focus()
	return ni
}

func (ni NoteInput) Text() string {
	return ni.textArea.Value()
}

func (ni NoteInput) Clear() NoteComponent {
	ni.textArea.SetValue("")

	return ni
}
