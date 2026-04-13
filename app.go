package main

import (
	"fmt"
	"sync"
)

type App struct {
	Projects             []*Project          `json:"projects"`
	ProjectMap           map[string]*Project `json:"-"`
	Mu                   sync.Mutex          `json:"-"`
	CMDMap               map[string]CMDFunc  `json:"-"`
	Running              bool                `json:"-"`
	Locked               bool                `json:"locked"`
	currentProjectIndex  int                 `json:"-"`
	currentCategoryIndex int                 `json:"-"`
	currentTaskIndex     int                 `json:"-"`
	deletedTasks         []*Task
}

func NewApp() *App {
	app := &App{
		Projects:   []*Project{},
		ProjectMap: map[string]*Project{},
		Running:    true,
		Locked:     true,
		deletedTasks: []*Task{},
	}
	app.InitCommands()
	return app
}

func (a *App) ErrorLog(err error) {
	fmt.Println("Unexpected error occured: " + err.Error())
}

func (a *App) Run() {
	for a.Running == true {
		a.Parser()
	}
	err := a.Save(nil)
	if err != nil {
		fmt.Println("Could not save tasks")
	}
}
