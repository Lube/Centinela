package domain

import "time"

const PedidoDeFix = "Pedido de Fix"
const Bug = "Bug"

// Issue is a Centinela's representation of a Jira Issue
type Issue struct {
	ID            string
	Description   string
	Type   		  string
	Assignee      string
	Status        string
	TimesNotified int
	DueDate       time.Time
}