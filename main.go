package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/john-marinelli/bon/bon"
	"github.com/john-marinelli/bon/cfg"
	"github.com/john-marinelli/bon/util"
)

func main() {
	logFile, _ := os.OpenFile("bon.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	log.SetOutput(logFile)
	args := os.Args
	screen, err := util.ParseArgs(args)
	cfg.Initialize()
	switch err {
	case util.ArgNumberErr:
		fmt.Println(
			"Incorrect number of arguments.\nUsage:\n\t" +
				"bon (for note input)\n\tbon bon (for archive viewing)",
		)
		os.Exit(0)
	case util.WrongArgErr:
		fmt.Println(
			"Incorrect arguments.\nUsage:\n\t" +
				"bon (for note input)\n\tbon bon (for archive viewing)",
		)
		os.Exit(0)
	}
	b := bon.NewBon(screen)
	p := tea.NewProgram(b, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		p.Quit()
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}
}
