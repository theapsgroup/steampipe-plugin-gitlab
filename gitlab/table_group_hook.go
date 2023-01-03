package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGroupHook() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_hook",
		Description: "Hooks for a GitLab Group",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    listGroupHooks,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "The ID of the hook.",
			},
			{
				Name:        "url",
				Type:        proto.ColumnType_STRING,
				Description: "The url the hook invokes.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Timestamp of when the hook was created.",
			},
			{
				Name:        "push_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if push events will be sent to the hook.",
			},
			{
				Name:        "push_events_branch_filter",
				Type:        proto.ColumnType_STRING,
				Description: "The filter for branches on which to send push events to the hook.",
			},
			{
				Name:        "issues_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if issue events will be sent to the hook.",
			},
			{
				Name:        "confidential_issues_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if confidential issue events will be sent to the hook.",
			},
			{
				Name:        "merge_requests_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if merge request events will be sent to the hook.",
			},
			{
				Name:        "tag_push_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if tag push events will be sent to the hook.",
			},
			{
				Name:        "note_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if note events will be sent to the hook.",
			},
			{
				Name:        "job_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if job events will be sent to the hook.",
			},
			{
				Name:        "pipeline_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if pipeline events will be sent to the hook.",
			},
			{
				Name:        "wiki_page_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if wiki events will be sent to the hook.",
			},
			{
				Name:        "deployment_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if deployment events will be sent to the hook.",
			},
			{
				Name:        "releases_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if release events will be sent to the hook.",
			},
			{
				Name:        "subgroup_events",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if events from sub-groups will be sent to the hook.",
			},
			{
				Name:        "group_id",
				Type:        proto.ColumnType_INT,
				Description: "The group id - link to gitlab_group.id`.",
			},
		},
	}
}

func listGroupHooks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())

	hooks, _, err := conn.Groups.ListGroupHooks(groupId)
	if err != nil {
		return nil, err
	}

	for _, hook := range hooks {
		d.StreamListItem(ctx, hook)
	}

	return nil, nil
}
