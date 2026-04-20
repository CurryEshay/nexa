package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app, err := Load()
	if err != nil {
		fmt.Println(err)
	}

	if len(app.Projects) == 0 {
		app.NewProject("init")
		app.NewCategory("init", "init")
	}

	m := NewModel(app)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
