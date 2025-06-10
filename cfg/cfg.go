package cfg

import "os"

type BonConfig struct {
	RootDir string
	ArchDir string
	MaxDays int
}

var Config BonConfig

func init() {
	Config.RootDir = os.ExpandEnv("$HOME/.bon")
	Config.ArchDir = Config.RootDir + "/archive"
}
