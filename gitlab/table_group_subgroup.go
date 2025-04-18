package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableGroupSubgroup() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_subgroup",
		Description: "Obtain information about subgroups for a specific group within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate:    listGroupSubgroups,
			KeyColumns: plugin.SingleColumn("parent_id"),
		},
		Columns: groupColumns(),
	}
}

// Hydrate Functions
func listGroupSubgroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupSubgroups", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupSubgroups", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["parent_id"].GetInt64Value())
	stats := true
	opt := &api.ListSubGroupsOptions{Statistics: &stats, ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		plugin.Logger(ctx).Debug("listGroupSubgroups", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)
		groups, resp, err := conn.Groups.ListSubGroups(groupId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroupSubgroups", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain subgroups for group_id %d\n%v", groupId, err)
		}

		for _, group := range groups {
			d.StreamListItem(ctx, group)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listGroupSubgroups", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroupSubgroups", "completed successfully")
	return nil, nil
}

// Column Function
func groupSubgroupColumns() []*plugin.Column {
	cols := groupColumns()
	gic := plugin.Column{
		Name:        "group_id",
		Type:        proto.ColumnType_INT,
		Description: "Group ID",
		Transform:   transform.FromQual("group_id"),
	}
	return append(cols, &gic)
}
