package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_issue",
		Description: "Obtain information about issues with the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listIssues,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "assignee", Require: plugin.Optional},
				{Name: "assignee_id", Require: plugin.Optional},
				{Name: "author_id", Require: plugin.Optional},
				{Name: "confidential", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
			},
		},
		Columns: issueColumns(),
	}
}

// Hydrate Functions
func listIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.EqualsQuals
	if q["assignee"] == nil &&
		q["assignee_id"] == nil &&
		q["author_id"] == nil &&
		q["project_id"] == nil &&
		isPublicGitLab(d) {
		plugin.Logger(ctx).Error("listIssues", "Public GitLab requires an '=' qualifier for at least one of the following columns 'assignee', 'assignee_id', 'author_id', 'project_id' - none was provided")
		return nil, fmt.Errorf("when using the gitlab_issue table with GitLab Cloud, `List` call requires an '=' qualifier for one or more of the following columns: 'assignee', 'assignee_id', 'author_id', 'project_id'")
	}

	if q["project_id"] != nil {
		plugin.Logger(ctx).Debug("listIssues", "project_id qualifier obtained, re-directing SDK call to ListProjectIssues")
		return listProjectIssues(ctx, d, h)
	}

	return listAllIssues(ctx, d, h)
}

func listProjectIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectIssues", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectIssues", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	if q["project_id"] == nil {
		return nil, nil
	}

	defaultScope := "all"
	opt := &api.ListProjectIssuesOptions{
		Scope: &defaultScope,
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	opt = addOptionalProjectIssueQualifiers(ctx, opt, q)
	projectId := int(q["project_id"].GetInt64Value())

	for {
		plugin.Logger(ctx).Debug("listProjectIssues", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		issues, resp, err := conn.Issues.ListProjectIssues(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectIssues", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain issues for project_id %d\n%v", projectId, err)
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectIssues", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectIssues", "completed successfully")
	return nil, nil
}

func listAllIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listAllIssues", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listAllIssues", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	defaultScope := "all"
	opt := &api.ListIssuesOptions{
		Scope: &defaultScope,
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}
	opt = addOptionalIssueQualifiers(ctx, opt, q)

	for {
		plugin.Logger(ctx).Debug("listAllIssues", "page", opt.Page, "perPage", opt.PerPage)
		issues, resp, err := conn.Issues.ListIssues(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listAllIssues", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain issues\n%v", err)
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listAllIssues", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listAllIssues", "completed successfully")
	return nil, nil
}

// Assist Functions
func addOptionalProjectIssueQualifiers(ctx context.Context, opts *api.ListProjectIssuesOptions, q map[string]*proto.QualValue) *api.ListProjectIssuesOptions {
	if q["assignee"] != nil {
		assignee := q["assignee"].GetStringValue()
		opts.AssigneeUsername = &assignee
		plugin.Logger(ctx).Debug("listProjectIssues", "filter[assignee]", assignee)
	}

	if q["assignee_id"] != nil {
		assigneeId := int(q["assignee_id"].GetInt64Value())
		opts.AssigneeID = &assigneeId
		plugin.Logger(ctx).Debug("listProjectIssues", "filter[assignee_id]", assigneeId)
	}

	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opts.AuthorID = &authorId
		plugin.Logger(ctx).Debug("listProjectIssues", "filter[author_id]", authorId)
	}

	if q["confidential"] != nil {
		confidential := q["confidential"].GetBoolValue()
		opts.Confidential = &confidential
		plugin.Logger(ctx).Debug("listProjectIssues", "filter[confidential]", confidential)
	}

	return opts
}

