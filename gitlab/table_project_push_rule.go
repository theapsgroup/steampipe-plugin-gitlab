package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableProjectPushRule() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_push_rule",
		Description: "Obtain information on push rules for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectPushRules,
		},
		Columns: projectPushRuleColumns(),
	}
}

func listProjectPushRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectPushRules", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectPushRules", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	plugin.Logger(ctx).Debug("listProjectPushRules", "projectId", projectId)

	pushRules, _, err := conn.Projects.GetProjectPushRules(projectId)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectPushRules", "projectId", projectId, "error", err)
		return nil, fmt.Errorf("unable to obtain push rules for project_id %d\n%v", projectId, err)
	}

	// A project without configured push rules returns an empty (nil) response.
	if pushRules == nil {
		plugin.Logger(ctx).Debug("listProjectPushRules", "no push rules configured", projectId)
		return nil, nil
	}

	d.StreamListItem(ctx, pushRules)

	plugin.Logger(ctx).Debug("listProjectPushRules", "completed successfully")
	return nil, nil
}

// Column Function
func projectPushRuleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the push rule.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the project push rule was created.",
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
			Description: "Indicates if the commit author is a GitLab member with a verified email.",
		},
		{
			Name:        "prevent_secrets",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if pushes that contain files likely to hold secrets are rejected.",
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
			Description: "Indicates if committers can only push commits to repositories that match their own email.",
		},
		{
			Name:        "reject_unsigned_commits",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if commits not signed by GPG will be rejected.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
