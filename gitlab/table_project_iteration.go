package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectIteration() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_iteration",
		Description: "Iterations for a specific project in the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listProjectIterations,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: iterationColumns(),
	}
}

func listProjectIterations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.KeyColumnQuals

	projectId := int(q["project_id"].GetInt64Value())

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opt := &api.ListProjectIterationsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		iterations, resp, err := conn.ProjectIterations.ListProjectIterations(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, iteration := range iterations {
			d.StreamListItem(ctx, iteration)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
