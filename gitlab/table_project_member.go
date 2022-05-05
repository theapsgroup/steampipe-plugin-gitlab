package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"time"
)

type ProjectMember struct {
	ID          int
	Username    string
	Name        string
	State       string
	CreatedAt   *time.Time
	ExpiresAt   *api.ISOTime
	AccessLevel int
	AccessDesc  string
	WebUrl      string
	AvatarUrl   string
	ProjectID   int
}

func tableProjectMember() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_member",
		Description: "Project Members for a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectMembers,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The id of the project member - link to `gitlab_user.id`"},
			{Name: "username", Type: proto.ColumnType_STRING, Description: "The username of the project member - link to `gitlab_user.username`."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the project member."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the project member active, blocked, etc"},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The url of the project members avatar.", Transform: transform.FromField("AvatarUrl")},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url for profile of the project member.", Transform: transform.FromField("WebUrl")},
			{Name: "expires_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date of expiration of access to the project.", Transform: transform.FromField("ExpiresAt").NullIfZero().Transform(isoTimeTransform)},
			{Name: "access_level", Type: proto.ColumnType_INT, Description: "The access level the project member holds within the project.", Transform: transform.FromGo()},
			{Name: "access_desc", Type: proto.ColumnType_STRING, Description: "The descriptive of the access level held by the project member."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The project id - link to gitlab_project.id`."},
		},
	}
}

func listProjectMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())

	opt := &api.ListProjectMembersOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		members, resp, err := conn.ProjectMembers.ListAllProjectMembers(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, member := range members {
			d.StreamListItem(ctx, &ProjectMember{
				ID:          member.ID,
				Username:    member.Username,
				Name:        member.Name,
				State:       member.State,
				AvatarUrl:   member.AvatarURL,
				WebUrl:      member.WebURL,
				ExpiresAt:   member.ExpiresAt,
				AccessLevel: int(member.AccessLevel),
				AccessDesc:  parseAccessLevel(int(member.AccessLevel)),
				ProjectID:   projectId,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
