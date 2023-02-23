package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGroupPushRule() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_push_rule",
		Description: "Obtain information on push rules for a specific group within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    listGroupPushRules,
		},
		Columns: groupPushRuleColumns(),
	}
}

func listGroupPushRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupPushRules", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupPushRules", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())
	plugin.Logger(ctx).Debug("listGroupPushRules", "groupId", groupId)

	pushRules, _, err := conn.Groups.GetGroupPushRules(groupId)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupPushRules", "groupId", groupId, "error", err)
		return nil, fmt.Errorf("unable to obtain push rules for group_id %d\n%v", groupId, err)
	}

	d.StreamListItem(ctx, pushRules)

	plugin.Logger(ctx).Debug("listGroupPushRules", "completed successfully")
	return nil, nil
}

// Column Function
func groupPushRuleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the push rule.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the group push rule was created.",
		},
		{
			Name:        "commit_message_regex",
			Type:        proto.ColumnType_STRING,
			Description: "The regex that a commit message must adhere to.",
		},
		{
			Name:        "commit_message_negative_regex",
			Type:        proto.ColumnType_STRING,
			Description: "The regex that a commit message can not adhere to.",
		},
		{
			Name:        "branch_name_regex",
			Type:        proto.ColumnType_STRING,
			Description: "The regex that a branch name must adhere to.",
		},
		{
			Name:        "deny_delete_tag",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if tag deletion will be denied.",
		},
		{
			Name:        "member_check",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if member checks are performed.",
		},
		{
			Name:        "prevent_secrets",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if push should be denied if contains secrets.",
		},
		{
			Name:        "author_email_regex",
			Type:        proto.ColumnType_STRING,
			Description: "The regex that commit authors email address must adhere to.",
		},
		{
			Name:        "file_name_regex",
			Type:        proto.ColumnType_STRING,
			Description: "The regex that file names must not adhere to.",
		},
		{
			Name:        "max_file_size",
			Type:        proto.ColumnType_INT,
			Description: "Length of maximum file size (MB).",
		},
		{
			Name:        "commit_committer_check",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if committer must have a verified email address.",
		},
		{
			Name:        "reject_unsigned_commits",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if commits not signed by GPG will be rejected.",
		},
		{
			Name:        "group_id",
			Type:        proto.ColumnType_INT,
			Description: "The group id - link to gitlab_group.id`.",
			Transform:   transform.FromQual("group_id"),
		},
	}
}
