package domain

import "time"

// Issue is a Centinela's representation of a Jira Issue
type Issue struct {
	ID            string
	Assignee      string
	Status        string
	TimesNotified int
	DueDate       time.Time
}
