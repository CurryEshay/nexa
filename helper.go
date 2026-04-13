package main

import (
	"errors"
)

// NOTE: Methods with lower case are UNPROTECTED and only for internal use

func (a *App) FindByPath(projectName string, categoryName string, taskName string) (*Project, *Category, *Task, error) {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	// 1. Jump to Project
	project, ok := a.ProjectMap[projectName]
	if !ok {
		return nil, nil, nil, errors.New("Project not found")
	}

	// If only project name was requested
	if categoryName == "" && taskName == "" {
		return project, nil, nil, nil
	}

	// 2. Jump to Category
	category, ok := project.CategoryMap[categoryName]
	if !ok {
		return nil, nil, nil, errors.New("Category not found")
	}

	// If only project name and category was requested
	if taskName == "" {
		return project, category, nil, nil
	}

	// 3. Jump to Task
	task, ok := category.TaskMap[taskName]
	if !ok {
		return nil, nil, nil, errors.New("Task not found")
	}

	return project, category, task, nil
}

func (a *App) ProjectExists(projectName string) bool {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	// We pass empty strings for category and task since we only care about the project
	project, _, _, err := a.FindByPath(projectName, "", "")

	// If there is no error and the pointer isn't nil, the project exists
	return err == nil && project != nil
}

func (a *App) CategoryExists(projectName string, categoryName string) bool {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	// We pass an empty string for the taskName
	_, category, _, err := a.FindByPath(projectName, categoryName, "")

	// Category exists only if the search returned a valid pointer without errors
	return err == nil && category != nil
}

func (a *App) TaskExists(projectName string, categoryName string, taskName string) bool {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	_, _, task, err := a.FindByPath(projectName, categoryName, taskName)

	// Task exists only if the final leaf in the tree was found
	return err == nil && task != nil
}

func (a *App) findByPath(projectName string, categoryName string, taskName string) (*Project, *Category, *Task, error) {

	// 1. Jump to Project
	project, ok := a.ProjectMap[projectName]
	if !ok {
		return nil, nil, nil, errors.New("Project not found")
	}

	// If only project name was requested
	if categoryName == "" && taskName == "" {
		return project, nil, nil, nil
	}

	// 2. Jump to Category
	category, ok := project.CategoryMap[categoryName]
	if !ok {
		return nil, nil, nil, errors.New("Category not found")
	}

	// If only project name and category was requested
	if taskName == "" {
		return project, category, nil, nil
	}

	// 3. Jump to Task
	task, ok := category.TaskMap[taskName]
	if !ok {
		return nil, nil, nil, errors.New("Task not found")
	}

	return project, category, task, nil
}

func (a *App) projectExists(projectName string) (bool, error) {
	// We pass empty strings for category and task since we only care about the project
	project, _, _, err := a.findByPath(projectName, "", "")

	// If there is no error and the pointer isn't nil, the project exists
	return err == nil && project != nil, err
}

func (a *App) categoryExists(projectName string, categoryName string) (bool, error) {
	// We pass an empty string for the taskName
	_, category, _, err := a.findByPath(projectName, categoryName, "")

	// Category exists only if the search returned a valid pointer without errors
	return err == nil && category != nil, err
}

func (a *App) taskExists(projectName string, categoryName string, taskName string) (bool, error) {
	_, _, task, err := a.findByPath(projectName, categoryName, taskName)

	// Task exists only if the final leaf in the tree was found
	return err == nil && task != nil, err
}

func (a *App) isValidTask(projectName string, categoryName string, taskName string) (bool, error) {
	if exists, _ := a.categoryExists(projectName, categoryName); exists == false || taskName == "" {
		return false, errors.New("Path to task is invalid")
	}
	if exists, _ := a.taskExists(projectName, categoryName, taskName); exists == true {
		return false, errors.New("Task already exists")
	}

	return true, nil
}

func (a *App) isValidCategory(projectName string, categoryName string) (bool, error) {
	if exists, _ := a.projectExists(projectName); exists == false || categoryName == "" {
		return false, errors.New("Path to category is invalid")
	}
	if exists, _ := a.categoryExists(projectName, categoryName); exists == true {
		return false, errors.New("Category already exists")
	}

	return true, nil
}

func (a *App) isValidProject(projectName string) (bool, error) {
	if exists, _ := a.projectExists(projectName); exists == true || projectName == "" {
		return false, errors.New("Project exists already or is blank")
	}
	return true, nil
}

func (a *App) findProjectIndex(name string) (projectIndex uint,err error) {
	// Iterate through projects to find 
	for i := range a.Projects {
		if (a.Projects)[i].Name == name {
			return uint(i), nil
		}
	}
	return 0, errors.New("no project with same name found")
}

func (a *App) findCategoryIndex(projectName string, categoryName string) ( uint, uint, error) {
	// We use the index 'i' to access the original slice memory
	projectIndex, err := a.findProjectIndex(projectName)
	if err != nil {
		return 0, 0, err
	}

	for index, category := range a.Projects[projectIndex].Categories {
		if category.Name == categoryName {
			return projectIndex, uint(index), nil
		}
	}

	return 0, 0, errors.New("no category with same name found")
}

func (a *App) findTaskIndex(projectName string, categoryName string, taskName string) (uint, uint, uint, error) {
	// Find relevant project and categories
	projectIndex, categoryIndex, err := a.findCategoryIndex(projectName, categoryName)
	if err != nil {
		return 0, 0,0, err
	}

	// Find if name of tasks from relevant category matches task name
	for index, task := range a.Projects[projectIndex].Categories[categoryIndex].Tasks {
		if task.Name == taskName {
			return projectIndex, categoryIndex, uint(index), nil
		}
	}
	return 0, 0,0, errors.New("no task with same name found")
}

