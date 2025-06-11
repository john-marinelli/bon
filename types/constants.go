package types

import "github.com/charmbracelet/lipgloss"

type InputMode int
type BonScreen int
type BonViewMode int
type BonListMode int

type UpdateFilesMsg bool

const (
	PathInput InputMode = 0
	NoteInput InputMode = 1
)

const (
	PickerMode BonViewMode = 0
	ListMode   BonViewMode = 1
	ViewerMode BonViewMode = 2
)

const (
	SaveMode   BonListMode = 0
	BrowseMode BonListMode = 1
)

const (
	ViewScreen  BonScreen = 0
	InputScreen BonScreen = 1
)

var (
	AboutToDelete string = lipgloss.NewStyle().Foreground(lipgloss.Color("#64b579")).Render("○")
	DayFromDelete string = lipgloss.NewStyle().Foreground(lipgloss.Color("#dbb05a")).Bold(true).Render("!")
	Safe          string = lipgloss.NewStyle().Foreground(lipgloss.Color("#e84646")).Bold(true).Render("X")
)
