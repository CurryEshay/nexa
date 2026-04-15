package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	// These are the Charmbracelet libraries you just downloaded
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// The "Selected" look - White text on a Purple background
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")).
			Bold(true).
			PaddingLeft(1)

	// Projects: Bold and underlined to distinguish as headers
	projStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Bold(true).
			Underline(true)

	// Categories: Dimmer than projects
	catStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).Bold(true)

	// Tasks: Clean and spaced
	taskStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).Bold(true)

	// Done/Urgent
	doneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
	urgentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

type model struct {
	app                *App
	input              textinput.Model // The actual typing box
	width              int             // For relative sizing
	height             int
	error              string
	commands           []string
	commandIndex       uint
	commandScrollIndex uint
}

// Initializer
func NewModel(a *App) model {
	ti := textinput.New()
	ti.Placeholder = "Type command. e.g. mk;project "
	ti.Focus()
	return model{
		app:                a,
		input:              ti,
		commandIndex:       0,
		commandScrollIndex: 0,
		commands:           []string{""},
	}
}

// Init is the first function called when the program starts.
// It can return a command (tea.Cmd) to perform an action (like a timer or file read).
func (m model) Init() tea.Cmd {
	// Just return nil if you have no startup actions.
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.app.Save(nil)
			m.app.QuitProgram(nil)
			return m, tea.Quit

		case "shift+up":

			// Check if there is space to move up within a project
			if m.app.currentCategoryIndex > -1 {
				m.app.currentCategoryIndex--

				if m.app.currentCategoryIndex != -1 {
					if len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks) == 0 {
						m.app.currentTaskIndex = -1
					} else {
						m.app.currentTaskIndex = 0
					}
				} else if len(m.app.Projects[m.app.currentProjectIndex].Tasks) == 0 {
					m.app.currentTaskIndex = -1
				} else {
					m.app.currentTaskIndex = 0
				}

				// If there is no space to move up, i.e. a project is selected, if there is another project above, move to it's last task
			} else if m.app.currentProjectIndex > 0 && m.app.currentCategoryIndex == -1 {
				m.app.currentProjectIndex--
				m.app.currentCategoryIndex = len(m.app.Projects[m.app.currentProjectIndex].Categories) - 1

				// If next project is not empty
				if m.app.currentCategoryIndex != -1 {
					m.app.currentTaskIndex = 0
					// If next project is empty
				} else {
					m.app.currentTaskIndex = -1
				}

			}
			return m, nil

		case "shift+down":

			// Check if there is space below to move
			if m.app.currentCategoryIndex < len(m.app.Projects[m.app.currentProjectIndex].Categories)-1 {
				m.app.currentCategoryIndex++

				if len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks) == 0 {
					m.app.currentTaskIndex = -1
				} else {
					m.app.currentTaskIndex = 0
				}

				// If at the end of a project's categories, check if there is a project below and move to it
			} else if m.app.currentProjectIndex < len(m.app.Projects)-1 && m.app.currentCategoryIndex == len(m.app.Projects[m.app.currentProjectIndex].Categories)-1 {
				m.app.currentProjectIndex++
				m.app.currentCategoryIndex = -1

				if len(m.app.Projects[m.app.currentProjectIndex].Tasks) == 0 {
					m.app.currentTaskIndex = -1
				} else {
					m.app.currentTaskIndex = 0
				}

			}
			return m, nil

		// done selected task
		case "ctrl+d":
			m.app.Save(nil)

			// Check if a task is currently highlighted
			if m.app.currentTaskIndex != -1 {
				// If on a category
				if m.app.currentCategoryIndex != -1 {
					err := m.app.doneTask(m.app.Projects[m.app.currentProjectIndex].Name, m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Name, m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks[m.app.currentTaskIndex].Name)
					if err != nil {
						m.error = err.Error()
					} else {
						m.error = "" // Clear it on success
					}
				} else {
					task := m.app.Projects[m.app.currentProjectIndex].Tasks[m.app.currentTaskIndex]
					err := m.app.doneTask(task.ProjectName, task.CategoryName, task.Name)
					if err != nil {
						m.error = err.Error()
					} else {
						m.error = "" // Clear it on success
					}
				}

			}
			if m.app.currentCategoryIndex != -1 {
				if m.app.currentTaskIndex > len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks)-1 {
					m.app.currentTaskIndex = len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks) - 1
				}
			} else {
				if m.app.currentTaskIndex > len(m.app.Projects[m.app.currentProjectIndex].Tasks)-1 {
					m.app.currentTaskIndex = len(m.app.Projects[m.app.currentProjectIndex].Tasks) - 1
				}
			}

			return m, nil

		case "ctrl+z":
			m.app.Save(nil)
			err := m.app.restoreTask()
			if err != nil {
				m.error = err.Error()
			} else {
				m.error = "" // Clear it on success
			}

		// scroll through previous commands
		case "up":
			if m.commandScrollIndex > 0 {
				m.commandScrollIndex--
			}
			m.input.SetValue(m.commands[m.commandScrollIndex])
			return m, nil

		// scroll through previous commands
		case "down":
			if m.commandScrollIndex < m.commandIndex {
				m.commandScrollIndex++
			}
			m.input.SetValue(m.commands[m.commandScrollIndex])
			return m, nil

		// scroll through tasks
		case "shift+right":

			// If on a category
			if m.app.currentCategoryIndex != -1 {

				// Check if there is space to scroll down and if and the length of the category is not 0
				if m.app.currentTaskIndex < len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks)-1 && len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks) != 0 {
					m.app.currentTaskIndex++

				}

				// If on a project
			} else {
				// Check if there is space to go down and the project is not empty
				if m.app.currentTaskIndex < len(m.app.Projects[m.app.currentProjectIndex].Tasks)-1 && len(m.app.Projects[m.app.currentProjectIndex].Tasks) != 0 {
					m.app.currentTaskIndex++
				}
			}
			return m, nil

		case "shift+left":

			// If on a category
			if m.app.currentCategoryIndex != -1 {
				// Check if there is space to go up and the category is not empty
				if m.app.currentTaskIndex > 0 && len(m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks) != 0 {
					m.app.currentTaskIndex--
				}
				// If on a project
			} else {
				// Check if there is there is space to scroll up and the project is not empty
				if m.app.currentTaskIndex > 0 && len(m.app.Projects[m.app.currentProjectIndex].Tasks) != 0 {
					m.app.currentTaskIndex--
				}
			}
			return m, nil

		case "enter":
			m.app.Save(nil)
			// Quit if user requests
			if strings.TrimSpace(m.input.Value()) == "q" || strings.TrimSpace(m.input.Value()) == "quit" {
				m.app.Save(nil)
				m.app.QuitProgram(nil)
				return m, tea.Quit
			}
			// 1. Run your command
			err := m.app.handleCommands(m.input.Value())
			if err != nil {
				m.error = err.Error()
			} else {
				m.error = "" // Clear it on success
			}

			// Make new command index for next command
			m.commands[m.commandIndex] = m.input.Value()
			m.commands = append(m.commands, "")
			m.commandIndex += 1
			m.commandScrollIndex = m.commandIndex
			m.app.Save(nil)

			// 3. Clear input
			m.input.Reset()

			// Update all project and categories rankings
			for _, project := range m.app.Projects {
				project.syncTasks()
			}
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Calculate sizes
	sidebarWidth := int(float64(m.width) * 0.3)
	mainWidth := m.width - sidebarWidth - 4 // -4 for borders/padding

	// If there are no projects, or the index is invalid, reset or show empty state
	if len(m.app.Projects) == 0 {
		m.app.currentProjectIndex = 0
		m.app.currentCategoryIndex = -1
		m.app.currentTaskIndex = -1
		return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("\n  No projects found. Use 'mk;project_name' to create one.") + "\n\n" + m.input.View()
	}

	// Ensure project index hasn't drifted out of bounds
	if m.app.currentProjectIndex >= len(m.app.Projects) {
		m.app.currentProjectIndex = len(m.app.Projects) - 1
	}

	// Inside View function, where you start rendering taskView
	if len(m.app.Projects) > 0 {
		p := m.app.Projects[m.app.currentProjectIndex]

		// Safety: Clamp index before looping to prevent "index out of range"
		var currentListLen int
		if m.app.currentCategoryIndex != -1 {
			currentListLen = len(p.Categories[m.app.currentCategoryIndex].Tasks)
		} else {
			currentListLen = len(p.Tasks)
		}

		if m.app.currentTaskIndex >= currentListLen {
			m.app.currentTaskIndex = currentListLen - 1
		}

		// ... rest of your for index, task := range ...
	}

	// 1. Render Sidebar (The Tree)
	sidebar := ""
	if len(m.app.Projects) > 0 {
		for projectIndex, project := range m.app.Projects {
			if m.app.currentProjectIndex == projectIndex && m.app.currentCategoryIndex == -1 {
				sidebar = sidebar + renderBlock(selectedStyle, fmt.Sprint(project.Name), sidebarWidth) + "\n"
			} else {
				sidebar = sidebar + renderBlock(projStyle, fmt.Sprint(project.Name), sidebarWidth) + "\n"
			}
			if len(m.app.Projects[projectIndex].Categories) > 0 {
				for categoryIndex, category := range m.app.ProjectMap[project.Name].Categories {
					if m.app.currentCategoryIndex == categoryIndex && m.app.currentProjectIndex == projectIndex {
						sidebar = sidebar + renderBlock(selectedStyle, fmt.Sprint("└─"+category.Name), sidebarWidth) + "\n"
					} else {
						sidebar = sidebar + renderBlock(catStyle, fmt.Sprint("└─"+category.Name), sidebarWidth) + "\n"
					}
				}
			}
		}
	}
	// Iterate through your app.Projects and app.Categories
	// Use m.cursor to highlight the selected one with a ">"

	// 2. Render Tasks (The Main View)
	taskView := ""

	// Check projects exist
	if len(m.app.Projects) > 0 {
		// If on category
		if m.app.currentCategoryIndex != -1 {

			// Loop through every task
			for index, task := range m.app.Projects[m.app.currentProjectIndex].Categories[m.app.currentCategoryIndex].Tasks {
				dueDays := time.Until(task.Deadline).Hours() / 24

				// If task is not selected
				if index != m.app.currentTaskIndex {

					if task.Repeating == false {
						// If task is due in less than a day
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours", task.Name, task.Priority, deadlineString, int(dueHours)), mainWidth) + "\n" + "\n"

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days", task.Name, task.Priority, deadlineString, int(dueDays)), mainWidth) + "\n" + "\n"
						}

					} else {
						var incrementString string
						if task.IncrementYears > 0 {
							amount := strconv.Itoa(int(task.IncrementYears))
							incrementString = incrementString + amount + " Years "
						}
						if task.IncrementMonths > 0 {
							amount := strconv.Itoa(int(task.IncrementMonths))
							incrementString = incrementString + amount + " Months "
						}
						if task.IncrementWeeks > 0 {
							amount := strconv.Itoa(int(task.IncrementWeeks))
							incrementString = incrementString + amount + " Weeks "
						}
						if task.IncrementDays > 0 {
							amount := strconv.Itoa(int(task.IncrementDays))
							incrementString = incrementString + amount + " Days "
						}
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueHours), incrementString), mainWidth) + "\n" + "\n"

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + taskStyle.Render(fmt.Sprintf("└─ %s    P%d    %s    %v days, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueDays), incrementString)) + "\n" + "\n"
						}

					}

					// If task is selected
				} else {
					if task.Repeating == false {
						// If task is due in less than a day
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours", task.Name, task.Priority, deadlineString, int(dueHours)), mainWidth) + "\n" + "\n"

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days", task.Name, task.Priority, deadlineString, int(dueDays)), mainWidth) + "\n" + "\n"
						}

					} else {
						var incrementString string
						if task.IncrementYears > 0 {
							amount := strconv.Itoa(int(task.IncrementYears))
							incrementString = incrementString + amount + " Years "
						}
						if task.IncrementMonths > 0 {
							amount := strconv.Itoa(int(task.IncrementMonths))
							incrementString = incrementString + amount + " Months "
						}
						if task.IncrementWeeks > 0 {
							amount := strconv.Itoa(int(task.IncrementWeeks))
							incrementString = incrementString + amount + " Weeks "
						}
						if task.IncrementDays > 0 {
							amount := strconv.Itoa(int(task.IncrementDays))
							incrementString = incrementString + amount + " Days "
						}
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueHours), incrementString), mainWidth) + "\n" + "\n"
							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueDays), incrementString), mainWidth) + "\n" + "\n"

						}

					}
				}
			}

			// If on project
		} else {
			// Loop through all tasks
			for index, task := range m.app.Projects[m.app.currentProjectIndex].Tasks {
				dueDays := time.Until(task.Deadline).Hours() / 24

				// If task is not selected
				if index != m.app.currentTaskIndex {

					if task.Repeating == false {
						// If task is due in less than a day
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours", task.Name, task.Priority, deadlineString, int(dueHours)), mainWidth) + "\n" + "\n"

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days", task.Name, task.Priority, deadlineString, int(dueDays)), mainWidth) + "\n" + "\n"
						}

					} else {
						var incrementString string
						if task.IncrementYears > 0 {
							amount := strconv.Itoa(int(task.IncrementYears))
							incrementString = incrementString + amount + " Years "
						}
						if task.IncrementMonths > 0 {
							amount := strconv.Itoa(int(task.IncrementMonths))
							incrementString = incrementString + amount + " Months "
						}
						if task.IncrementWeeks > 0 {
							amount := strconv.Itoa(int(task.IncrementWeeks))
							incrementString = incrementString + amount + " Weeks "
						}
						if task.IncrementDays > 0 {
							amount := strconv.Itoa(int(task.IncrementDays))
							incrementString = incrementString + amount + " Days "
						}
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueHours), incrementString), mainWidth) + "\n" + "\n"

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueDays), incrementString), mainWidth) + "\n" + "\n"
						}

					}

					// If task is selected
				} else {
					if task.Repeating == false {
						// If task is due in less than a day
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(taskStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours", task.Name, task.Priority, deadlineString, int(dueHours)), mainWidth) + "\n" + "\n"

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days", task.Name, task.Priority, deadlineString, int(dueDays)), mainWidth) + "\n" + "\n"
						}

					} else {
						var incrementString string
						if task.IncrementYears > 0 {
							amount := strconv.Itoa(int(task.IncrementYears))
							incrementString = incrementString + amount + " Years "
						}
						if task.IncrementMonths > 0 {
							amount := strconv.Itoa(int(task.IncrementMonths))
							incrementString = incrementString + amount + " Months "
						}
						if task.IncrementWeeks > 0 {
							amount := strconv.Itoa(int(task.IncrementWeeks))
							incrementString = incrementString + amount + " Weeks "
						}
						if task.IncrementDays > 0 {
							amount := strconv.Itoa(int(task.IncrementDays))
							incrementString = incrementString + amount + " Days "
						}
						if math.Max(1, dueDays) == 1 {
							dueHours := int(time.Until(task.Deadline).Hours())
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v hours, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueHours), incrementString), mainWidth)

							// If task is due in more than a day
						} else {
							deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
							taskView = taskView + renderBlock(selectedStyle, fmt.Sprintf("└─ %s    P%d    %s    %v days, Repeats Every: %s", task.Name, task.Priority, deadlineString, int(dueDays), incrementString), mainWidth) + "\n" + "\n"
						}

					}
				}
			}
		}
	}
	// Logic: If selection is a Category, show category tasks.
	// If selection is a Project, show ALL tasks in that project.

	// 3. Styles
	sidebarStyle := lipgloss.NewStyle().Width(sidebarWidth).Border(lipgloss.RoundedBorder())
	mainStyle := lipgloss.NewStyle().Width(mainWidth).Border(lipgloss.RoundedBorder())

	// 4. Combine
	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		sidebarStyle.Render(sidebar),
		mainStyle.Render(taskView),
	)

	// Errors

	var errorBar string
	if m.error != "" {
		errorStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("9")).  // Red background
			Foreground(lipgloss.Color("15")). // White text
			Width(m.width).
			Padding(0, 1)

		errorBar = errorStyle.Render("ERROR: " + m.error)
	}

	return panels + "\n" + errorBar + "\n" + m.input.View()
}

// Helper function to render a "tight" block
func renderBlock(style lipgloss.Style, text string, maxWidth int) string {
	// 1. Force wrap the text to the boundary
	w := lipgloss.NewStyle().Width(maxWidth).Render(text)

	// 2. Find the widest line in that wrap
	actualWidth := lipgloss.Width(w)

	// 3. Render with the background color filling only that width
	return style.Width(actualWidth).Render(text)
}
