package bon

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/john-marinelli/bon/cfg"
	"github.com/john-marinelli/bon/components"
	"github.com/john-marinelli/bon/dstructs"
	"github.com/john-marinelli/bon/types"
)

var (
	archTitleStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Top = ""
		b.Left = ""
		b.Right = ""
		st := lipgloss.NewStyle().BorderStyle(b).MarginBottom(1)
		st = st.Foreground(
			lipgloss.Color("#dbc89a"),
		).BorderForeground(
			lipgloss.Color("#779ed1"),
		).Margin(0, 0, 0, 2)
		return st
	}()
	selScreenStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		st := lipgloss.NewStyle().BorderStyle(
			b,
		).BorderForeground(
			lipgloss.Color("#9ac3db"),
		)
		return st
	}
	disScreenStyle = func() lipgloss.Style {
		st := lipgloss.NewStyle().Padding(1)
		return st
	}
)

type BonView struct {
	picker      filepicker.Model
	pickerStyle lipgloss.Style
	viewer      components.NoteViewer
	list        components.NoteList
	listStyle   lipgloss.Style
	noteList    components.NoteList
	focused     types.BonViewMode
	prevFocused types.BonViewMode
	selected    string
	quitting    bool
	baseDir     string
	err         error
}

type clearErrMsg struct{}

func clearErrAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrMsg{}
	})
}

func NewBonView(notes []dstructs.Note, err error) BonView {
	width, height, _ := term.GetSize(uintptr(os.Stdout.Fd()))
	p := filepicker.New()

	p.AllowedTypes = []string{".md"}
	p.ShowPermissions = false
	p.ShowSize = false
	p.CurrentDirectory = cfg.Config.ArchDir
	p.Styles.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#dbc89a"))
	p.Styles.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color("#dbc89a"))
	p.Styles.Directory = lipgloss.NewStyle().Foreground(lipgloss.Color("#779ed1"))

	v := components.NewNoteViewer()

	nl := components.NewNoteList(notes, width/2, height)

	arch := BonView{
		picker:      p,
		viewer:      v,
		baseDir:     cfg.Config.ArchDir,
		list:        nl,
		focused:     types.PickerMode,
		listStyle:   lipgloss.NewStyle().Width((width / 2) - 2).Height(height - 2),
		pickerStyle: lipgloss.NewStyle().Width((width / 2) - 2).Height(height - 2),
		err:         err,
	}

	if err != nil {
		arch.err = err
	}

	return arch
}

func (a BonView) Init() tea.Cmd {
	return a.picker.Init()
}

func (a BonView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			a.quitting = true
			return a, tea.Quit
		case "esc":
			if a.focused == types.ViewerMode {
				a.focused = types.PickerMode
			}
		case "ctrl+b", "2":
			if a.focused == types.PickerMode {
				a.focused = types.ListMode
			}
		case "ctrl+a", "1":
			if a.focused == types.ListMode {
				a.focused = types.PickerMode
			}
		case "enter":
			if a.focused == types.ListMode && !a.list.Saving() {
				title, note, err := a.list.GetSelectedContent()
				if err != nil {
					a.err = err
				} else {
					a.viewer.SetContent(note)
					a.viewer.SetTitle(title)
					a.focused = types.ViewerMode
					width, height, err := term.GetSize(uintptr(os.Stdout.Fd()))
					if err != nil {
						a.err = err
					}
					return a, func() tea.Msg {
						return tea.WindowSizeMsg{
							Width:  width,
							Height: height,
						}
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		a.pickerStyle = a.pickerStyle.Height(msg.Height)
		a.listStyle = a.listStyle.Height(msg.Height)

		a.pickerStyle = a.pickerStyle.Width(msg.Width / 2)
		a.listStyle = a.listStyle.Width(msg.Width / 2)
	case clearErrMsg:
		a.err = nil
	}

	var cmd tea.Cmd
	switch a.focused {
	case types.PickerMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			ks := msg.String()
			if (ks == "esc" || ks == "h") &&
				a.picker.CurrentDirectory == cfg.Config.ArchDir {
				return a, nil
			}
		}
		a.picker, cmd = a.picker.Update(msg)
	case types.ViewerMode:
		a.viewer, cmd = a.viewer.Update(msg)
	case types.ListMode:
		a.list, cmd = a.list.Update(msg)
	}

	if didSelect, path := a.picker.DidSelectFile(msg); didSelect {
		cnt, err := a.openSelected(path)
		if err != nil {
			a.err = err
			return a, tea.Batch(cmd, clearErrAfter(2*time.Second))
		}
		t := strings.Split(path, "/")

		a.viewer.SetTitle(strings.Split(t[len(t)-1], ".")[0])
		a.viewer.SetContent(cnt)
		a.focused = types.ViewerMode
		width, height, err := term.GetSize(uintptr(os.Stdout.Fd()))
		if err != nil {
			a.err = err
		}

		return a, func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  width,
				Height: height,
			}
		}
	}

	if didSelect, path := a.picker.DidSelectDisabledFile(msg); didSelect {
		a.err = errors.New(path + " is not valid.")
		a.selected = ""
		return a, tea.Batch(cmd, clearErrAfter(2*time.Second))
	}

	return a, cmd
}

func (a BonView) View() string {
	if a.quitting {
		return "Goodby! :)"
	}

	if a.focused == types.ViewerMode {
		return a.viewer.View()
	}

	var s strings.Builder
	s.WriteString("\n")

	if a.err != nil {
		s.WriteString(a.picker.Styles.DisabledFile.Render(a.err.Error()))
	} else {
		s.WriteString(archTitleStyle.Render(".bon archive"))
	}
	s.WriteString("\n" + a.picker.View() + "\n")
	p := a.pickerStyle.Render(s.String())
	v := a.list.View()
	if a.focused == types.PickerMode {
		p = selScreenStyle().Render(p)
	} else {
		p = disScreenStyle().Render(p)
	}
	if a.focused == types.ListMode {
		v = selScreenStyle().Render(v)
	} else {
		v = disScreenStyle().Render(v)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		p,
		v,
	)
}

func (a BonView) openSelected(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}
