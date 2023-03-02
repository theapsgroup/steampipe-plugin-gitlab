package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableGroupAccessRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_access_request",
		Description: "Obtain access requests for a specific group in the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listGroupAccessRequests,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "group_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: groupAccessRequestColumns(),
	}
}

// Hydrate Function
func listGroupAccessRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupAccessRequests", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupAccessRequests", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())
	opt := &api.ListAccessRequestsOptions{
		Page:    1,
		PerPage: 20,
	}

	for {
		plugin.Logger(ctx).Debug("listGroupAccessRequests", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)
		reqs, resp, err := conn.AccessRequests.ListGroupAccessRequests(groupId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroupAccessRequests", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain access requests for group_id %d\n%v", groupId, err)
		}

		for _, req := range reqs {
			d.StreamListItem(ctx, req)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroupAccessRequests", "completed successfully")
	return nil, nil
}

// Column Function
func groupAccessRequestColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the access request.",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user requesting access.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the user requesting access.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the access request.",
		},
		{
			Name:        "access_level",
			Type:        proto.ColumnType_INT,
			Description: "The numeric value of the access level requested by the user.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of access request creation.",
		},
		{
			Name:        "requested_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of access request submission.",
		},
		{
			Name:        "group_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the group - link to `gitlab_group.id",
			Transform:   transform.FromQual("group_id"),
		},
	}
}
