package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"strings"
)

func tableGroup() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group",
		Description: "Obtain information about groups within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listGroups,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getGroup,
		},
		Columns: groupColumns(),
	}
}

// Hydrate Functions
func listGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroups", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroups", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	stats := true
	opt := &api.ListGroupsOptions{Statistics: &stats, ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		plugin.Logger(ctx).Debug("listGroups", "page", opt.Page, "perPage", opt.PerPage)

		groups, resp, err := conn.Groups.ListGroups(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroups", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain groups\n%v", err)
		}

		for _, group := range groups {
			d.StreamListItem(ctx, group)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listGroups", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroups", "completed successfully")
	return nil, nil
}

func getGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getGroup", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getGroup", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["id"].GetInt64Value())
	opts := &api.GetGroupOptions{}
	plugin.Logger(ctx).Debug("getGroup", "groupId", groupId)

	group, _, err := conn.Groups.GetGroup(groupId, opts)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			plugin.Logger(ctx).Warn("getGroup", "groupId", groupId, "no group was found, returning empty result set")
			return nil, nil
		}
		plugin.Logger(ctx).Error("getGroup", "groupId", groupId, "error", err)
		return nil, fmt.Errorf("unable to obtain branches for group_id %d\n%v", groupId, err)
	}

	plugin.Logger(ctx).Debug("getGroup", "completed successfully")
	return group, nil
}

// Column Function
func groupColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the group.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The group name.",
		},
		{
			Name:        "path",
			Type:        proto.ColumnType_STRING,
			Description: "The group path.",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Description: "The groups description.",
		},
		{
			Name:        "membership_lock",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if membership of the group is locked.",
		},
		{
			Name:        "visibility",
			Type:        proto.ColumnType_STRING,
			Description: "The groups visibility (private/internal/public)",
		},
		{
			Name:        "lfs_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Does the group have Large File System enabled.",
			Transform:   transform.FromField("LFSEnabled"),
		},
		{
			Name:        "default_branch_protection",
			Type:        proto.ColumnType_INT,
			Description: "Indicates level of protection applied to default branch see: https://docs.gitlab.com/ee/api/groups.html#options-for-default_branch_protection for details.",
		},
		{
			Name:        "avatar_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for the groups avatar.",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for the group.",
		},
		{
			Name:        "request_access_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the group allows access requests.",
		},
		{
			Name:        "full_name",
			Type:        proto.ColumnType_STRING,
			Description: "The full name of the group.",
		},
		{
			Name:        "full_path",
			Type:        proto.ColumnType_STRING,
			Description: "The full path of the group",
		},
		{
			Name:        "file_template_project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project template used (if any).",
		},
		{
			Name:        "parent_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the groups parent group (for sub-groups)",
		},
		{
			Name:        "custom_attributes",
			Type:        proto.ColumnType_JSON,
			Description: "An array of custom attributes.",
		},
		{
			Name:        "share_with_group_lock",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if this group can be shared with other groups",
		},
		{
			Name:        "require_two_factor_authentication",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if this group requires 2fa.",
			Transform:   transform.FromField("RequireTwoFactorAuth"),
		},
		{
			Name:        "two_factor_grace_period",
			Type:        proto.ColumnType_INT,
			Description: "The grace period (in hours) for 2fa.",
		},
		{
			Name:        "project_creation_level",
			Type:        proto.ColumnType_STRING,
			Description: "The level at which project creation is permitted developer/maintainer/owner",
		},
		{
			Name:        "auto_devops_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the group has auto devops enabled.",
		},
		{
			Name:        "subgroup_creation_level",
			Type:        proto.ColumnType_STRING,
			Description: "The level at which sub-group creation is permitted developer/maintainer/owner",
			Transform:   transform.FromField("SubGroupCreationLevel"),
		},
		{
			Name:        "emails_disabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if this group has email notifications disabled.",
		},
		{
			Name:        "mentions_disabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if this group has mention notifications disabled.",
		},
		{
			Name:        "runners_token",
			Type:        proto.ColumnType_STRING,
			Description: "The groups runner token.",
		},
		{
			Name:        "ldap_cn",
			Type:        proto.ColumnType_STRING,
			Description: "The LDAP CN associated with group.",
			Transform:   transform.FromField("LDAPCN"),
		},
		{
			Name:        "ldap_access",
			Type:        proto.ColumnType_INT,
			Description: "The LDAP Access associated with group.",
			Transform:   transform.FromField("LDAPAccess"),
		},
		{
			Name:        "ldap_group_links",
			Type:        proto.ColumnType_JSON,
			Description: "The LDAP groups linked to the group.",
			Transform:   transform.FromField("LDAPGroupLinks"),
		},
		{
			Name:        "shared_runners_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if shared runners are enabled for the group.",
		},
		{
			Name:        "shared_runners_minutes_limit",
			Type:        proto.ColumnType_INT,
			Description: "The limit in minutes of time the group can utilise shared runner resources.",
		},
		{
			Name:        "extra_shared_runners_minutes_limit",
			Type:        proto.ColumnType_INT,
			Description: "The limit in minutes of extra time the group can utilise shared runner resources.",
		},
		{
			Name:        "marked_for_deletion_on",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp for when the group was marked to be deleted.",
			Transform:   transform.FromField("MarkedForDeletionOn").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp for when the group was created.",
		},
		{
			Name:        "prevent_forking_outside_group",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if forking is prevented for projects not within the group.",
		},
		// Group Statistics
		{
			Name:        "commit_count",
			Type:        proto.ColumnType_INT,
			Description: "The number of commits in the group.",
			Transform:   transform.FromField("Statistics.CommitCount"),
		},
		{
			Name:        "storage_size",
			Type:        proto.ColumnType_INT,
			Description: "The storage size of the group on disk.",
			Transform:   transform.FromField("Statistics.StorageSize"),
		},
		{
			Name:        "repository_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of repositories in the group.",
			Transform:   transform.FromField("Statistics.RepositorySize"),
		},
		{
			Name:        "wiki_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of wikis within the group.",
			Transform:   transform.FromField("Statistics.WikiSize"),
		},
		{
			Name:        "lfs_objects_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of LFS objects within the group.",
			Transform:   transform.FromField("Statistics.LFSObjectsSize"),
		},
		{
			Name:        "job_artifacts_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of job artifacts within the group.",
			Transform:   transform.FromField("Statistics.JobArtifactsSize"),
		},
		{
			Name:        "pipeline_artifacts_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of pipeline artifacts within the group.",
			Transform:   transform.FromField("Statistics.PipelineArtifactsSize"),
		},
		{
			Name:        "packages_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of packages within the group.",
			Transform:   transform.FromField("Statistics.PackagesSize"),
		},
		{
			Name:        "snippets_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of snippets within the group.",
			Transform:   transform.FromField("Statistics.SnippetsSize"),
		},
		{
			Name:        "uploads_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of uploads within the group.",
			Transform:   transform.FromField("Statistics.UploadsSize"),
		},
	}
}
