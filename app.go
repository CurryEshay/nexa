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
	Locked               bool                `json:"-"`
	currentProjectIndex  int                 `json:"-"`
	currentCategoryIndex int                 `json:"-"`
	currentTaskIndex     int                 `json:"-"`
	deletedTasks         []*Task             `json:"-"`
	PriorityCompression  float64             `json:"priority_compression"`
	SmoothingConstant    float64             `json:"smoothing_constant"`
	TimeAggression       float64             `json:"time_aggression"`
	OverdueConstant      float64             `json:"overdue_constant"`
	OverdueAggression    float64             `json:"overdue_aggression"`
}

func NewApp() *App {
	app := &App{
		Projects:            []*Project{},
		ProjectMap:          map[string]*Project{},
		Running:             true,
		Locked:              true,
		deletedTasks:        []*Task{},
		PriorityCompression: 0.9,
		SmoothingConstant:   24,
		TimeAggression:      1.15,
		OverdueConstant:     0.75,
		OverdueAggression:   1.5,
	}
	app.InitCommands()
	return app
}

func (a *App) ErrorLog(err error) {
	fmt.Println("Unexpected error occured: " + err.Error())
}
