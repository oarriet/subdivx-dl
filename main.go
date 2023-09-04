package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oarriet/subdivx-dl/tui"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	p := tea.NewProgram(tui.NewModel())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
