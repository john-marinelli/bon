package dstructs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/john-marinelli/bon/cfg"
)

type FTree struct {
	AllPaths []string
}

func NewFTree() (FTree, error) {
	ft := FTree{}
	ft.AllPaths = []string{}

	err := filepath.Walk(cfg.Config.ArchDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			p := strings.Trim(strings.Replace(path, cfg.Config.ArchDir, "", -1), "/")
			ft.AllPaths = append(ft.AllPaths, p)
		}
		return nil
	})

	return ft, err
}
