package components

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type NoteViewer struct {
	viewer  viewport.Model
	content string
	title   string
	md      *glamour.TermRenderer
}

func NewNoteViewer() NoteViewer {
	vp := viewport.New(0, 0)
	md, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		fmt.Println("renderer failed to init")
		os.Exit(1)
	}
	return NoteViewer{
		viewer: vp,
		md:     md,
	}
}

func (nv NoteViewer) GetContent() string {
	return nv.viewer.View()
}

func (nv NoteViewer) Init() tea.Cmd {
	return nil
}

func (nv NoteViewer) Update(msg tea.Msg) (NoteViewer, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(nv.headerView())
		footerHeight := lipgloss.Height(nv.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		nv.viewer.YPosition = headerHeight
		nv.viewer.Width = msg.Width
		nv.viewer.Height = msg.Height - verticalMarginHeight
	}

	nv.viewer, cmd = nv.viewer.Update(msg)

	return nv, cmd
}

func (nv NoteViewer) View() string {

	return fmt.Sprintf(
		"%s\n%s\n%s",
		nv.headerView(),
		nv.viewer.View(),
		nv.footerView(),
	)
}

func (nv *NoteViewer) SetTitle(title string) {
	nv.title = title
}

func (nv *NoteViewer) SetContent(cnt string) {
	txt, err := nv.md.Render(cnt)
	if err != nil {
		fmt.Println("content failed to render")
		os.Exit(1)
	}
	nv.viewer.SetContent(txt)
	headerHeight := lipgloss.Height(nv.headerView())
	nv.viewer.YPosition = headerHeight
}

func (nv NoteViewer) headerView() string {
	t := titleStyle.Render(nv.title)
	l := strings.Repeat("-", max(0, nv.viewer.Width-lipgloss.Width(t)))
	return lipgloss.JoinHorizontal(lipgloss.Center, t, l)
}

func (nv NoteViewer) footerView() string {
	ft := infoStyle.Render(fmt.Sprintf("%3.f%%", nv.viewer.ScrollPercent()*100))
	l := strings.Repeat("-", max(0, nv.viewer.Width-lipgloss.Width(ft)))
	return lipgloss.JoinHorizontal(lipgloss.Center, l, ft)
}
