package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectHook() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_hook",
		Description: "Obtain information about the hooks for a specific project in the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectHooks,
		},
		Columns: projectHookColumns(),
	}
}

// Hydrate Functions
func listProjectHooks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectHooks", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectHooks", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectHooksOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectHooks", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		hooks, resp, err := conn.Projects.ListProjectHooks(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectHooks", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain hooks for project_id %d\n%v", projectId, err)
		}

		for _, hook := range hooks {
			d.StreamListItem(ctx, hook)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectHooks", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectHooks", "completed successfully")
	return nil, nil
}

// Column Function
func projectHookColumns() []*plugin.Column {
	return []*plugin.Column{
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
			Name:        "confidential_note_events",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if confidential note events will be sent to the hook.",
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
			Name:        "enable_ssl_verification",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if SSL verification is enabled for the hook.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
		},
	}
}
