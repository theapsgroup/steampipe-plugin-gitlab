// Package gitlab implements gitlab api calls for steampipe.
package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableMyEvents() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_my_event",
		Description: "Obtain information about my events.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
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
			Hydrate: listMyEvents,
		},
		Columns: myEventColumns(),
	}
}

// Hydrate Functions
func listMyEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listMyEvents", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listMyEvents", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

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
		plugin.Logger(ctx).Debug("listMyEvents", "page", opt.Page, "perPage", opt.PerPage)
		events, resp, err := conn.Events.ListCurrentUserContributionEvents(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listMyEvents", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain my events\n%v", err)
		}

		for _, event := range events {
			plugin.Logger(ctx).Debug("listMyEvents", "event", event)
			d.StreamListItem(ctx, event)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listMyEvents", "completed successfully")
	return nil, nil
}

// Column Function
func myEventColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The event ID",
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "Possibly abandoned title of the event",
		},
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
