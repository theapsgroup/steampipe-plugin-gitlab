package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

type GroupMember struct {
	ID          int
	Username    string
	Name        string
	State       string
	AvatarUrl   string
	WebUrl      string
	ExpiresAt   *api.ISOTime
	AccessLevel int
	AccessDesc  string
	GroupID     int
}

func tableGroupMember() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_member",
		Description: "Obtain information about members of a specific group within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    listGroupMembers,
		},
		Columns: groupMemberColumns(),
	}
}

// Hydrate Functions
func listGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupMembers", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupMembers", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())
	opt := &api.ListGroupMembersOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		plugin.Logger(ctx).Debug("listGroupMembers", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)
		members, resp, err := conn.Groups.ListAllGroupMembers(groupId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroupMembers", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain members for group_id %d\n%v", groupId, err)
		}

		for _, member := range members {
			d.StreamListItem(ctx, &GroupMember{
				ID:          member.ID,
				Username:    member.Username,
				Name:        member.Name,
				State:       member.State,
				AvatarUrl:   member.AvatarURL,
				WebUrl:      member.WebURL,
				ExpiresAt:   member.ExpiresAt,
				AccessLevel: int(member.AccessLevel),
				AccessDesc:  parseAccessLevel(int(member.AccessLevel)),
				GroupID:     groupId,
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listGroupMembers", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroupMembers", "completed successfully")
	return nil, nil
}

// Column Functions
func groupMemberColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The id of the group member - link to `gitlab_user.id`",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the group member - link to `gitlab_user.username`.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the group member.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the group member active, blocked, etc",
		},
		{
			Name:        "avatar_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the group members avatar.",
			Transform:   transform.FromField("AvatarUrl"),
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for profile of the group member.",
			Transform:   transform.FromField("WebUrl"),
		},
		{
			Name:        "expires_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The date of expiration from the group for the group member.",
			Transform:   transform.FromField("ExpiresAt").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "access_level",
			Type:        proto.ColumnType_INT,
			Description: "The access level the group member holds within the group.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "access_desc",
			Type:        proto.ColumnType_STRING,
			Description: "The descriptive of the access level held by the group member.",
		},
		{
			Name:        "group_id",
			Type:        proto.ColumnType_INT,
			Description: "The group id - link to gitlab_group.id`.",
		},
	}
}
