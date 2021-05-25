package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectMergeRequest() *plugin.Table {
	return &plugin.Table{
		Name: "gitlab_project_merge_request",
		Description: "GitLab Merge Requests for a specific Project",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate: listProjectMergeRequests,
		},
		Columns: gitlabMergeRequestColumns(),
	}
}

func listProjectMergeRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())

	opt := &api.ListProjectMergeRequestsOptions{
		ListOptions: api.ListOptions{
			Page: 1,
			PerPage: 30,
		},
	}

	for {
		mergeRequests, response, err := conn.MergeRequests.ListProjectMergeRequests(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, mergeRequest := range mergeRequests {
			d.StreamListItem(ctx, mergeRequest)
		}

		if response.CurrentPage >= response.TotalPages {
			break
		}

		opt.Page = response.NextPage
	}

	return nil, nil
}