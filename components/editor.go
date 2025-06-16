package components

import (
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/john-marinelli/bon/cfg"
	"github.com/john-marinelli/bon/data"
	"github.com/john-marinelli/bon/types"
)

type EditorDoneMsg struct {
	Err error
}

type Editor struct {
	program string
	Err     error
	tmpFile string
	isBon   bool
}

func NewEditor() Editor {
	return Editor{
		program: cfg.Config.Editor,
		tmpFile: cfg.Config.RootDir + "/temp.md",
	}
}

func (e *Editor) rmTempFile() {
	err := os.Remove(e.tmpFile)
	e.Err = err
}

func (e *Editor) getTemp() (string, error) {
	f, err := os.ReadFile(e.tmpFile)
	return string(f), err
}

func (e *Editor) transferTemp() {
	tmp, err := os.ReadFile(e.tmpFile)
	if err != nil {
		e.Err = err
		return
	}
	data.SaveNote("", string(tmp))
}

func (e Editor) Open(path string, screen types.BonScreen) tea.Cmd {
	var p string

	if screen == types.InputScreen && path != "" {
		p = cfg.Config.ArchDir + "/" + path + ".md"
		dPath := filepath.Dir(p)
		if _, err := os.Stat(dPath); os.IsNotExist(err) {
			err := os.MkdirAll(dPath, os.ModePerm)
			if err != nil {
				return func() tea.Msg {
					return EditorDoneMsg{
						Err: err,
					}
				}
			}
		}
	} else if screen == types.InputScreen {
		p = e.tmpFile
	} else {
		p = path
	}

	cmd := exec.Command(e.program, p)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err == nil && path == "" && screen == types.InputScreen {
			e.transferTemp()
			e.rmTempFile()
		}
		return EditorDoneMsg{
			Err: err,
		}
	})

}

func (e Editor) OpenBon(id int, note string) tea.Cmd {
	f, err := os.Create(e.tmpFile)
	if err != nil {
		return func() tea.Msg {
			return EditorDoneMsg{
				Err: err,
			}
		}
	}
	f.WriteString(note)
	f.Close()
	cmd := exec.Command(e.program, e.tmpFile)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		content, err := e.getTemp()
		if err != nil {
			return EditorDoneMsg{
				Err: err,
			}
		}

		err = data.EditBonNote(id, content)
		e.rmTempFile()
		return EditorDoneMsg{
			Err: err,
		}
	})

}
