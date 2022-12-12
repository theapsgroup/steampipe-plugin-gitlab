package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableGroupPushRule() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_push_rule",
		Description: "Push Rules for a GitLab Group",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("group_id"),
			Hydrate:    listGroupPushRules,
		},
		Columns: []*plugin.Column{
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
			},
		},
	}
}

func listGroupPushRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	groupId := int(d.KeyColumnQuals["group_id"].GetInt64Value())

	pushRules, _, err := conn.Groups.GetGroupPushRules(groupId)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, pushRules)

	return nil, nil
}
