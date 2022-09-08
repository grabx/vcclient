package vcclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Jobs struct {
	ResetTriggerDepedency bool   `json:"ResetTriggerDepedency"`
	Stats                 Stats  `json:"Stats"`
	ID                    string `json:"Id"`
	Test                  bool   `json:"Test"`
	IsTaskRepository      bool   `json:"IsTaskRepository"`
	UseRunRandomValue     bool   `json:"UseRunRandomValue"`
	RunRandomValue        int    `json:"RunRandomValue"`
	Missed                bool   `json:"Missed"`
	MissedDate            string `json:"MissedDate"`
	Name                  string `json:"Name"`
	Description           string `json:"Description"`
	Group                 string `json:"Group"`
	JobDebugging          bool   `json:"JobDebugging"`
	RunMissed             bool   `json:"RunMissed"`
	RunOnce               bool   `json:"RunOnce"`
	RemoveAfterExecution  bool   `json:"RemoveAfterExecution"`
	RunTasksInOrder       bool   `json:"RunTasksInOrder"`
	NotStartIfRunning     bool   `json:"NotStartIfRunning"`
	QueueJobs             bool   `json:"QueueJobs"`
	UniqueRunID           int    `json:"UniqueRunId"`
}

type Stats struct {
	JobID               string    `json:"JobId"`
	Active              bool      `json:"Active"`
	ExitCode            int       `json:"ExitCode"`
	ExitCodeResult      int       `json:"ExitCodeResult"`
	DateLastExecution   time.Time `json:"DateLastExecution"`
	DateNextExecution   time.Time `json:"DateNextExecution"`
	DateLastExited      time.Time `json:"DateLastExited"`
	DateCreated         time.Time `json:"DateCreated"`
	DateModified        time.Time `json:"DateModified"`
	NoExecutes          int       `json:"NoExecutes"`
	ExecutionTime       float64   `json:"ExecutionTime"`
	LastTriggerID       string    `json:"LastTriggerId"`
	Status              int       `json:"Status"`
	CPUTime             int       `json:"CPUTime"`
	TriggerCPUTime      int       `json:"TriggerCPUTime"`
	TasksCPUTime        float64   `json:"TasksCPUTime"`
	NotificationCPUTime int       `json:"NotificationCPUTime"`
}

/*
Get Visual Cron Jobs
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
	log.Println(c.Token)

	// Get API Endpoint for all Jobs
	log.Printf("Requesting url: %s", fmt.Sprintf("%s/Job/List?token=%s", c.BaseURL, c.Token))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Job/List?token=%s", c.BaseURL, c.Token), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res := []Jobs{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
