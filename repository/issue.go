package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"jira/domain"
	"jira/lib"
	"time"
)

const jiraTimeLayout = "2006-01-02T15:04:05.000Z0700"

const activeBugJQLQuery = `project = AC AND issuetype in ("Bug de GIN", "Incidencia de GIN", "Technical Bug") AND status in (Activo, "En Proceso", "En progreso", "Esperando Deploy", "Pendiente de Fix")`
const issueByKeyJQLQuery = `project = AC AND key = %s`
const activePedidoDeFixJQLQuery = `project = AC AND issuetype = "Pedido de fix GIN" AND status in (Activo) AND labels in (EMPTY) ORDER BY cf[11400] ASC, created ASC`

func GetIssue(ctx context.Context, dataStoreClient *datastore.Client, ID string) (domain.Issue, error) {

	var issue domain.Issue
	err := dataStoreClient.Get(ctx, datastore.NameKey("Issue", ID, nil), &issue)
	if err != nil {
		return domain.Issue{}, err
	}

	return issue, nil
}

func GetAllIssues(ctx context.Context, dataStoreClient *datastore.Client, issueType string) ([]*domain.Issue, error) {

	q := datastore.NewQuery("Issue").Filter("Type =", issueType)

	var issues []*domain.Issue
	_, err := dataStoreClient.GetAll(ctx, q, &issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func GetIssuesToNotify(ctx context.Context, dataStoreClient *datastore.Client, config lib.Config, issueType string) ([]*domain.Issue, error) {

	var deadLine time.Duration
	if issueType == domain.Bug {
		deadLine = config.BugDeadline
	} else {
		deadLine = config.PedidoDeFixDeadline
	}

	q := datastore.NewQuery("Issue").
		Filter("DueDate <", time.Now().Add(deadLine)).
		Filter("DueDate >", time.Now())

	var issuesToCheck []*domain.Issue
	_, err := dataStoreClient.GetAll(ctx, q, &issuesToCheck)
	if err != nil {
		return nil, err
	}

	var issuesToNotify []*domain.Issue
	for i := range issuesToCheck {
		if issuesToCheck[i].TimesNotified < config.MaxTimesToNotify && issuesToCheck[i].Type == issueType {
			issuesToNotify = append(issuesToNotify, issuesToCheck[i])
		}
	}

	return issuesToNotify, nil
}

func UpdateIssuesNotifications(ctx context.Context, dataStoreClient *datastore.Client, issues []*domain.Issue) error {

	var keys []*datastore.Key
	for i := range issues {
		issues[i].TimesNotified = issues[i].TimesNotified + 1
		keys = append(keys, datastore.NameKey("Issue", issues[i].ID, nil))
	}

	if _, err := dataStoreClient.PutMulti(ctx, keys, issues); err != nil {
		return err
	}

	return nil
}

func IndexActiveBugs(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client) (newBugs []*domain.Issue, err error) {

	issues, _, err := jiraClient.Issue.Search(activeBugJQLQuery, &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	})
	if err != nil {
		return newBugs, err
	}

	for _, issue := range issues {
		k := datastore.NameKey("Issue", issue.Key, nil)

		var issueToLookup domain.Issue
		err := dataStoreClient.Get(ctx, k, &issueToLookup)
		if err != nil {
			if err != datastore.ErrNoSuchEntity {
				return newBugs, err
			}
		}

		if err == nil {
			continue
		}

		i := new(domain.Issue)
		i.ID = k.Name

		if issue.Fields.Assignee != nil {
			i.Assignee = issue.Fields.Assignee.DisplayName
		}

		stringDate, _ := issue.Fields.Unknowns["customfield_11400"].(string)
		t, err := time.Parse(jiraTimeLayout, stringDate)
		if err != nil {
			return newBugs, err
		}
		i.DueDate = t
		i.Status = issue.Fields.Status.Name
		i.Priority = issue.Fields.Priority.Name
		i.URL = fmt.Sprintf("https://mercadolibre.atlassian.net/browse/%s", i.ID)
		i.ListURL = fmt.Sprintf("https://mercadolibre.atlassian.net/browse/%s?filter=18341", i.ID)
		i.Description = issue.Fields.Summary
		i.Type = domain.Bug

		if _, err := dataStoreClient.Put(ctx, k, i); err != nil {
			return newBugs, err
		}

		newBugs = append(newBugs, i)
	}

	return newBugs,nil
}

func UpdateCurrentIssues(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client, statesToCheck []string) error {

	q := datastore.NewQuery("Issue").KeysOnly()

	var activeIssues []domain.Issue
	keys, err := dataStoreClient.GetAll(ctx, q, activeIssues)
	if err != nil {
		return err
	}

	for _, key := range keys {
		var issueToUpdate domain.Issue
		err := dataStoreClient.Get(ctx, key, &issueToUpdate)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		issues, _, err := jiraClient.Issue.Search(fmt.Sprintf(issueByKeyJQLQuery, issueToUpdate.ID), &jira.SearchOptions{
			StartAt:    0,
			MaxResults: 1,
		})

		if lib.ContainsString(issues[0].Fields.Status.Name, statesToCheck) {
			issueToUpdate.Status = issues[0].Fields.Status.Name
			issueToUpdate.Priority = issues[0].Fields.Priority.Name

			if issues[0].Fields.Assignee != nil {
				issueToUpdate.Assignee = issues[0].Fields.Assignee.DisplayName
			} else {
				issueToUpdate.Assignee = ""
			}

			stringDate, _ := issues[0].Fields.Unknowns["customfield_11400"].(string)
			t, err := time.Parse(jiraTimeLayout, stringDate)
			if err != nil {
				return err
			}
			issueToUpdate.DueDate = t

			if _, err := dataStoreClient.Put(ctx, key, &issueToUpdate); err != nil {
				return err
			}
		} else {
			if err := dataStoreClient.Delete(ctx, key); err != nil {
				return err
			}
		}
	}

	return nil
}

func IndexActivePedidosDeFix(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client) (newPedidosDeFix []*domain.Issue, err error) {

	issues, _, err := jiraClient.Issue.Search(activePedidoDeFixJQLQuery, &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	})
	if err != nil {
		return newPedidosDeFix, err
	}

	for _, issue := range issues {
		k := datastore.NameKey("Issue", issue.Key, nil)

		var issueToLookup domain.Issue
		err := dataStoreClient.Get(ctx, k, &issueToLookup)
		if err != nil {
			if err != datastore.ErrNoSuchEntity {
				return newPedidosDeFix, err
			}
		}

		if err == nil {
			continue
		}

		i := new(domain.Issue)
		i.ID = k.Name

		if issue.Fields.Assignee != nil {
			i.Assignee = issue.Fields.Assignee.DisplayName
		}

		stringDate, _ := issue.Fields.Unknowns["customfield_11400"].(string)
		t, err := time.Parse(jiraTimeLayout, stringDate)
		if err != nil {
			return newPedidosDeFix, err
		}

		i.DueDate = t
		i.Status = issue.Fields.Status.Name
		i.Priority = issue.Fields.Priority.Name
		i.URL = fmt.Sprintf("https://mercadolibre.atlassian.net/browse/%s", i.ID)
		i.ListURL = fmt.Sprintf("https://mercadolibre.atlassian.net/browse/%s?filter=18342", i.ID)
		i.Description = issue.Fields.Summary
		i.Type = domain.PedidoDeFix

		if _, err := dataStoreClient.Put(ctx, k, i); err != nil {
			return newPedidosDeFix, err
		}

		newPedidosDeFix = append(newPedidosDeFix, i)
	}

	return newPedidosDeFix, err
}

func ResetIssuesNotifications(ctx context.Context, dataStoreClient *datastore.Client, issues []*domain.Issue) error {

	var keys []*datastore.Key
	for i := range issues {
		issues[i].TimesNotified = 0
		keys = append(keys, datastore.NameKey("Issue", issues[i].ID, nil))
	}

	if _, err := dataStoreClient.PutMulti(ctx, keys, issues); err != nil {
		return err
	}

	return nil
}

func Take(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client, user lib.JiraUser, issueID string) error {

	issue, err := GetIssue(ctx, dataStoreClient, issueID )
	if err != nil {
		return err
	}

	_, err = jiraClient.Issue.UpdateAssignee(issue.ID, &jira.User{AccountID: string(user.Username)})
	if err != nil {
		return err
	}

	return nil
}

func Release(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client, issueID string) error {

	issue, err := GetIssue(ctx, dataStoreClient, issueID )
	if err != nil {
		return err
	}

	_, err = jiraClient.Issue.UpdateAssignee(issue.ID, &jira.User{})
	if err != nil {
		return err
	}

	return nil
}
