package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type (
	OnCallInfo struct {
		Level int64
	}

	User struct {
		AvatarURL  string
		Email      string
		Name       string
		OnCallInfo []OnCallInfo `json:"on_call"`
	}
	APIResponse struct {
		Users     []User     `json:"users"`
		Incidents []Incident `json:"incidents"`
	}
	Incident struct {
		Id                  string `json:"id"`
		IncidentKey         string `json:"incident_key"`
		NumberOfEscalations int64  `json:"number_of_escalations"`
		ResolvedByUser      User   `json:"resolved_by_user"`
		Status              string `json:"status"`
		URL                 string `json:"html_url"`
		TriggerType         string `json:"trigger_type"`
		TriggerDetails      string `json:"trigger_details_html_url"`
	}
)

var HAPPY_GIFS = [...]string{
	"http://replygif.net/i/726.gif",
	"http://replygif.net/i/765.gif",
	"http://replygif.net/i/321.gif",
	"http://replygif.net/i/104.gif",
	"http://replygif.net/i/720.gif",
	"http://replygif.net/i/769.gif",
	"http://replygif.net/i/1033.gif",
	"http://replygif.net/i/1116.gif",
	"http://replygif.net/i/97.gif",
	"http://replygif.net/i/840.gif",
	"http://replygif.net/i/237.gif",
}

func getAccountURL() string {
	return fmt.Sprintf("https://%s.pagerduty.com/api/v1", account)
}

func doRequest(path string) (*http.Response, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", getAccountURL(), path)
	// for now just use GET since we are read-only
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error connecting to PagerDuty: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", token))
	return client.Do(req)
}

// Ideally this would be abstracted into a PagerDuty module
// Returns current person on-call
func getOnCall() (string, error) {
	resp, err := doRequest("/users/on_call")
	if err != nil {
		log.Printf("Error requesting on call from PagerDuty: %s", err)
		return "", err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from PagerDuty: %s", err)
		return "", err
	}
	var response APIResponse
	c := bytes.NewBufferString(string(contents))
	d := json.NewDecoder(c)
	if err := d.Decode(&response); err != nil {
		log.Printf("Error decoding JSON from PagerDuty: %s", err)
		return "", err
	}
	resp.Body.Close()
	var oncall string
	for _, u := range response.Users {
		for _, x := range u.OnCallInfo {
			if x.Level == 1 {
				oncall = u.Name
				break
			}
		}
	}
	return fmt.Sprintf("Current on call: %s", oncall), nil
}

// Returns current incidents (only shows triggered and acknowledged)
func getIncidents() (string, error) {
	resp, err := doRequest("/incidents?status=triggered,acknowledged")
	if err != nil {
		log.Printf("Error requesting incidents from PagerDuty: %s", err)
		return "", err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from PagerDuty: %s", err)
		return "", err
	}
	var response APIResponse
	c := bytes.NewBufferString(string(contents))
	d := json.NewDecoder(c)
	if err := d.Decode(&response); err != nil {
		log.Printf("Error decoding JSON from PagerDuty: %s", err)
		return "", err
	}
	resp.Body.Close()
	data := ""
	if len(response.Incidents) > 0 {
		for _, i := range response.Incidents {
			data += fmt.Sprintf("%s: %s\n  Status: %s\n  Escalations: %v\n  Details: %s\n", i.Id, i.IncidentKey, i.Status, i.NumberOfEscalations, i.URL)
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		gif := HAPPY_GIFS[rand.Intn(len(HAPPY_GIFS))]
		data = fmt.Sprintf("No incidents... \n  %s", gif)
	}
	return data, nil
}
