package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
		Columns: projectIterationColumns(),
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

func projectIterationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "The ID of the iteration.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "iid",
			Description: "The instance ID of the iteration.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "sequence",
			Description: "The sequence number of the iteration.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "project_id",
			Description: "The ID of the project to which this iteration belongs.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("project_id"),
		},
		{
			Name:        "title",
			Description: "The title of the iteration.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "The description of the iteration.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "state",
			Description: "The state of the iteration.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Timestamp of when the iteration was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp of when the iteration was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "start_date",
			Description: "Timestamp indicating the start of the iteration.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "due_date",
			Description: "Timestamp indicating the due date of the iteration.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "web_url",
			Description: "The web url of the iteration.",
			Type:        proto.ColumnType_STRING,
		},
	}
}
