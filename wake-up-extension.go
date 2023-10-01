package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func checkIfReviewAppIsAsleep(namespace string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", os.Getenv("THANOS_URL")+"/api/v1/query", nil)
		if err != nil {
			return errMsg(err)
		}
		q := req.URL.Query()
		query := os.Getenv("THANOS_QUERY")
		q.Add("query", strings.Replace(query, "%NS%", namespace, 1))
		req.URL.RawQuery = q.Encode()
		resp, err := client.Do(req)
		if err != nil {
			return errMsg(err)
		}
		var thResponse ThanosResponse
		err = json.NewDecoder(resp.Body).Decode(&thResponse)
		if err != nil {
			return errMsg(err)
		}
		if thResponse.IsAsleep() {
			return noOutputResultMsg{false, "Review app is sleeping"}
		} else {
			return noOutputResultMsg{true, "Review app is awake"}
		}
	}
}

func getReviewAppsSleepingStatus() ([]ThanosResult, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("THANOS_URL")+"/api/v1/query", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	query := os.Getenv("THANOS_QUERY_ALL_NS")
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var thResponse ThanosResponse
	err = json.NewDecoder(resp.Body).Decode(&thResponse)
	if err != nil {
		return nil, err
	}
	return thResponse.Data.Result, nil
}
