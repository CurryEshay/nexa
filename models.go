package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Project struct {
	Name        string               `json:"name"`
	Categories  []*Category          `json:"categories"`
	CategoryMap map[string]*Category `json:"-"`
	Tasks       []*Task              `json:"tasks"`
}

func (a *App) NewProject(name string) error {
	a.Mu.Lock()

	// Create new project
	project := &Project{Name: name, Categories: []*Category{}, CategoryMap: map[string]*Category{}}

	// Update app struct
	a.Projects = append(a.Projects, project)
	a.ProjectMap[name] = project

	// // Add "All" category
	// category := &Category{Name: "All", ProjectName: name, Tasks: []*Task{}, TaskMap: map[string]*Task{}}

	// // Update relevant projects slice and map
	// (*project).Categories = append((*project).Categories, category)
	// (*project).CategoryMap["All"] = category

	a.Mu.Unlock()
	return nil
}

func (a *App) removeProject(name string) error {

	// Find project index
	projectIndex, err := a.findProjectIndex(name)
	if err != nil {
		return err
	}

	delete((*a).ProjectMap, name)

	// Slice if more than 1 element, else make empty slice
	if len(a.Projects) > 1 {
		a.Projects = append(a.Projects[:projectIndex], a.Projects[projectIndex+1:]...)
	} else {
		a.Projects = []*Project{}
	}

	if len(a.Projects) == 0 {
		a.currentProjectIndex = 0
		a.currentCategoryIndex = -1
		a.currentTaskIndex = -1
	} else {
		// If we were looking at the deleted project, or our index is now too high
		if a.currentProjectIndex >= len(a.Projects) {
			a.currentProjectIndex = len(a.Projects) - 1
		}
		// Force reset sub-selections to prevent orphans
		a.currentCategoryIndex = -1
		a.currentTaskIndex = -1
	}
	return nil

}

func (a *App) newProject(name string) error {
	// Create new project
	project := &Project{Name: name, Categories: []*Category{}, CategoryMap: map[string]*Category{}}

	// Update app struct
	(*a).Projects = append((*a).Projects, project)
	(*a).ProjectMap[name] = project

	// // Add "All" category
	// category := &Category{Name: "All", ProjectName: name, Tasks: []*Task{}, TaskMap: map[string]*Task{}}

	// // Update relevant projects slice and map
	// (*project).Categories = append((*project).Categories, category)
	// (*project).CategoryMap["All"] = category

	return nil
}

func (p *Project) printProjects() {
	fmt.Println("--- Project: " + p.Name + " ----" + "\n")
	for _, category := range p.Categories {
		category.printCategories()
	}
}

func (p *Project) syncTasks() {
	// Reset task lists
	p.Tasks = []*Task{}

	// Add all tasks
	for _, category := range p.Categories {
		for _, task := range category.Tasks {
			p.Tasks = append(p.Tasks, task)
		}
	}

	p.SortTasksByUrgency()
}

func (p *Project) SortTasksByUrgency() {
	slices.SortFunc(p.Tasks, func(a, b *Task) int {
		scoreA, _ := ScoreTask(a)
		scoreB, _ := ScoreTask(b)

		// If a has a higher score, it should come BEFORE b.
		// In Go SortFunc:
		// Return -1 if a < b (a comes first)
		// Return 1  if a > b (b comes first)
		// Return 0  if equal

		if scoreA > scoreB {
			return -1 // a is "more urgent", move it to the start
		}
		if scoreA < scoreB {
			return 1 // b is "more urgent", move it to the start
		}
		return 0
	})
}

type Category struct {
	Name        string           `json:"name"`
	ProjectName string           `json:"project_name"`
	Tasks       []*Task          `json:"tasks"`
	TaskMap     map[string]*Task `json:"-"`
}

func (c *Category) SortTasksByUrgency() {
	slices.SortFunc(c.Tasks, func(a, b *Task) int {
		scoreA, _ := ScoreTask(a)
		scoreB, _ := ScoreTask(b)

		// If a has a higher score, it should come BEFORE b.
		// In Go SortFunc:
		// Return -1 if a < b (a comes first)
		// Return 1  if a > b (b comes first)
		// Return 0  if equal

		if scoreA > scoreB {
			return -1 // a is "more urgent", move it to the start
		}
		if scoreA < scoreB {
			return 1 // b is "more urgent", move it to the start
		}
		return 0
	})
}

