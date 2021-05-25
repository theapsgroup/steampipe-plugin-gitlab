package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableMergeRequest() *plugin.Table {
	return &plugin.Table{
		Name: "gitlab_merge_request",
		Description: "All GitLab Merge Requests",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"iid", "project_id"}),
			Hydrate: getMergeRequest,
		},
		List: &plugin.ListConfig{
			Hydrate: listMergeRequests,
		},
		Columns: gitlabMergeRequestColumns(),
	}
}

func getMergeRequest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	q := d.KeyColumnQuals
	iid := int(q["iid"].GetInt64Value())
	projectId := int(q["project_id"].GetInt64Value())

	mergeRequest, _, err := conn.MergeRequests.GetMergeRequest(projectId, iid, &api.GetMergeRequestsOptions{})
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func listMergeRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	defaultScope := "all"

	opt := &api.ListMergeRequestsOptions{
		Scope: &defaultScope,
		ListOptions: api.ListOptions{
			Page: 1,
			PerPage: 30,
		},
	}

	for {
		mergeRequests, response, err := conn.MergeRequests.ListMergeRequests(opt)
		if err != nil {
			return nil, err
		}

		for _, mergeRequest := range mergeRequests {
			d.StreamListItem(ctx, mergeRequest)
		}

		if response.CurrentPage >= response.TotalPages {
			break
		}

		opt.Page = response.NextPage
	}

	return nil, nil
}

// Transform Functions
func parseBasicUserCollection(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	var output []string
	if input.Value == nil {
		return nil, nil
	}

	users := input.Value.([]*api.BasicUser)

	for _, user := range users {
		output = append(output, user.Username)
	}

	return output, nil
}

