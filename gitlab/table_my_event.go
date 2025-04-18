package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "gitlab.com/gitlab-org/api/client-go"
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
		Columns: eventColumns(),
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
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listMyEvents", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listMyEvents", "completed successfully")
	return nil, nil
}
