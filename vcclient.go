package vcclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8001/VisualCron/json"
)

type VCClient struct {
	BaseURL    string
	UserName   string
	Password   string
	HTTPClient *http.Client
	Token      string
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

/*
Invoke new VC Client
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
Get auth token from api
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
	body, err := ioutil.ReadAll(res.Body)
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

func (c *VCClient) sendRequest(req *http.Request, v interface{}) error {
	// Get API Token before actual request to api
	token, err := GetToken(c)
	// If token was retrieved successfully continue with api request else log fatal
	if err != nil {
		log.Fatalln(err)
		return err
	}
	// Set token
	c.Token = token
	log.Println(c.Token)
	// Do the request for the requested API endpoint
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	log.Println(res.Body)
	log.Println(res.StatusCode)
	// Handle HTTP return codes
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	log.Println(v)
	// Unmarshall and populate interface
	fullResponse := successResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}
	return nil
}
