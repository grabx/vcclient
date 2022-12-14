package vcclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8001/VisualCron/json"
)

/*
The VisualCron API Client
*/
type VCClient struct {
	BaseURL    string
	UserName   string
	Password   string
	HTTPClient *http.Client
	Token      string
}

/*
Contains error Response data for failed API requests
*/
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*
Invoke new VisualCron API Client
*/
func NewClient(userName string, password string) *VCClient {
	return &VCClient{
		BaseURL:  baseURL,
		UserName: userName,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

/*
Get an auth token from the VisualCron API /VisualCron/json/logon?username=<name>&password=<password>
*/
func GetToken(c *VCClient) (string, error) {
	loginClient := http.Client{
		Timeout: time.Minute,
	}
	loginRequest, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s?username=%s&password=%s",
			c.BaseURL,
			"/logon",
			c.UserName,
			c.Password),
		nil)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	res, err := loginClient.Do(loginRequest)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	loginBody := Login{}
	jsonErr := json.Unmarshal(body, &loginBody)
	if jsonErr != nil {
		log.Fatalln(jsonErr)
		return "", err
	}
	return loginBody.Token, nil
}

/*
Helper function that parses the json to the provided interface to avoid code duplication
*/
func (c *VCClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	// Do the request for the requested API endpoint
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println(res.StatusCode)
	// Handle HTTP return codes
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	// Decode json into the interface passed from the endpoint's method
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		log.Println("Something went wrong")
		return err
	}
	return nil
}
