package main

import (
        "bytes"
        "encoding/json"
	"fmt"
        "net/http"
        "log"
        "io/ioutil"
)

type (
    OnCallInfo struct {
        Level   int64
    }

    User struct {
        AvatarURL   string
        Email       string
        Name        string
        OnCallInfo      []OnCallInfo `json:"on_call"`
    }
    OnCallAPIResponse struct {
        Users       []User
    }
)

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
    var response OnCallAPIResponse
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

