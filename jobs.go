package vcclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

/*
Contains all relevant info for a Job
*/
type Jobs struct {
	ResetTriggerDepedency bool    `json:"ResetTriggerDepedency"`
	Stats                 Stats   `json:"Stats"`
	TimeOut               TimeOut `json:"TimeOut"`
	ID                    string  `json:"Id"`
	Test                  bool    `json:"Test"`
	IsTaskRepository      bool    `json:"IsTaskRepository"`
	UseRunRandomValue     bool    `json:"UseRunRandomValue"`
	RunRandomValue        int     `json:"RunRandomValue"`
	Missed                bool    `json:"Missed"`
	MissedDate            string  `json:"MissedDate"`
	Name                  string  `json:"Name"`
	Description           string  `json:"Description"`
	Group                 string  `json:"Group"`
	JobDebugging          bool    `json:"JobDebugging"`
	RunMissed             bool    `json:"RunMissed"`
	RunOnce               bool    `json:"RunOnce"`
	RemoveAfterExecution  bool    `json:"RemoveAfterExecution"`
	RunTasksInOrder       bool    `json:"RunTasksInOrder"`
	NotStartIfRunning     bool    `json:"NotStartIfRunning"`
	QueueJobs             bool    `json:"QueueJobs"`
	UniqueRunID           int     `json:"UniqueRunId"`
}

/*
Contains all metadata for a job's stats
*/
type Stats struct {
	JobID               string  `json:"JobId"`
	Active              bool    `json:"Active"`
	ExitCode            int     `json:"ExitCode"`
	ExitCodeResult      int     `json:"ExitCodeResult"`
	DateLastExecution   string  `json:"DateLastExecution"`
	DateNextExecution   string  `json:"DateNextExecution"`
	DateLastExited      string  `json:"DateLastExited"`
	DateCreated         string  `json:"DateCreated"`
	DateModified        string  `json:"DateModified"`
	NoExecutes          int     `json:"NoExecutes"`
	ExecutionTime       float64 `json:"ExecutionTime"`
	LastTriggerID       string  `json:"LastTriggerId"`
	Status              int     `json:"Status"`
	CPUTime             float64 `json:"CPUTime"`
	TriggerCPUTime      float64 `json:"TriggerCPUTime"`
	TasksCPUTime        float64 `json:"TasksCPUTime"`
	NotificationCPUTime float64 `json:"NotificationCPUTime"`
}

/*
Contains all configuration about timeout settings per job
*/
type TimeOut struct {
	Use           bool `json:"Use"`
	Hours         int  `json:"Hours"`
	Minutes       int  `json:"Minutes"`
	Seconds       int  `json:"Seconds"`
	ConsiderError bool `json:"ConsiderError"`
}

/*
Get all Jobs configured in VisualCron /VisualCron/json/Job/List
*/
func (c *VCClient) GetJobs(ctx context.Context) (*[]Jobs, error) {
	// Get API Token before actual request to api
	token, err := GetToken(c)
	// If token was retrieved successfully continue with api request else log fatal
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// Set token
	c.Token = token
	// Get API Endpoint for all Jobs
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Job/List?token=%s", c.BaseURL, c.Token), nil)
	if err != nil {
		return nil, err
	}
	// Execute the request
	req = req.WithContext(ctx)

	// Send the request to the endpoint
	var res []Jobs
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
