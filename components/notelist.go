package components

import (
	"errors"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/john-marinelli/bon/dstructs"
	"github.com/john-marinelli/bon/types"
	"github.com/john-marinelli/bon/util"
)

var acStyle lipgloss.Style = lipgloss.NewStyle().Margin(1, 0, 0, 0)

func noteStatus(dl int) types.NoteStatus {
	if dl < 1 {
		return types.AboutToDelete
	}
	if dl > 2 {
		return types.Safe
	}

	return types.DayFromDelete
}

type item struct {
	date  string
	short string
	full  string
	id    int
}

func (li item) Title() string {
	return li.short
}

func (li item) Description() string {
	return "Date: " + li.date + " Id: " + strconv.Itoa(li.id)
}

func (li item) FilterValue() string {
	return li.full
}

func (li item) GetId() int {
	return li.id
}

type NoteList struct {
	list      list.Model
	notes     []dstructs.Note
	input     NoteComponent
	height    int
	calcWidth func(w int) int
	focused   types.BonListMode
	selNote   dstructs.Note
}

func NewNoteList(
	notes []dstructs.Note,
	width int,
	height int,
	calcWidth func(w int) int,
) NoteList {
	li := []list.Item{}
	for _, n := range notes {
		li = append(li, item{
			date:  n.Date.Format(time.RFC3339),
			short: string(noteStatus(n.DaysLeft)) + " " + util.Truncate(n.Content, 20),
			full:  n.Content,
			id:    n.Id,
		})
	}
	ld := list.NewDefaultDelegate()
	ld.Styles.NormalDesc = ld.Styles.NormalDesc.Foreground(lipgloss.Color("#779ed1"))
	ld.Styles.NormalTitle = ld.Styles.NormalTitle.Foreground(lipgloss.Color("#779ed1"))
	ld.Styles.SelectedTitle = ld.Styles.SelectedTitle.BorderForeground(lipgloss.Color("#dbc89a"))
	ld.Styles.SelectedTitle = ld.Styles.SelectedTitle.Foreground(lipgloss.Color("#dbc89a"))
	ld.Styles.SelectedDesc = ld.Styles.SelectedDesc.BorderForeground(lipgloss.Color("#dbc89a"))
	ld.Styles.SelectedDesc = ld.Styles.SelectedDesc.Foreground(lipgloss.Color("#dbc89a"))
	l := list.New(li, ld, 0, 0)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	//TODO change style of selected note
	l.SetSize(width, height)

	ft, err := dstructs.NewFTree()
	if err != nil {
		panic(err)
	}

	ac := NewAutoComplete(ft.AllPaths, width)

	return NoteList{
		list:    l,
		notes:   notes,
		input:   ac,
		focused: types.BrowseMode,
	}
}

func (nl NoteList) Init() tea.Cmd {
	return nil
}

func (nl NoteList) Update(msg tea.Msg) (NoteList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			if nl.focused == types.BrowseMode {
				nl.focused = types.SaveMode
				nl.selNote = nl.notes[nl.list.Index()]
			}
		case "enter":
			if nl.focused == types.SaveMode {
				nl = nl.saveSelected()
			}
		case "esc":
			if nl.focused == types.SaveMode {
				nl.input = nl.input.Clear()
				nl.focused = types.BrowseMode
			}
		}
	case tea.WindowSizeMsg:
		nl.list.SetSize(nl.calcWidth(msg.Width), msg.Height-1)
	}
	var cmd tea.Cmd
	if nl.focused == types.BrowseMode {
		nl.list, cmd = nl.list.Update(msg)
	} else {
		nl.input, cmd = nl.input.Update(msg)
	}

	return nl, cmd
}

func (nl NoteList) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		nl.list.View(),
		acStyle.Render(nl.input.View()),
	)
}

func (nl NoteList) GetSelectedContent() (string, string, error) {
	sel, ok := nl.list.SelectedItem().(item)
	if ok {
		return sel.date, sel.full, nil
	}
	return "", "", errors.New("failed to open bon note")
}

func (nl NoteList) saveSelected() NoteList {
	dstructs.SaveNote(nl.input.Text(), nl.selNote.Content)
	n, err := dstructs.DeleteBonNote(nl.selNote.Id)
	if err != nil {
		panic(err)
	}

	nl.setItems(n)
	return nl
}

func (nl NoteList) Saving() bool {
	if nl.focused == types.SaveMode {
		return true
	}

	return false
}

func (nl *NoteList) setItems(notes []dstructs.Note) {
	li := []list.Item{}
	for _, n := range notes {
		li = append(li, item{
			date:  n.Date.Format(time.RFC3339),
			short: string(noteStatus(n.DaysLeft)) + " " + util.Truncate(n.Content, 20),
			full:  n.Content,
			id:    n.Id,
		})
	}
	nl.notes = notes
	nl.list.SetItems(li)
}
