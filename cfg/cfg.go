package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type BonConfig struct {
	RootDir string
	ArchDir string
	BonFile string
	BonCfg  string
	Editor  string
	MaxDays int `json:"max_note_days"`
}

var Config *BonConfig = &BonConfig{
	MaxDays: 7,
	Editor:  "lvim",
}

func Initialize() {
	var home string
	var err error
	home, exist := os.LookupEnv("BON_HOME")
	if !exist {
		home, err = os.UserHomeDir()
		if err != nil {
			fmt.Printf("failed to load home dir: %s", err.Error())
			os.Exit(1)
		}
	}
	Config.RootDir = home + "/.bon"
	Config.ArchDir = Config.RootDir + "/archive"
	Config.BonFile = Config.RootDir + "/bon.json"
	Config.BonCfg = Config.RootDir + "/boncfg.json"

	createMissingDirs([]string{Config.RootDir, Config.ArchDir})
	createMissingFiles(Config.BonFile, Config.BonCfg)

	cfg, err := os.ReadFile(Config.BonCfg)
	if err != nil {
		fmt.Printf("error in reading config: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(cfg, Config)
	if err != nil {
		fmt.Printf("error in parsing config: %s", err.Error())
	}
}

func createMissingDirs(dirs []string) {
	for _, d := range dirs {
		if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(d, os.ModePerm)
		}
	}
}

func createMissingFiles(bonFile string, bonCfg string) {
	if _, err := os.Stat(bonFile); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(bonFile)
		defer f.Close()
		if err != nil {
			fmt.Printf("failed to create bon file: %s", err.Error())
			os.Exit(1)
		}
		f.WriteString("[]")
	}

	if _, err := os.Stat(bonCfg); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(bonCfg)
		defer f.Close()
		if err != nil {
			fmt.Printf("failed to create bon file: %s", err.Error())
			os.Exit(1)
		}
		f.WriteString("{}")
	}
}
