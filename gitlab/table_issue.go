package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableIssue() *plugin.Table {
	return &plugin.Table{
		Name: "gitlab_issue",
		Description: "GitLab Issues",
		List: &plugin.ListConfig{
			Hydrate: listIssues,
		},
		Columns: []*plugin.Column{
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
		},
	}
}

func listIssues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	defaultScope := "all"

	opt := &api.ListIssuesOptions{
		Scope: &defaultScope,
		ListOptions: api.ListOptions{
			Page: 1,
			PerPage: 20,
		},
	}

	for {
		issues, resp, err := conn.Issues.ListIssues(opt)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			d.StreamListItem(ctx, issue)
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
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