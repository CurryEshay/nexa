package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (a *App) HandlePath(pathString string) (string, string, string, error) {
	pathString = strings.ReplaceAll(pathString, ",", "/")
	path := strings.Split(pathString, "/")

	if len(path) == 0 {
		return "", "", "", errors.New("Path not properly deliminated (Please use , or /)")
	}

	// Trim all path elemets
	for i, arg := range path {
		path[i] = strings.TrimSpace(arg)
	}

	// Replace * or . elements with selected
	for index, element := range path {
		if element == "*" || element == "." {
			switch index {
			case 0: // Project level
				if len(a.Projects) == 0 {
					return "", "", "", errors.New("no projects exist to reference")
				}
				path[index] = a.Projects[a.currentProjectIndex].Name

			case 1: // Category level
				if a.currentCategoryIndex == -1 && (path[2] != "*" && path[2] != ".") {
					return "", "", "", errors.New("cannot reference null category")
				} else if path[2] == "*" || path[2] == "." {
					continue
				}
				// Check if project actually has categories
				proj := a.Projects[a.currentProjectIndex]
				if len(proj.Categories) == 0 || a.currentCategoryIndex >= len(proj.Categories) {
					return "", "", "", errors.New("category selection out of sync")
				}
				path[index] = proj.Categories[a.currentCategoryIndex].Name

			case 2: // Task level
				if a.currentTaskIndex == -1 {
					return "", "", "", errors.New("no task selected")
				}
				// Determine if we are in Category view or Project view
				var tasks []*Task
				if a.currentCategoryIndex != -1 {
					tasks = a.Projects[a.currentProjectIndex].Categories[a.currentCategoryIndex].Tasks
					path[index] = tasks[a.currentTaskIndex].Name
				} else {
					tasks = a.Projects[a.currentProjectIndex].Tasks
					path[index] = tasks[a.currentTaskIndex].Name
					path[1] = tasks[a.currentTaskIndex].CategoryName
				}

				if len(tasks) == 0 || a.currentTaskIndex >= len(tasks) {
					return "", "", "", errors.New("task selection out of sync")
				}

			}
		}
	}

	// Just project
	if len(path) == 1 {
		return path[0], "", "", nil
	}

	// Project + category
	if len(path) == 2 {
		return path[0], path[1], "", nil
	}

	// Full path, project + category + task
	if len(path) == 3 {
		return path[0], path[1], path[2], nil
	}

	return "", "", "", errors.New("Inappropriate number of path elements, please input {projct name}/{category name}/{}")
}

func parseIncrement(input string) (uint, uint, uint, uint, error) {
	bits := strings.Split(input, "_")
	if len(bits) != 2 {
		return 0, 0, 0, 0, fmt.Errorf("invalid shorthand format")
	}

	amount, err := strconv.Atoi(bits[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid shortcut number: %s", bits[0])
	}
	switch bits[1] {
	case "d":
		return 0, 0, 0, uint(amount), nil
	case "w":
		return 0, 0, uint(amount), 0, nil
	case "m":
		return 0, uint(amount), 0, 0, nil
	case "y":
		return uint(amount), 0, 0, 0, nil
	default:
		return 0, 0, 0, 0, fmt.Errorf("unknown unit: %s", bits[1])
	}
}

// parseFlexibleDeadline is the entry point that coordinates the date and time parsing.
func parseFlexibleDeadline(input string) (time.Time, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	input = strings.ReplaceAll(input, "/", "-")

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return time.Time{}, fmt.Errorf("empty input")
	}

	var targetDate time.Time
	var err error

	// 1. Resolve the Date (Absolute or Shorthand)
	if strings.Contains(parts[0], "_") {
		targetDate, err = processShorthand(parts[0])
	} else {
		targetDate, err = processAbsoluteDate(parts[0])
	}

	if err != nil {
		return time.Time{}, err
	}

	// 2. Resolve the Time (Provided or Default)
	timePart := "11:59pm"
	if len(parts) > 1 {
		timePart = parts[1]
	}

	return applyTimeToDate(targetDate, timePart)
}

// processShorthand calculates a date based on relative jumps (1_d, 2_w, etc.)
func processShorthand(shorthand string) (time.Time, error) {
	now := time.Now()
	bits := strings.Split(shorthand, "_")
	if len(bits) != 2 {
		return time.Time{}, errors.New("invalid shorthand format")
	}

	amount, err := strconv.Atoi(bits[0])
	if err != nil {
		return time.Time{}, errors.New("invalid shorthand format")
	}

	switch bits[1] {
	case "d":
		return now.AddDate(0, 0, amount), nil
	case "w":
		return now.AddDate(0, 0, amount*7), nil
	case "m":
		return now.AddDate(0, amount, 0), nil
	case "y":
		return now.AddDate(amount, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown unit: %s", bits[1])
	}
}

// processAbsoluteDate handles DD-MM-YYYY strings
func processAbsoluteDate(dateStr string) (time.Time, error) {
	dParts := strings.Split(dateStr, "-")
	if len(dParts) != 3 {
		return time.Time{}, fmt.Errorf("use DD-MM-YYYY or shorthand (1_d)")
	}

	year, _ := strconv.Atoi(dParts[2])  // Year is at index 2
	month, _ := strconv.Atoi(dParts[1]) // Month is at index 1
	day, _ := strconv.Atoi(dParts[0])   // Day is at index 0

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// applyTimeToDate parses the time string and merges it with the target date
func applyTimeToDate(baseDate time.Time, timeStr string) (time.Time, error) {
	var hour, min int
	var t time.Time
	var err error

	if strings.Contains(timeStr, "am") || strings.Contains(timeStr, "pm") {
		t, err = time.Parse("3:04pm", timeStr)
	} else {
		t, err = time.Parse("15:04", timeStr)
	}

	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, min = t.Hour(), t.Minute()
	return time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, min, 0, 0, time.Local), nil
}
