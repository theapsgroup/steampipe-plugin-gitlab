package gitlab

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
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
		Description: "Obtain information about members of a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectMembers,
		},
		Columns: projectMemberColumns(),
	}
}

// Hydrate Functions
func listProjectMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectMembers", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectMembers", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectMembersOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		plugin.Logger(ctx).Debug("listProjectMembers", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		members, resp, err := conn.ProjectMembers.ListAllProjectMembers(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectMembers", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain members for project_id %d\n%v", projectId, err)
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
				CreatedAt:   member.CreatedAt,
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectMembers", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectMembers", "completed successfully")
	return nil, nil
}

// Column Function
func projectMemberColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The id of the project member - link to `gitlab_user.id`",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the project member - link to `gitlab_user.username`.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the project member.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the project member active, blocked, etc",
		},
		{
			Name:        "avatar_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the project members avatar.",
			Transform:   transform.FromField("AvatarUrl"),
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for profile of the project member.",
			Transform:   transform.FromField("WebUrl"),
		},
		{
			Name:        "expires_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The date of expiration of access to the project.",
			Transform:   transform.FromField("ExpiresAt").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "access_level",
			Type:        proto.ColumnType_INT,
			Description: "The access level the project member holds within the project.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "access_desc",
			Type:        proto.ColumnType_STRING,
			Description: "The descriptive of the access level held by the project member.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to gitlab_project.id`.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp at which the user was created.",
		},
	}
}
