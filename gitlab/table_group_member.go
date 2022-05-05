package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
	api "github.com/xanzy/go-gitlab"
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
		Description: "Group Members for a GitLab Group",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    listGroupMembers,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The id of the group member - link to `gitlab_user.id`"},
			{Name: "username", Type: proto.ColumnType_STRING, Description: "The username of the group member - link to `gitlab_user.username`."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the group member."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the group member active, blocked, etc"},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The url of the group members avatar.", Transform: transform.FromField("AvatarUrl")},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url for profile of the group member.", Transform: transform.FromField("WebUrl")},
			{Name: "expires_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date of expiration from the group for the group member.", Transform: transform.FromField("ExpiresAt").NullIfZero().Transform(isoTimeTransform)},
			{Name: "access_level", Type: proto.ColumnType_INT, Description: "The access level the group member holds within the group.", Transform: transform.FromGo()},
			{Name: "access_desc", Type: proto.ColumnType_STRING, Description: "The descriptive of the access level held by the group member."},
			{Name: "group_id", Type: proto.ColumnType_INT, Description: "The group id - link to gitlab_group.id`."},
		},
	}
}

func listGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	groupId := int(d.KeyColumnQuals["group_id"].GetInt64Value())

	opt := &api.ListGroupMembersOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		members, resp, err := conn.Groups.ListAllGroupMembers(groupId, opt)
		if err != nil {
			return nil, err
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
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
