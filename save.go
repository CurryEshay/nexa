package main

import (
	"bytes"
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
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}

	// WriteFile creates or overwrites the file with the JSON data
	return os.WriteFile(configPath, data, 0644)
}
func Load() (*App, error) {
	configPath := getStoragePath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Initialize with defaults (aggression constants, etc.)
	loadedApp := NewApp()

	// Trim whitespace to check the first actual character
	trimmedData := bytes.TrimSpace(data)
	if len(trimmedData) == 0 {
		return loadedApp, nil
	}

	if trimmedData[0] == '[' {
		// --- OLD SYSTEM: The file is just [project1, project2, ...] ---
		var oldProjects []*Project
		if err := json.Unmarshal(data, &oldProjects); err != nil {
			return nil, err
		}
		loadedApp.Projects = oldProjects
	} else {
		// --- NEW SYSTEM: The file is {"projects": [...], "time_aggression": 1.5, ...} ---
		if err := json.Unmarshal(data, loadedApp); err != nil {
			return nil, err
		}
	}

	// Rebuild the maps and pointers for both cases
	loadedApp.reconstructInternalState()

	return loadedApp, nil
}

// Helper to keep Load() clean
func (a *App) reconstructInternalState() {
	a.ProjectMap = make(map[string]*Project)
	for _, p := range a.Projects {
		if p == nil {
			continue
		}
		a.ProjectMap[p.Name] = p

		p.CategoryMap = make(map[string]*Category)
		for _, c := range p.Categories {
			if c == nil {
				continue
			}
			p.CategoryMap[c.Name] = c

			c.TaskMap = make(map[string]*Task)
			for _, t := range c.Tasks {
				if t == nil {
					continue
				}
				c.TaskMap[t.Name] = t
			}
		}
		p.syncTasks(a)
	}
	a.InitCommands()
}
