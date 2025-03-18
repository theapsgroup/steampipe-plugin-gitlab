package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableEpic() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_epic",
		Description: "Obtain information about epics for a specific group within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listEpics,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "group_id",
					Require: plugin.Required,
				},
				{
					Name:      "author_id",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
				{
					Name:      "state",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
			},
		},
		Columns: epicColumns(),
	}
}

// Hydrate Functions
func listEpics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listEpics", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listEpics", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	groupId := int(q["group_id"].GetInt64Value())
	opt := &api.ListGroupEpicsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	// Optional Qualifiers
	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opt.AuthorID = &authorId
		plugin.Logger(ctx).Debug("listEpics", "filter[author_id]", authorId)
	}

	if q["state"] != nil {
		state := q["state"].GetStringValue()
		opt.State = &state
		plugin.Logger(ctx).Debug("listEpics", "filter[state]", state)
	}

	for {
		plugin.Logger(ctx).Debug("listEpics", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)
		epics, resp, err := conn.Epics.ListGroupEpics(groupId, opt)
		if err != nil {
			if resp.StatusCode == 403 {
				return nil, nil
			}
			plugin.Logger(ctx).Error("listEpics", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain branches for group_id %d\n%v", groupId, err)
		}

		for _, epic := range epics {
			d.StreamListItem(ctx, epic)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listEpics", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listEpics", "completed successfully")
	return nil, nil
}

// Column Function
func epicColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "The ID of the epic.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "iid",
			Description: "The instance ID of the epic.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("IID").NullIfZero(),
		},
		{
			Name:        "group_id",
			Description: "The ID of the parent group for this epic.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "parent_id",
			Description: "The ID of the parent for the epic.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "title",
			Description: "The epics title.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "The epics description.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "state",
			Description: "The state of the epic.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "web_url",
			Description: "The web url of the epic.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "author_id",
			Description: "The ID of the author of the epic.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Author.ID"),
		},
		{
			Name:        "author",
			Description: "The username of the author of the epic.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Author.Username"),
		},
		{
			Name:        "author_name",
			Description: "The display name of the author of the epic.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Author.Name"),
		},
		{
			Name:        "author_url",
			Description: "The url for the authors profile.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Author.WebURL"),
		},
		{
			Name:        "start_date",
			Description: "Timestamp indicating the start of the epic.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("StartDate").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "due_date",
			Description: "Timestamp indicating the due date of the epic.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("DueDate").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "end_date",
			Description: "Timestamp indicating the end of the epic.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "labels",
			Description: "An array of labels associated with the epic.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp of when the epic was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp of when the epic was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "closed_at",
			Description: "Timestamp of when the epic was closed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "upvotes",
			Description: "The number of up votes for the epic.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "downvotes",
			Description: "The number of down votes for the epic.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "user_notes_count",
			Description: "A count of user notes on the epic.",
			Type:        proto.ColumnType_INT,
		},
	}
}