func addOptionalIssueQualifiers(ctx context.Context, opts *api.ListIssuesOptions, q map[string]*proto.QualValue) *api.ListIssuesOptions {
	if q["assignee"] != nil {
		assignee := q["assignee"].GetStringValue()
		opts.AssigneeUsername = &assignee
		plugin.Logger(ctx).Debug("listAllIssues", "filter[assignee]", assignee)
	}

	if q["assignee_id"] != nil {
		assigneeId := int(q["assignee_id"].GetInt64Value())
		opts.AssigneeID = api.AssigneeID(assigneeId)
		plugin.Logger(ctx).Debug("listAllIssues", "filter[assignee_id]", assigneeId)
	}

	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opts.AuthorID = &authorId
		plugin.Logger(ctx).Debug("listAllIssues", "filter[author_id]", authorId)
	}

	if q["confidential"] != nil {
		confidential := q["confidential"].GetBoolValue()
		opts.Confidential = &confidential
		plugin.Logger(ctx).Debug("listAllIssues", "filter[confidential]", confidential)
	}

	return opts
}

// Transform Functions
func parseAssignees(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return nil, nil
	}

	assignees := input.Value.([]*api.IssueAssignee)
	var output []string

	for _, assignee := range assignees {
		output = append(output, assignee.Username)
	}

	return output, nil
}

