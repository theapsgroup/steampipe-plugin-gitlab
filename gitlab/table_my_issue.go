package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "github.com/xanzy/go-gitlab"
)

func tableMyIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_my_issue",
		Description: "GitLab Issues that are Created By or Assigned To the authenticated user.",
		List: &plugin.ListConfig{
			Hydrate: listMyIssues,
		},
		Columns: issueColumns(),
	}
}

func listMyIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	createdByScope := "created_by_me"
	assignedToScope := "assigned_to_me"
	createdByOptions := &api.ListIssuesOptions{Scope: &createdByScope, ListOptions: api.ListOptions{Page: 1, PerPage: 50}}
	assignedToOptions := &api.ListIssuesOptions{Scope: &assignedToScope, ListOptions: api.ListOptions{Page: 1, PerPage: 50}}

	for {
		issues, resp, err := conn.Issues.ListIssues(createdByOptions)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
		}

		if resp.NextPage == 0 {
			break
		}

		createdByOptions.Page = resp.NextPage
	}

	for {
		issues, resp, err := conn.Issues.ListIssues(assignedToOptions)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
		}

		if resp.NextPage == 0 {
			break
		}

		assignedToOptions.Page = resp.NextPage
	}

	return nil, nil
}
