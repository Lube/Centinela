package commands

import (
	"cloud.google.com/go/datastore"
	"context"
	"jira/domain"
	"jira/repository"
)

// ExpireNotifications pools Jira's api and syncs issues with Centinela's dataStore.
func ExpireNotifications() error {

	ctx := context.Background()

	dataStoreClient, err := datastore.NewClient(ctx, "centinela-258804")
	if err != nil {
		return err
	}

	bugs, err := repository.GetAllIssues(ctx, dataStoreClient, domain.Bug)
	if err != nil {
		return err
	}

	pedidosDeFix, err := repository.GetAllIssues(ctx, dataStoreClient, domain.PedidoDeFix)
	if err != nil {
		return err
	}

	err = repository.ResetIssuesNotifications(ctx, dataStoreClient, append(bugs, pedidosDeFix...))
	if err != nil {
		return err
	}

	return nil
}