func (c *Category) printCategories() {
	c.SortTasksByUrgency()
	fmt.Printf("--- Category: %s --- \n", c.Name)
	fmt.Println("")
	for _, task := range c.Tasks {

		dueDays := time.Until(task.Deadline).Hours() / 24

		if math.Max(1, dueDays) == 1 {
			dueHours := int(time.Until(task.Deadline).Hours())
			deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
			fmt.Printf("- %s    P%d    %s    %v hours \n", strings.ReplaceAll(task.Name, "_", " "), task.Priority, deadlineString, int(dueHours))
		} else {
			deadlineString := time.Time.Format(task.Deadline, "Monday, Jan 02 3:04 PM")
			fmt.Printf("- %s    P%d    %s    %v days \n", task.Name, task.Priority, deadlineString, int(dueDays))
		}
		fmt.Println("")

	}
}

func (a *App) NewCategory(projectName string, categoryName string) error {

	a.Mu.Lock()
	defer a.Mu.Unlock()
	if exists, err := a.projectExists(projectName); exists == false || err != nil {
		return errors.New("Project name does not exist")
	}
	// Create blank category
	category := &Category{Name: categoryName, ProjectName: projectName, Tasks: []*Task{}, TaskMap: map[string]*Task{}}

	// Find relevant project
	project := a.ProjectMap[projectName]

	// Update relevant projects slice and map
	(*project).Categories = append((*project).Categories, category)
	(*project).CategoryMap[categoryName] = category

	return nil
}

func (a *App) newCategory(projectName string, categoryName string) error {
	if exists, err := a.projectExists(projectName); exists == false || err != nil {
		return errors.New("Project name does not exist")
	}
	// Create blank category
	category := &Category{Name: categoryName, ProjectName: projectName, Tasks: []*Task{}, TaskMap: map[string]*Task{}}

	// Find relevant project
	project := a.ProjectMap[projectName]

	// Update relevant projects slice and map
	(*project).Categories = append((*project).Categories, category)
	(*project).CategoryMap[categoryName] = category

	return nil
}

func (a *App) removeCategory(projectName string, categoryName string) error {

	// Find project index
	projectIndex, categoryIndex, err := a.findCategoryIndex(projectName, categoryName)
	if err != nil {
		return err
	}

	delete((*a).ProjectMap[projectName].CategoryMap, categoryName)

	// Slice if more than 1 element, else make empty slice
	if len(a.Projects[projectIndex].Categories) > 1 {
		a.Projects[projectIndex].Categories = append(a.Projects[projectIndex].Categories[:categoryIndex], a.Projects[projectIndex].Categories[categoryIndex+1:]...)
	} else {
		a.Projects[projectIndex].Categories = []*Category{}
	}

	a.Projects[projectIndex].syncTasks()

	// If the deleted category was the one currently selected
	if len(a.Projects[projectIndex].Categories) == 0 {
		a.currentCategoryIndex = -1 // Drop back to "All Tasks" view
		a.currentTaskIndex = -1     // Reset task focus
	} else if a.currentCategoryIndex >= len(a.Projects[projectIndex].Categories) {
		// If the selection is now out of bounds, snap it to the last valid category
		a.currentCategoryIndex = len(a.Projects[projectIndex].Categories) - 1
	}

	return nil

}

type Task struct {
	Name            string    `json:"name"`
	Deadline        time.Time `json:"deadline"`
	ProjectName     string    `json:"project_name"`
	CategoryName    string    `json:"category_name"`
	Priority        uint64    `json:"priority"`
	Repeating       bool      `json:"repeating"`
	IncrementDays   uint      `json:"increment_days"`
	IncrementWeeks  uint      `json:"increment_weeks"`
	IncrementMonths uint      `json:"increment_months"`
	IncrementYears  uint      `json:"increment_years"`
}

