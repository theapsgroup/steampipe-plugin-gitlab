package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableGroupIteration() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_iteration",
		Description: "Obtain information about iterations for a specific group within the GitLab instance.",
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

// Hydrate Functions
func listGroupIterations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupIterations", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupIterations", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	groupId := int(q["group_id"].GetInt64Value())
	opt := &api.ListGroupIterationsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		plugin.Logger(ctx).Debug("listGroupIterations", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)

		iterations, resp, err := conn.GroupIterations.ListGroupIterations(groupId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroupIterations", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain iterations for group_id %d\n%v", groupId, err)
		}

		for _, iteration := range iterations {
			d.StreamListItem(ctx, iteration)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listGroupIterations", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroupIterations", "completed successfully")
	return nil, nil
}

// Column Function
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
			Transform:   transform.FromField("IID").NullIfZero(),
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
