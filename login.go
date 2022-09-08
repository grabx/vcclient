package vcclient

/*
Contains the login data for the VisualCron API
*/
type Login struct {
	Result int    `json:"Result"`
	Token  string `json:"Token"`
}
