package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getStoragePath() string {
	// Find the user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to local directory if home is inaccessible
		return "nexa_data.json"
	}

	// Define the folder path (~/.nexa)
	nexaDir := filepath.Join(home, ".nexa")

	// Create the folder if it doesn't exist
	// 0755 gives the user read/write/execute permissions
	err = os.MkdirAll(nexaDir, 0755)
	if err != nil {
		// If we can't create the folder, fallback to local
		return "nexa_data.json"
	}

	// Return the full path to the file
	return filepath.Join(nexaDir, "nexa_data.json")
}

func (a *App) Save(args []string) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	configPath := getStoragePath()

	// Marshal transforms the Projects slice into a JSON byte array
	data, err := json.MarshalIndent(a.Projects, "", "  ")
	if err != nil {
		return err
	}

	// WriteFile creates or overwrites the file with the JSON data
	return os.WriteFile(configPath, data, 0644)
}
func (a *App) Load(args []string) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	configPath := getStoragePath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err // Usually means file doesn't exist yet
	}

	var loadedProjects []*Project
	if err := json.Unmarshal(data, &loadedProjects); err != nil {
		return err
	}

	// Replace old slice and REBUILD MAPS
	a.Projects = loadedProjects
	a.ProjectMap = make(map[string]*Project)

	for _, p := range a.Projects {
		a.ProjectMap[p.Name] = p
		p.CategoryMap = make(map[string]*Category)

		for _, c := range p.Categories {
			p.CategoryMap[c.Name] = c
			c.TaskMap = make(map[string]*Task)

			for _, t := range c.Tasks {
				c.TaskMap[t.Name] = t
			}
		}
		p.syncTasks()
	}
	a.InitCommands()
	return nil
}
