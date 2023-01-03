package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "github.com/xanzy/go-gitlab"
)

func tableGroupIteration() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_iteration",
		Description: "Iterations for a specific group in the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listGroupIterations,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "group_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: iterationColumns(),
	}
}

func iterationColumns() []*plugin.Column {
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
			Name:        "group_id",
			Description: "The ID of the group to which this iteration belongs.",
			Type:        proto.ColumnType_INT,
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

func listGroupIterations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.EqualsQuals

	groupId := int(q["group_id"].GetInt64Value())

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opt := &api.ListGroupIterationsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		iterations, resp, err := conn.GroupIterations.ListGroupIterations(groupId, opt)
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
