package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableMyIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_my_issue",
		Description: "Obtain information about issues that are created by or assigned to the authenticated user within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listMyIssues,
		},
		Columns: issueColumns(),
	}
}

func listMyIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listMyIssues", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listMyIssues", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	createdByScope := "created_by_me"
	assignedToScope := "assigned_to_me"
	createdByOptions := &api.ListIssuesOptions{Scope: &createdByScope, ListOptions: api.ListOptions{Page: 1, PerPage: 50}}
	assignedToOptions := &api.ListIssuesOptions{Scope: &assignedToScope, ListOptions: api.ListOptions{Page: 1, PerPage: 50}}

	for {
		plugin.Logger(ctx).Debug("listMyIssues", "type", createdByScope, "page", createdByOptions.Page, "perPage", createdByOptions.PerPage)
		issues, resp, err := conn.Issues.ListIssues(createdByOptions)
		if err != nil {
			plugin.Logger(ctx).Error("listMyIssues", "type", createdByScope, "page", createdByOptions.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain issues created by the current user\n%v", err)
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listMyIssues", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		createdByOptions.Page = resp.NextPage
	}

	for {
		plugin.Logger(ctx).Debug("listMyIssues", "type", assignedToScope, "page", assignedToOptions.Page, "perPage", assignedToOptions.PerPage)
		issues, resp, err := conn.Issues.ListIssues(assignedToOptions)
		if err != nil {
			plugin.Logger(ctx).Error("listMyIssues", "type", assignedToScope, "page", assignedToOptions.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain issues assigned to the current user\n%v", err)
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listMyIssues", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		assignedToOptions.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listMyIssues", "completed successfully")
	return nil, nil
}
