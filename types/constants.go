package types

type NoteStatus rune
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

const (
	AboutToDelete NoteStatus = '❌'
	DayFromDelete NoteStatus = '❗'
	Safe          NoteStatus = '🟢'
)
