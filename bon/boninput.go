package bon

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/john-marinelli/bon/components"
	"github.com/john-marinelli/bon/data"
	"github.com/john-marinelli/bon/types"
)

var pathBorder = lipgloss.NewStyle().Border(lipgloss.ASCIIBorder())

type errMsg error

type BonInput struct {
	pathInput components.AutoComplete
	teaProg   *tea.Program
	editor    components.Editor
	ft        data.FTree
	err       error
}

func NewBonInput() BonInput {
	bi := BonInput{}
	width, _, _ := term.GetSize(uintptr(os.Stdout.Fd()))

	ft, err := data.NewFTree()
	bi.err = err

	p := components.NewAutoComplete(ft.AllPaths, width)
	bi.pathInput = p

	bi.editor = components.NewEditor()

	return bi
}

func (bi BonInput) Init() tea.Cmd {
	return nil
}

func (bi BonInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case errMsg:
		bi.err = msg
		return bi, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd = bi.editor.Open(bi.pathInput.Text(), types.InputScreen)
			return bi, cmd
		case "ctrl+c":
			return bi, tea.Quit
		}
	case components.EditorDoneMsg:
		if msg.Err != nil {
			fmt.Println(msg.Err.Error())
		}
		return bi, tea.Quit
	}

	bi.pathInput, cmd = bi.pathInput.Update(msg)

	return bi, cmd
}

func (bi BonInput) View() string {
	return fmt.Sprintf(
		"%s",
		bi.pathInput.View(),
	)
}