func ScoreTask(task *Task) (float64, error) {
	now := time.Now()
	dueHours := task.Deadline.Sub(now).Hours()
	if dueHours > 0 {
		timeScore := 2 / (float64(dueHours + 1))
		priorityScore := float64(task.Priority) * 0.1
		score := timeScore + priorityScore
		return score, nil
	} else {
		timeScore := float64(2)
		priorityScore := float64(task.Priority) * 0.1
		overdueScore := math.Min(2.0, math.Abs(dueHours)*0.1)
		score := timeScore + priorityScore + overdueScore
		return score, nil
	}

}

func (t *Task) incrementTask() error {
	if t.Repeating == false {
		return errors.New("Can only increment repeating tasks")
	}

	t.Deadline = t.Deadline.AddDate(int(t.IncrementYears), int(t.IncrementMonths), int(t.IncrementWeeks)*7+int(t.IncrementDays))
	return nil
}

func (a *App) NewTask(projectName string, categoryName string, taskName string, priority uint, deadline time.Time) error {

	a.Mu.Lock()
	defer a.Mu.Unlock()
	if exists, err := a.categoryExists(projectName, categoryName); exists == false || err != nil {
		return errors.New("Category name does not exist")
	}
	// Create blank task
	task := &Task{Name: taskName, ProjectName: projectName, CategoryName: categoryName, Priority: uint64(priority), Deadline: deadline, Repeating: false}

	// Find relevent category
	category := a.ProjectMap[projectName].CategoryMap[categoryName]

	// Update relevant category's slice and map
	(*category).Tasks = append((*category).Tasks, task)
	(*category).TaskMap[taskName] = task

	// Find relevant project and sync
	project := a.ProjectMap[projectName]

	project.syncTasks()
	return nil
}

func (a *App) newTask(projectName string, categoryName string, taskName string, priority uint, deadline time.Time) error {
	if exists, err := a.categoryExists(projectName, categoryName); exists == false || err != nil {
		return errors.New("Category name does not exist")
	}

	if valid, err := a.isValidTask(projectName, categoryName, taskName); valid == false || err != nil {
		if err != nil {
			return err
		}
		return errors.New("Not a valid task name to create")
	}
	// Create blank task
	task := &Task{Name: taskName, ProjectName: projectName, CategoryName: categoryName, Priority: uint64(priority), Deadline: deadline}

	// Find relevent category
	category := a.ProjectMap[projectName].CategoryMap[categoryName]

	// Update relevant category's slice and map
	(*category).Tasks = append((*category).Tasks, task)
	(*category).TaskMap[taskName] = task

	// Find relevant project and sync
	project := a.ProjectMap[projectName]

	project.syncTasks()

	projectIndex, categoryIndex, err := a.findCategoryIndex(projectName, categoryName)
	if err != nil {
		return err
	}

	if a.currentCategoryIndex > 0 {
		if projectIndex == uint(a.currentProjectIndex) && categoryIndex == uint(a.currentCategoryIndex) {
			a.currentTaskIndex = 0
		}
	} else {
		if projectIndex == uint(a.currentProjectIndex) && a.currentCategoryIndex == -1 {
			a.currentTaskIndex = 0
		}
	}

	return nil
}

