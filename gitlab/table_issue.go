package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_issue",
		Description: "All GitLab Issues",
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
	q := d.KeyColumnQuals

	if q["assignee"] == nil &&
		q["assignee_id"] == nil &&
		q["author_id"] == nil &&
		q["project_id"] == nil &&
		isPublicGitLab(d) {
		return nil, fmt.Errorf("when using the gitlab_issue table with GitLab Cloud, `List` call requires an '=' qualifier for one or more of the following columns: 'assignee', 'assignee_id', 'author_id', 'project_id'")
	}

	if q["project_id"] != nil {
		return listProjectIssues(ctx, d, h)
	}

	return listAllIssues(ctx, d, h)
}

func listProjectIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	q := d.KeyColumnQuals
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

	opt = addOptionalProjectIssueQualifiers(opt, q)
	projectId := int(q["project_id"].GetInt64Value())

	for {
		issues, resp, err := conn.Issues.ListProjectIssues(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

func listAllIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	q := d.KeyColumnQuals
	defaultScope := "all"
	opt := &api.ListIssuesOptions{
		Scope: &defaultScope,
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}
	opt = addOptionalIssueQualifiers(opt, q)

	for {
		issues, resp, err := conn.Issues.ListIssues(opt)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return nil, nil
}

// Assist Functions
func addOptionalProjectIssueQualifiers(opts *api.ListProjectIssuesOptions, q map[string]*proto.QualValue) *api.ListProjectIssuesOptions {
	if q["assignee"] != nil {
		assignee := q["assignee"].GetStringValue()
		opts.AssigneeUsername = &assignee
	}

	if q["assignee_id"] != nil {
		assigneeId := int(q["assignee_id"].GetInt64Value())
		opts.AssigneeID = &assigneeId
	}

	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opts.AuthorID = &authorId
	}

	if q["confidential"] != nil {
		confidential := q["confidential"].GetBoolValue()
		opts.Confidential = &confidential
	}

	return opts
}

func addOptionalIssueQualifiers(opts *api.ListIssuesOptions, q map[string]*proto.QualValue) *api.ListIssuesOptions {
	if q["assignee"] != nil {
		assignee := q["assignee"].GetStringValue()
		opts.AssigneeUsername = &assignee
	}

	if q["assignee_id"] != nil {
		assigneeId := int(q["assignee_id"].GetInt64Value())
		opts.AssigneeID = &assigneeId
	}

	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opts.AuthorID = &authorId
	}

	if q["confidential"] != nil {
		confidential := q["confidential"].GetBoolValue()
		opts.Confidential = &confidential
	}

	return opts
}

// Transform Functions
func parseAssignees(ctx context.Context, input *transform.TransformData) (interface{}, error) {
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
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the Issue."},
		{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the Issue."},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the Issue."},
		{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the Issue (opened, closed, etc)."},
		{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project - link to `gitlab_project.id`."},
		{Name: "external_id", Type: proto.ColumnType_STRING, Description: "The external ID of the issue."},
		{Name: "author_id", Type: proto.ColumnType_INT, Description: "The ID of the author - link to `gitlab_user.id`.", Transform: transform.FromField("Author.ID")},
		{Name: "author", Type: proto.ColumnType_STRING, Description: "The username of the author - link to `gitlab_user.username`.", Transform: transform.FromField("Author.Username")},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of issue creation."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of last update to the issue."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when issue was closed. (null if not closed)."},
		{Name: "closed_by_id", Type: proto.ColumnType_INT, Description: "The ID of the user whom closed the issue - link to `gitlab_user.id`.", Transform: transform.FromField("ClosedBy.ID")},
		{Name: "closed_by", Type: proto.ColumnType_STRING, Description: "The username of the user whom closed the issue - link to `gitlab_user.username`.", Transform: transform.FromField("ClosedBy.Username")},
		{Name: "assignee_id", Type: proto.ColumnType_INT, Description: "The ID of the user assigned to the issue - link to `gitlab_user.id`.", Transform: transform.FromField("Assignee.ID")},
		{Name: "assignee", Type: proto.ColumnType_STRING, Description: "The username of the user assigned to the issue - link to `gitlab_user.username`", Transform: transform.FromField("Assignee.Username")},
		{Name: "assignees", Type: proto.ColumnType_JSON, Description: "An array of assigned usernames, for when more than one user is assigned.", Transform: transform.FromField("Assignees").Transform(parseAssignees)},
		{Name: "upvotes", Type: proto.ColumnType_INT, Description: "Count of up-votes received on the issue.", Transform: transform.FromGo()},
		{Name: "downvotes", Type: proto.ColumnType_INT, Description: "Count of down-votes received on the issue.", Transform: transform.FromGo()},
		{Name: "due_date", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of due date for the issue to be completed by.", Transform: transform.FromField("DueDate").NullIfZero().Transform(isoTimeTransform)},
		{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to access the issue.", Transform: transform.FromField("WebURL")},
		{Name: "confidential", Type: proto.ColumnType_BOOL, Description: "Indicates if the issue is marked as confidential."},
		{Name: "discussion_locked", Type: proto.ColumnType_BOOL, Description: "Indicates if the issue has the discussions locked against new input."},
	}
}
