package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableGroupProject() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_project",
		Description: "Obtain information about the project(s) that reside within a group.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    listGroupProjects,
		},
		Columns: groupProjectColumns(),
	}
}

// Hydrate Functions
func listGroupProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupProjects", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupProjects", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())
	includeSubGroups := true
	opt := &api.ListGroupProjectsOptions{
		IncludeSubGroups: &includeSubGroups,
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		plugin.Logger(ctx).Debug("listGroupProjects", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)
		groups, resp, err := conn.Groups.ListGroupProjects(groupId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroupProjects", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain projects for group_id %d\n%v", groupId, err)
		}

		for _, group := range groups {
			d.StreamListItem(ctx, group)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listGroupProjects", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroupProjects", "completed successfully")
	return nil, nil
}

// Column Function
func groupProjectColumns() []*plugin.Column {
	cols := projectColumns()
	gic := plugin.Column{
		Name:        "group_id",
		Type:        proto.ColumnType_INT,
		Description: "Group ID",
		Transform:   transform.FromQual("group_id"),
	}
	return append(cols, &gic)
}
