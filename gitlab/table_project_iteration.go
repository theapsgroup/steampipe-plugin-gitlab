package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectIteration() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_iteration",
		Description: "Obtain information about iterations for a specific project in the GitLab instance.",
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

// Hydrate Functions
func listProjectIterations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectIterations", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectIterations", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	projectId := int(q["project_id"].GetInt64Value())
	opt := &api.ListProjectIterationsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		plugin.Logger(ctx).Debug("listProjectIterations", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		iterations, resp, err := conn.ProjectIterations.ListProjectIterations(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectIterations", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain iterations for project_id %d\n%v", projectId, err)
		}

		for _, iteration := range iterations {
			d.StreamListItem(ctx, iteration)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectIterations", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectIterations", "completed successfully")
	return nil, nil
}

// Column Function
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