func (a *App) removeTask(projectName string, categoryName string, taskName string) error {

	// Find project index
	projectIndex, categoryIndex, taskIndex, err := a.findTaskIndex(projectName, categoryName, taskName)
	if err != nil {
		return err
	}
	task := a.Projects[projectIndex].Categories[categoryIndex].Tasks[taskIndex]
	a.deletedTasks = append(a.deletedTasks, task)
	delete((*a).ProjectMap[projectName].CategoryMap[categoryName].TaskMap, taskName)

	if len(a.Projects[projectIndex].Categories[categoryIndex].Tasks) > 1 {
		a.Projects[projectIndex].Categories[categoryIndex].Tasks = append(a.Projects[projectIndex].Categories[categoryIndex].Tasks[:taskIndex], a.Projects[projectIndex].Categories[categoryIndex].Tasks[taskIndex+1:]...)
	} else {
		a.Projects[projectIndex].Categories[categoryIndex].Tasks = []*Task{}
	}
	// slice if there is more than 1 task

	a.Projects[projectIndex].syncTasks()

	if a.currentCategoryIndex != -1 {
		if a.currentTaskIndex > len(a.Projects[projectIndex].Categories[categoryIndex].Tasks)-1 {
			a.currentTaskIndex = len(a.Projects[projectIndex].Categories[categoryIndex].Tasks) - 1
		}
	} else {
		if a.currentTaskIndex > len(a.Projects[projectIndex].Tasks)-1 {
			a.currentTaskIndex = len(a.Projects[projectIndex].Tasks) - 1
		}
	}

	return nil

}
func (a *App) restoreTask() error {
	if len(a.deletedTasks) == 0 {
		return errors.New("No tasks were deleted")
	}
	lastTaskIndex := len(a.deletedTasks) - 1
	task := a.deletedTasks[lastTaskIndex]

	if len(a.deletedTasks) > 0 {
		a.deletedTasks = append(a.deletedTasks[:lastTaskIndex], a.deletedTasks[lastTaskIndex+1:]...)
	} else {
		a.deletedTasks = []*Task{}
	}

	if valid, _ := a.isValidTask(task.ProjectName, task.CategoryName, task.Name); valid == false {

		return errors.New("Task no longer valid (Project or Category may have been deleted)")
	}

	err := a.newTask(task.ProjectName, task.CategoryName, task.Name, uint(task.Priority), task.Deadline)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) doneTask(projectName string, categoryName string, taskName string) error {
	_, _, task, err := a.findByPath(projectName, categoryName, taskName)
	if err != nil {
		return err
	}

	if task.Repeating == false {
		err := a.removeTask(projectName, categoryName, taskName)
		if err != nil {
			return err
		} else {
			return nil // Clear it on success
		}
		// If on a project
	} else {
		task.incrementTask()
		return nil
	}

}

func (a *App) newRepeatingTask(projectName string, categoryName string, taskName string, priority uint, referenceTime time.Time, incrementString string) error {
	if exists, err := a.categoryExists(projectName, categoryName); exists == false || err != nil {
		return errors.New("Category name does not exist")
	}

	if valid, err := a.isValidTask(projectName, categoryName, taskName); valid == false || err != nil {
		if err != nil {
			return err
		}
		return errors.New("Not a valid task name to create")
	}
	incrementString = strings.TrimSpace(incrementString)

	incrementsString := strings.Split(incrementString, " ")

	if len(incrementsString) != 4 {
		return errors.New("Invalid increment format")
	}

	var increments [4]uint
	for index := range incrementsString {
		increment, err := strconv.Atoi(incrementsString[index])
		if err != nil {
			return errors.New("Increments are not integers")
		}
		increments[index] = uint(increment)
	}
	incrementYears, incrementMonths, incrementWeeks, incrementDays := increments[0], increments[1], increments[2], increments[3]

	if incrementYears == 0 && incrementMonths == 0 && incrementWeeks == 0 && incrementDays == 0 {
		return errors.New("repeating increment cannot be zero")
	}

	for referenceTime.Before(time.Now()) {
		referenceTime = referenceTime.AddDate(int(incrementYears), int(incrementMonths), int(incrementWeeks)*7+int(incrementDays))
	}
	// Create blank task
	task := &Task{Name: taskName, ProjectName: projectName, CategoryName: categoryName, Priority: uint64(priority), Deadline: referenceTime, Repeating: true,
		IncrementYears: incrementYears, IncrementMonths: incrementMonths, IncrementWeeks: incrementWeeks, IncrementDays: incrementDays}
	// Find relevent category
	category := a.ProjectMap[projectName].CategoryMap[categoryName]

	// Update relevant category's slice and map
	(*category).Tasks = append((*category).Tasks, task)
	(*category).TaskMap[taskName] = task

	// Find relevant project and sync
	project := a.ProjectMap[projectName]

	project.syncTasks()

	projectIndex, categoryIndex, err := a.findCategoryIndex(projectName, categoryName)
	if err != nil {
		return err
	}

	if a.currentCategoryIndex > 0 {
		if projectIndex == uint(a.currentProjectIndex) && categoryIndex == uint(a.currentCategoryIndex) {
			a.currentTaskIndex = 0
		}
	} else {
		if projectIndex == uint(a.currentProjectIndex) && a.currentCategoryIndex == -1 {
			a.currentTaskIndex = 0
		}
	}

	return nil
}
