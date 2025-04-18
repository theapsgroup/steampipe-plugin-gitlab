package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableUserEvents() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_user_event",
		Description: "Obtain information about a user's events.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "author_id",
					Require: plugin.Required,
				},
				{
					Name:      "created_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">=", "=", "<", "<="},
				},
				{
					Name:      "target_type",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
				{
					Name:      "action_name",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
			},
			Hydrate: listUserEvents,
		},
		Columns: eventColumns(),
	}
}

// Hydrate Functions
func listUserEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listUserEvents", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserEvents", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	userID := int(d.EqualsQuals["author_id"].GetInt64Value())
	opt := &api.ListContributionEventsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	if d.Quals["created_at"] != nil {
		for _, q := range d.Quals["created_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			beforeTime := givenTime.Add(time.Duration(-1) * time.Second)
			afterTime := givenTime.Add(time.Second * 1)
			givenISOTime := api.ISOTime(givenTime)
			beforeISOTime := api.ISOTime(beforeTime)
			afterISOTime := api.ISOTime(afterTime)

			switch q.Operator {
			case ">":
				opt.After = &afterISOTime
			case ">=":
				opt.After = &givenISOTime
			case "=":
				opt.After = &beforeISOTime
				opt.Before = &afterISOTime
			case "<=":
				opt.Before = &givenISOTime
			case "<":
				opt.Before = &beforeISOTime
			}
		}
	}

	if d.Quals["target_type"] != nil {
		targetType := api.EventTargetTypeValue(
			strings.ToLower(d.EqualsQuals["target_type"].GetStringValue()))
		opt.TargetType = &targetType
	}

	if d.Quals["action_name"] != nil {
		action := api.EventTypeValue(
			strings.ToLower(d.EqualsQuals["action_name"].GetStringValue()))
		opt.Action = &action
	}

	for {
		plugin.Logger(ctx).Debug("listUserEvents", "userID", userID, "page", opt.Page, "perPage", opt.PerPage)
		events, resp, err := conn.Users.ListUserContributionEvents(userID, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listUserEvents", "userID", userID, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain events for user_id %d\n%v", userID, err)
		}

		for _, event := range events {
			plugin.Logger(ctx).Debug("listMyEvents", "event", event)
			d.StreamListItem(ctx, event)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listUserEvents", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listUserEvents", "completed successfully")
	return nil, nil
}

// Column Function
func eventColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The event ID",
		},
		// The Go gitlab library has a Title attribute but the api no longer provides it.
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project ID",
		},
		{
			Name:        "action_name",
			Type:        proto.ColumnType_STRING,
			Description: "The action this event tracks: approved, closed, commented on, created, destroyed, expired, joined, left, merged, pushed to, reopened, updated",
		},
		{
			Name:        "target_id",
			Type:        proto.ColumnType_INT,
			Description: "The target ID",
		},
		{
			Name:        "target_iid",
			Type:        proto.ColumnType_INT,
			Description: "The target IID",
			Transform:   transform.FromField("TargetIID").NullIfZero(),
		},
		{
			Name:        "target_type",
			Type:        proto.ColumnType_STRING,
			Description: "What the event was: issue, milestone, merge_request, note, project, snippet, user",
		},
		{
			Name:        "author_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user who created the event",
		},
		{
			Name:        "target_title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the target",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "When the event was created",
		},
		{
			Name:        "push_data",
			Type:        proto.ColumnType_JSON,
			Description: "JSON struct if there's push data",
		},
		{
			Name:        "note",
			Type:        proto.ColumnType_JSON,
			Description: "JSON struct if there's a note",
		},
		{
			Name:        "author",
			Type:        proto.ColumnType_JSON,
			Description: "JSON struct describing the user",
		},
		{
			Name:        "author_username",
			Type:        proto.ColumnType_STRING,
			Description: "author_username",
		},
	}
}
