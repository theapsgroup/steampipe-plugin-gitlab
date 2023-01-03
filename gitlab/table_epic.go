package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableEpic() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_epic",
		Description: "Epics for a specific group in the GitLab instance.",
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
		},
		{
			Name:        "group_id",
			Description: "The ID of the parent group for this epic.",
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
			Name:        "reference",
			Description: "The epics reference identifier.",
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
			Name:        "start_date",
			Description: "Timestamp indicating the start of the epic.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "due_date",
			Description: "Timestamp indicating the due date of the epic.",
			Type:        proto.ColumnType_TIMESTAMP,
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
	}
}

func listEpics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.EqualsQuals

	groupId := int(q["group_id"].GetInt64Value())

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

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
	}

	if q["state"] != nil {
		state := q["state"].GetStringValue()
		opt.State = &state
	}

	for {
		epics, resp, err := conn.Epics.ListGroupEpics(groupId, opt)
		if err != nil {
			return nil, err
		}

		for _, epic := range epics {
			d.StreamListItem(ctx, epic)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
