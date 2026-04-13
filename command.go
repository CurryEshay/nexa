package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type CMDFunc func(args []string) error

func (a *App) InitCommands() {
	a.CMDMap = map[string]CMDFunc{
		"mk": a.HandleMake,
		// "c":     a.ClearTerminal, // removed for TUI
		"rm":   a.HandleRemove,
		"lk":   a.LockApp,
		"ulk":  a.UnlockApp,
		"mv":   a.HandleMove,
		"s":    a.Save,
		"save": a.Save,
		"rs":   a.handleReschedule,
		"dn":   a.handleDone,
		"rp":   a.handleMakeRepeatingTask,
		"p":    a.handleReprioritise,
	}
}

func (a *App) handleCommands(input string) error {

	input = strings.TrimSpace(input)
	args, err := a.SplitArgString(input)
	if err != nil {
		return err

	}
	args[0] = strings.ToLower(args[0])
	handler, ok := a.CMDMap[args[0]]
	if !ok {
		return errors.New("Invalid command")
	}
	err = handler([]string(args[1:]))
	if err != nil {
		return err
	}

	return nil
}

func (a *App) UnlockApp(args []string) error {
	a.Locked = false
	fmt.Println("App is now unlocked")
	return nil
}

func (a *App) LockApp(args []string) error {
	a.Locked = true
	fmt.Println("App is now locked")
	return nil
}
func (a *App) SplitArgString(argString string) ([]string, error) {

	// Seperate by ; deliminator
	args := strings.Split(argString, ";")

	// Make sure args slice is not empty and raise error if it is empty
	if len(args) == 0 {
		return nil, errors.New("No argument in string, likely improper deliminator used")
	}

	// Make sure no extra whitespace exists in each argument
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}

	return args, nil
}

func (a *App) Parser() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(">> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		a.ErrorLog(err)
		return
	}
	input = strings.TrimSpace(input)
	args, err := a.SplitArgString(input)
	if err != nil {
		a.ErrorLog(err)
		return

	}
	args[0] = strings.ToLower(args[0])
	handler, ok := a.CMDMap[args[0]]
	if !ok {
		a.CMDMap["h"]([]string{})
		return
	}
	err = handler([]string(args[1:]))
	if err != nil {
		a.ErrorLog(err)
		return
	}

}

// Display help command
func (a *App) ShowHelpMenu(args []string) error {
	fmt.Println(`--- HELP MENU ---
	mk\{proj}\{proj name}  - makes project
	mk\{proj/category}\{category name} - makes category under project
	mk\{proj/category/task}\{priority}\{date}\{time - default 11:59PM} - makes task under category
	rm\{path/to/object} - removes an object (project, category or task, be careful)
	dn\{path/to/task} - sets task to done
	mv\{old/path/to/object}\{new/path/to/object} - move object to other spot or rename
	rd\{path/to/task}\{name}\{priority}\{date}\{time - default 11:59PM} - re-define task, use * to keep the same
	q - quit program`)
	return nil

}

// Quit program
func (a *App) QuitProgram(args []string) error {
	a.Running = false
	return nil
}

// Handle make command
func (a *App) HandleMake(args []string) error {

	// Get path from first argument
	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	a.Mu.Lock()
	defer a.Mu.Unlock()

	// Check if only project was provided and given project does not yet exist
	if valid, err := a.isValidProject(projectName); valid == true && err == nil && categoryName == "" {
		err := a.handleMakeProject(projectName)
		return err
	}

	// Get the references for the path

	// if only project exists, create category

	if valid, err := a.isValidCategory(projectName, categoryName); valid == true && err == nil && taskName == "" {
		err := a.handleMakeCategory(projectName, categoryName)
		return err

	}

	// If project and category exist but not task, create task
	if valid, err := a.isValidTask(projectName, categoryName, taskName); valid == true && err == nil {
		err = a.handleMakeTask(projectName, categoryName, taskName, args[1:])
		return err
	}
	return errors.New("Path does not exist")
}