// Column Functions
func issueColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the Issue.",
		},
		{
			Name:        "iid",
			Type:        proto.ColumnType_INT,
			Description: "The instance ID of the Issue.",
			Transform:   transform.FromField("IID"),
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the Issue.",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Description: "The description of the Issue.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the Issue (opened, closed, etc).",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project - link to `gitlab_project.id`.",
		},
		{
			Name:        "external_id",
			Type:        proto.ColumnType_STRING,
			Description: "The external ID of the issue.",
		},
		{
			Name:        "author_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the author - link to `gitlab_user.id`.",
			Transform:   transform.FromField("Author.ID"),
		},
		{
			Name:        "author",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the author - link to `gitlab_user.username`.",
			Transform:   transform.FromField("Author.Username"),
		},
		{
			Name:        "author_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the author",
			Transform:   transform.FromField("Author.Name"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of issue creation.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of last update to the issue.",
		},
		{
			Name:        "closed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when issue was closed. (null if not closed).",
		},
		{
			Name:        "closed_by_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user whom closed the issue - link to `gitlab_user.id`.",
			Transform:   transform.FromField("ClosedBy.ID"),
		},
		{
			Name:        "closed_by",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user whom closed the issue - link to `gitlab_user.username`.",
			Transform:   transform.FromField("ClosedBy.Username"),
		},
		{
			Name:        "assignee_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user assigned to the issue - link to `gitlab_user.id`.",
			Transform:   transform.FromField("Assignee.ID"),
		},
		{
			Name:        "assignee",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user assigned to the issue - link to `gitlab_user.username`",
			Transform:   transform.FromField("Assignee.Username"),
		},
		{
			Name:        "assignees",
			Type:        proto.ColumnType_JSON,
			Description: "An array of assigned usernames, for when more than one user is assigned.",
			Transform:   transform.FromField("Assignees").Transform(parseAssignees),
		},
		{
			Name:        "upvotes",
			Type:        proto.ColumnType_INT,
			Description: "Count of up-votes received on the issue.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "downvotes",
			Type:        proto.ColumnType_INT,
			Description: "Count of down-votes received on the issue.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "due_date",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of due date for the issue to be completed by.",
			Transform:   transform.FromField("DueDate").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to access the issue.",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "confidential",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the issue is marked as confidential.",
		},
		{
			Name:        "discussion_locked",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the issue has the discussions locked against new input.",
		},
		{
			Name:        "weight",
			Type:        proto.ColumnType_INT,
			Description: "The weight assigned to the issue.",
		},
		{
			Name:        "issue_type",
			Type:        proto.ColumnType_STRING,
			Description: "The type of issue.",
		},
		{
			Name:        "subscribed",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if current user is subscribed to the issue.",
		},
		{
			Name:        "user_notes_count",
			Type:        proto.ColumnType_INT,
			Description: "Count of user notes on the issue.",
		},
		{
			Name:        "merge_requests_count",
			Type:        proto.ColumnType_INT,
			Description: "Count of merge requests associated with the issue.",
		},
		// Milestone
		{
			Name:        "milestone_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the milestone the issues is placed into.",
			Transform:   transform.FromField("Milestone.ID"),
		},
		{
			Name:        "milestone_iid",
			Type:        proto.ColumnType_INT,
			Description: "The instance id of the milestone",
			Transform:   transform.FromField("Milestone.IID"),
		},
		{
			Name:        "milestone_title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the milestone.",
			Transform:   transform.FromField("Milestone.Title"),
		},
		{
			Name:        "milestone_description",
			Type:        proto.ColumnType_STRING,
			Description: "The description of the milestone.",
			Transform:   transform.FromField("Milestone.Description"),
		},
		{
			Name:        "milestone_created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp at which the milestone was created.",
			Transform:   transform.FromField("Milestone.CreatedAt"),
		},
		{
			Name:        "milestone_updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp at which the milestone was updated.",
			Transform:   transform.FromField("Milestone.UpdatedAt"),
		},
		{
			Name:        "milestone_start_date",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the milestone was started.",
			Transform:   transform.FromField("Milestone.StartDate").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "milestone_due_date",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of due date for the milestone to be completed by.",
			Transform:   transform.FromField("Milestone.DueDate").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "milestone_state",
			Type:        proto.ColumnType_STRING,
			Description: "The current state of the milestone.",
			Transform:   transform.FromField("Milestone.State"),
		},
		{
			Name:        "milestone_expired",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the milestone is expired.",
			Transform:   transform.FromField("Milestone.Expired"),
		},
		// Labels
		{
			Name:        "labels",
			Type:        proto.ColumnType_JSON,
			Description: "An array of strings for the textual labels applied to the issue.",
		},
		// Refs
		{
			Name:        "short_ref",
			Type:        proto.ColumnType_STRING,
			Description: "Short reference of the issue.",
			Transform:   transform.FromField("References.Short"),
		},
		{
			Name:        "rel_ref",
			Type:        proto.ColumnType_STRING,
			Description: "Relative reference of the issue.",
			Transform:   transform.FromField("References.Relative"),
		},
		{
			Name:        "full_ref",
			Type:        proto.ColumnType_STRING,
			Description: "Full reference of the issue.",
			Transform:   transform.FromField("References.Full"),
		},
		// Time Stats
		{
			Name:        "time_estimate",
			Type:        proto.ColumnType_INT,
			Description: "Time estimated against the issue.",
			Transform:   transform.FromField("TimeStats.TimeEstimate"),
		},
		{
			Name:        "total_time_spent",
			Type:        proto.ColumnType_INT,
			Description: "Total time spent on the issue.",
			Transform:   transform.FromField("TimeStats.TotalTimeSpent"),
		},
		// IDs
		{
			Name:        "issue_link_id",
			Type:        proto.ColumnType_INT,
			Description: "Issue link id.",
		},
		{
			Name:        "epic_issue_id",
			Type:        proto.ColumnType_INT,
			Description: "Epic issue id.",
		},
		// Epic Fields (Not all on SDK object are turned by this API call.)
		{
			Name:        "epic_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the associated epic.",
			Transform:   transform.FromField("Epic.ID"),
		},
		{
			Name:        "epic_iid",
			Type:        proto.ColumnType_INT,
			Description: "The IID of the associated epic.",
			Transform:   transform.FromField("Epic.IID"),
		},
		{
			Name:        "epic_title",
			Type:        proto.ColumnType_STRING,
			Description: "Title of the associated epic.",
			Transform:   transform.FromField("Epic.Title"),
		},
		{
			Name:        "epic_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the associated epic.",
			Transform:   transform.FromField("Epic.URL"),
		},
		{
			Name:        "epic_group_id",
			Type:        proto.ColumnType_INT,
			Description: "The group ID of the associated epic.",
			Transform:   transform.FromField("Epic.GroupID"),
		},
	}
}