// Column Functions
func gitlabMergeRequestColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "id", Type: proto.ColumnType_INT, Description: "The global ID for the merge request."},
		{Name: "iid", Type: proto.ColumnType_INT, Description: "The internal ID to the project for the merge request", Transform: transform.FromField("IID").NullIfZero()},
		{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project containing the merge request - link to `gitlab_project.ID`"},
		{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the merge request."},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the merge request."},
		{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the merge request (open, closed, merged)"},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of merge request creation."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of last update to the issue."},
		{Name: "merged_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when merge request was merged."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when merge request was closed."},
		{Name: "target_branch", Type: proto.ColumnType_STRING, Description: "The target branch for the merge request."},
		{Name: "source_branch", Type: proto.ColumnType_STRING, Description: "The source branch of the merge request."},
		{Name: "author_id", Type: proto.ColumnType_INT, Description: "The ID of the author - link to `gitlab_user.id`.", Transform: transform.FromField("Author.ID").NullIfZero()},
		{Name: "author_username", Type: proto.ColumnType_STRING, Description: "The username of the author - link to `gitlab_user.username`.", Transform: transform.FromField("Author.Username").NullIfZero()},
		{Name: "upvotes", Type: proto.ColumnType_INT, Description: "Count of up-votes received on the merge request.", Transform: transform.FromGo()},
		{Name: "downvotes", Type: proto.ColumnType_INT, Description: "Count of down-votes received on the issue.", Transform: transform.FromGo()},
		{Name: "assignee_id", Type: proto.ColumnType_INT, Description: "The ID of the assignee - link to `gitlab_user.id`.", Transform: transform.FromField("Assignee.ID").NullIfZero()},
		{Name: "assignee_username", Type: proto.ColumnType_STRING, Description: "The username of the assignee - link to `gitlab_user.username`.", Transform: transform.FromField("Assignee.Username").NullIfZero()},
		{Name: "assignees", Type: proto.ColumnType_JSON, Description: "An array of assigned usernames, for when more than one user is assigned.", Transform: transform.FromField("Assignees").NullIfZero().Transform(parseBasicUserCollection)},
		{Name: "reviewers", Type: proto.ColumnType_JSON, Description: "An array of usernames who've been asked to review the merge request.", Transform: transform.FromField("Reviewers").NullIfZero().Transform(parseBasicUserCollection)},
		{Name: "work_in_progress", Type: proto.ColumnType_BOOL, Description: "Indicates if the merge request is a work in progress."},
		{Name: "merge_when_pipeline_succeeds", Type: proto.ColumnType_BOOL, Description: "Indicates if the merge request will be merged upon completion of CI/CD pipeline."},
		{Name: "merge_status", Type: proto.ColumnType_STRING, Description: "Descriptive status about the ability of being able to merge the merge request."},
		{Name: "merge_error", Type: proto.ColumnType_STRING, Description: "Error message if the merge request can not be merged."},
		{Name: "merged_by_id", Type: proto.ColumnType_INT, Description: "The ID of the user who merged the merge request - link to `gitlab_user.id`.", Transform: transform.FromField("MergedBy.ID").NullIfZero()},
		{Name: "merged_by_username", Type: proto.ColumnType_STRING, Description: "The username of the user who merged the merge request - link to `gitlab_user.username`.", Transform: transform.FromField("MergedBy.Username").NullIfZero()},
		{Name: "closed_by_id", Type: proto.ColumnType_INT, Description: "The ID of the user who closed the merge request - link to `gitlab_user.id`.", Transform: transform.FromField("ClosedBy.ID").NullIfZero()},
		{Name: "closed_by_username", Type: proto.ColumnType_STRING, Description: "The username of the user who closed the merge request - link to `gitlab_user.username`.", Transform: transform.FromField("ClosedBy.Username").NullIfZero()},
		{Name: "subscribed", Type: proto.ColumnType_BOOL, Description: "Indicates if the user associated to the token used to access the data is subscribed to the merge request."},
		{Name: "sha", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("SHA")},
		{Name: "merge_commit_sha", Type: proto.ColumnType_STRING, Description: "The hash of the merge commit.", Transform: transform.FromField("MergeCommitSHA")},
		{Name: "squash_commit_sha", Type: proto.ColumnType_STRING, Description: "The hash of the squashed merge commit.", Transform: transform.FromField("SquashCommitSHA")},
		{Name: "user_notes_count", Type: proto.ColumnType_INT, Description: "A count of user notes on the merge request.", Transform: transform.FromGo()},
		{Name: "changes_count", Type: proto.ColumnType_INT, Description: "A count of changes contained within the merge request."},
		{Name: "should_remove_source_branch", Type: proto.ColumnType_BOOL, Description: "Indicates if source_branch should be deleted on merge."},
		{Name: "force_remove_source_branch", Type: proto.ColumnType_BOOL, Description: "Indicates if source_branch will be force deleted on merge."},
		{Name: "allow_collaboration", Type: proto.ColumnType_BOOL, Description: "Indicates if collaboration is allowed on the merge request."},
		{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to access the merge request.", Transform: transform.FromField("WebURL")},
		{Name: "discussion_locked", Type: proto.ColumnType_BOOL, Description: "Indicates if the merge request has the discussions locked against new input."},
		{Name: "squash", Type: proto.ColumnType_BOOL, Description: "Indicates if a squash is requested."},
		{Name: "diverged_commits_count", Type: proto.ColumnType_INT, Description: "A count of commits diverged from target_branch.", Transform: transform.FromGo()},
		{Name: "rebase_in_progress", Type: proto.ColumnType_BOOL, Description: "Indicates if a rebase is in progress."},
		{Name: "approvals_before_merge", Type: proto.ColumnType_INT, Description: "The number of approvals required before merge can proceed.", Transform: transform.FromGo()},
		{Name: "reference", Type: proto.ColumnType_STRING, Description: "The reference code of the merge request (example: `!4`."},
		{Name: "has_conflicts", Type: proto.ColumnType_BOOL, Description: "Indicates if the merge request has conflicts with the target_branch."},
	}
}