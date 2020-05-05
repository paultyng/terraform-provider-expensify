package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Report struct {
	Title string `json:"title"`
}

type Expense struct {
	Merchant string `json:"merchant"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
	Amount   int    `json:"amount"`
}

func (c *Client) Report(ctx context.Context, employeeEmail, policyID string, report Report, expenses []Expense) (string, error) {
	reqStructure := struct {
		requestJobDescription

		InputSettings struct {
			Type          string    `json:"type"`
			EmployeeEmail string    `json:"employeeEmail"`
			PolicyID      string    `json:"policyID"`
			Report        Report    `json:"report"`
			Expenses      []Expense `json:"expenses"`
		} `json:"inputSettings"`
	}{
		requestJobDescription: requestJobDescription{
			Type: "create",
			Credentials: credentials{
				PartnerUserID:     c.userID,
				PartnerUserSecret: c.userSecret,
			},
		},
		InputSettings: struct {
			Type          string    "json:\"type\""
			EmployeeEmail string    "json:\"employeeEmail\""
			PolicyID      string    "json:\"policyID\""
			Report        Report    "json:\"report\""
			Expenses      []Expense "json:\"expenses\""
		}{
			Type:          "report",
			EmployeeEmail: employeeEmail,
			PolicyID:      policyID,
			Report:        report,
			Expenses:      expenses,
		},
	}
	reqJSON, err := json.Marshal(reqStructure)
	if err != nil {
		return "", fmt.Errorf("unable to marshal request body: %w", err)
	}

	data := bytes.NewBuffer([]byte("requestJobDescription="))
	data.Write(reqJSON)

	url := c.serverURL.ResolveReference(relativeURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), data)
	if err != nil {
		return "", fmt.Errorf("unable to build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to perform request: %w", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %w", err)
	}

	respStructure := &struct {
		ReportID   string `json:"reportID"`
		ReportName string `json:"reportName"`
	}{}

	err = json.Unmarshal(bodyBytes, &respStructure)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return respStructure.ReportID, nil
}