func (a *App) handleMake(args []string) error {

	// Get path from first argument
	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	// Check if only project was provided and given project does not yet exist
	if valid, err := a.isValidProject(projectName); valid == true && err == nil && categoryName == "" {
		err := a.handleMakeProject(projectName)
		return err
	}

	// Get the references for the path

	// if only project exists, create category

	if valid, err := a.isValidCategory(projectName, categoryName); valid == true && err == nil && taskName == "" {
		err := a.handleMakeCategory(projectName, categoryName)
		return err

	}

	// If project and category exist but not task, create task
	if valid, err := a.isValidTask(projectName, categoryName, taskName); valid == true && err == nil {
		err = a.handleMakeTask(projectName, categoryName, taskName, args[1:])
		return err
	}
	return errors.New("Path does not exist")
}

func (a *App) handleMakeProject(projectName string) error {
	err := a.newProject(projectName)
	return err
}

func (a *App) handleMakeCategory(projectName string, categoryName string) error {
	err := a.newCategory(projectName, categoryName)
	return err
}

func (a *App) handleMakeTask(projectName string, categoryName string, taskName string, args []string) error {
	if len(args) < 2 {
		return errors.New("Invalid argument format, please refer to help guide")
	}

	priority, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	if len(args) == 2 {
		deadline, err := parseFlexibleDeadline(args[1])
		if err != nil {
			return err
		}
		err = a.newTask(projectName, categoryName, taskName, uint(priority), deadline)
		return err
	}
	if len(args) == 3 {
		deadlineString := fmt.Sprint(args[1] + " " + args[2])
		deadline, err := parseFlexibleDeadline(deadlineString)
		if err != nil {
			return err
		}
		err = a.newTask(projectName, categoryName, taskName, uint(priority), deadline)
		return err
	}

	return errors.New("Path already exists")

}
func (a *App) ClearTerminal(args []string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func (a *App) HandleRemove(args []string) error {
	if a.Locked == true {
		return errors.New("Error, project is locked, type ulk to unlock")
	}
	if len(args) != 1 {
		return errors.New("Invalid argument format, please refer to help guide")
	}

	a.Mu.Lock()
	defer a.Mu.Unlock()
	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	if projectName != "" && categoryName == "" && taskName == "" {
		err := a.removeProject(projectName)
		if err != nil {
			return err
		}
		return nil
	}
	if projectName != "" && categoryName != "" && taskName == "" {
		err := a.removeCategory(projectName, categoryName)
		if err != nil {
			return err
		}
		return nil
	}

	if projectName != "" && categoryName != "" && taskName != "" {
		err := a.removeTask(projectName, categoryName, taskName)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Invalid path in remove function")
}

func (a *App) handleRemove(args []string) error {
	if a.Locked == true {
		return errors.New("Error, project is locked, type ulk to unlock")
	}
	if len(args) != 1 {
		return errors.New("Invalid argument format, please refer to help guide")
	}
	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	if projectName != "" && categoryName == "" && taskName == "" {
		err := a.removeProject(projectName)
		if err != nil {
			return err
		}
		return nil
	}
	if projectName != "" && categoryName != "" && taskName == "" {
		err := a.removeCategory(projectName, categoryName)
		if err != nil {
			return err
		}
		return nil
	}

	if projectName != "" && categoryName != "" && taskName != "" {
		err := a.removeTask(projectName, categoryName, taskName)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Invalid path in remove function")
}

func (a *App) handleRemoveProject(projectName string) error {
	err := a.removeProject(projectName)
	return err
}

func (a *App) handleRemoveCategory(projectName string, categoryName string) error {
	err := a.removeCategory(projectName, categoryName)
	return err
}

func (a *App) handleRemoveTask(projectName string, categoryName string, taskName string) error {
	err := a.removeTask(projectName, categoryName, taskName)
	return err
}

func (a *App) HandleDisplay(Args []string) error {
	for _, project := range a.Projects {
		project.printProjects()
	}
	return nil
}

func (a *App) HandleMove(args []string) error {

	if a.Locked == true {
		return errors.New("Error, project is locked, type ulk to unlock")
	}

	if len(args) != 2 {
		return errors.New("Please enter 2 arguments for move function")
	}
	// Find old project / category / task
	oldProjectName, oldCategoryName, oldTaskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	newProjectName, newCategoryName, newTaskName, err := a.HandlePath(args[1])
	if err != nil {
		return err
	}

	a.Mu.Lock()
	defer a.Mu.Unlock()

	// Renaming project
	if oldProjectName != "" && oldCategoryName == "" && oldTaskName == "" && newProjectName != "" && newCategoryName == "" && newTaskName == "" {
		// Get old project path exists
		oldProjectExists, err := a.projectExists(oldProjectName)
		if err != nil {
			return err
		}

		// Get new project path is available
		newProjectValid, err := a.isValidProject(newProjectName)
		if err != nil {
			return err
		}

		// Verify old path exists and new path avaliable
		if oldProjectExists == true && newProjectValid == true {

			// Save reference to old project
			project := a.ProjectMap[oldProjectName]

			// Update map to reflect name change
			a.removeProject(oldProjectName)

			// Modify name
			project.Name = newProjectName

			// Update map and slice
			a.ProjectMap[newProjectName] = project
			a.Projects = append(a.Projects, project)

			// Update categories
			for _, category := range project.Categories {
				category.ProjectName = newProjectName
				for _, task := range category.Tasks {
					task.ProjectName = newProjectName
				}
			}
			return nil
		}
		return errors.New("Unexpected error in handle move project function")

	}
	// Moving category
	if oldProjectName != "" && oldCategoryName != "" && oldTaskName == "" && newProjectName != "" && newCategoryName != "" && newTaskName == "" {
		// Get old category path exists
		oldCategoryExists, err := a.categoryExists(oldProjectName, oldCategoryName)
		if err != nil {
			return err
		}

		// Get new category path is avaliable
		newCategoryValid, err := a.isValidCategory(newProjectName, newCategoryName)
		if err != nil {
			return err
		}

		// Verify old path exists and new path avaliable
		if oldCategoryExists == true && newCategoryValid == true {

			// Save reference to old category
			category := a.ProjectMap[oldProjectName].CategoryMap[oldCategoryName]

			// Remove category from old category map
			a.removeCategory(oldProjectName, oldCategoryName)

			// Modify name
			category.Name = newCategoryName

			// Move category and update slice and map
			a.ProjectMap[newProjectName].Categories = append(a.ProjectMap[newProjectName].Categories, category)
			a.ProjectMap[newProjectName].CategoryMap[newCategoryName] = category

			// Update tasks with new project
			for _, task := range category.Tasks {
				task.ProjectName = newProjectName
			}

			// Update new project and old project tasks
			a.ProjectMap[oldProjectName].syncTasks()
			a.ProjectMap[newProjectName].syncTasks()

			return nil
		}
		return errors.New("Unexpected error in handle move category function")

	}

	if oldProjectName != "" && oldCategoryName != "" && oldTaskName != "" && newProjectName != "" && newCategoryName != "" && newTaskName != "" {
		// Get old task path exists
		oldTaskExists, err := a.taskExists(oldProjectName, oldCategoryName, oldTaskName)
		if err != nil {
			return err
		}

		// Get new category path is avaliable
		newTaskValid, err := a.isValidTask(newProjectName, newCategoryName, newTaskName)
		if err != nil {
			return err
		}

		// Verify old path exists and new path avaliable
		if oldTaskExists == true && newTaskValid == true {

			// Save reference to old task
			task := a.ProjectMap[oldProjectName].CategoryMap[oldCategoryName].TaskMap[oldTaskName]

			// Remove task
			a.removeTask(oldProjectName, oldCategoryName, oldTaskName)

			// Modify name
			task.Name = newTaskName

			// Move task and update slice and map
			a.ProjectMap[newProjectName].CategoryMap[newCategoryName].Tasks = append(a.ProjectMap[newProjectName].CategoryMap[newCategoryName].Tasks, task)
			a.ProjectMap[newProjectName].CategoryMap[newCategoryName].TaskMap[newTaskName] = task

			// Update project name of task
			task.ProjectName = newProjectName

			// Update new project and old project tasks
			a.ProjectMap[oldProjectName].syncTasks()
			a.ProjectMap[newProjectName].syncTasks()

			return nil
		}
		return errors.New("Unexpected error in handle move task function")

	}

	return errors.New("Unexpected error in handle move function")
}

func (a *App) handleReschedule(args []string) error {

	if len(args) != 2 && len(args) != 3 {
		return errors.New("Invalid number of args")
	}
	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	if exists, err := a.taskExists(projectName, categoryName, taskName); exists == false {
		if err != nil {
			return err
		}
		return errors.New("Invalid task to reschedule, does not exist")
	}

	_, _, task, err := a.findByPath(projectName, categoryName, taskName)
	if err != nil {
		return err
	}

	if task.Repeating == false {
		if len(args) == 2 {
			deadline, err := parseFlexibleDeadline(args[1])
			if err != nil {
				return err
			}

			_, _, task, err := a.FindByPath(projectName, categoryName, taskName)
			if err != nil {
				return err
			}

			task.Deadline = deadline
			// err = a.removeTask(projectName, categoryName, taskName)
			// if err != nil {
			// 	return err
			// }
			// err = a.newTask(projectName, categoryName, taskName, uint(task.Priority), deadline)
			// if err != nil {
			// 	return err
			// }
		}
		if len(args) == 3 {
			deadline, err := parseFlexibleDeadline(args[1] + " " + args[2])
			if err != nil {
				return err
			}

			_, _, task, err := a.FindByPath(projectName, categoryName, taskName)
			if err != nil {
				return err
			}
			task.Deadline = deadline
			// err = a.removeTask(projectName, categoryName, taskName)
			// if err != nil {
			// 	return err
			// }
			// err = a.newTask(projectName, categoryName, taskName, uint(task.Priority), deadline)
			// if err != nil {
			// 	return err
			// }
		}

		return nil
	} else {
		if len(args) != 3 && len(args) != 4 {
			return errors.New("Invalid number of args")
		}
		if len(args) == 3 {
			referenceTime, err := parseFlexibleDeadline(args[2])
			if err != nil {
				return err
			}

			err = a.removeTask(projectName, categoryName, taskName)
			if err != nil {
				return err
			}
			err = a.newRepeatingTask(projectName, categoryName, taskName, uint(task.Priority), referenceTime, args[1])
			if err != nil {
				return err
			}

			return nil
		} else if len(args) == 4 {
			referenceTime, err := parseFlexibleDeadline(args[2] + " " + args[3])
			if err != nil {
				return err
			}
			err = a.removeTask(projectName, categoryName, taskName)
			if err != nil {
				return err
			}
			err = a.newRepeatingTask(projectName, categoryName, taskName, uint(task.Priority), referenceTime, args[1])
			if err != nil {
				return err
			}

			return nil
		}

	}
	return errors.New("Unexpected error in handleReschedule function")

}

func (a *App) handleDone(args []string) error {
	if len(args) != 1 {
		return errors.New("Inappropriate number of arguments in done function")
	}
	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}
	err = a.doneTask(projectName, categoryName, taskName)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) handleMakeRepeatingTask(args []string) error {
	projectName, categoryName, taskName, err := a.HandlePath(args[0])

	if err != nil {
		return err
	}

	if projectName == "" || categoryName == "" {
		return errors.New("Not a path to a task")
	}
	if len(args) != 4 && len(args) != 5 {
		return errors.New("Invalid argument format to make repeating task")
	}

	priority, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return err
	}

	if len(args) == 4 {
		deadline, err := parseFlexibleDeadline(args[3])
		if err != nil {
			return err
		}
		err = a.newRepeatingTask(projectName, categoryName, taskName, uint(priority), deadline, args[2])
		return err
	}
	if len(args) == 5 {
		deadline, err := parseFlexibleDeadline(args[3] + " " + args[4])
		if err != nil {
			return err
		}
		err = a.newRepeatingTask(projectName, categoryName, taskName, uint(priority), deadline, args[2])
		return err
	}

	return errors.New("Path already exists")

}

func (a *App) handleReprioritise(args []string) error {
	if len(args) != 2 {
		return errors.New("Innappropriate number of argumenets in reprioritise function")
	}

	projectName, categoryName, taskName, err := a.HandlePath(args[0])
	if err != nil {
		return err
	}

	priority, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	project, _, task, err := a.findByPath(projectName, categoryName, taskName)
	if err != nil {
		return err
	}

	task.Priority = uint64(priority)
	project.syncTasks()

	return nil
}
