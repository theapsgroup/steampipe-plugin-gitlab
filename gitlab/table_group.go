package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"strings"
)

func tableGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group",
		Description: "Groups within GitLab",
		List: &plugin.ListConfig{
			Hydrate: listGroups,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getGroup,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The group name."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The group path."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The groups description."},
			{Name: "membership_lock", Type: proto.ColumnType_BOOL, Description: "Indicates if membership of the group is locked."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The groups visibility (private/internal/public)"},
			{Name: "lfs_enabled", Type: proto.ColumnType_BOOL, Description: "Does the group have Large File System enabled.", Transform: transform.FromField("LFSEnabled")},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The url for the groups avatar."},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url for the group."},
			{Name: "request_access_enabled", Type: proto.ColumnType_BOOL, Description: "Indicates if the group allows access requests."},
			{Name: "full_name", Type: proto.ColumnType_STRING, Description: "The full name of the group."},
			{Name: "full_path", Type: proto.ColumnType_STRING, Description: "The full path of the group"},
			{Name: "parent_id", Type: proto.ColumnType_INT, Description: "The ID of the groups parent group (for sub-groups)"},
			{Name: "custom_attributes", Type: proto.ColumnType_JSON, Description: "An array of custom attributes."},
			{Name: "share_with_group_lock", Type: proto.ColumnType_BOOL, Description: "Indicates if this group can be shared with other groups"},
			{Name: "require_two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "Indicates if this group requires 2fa.", Transform: transform.FromField("RequireTwoFactorAuth")},
			{Name: "two_factor_grace_period", Type: proto.ColumnType_INT, Description: "The grace period (in hours) for 2fa."},
			{Name: "project_creation_level", Type: proto.ColumnType_STRING, Description: "The level at which project creation is permitted developer/maintainer/owner"},
			{Name: "auto_devops_enabled", Type: proto.ColumnType_BOOL, Description: "Indicates if the group has auto devops enabled."},
			{Name: "subgroup_creation_level", Type: proto.ColumnType_STRING, Description: "The level at which sub-group creation is permitted developer/maintainer/owner", Transform: transform.FromField("SubGroupCreationLevel")},
			{Name: "emails_disabled", Type: proto.ColumnType_BOOL, Description: "Indicates if this group has email notifications disabled."},
			{Name: "mentions_disabled", Type: proto.ColumnType_BOOL, Description: "Indicates if this group has mention notifications disabled."},
			{Name: "runners_token", Type: proto.ColumnType_STRING, Description: "The groups runner token."},
			{Name: "ldap_cn", Type: proto.ColumnType_STRING, Description: "The LDAP CN associated with group.", Transform: transform.FromField("LDAPCN")},
			{Name: "ldap_access", Type: proto.ColumnType_INT, Description: "The LDAP Access associated with group.", Transform: transform.FromField("LDAPAccess")},
			{Name: "ldap_group_links", Type: proto.ColumnType_JSON, Description: "The LDAP groups linked to the group.", Transform: transform.FromField("LDAPGroupLinks")},
			{Name: "shared_runners_minutes_limit", Type: proto.ColumnType_INT, Description: "The limit in minutes of time the group can utilise shared runner resources."},
			{Name: "extra_shared_runners_minutes_limit", Type: proto.ColumnType_INT, Description: "The limit in minutes of extra time the group can utilise shared runner resources."},
			{Name: "marked_for_deletion_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp for when the group was marked to be deleted.", Transform: transform.FromField("MarkedForDeletionOn").NullIfZero().Transform(isoTimeTransform)},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp for when the group was created."},
		},
	}
}

func listGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opt := &api.ListGroupsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 30,
	}}

	for {
		groups, resp, err := conn.Groups.ListGroups(opt)
		if err != nil {
			return nil, err
		}

		for _, group := range groups {
			d.StreamListItem(ctx, group)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

func getGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	groupId := int(d.KeyColumnQuals["id"].GetInt64Value())

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	group, _, err := conn.Groups.GetGroup(groupId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}

	return group, nil
}
