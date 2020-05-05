package sdk

import "context"

type Report struct {
}

type Expense struct {
}

func (c *Client) Report(ctx context.Context, employeeEmail, policyID string, report Report, expenses []Expense) error {
	panic("not implemented")
}
