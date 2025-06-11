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

type NoteInputState int

const (
	PathInput = iota
	NoteInput
)

type BonInput struct {
	inputs  []components.NoteComponent
	current types.InputMode
	ft      data.FTree
	err     error
}

func NewBonInput() BonInput {
	bi := BonInput{}
	in := make([]components.NoteComponent, 2)
	width, height, _ := term.GetSize(uintptr(os.Stdout.Fd()))

	ft, err := data.NewFTree()
	bi.err = err

	p := components.NewAutoComplete(ft.AllPaths, width)
	in[types.PathInput] = p

	ni := components.NewNoteInput(width, height, 3)
	in[types.NoteInput] = ni

	bi.inputs = in

	return bi
}

func (bi BonInput) Init() tea.Cmd {
	return nil
}

func (bi BonInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for i := range bi.inputs {
			bi.inputs[i], _ = bi.inputs[i].Update(msg)
		}
	case errMsg:
		bi.err = msg
		return bi, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if bi.current == types.PathInput {
				bi.inputs[bi.current] = bi.inputs[bi.current].Blur()
				bi.current = types.NoteInput
				bi.inputs[bi.current] = bi.inputs[bi.current].Focus()
				return bi, nil
			}
		case "ctrl+c":
			return bi, tea.Quit
		case "ctrl+s":
			data.SaveNote(bi.inputs[types.PathInput].Text(), bi.inputs[types.NoteInput].Text())
			if bi.current == types.NoteInput {
				return bi, tea.Quit
			}
		}
	}

	bi.inputs[bi.current], _ = bi.inputs[bi.current].Update(msg)

	return bi, nil
}

func (bi BonInput) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		bi.inputs[types.PathInput].View(),
		"\n",
		bi.inputs[types.NoteInput].View(),
	)
}
