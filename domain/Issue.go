package domain

import (
	_ "cloud.google.com/go/datastore"

	"time"
)

const PedidoDeFix = "Pedido de Fix"
const Bug = "Bug"

// Issue is a Centinela's representation of a Jira Issue
type Issue struct {
	ID            string
	Summary   string
	Description   string `datastore:",noindex"`
	Type   		  string
	Assignee      string
	Status        string
	TimesNotified int
	Priority      string
	URL           string
	ListURL		  string
	DueDate       time.Time
}

