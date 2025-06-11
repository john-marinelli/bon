package bon

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/john-marinelli/bon/data"
	"github.com/john-marinelli/bon/types"
)

type Bon struct {
	model tea.Model
}

func NewBon(screen types.BonScreen) Bon {
	notes, err := data.LoadAndClearNotes()
	var m tea.Model
	if screen == types.InputScreen {
		m = NewBonInput()
	} else {
		m = NewBonView(notes, err)
	}
	return Bon{
		model: m,
	}
}

func (b Bon) Init() tea.Cmd {

	return b.model.Init()
}

func (b Bon) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b.model.Update(msg)
}

func (b Bon) View() string {
	return b.model.View()
}
