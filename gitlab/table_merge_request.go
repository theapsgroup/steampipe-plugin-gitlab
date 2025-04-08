package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableMergeRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_merge_request",
		Description: "Obtain information about merge requests within the GitLab instance.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"iid", "project_id"}),
			Hydrate:    getMergeRequest,
		},
		List: &plugin.ListConfig{
			Hydrate: listMergeRequests,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project_id", Require: plugin.Optional},
				{Name: "author_id", Require: plugin.Optional},
				{Name: "assignee_id", Require: plugin.Optional},
				{Name: "reviewer_id", Require: plugin.Optional},
			},
		},
		Columns: gitlabMergeRequestColumns(),
	}
}

// Hydrate Functions
func getMergeRequest(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getMergeRequest", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getMergeRequest", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	iid := int(q["iid"].GetInt64Value())
	projectId := int(q["project_id"].GetInt64Value())
	plugin.Logger(ctx).Debug("getMergeRequest", "projectId", projectId, "iid", iid)

	mergeRequest, _, err := conn.MergeRequests.GetMergeRequest(projectId, iid, &api.GetMergeRequestsOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("getMergeRequest", "projectId", projectId, "iid", iid, "error", err)
		return nil, fmt.Errorf("unable to obtain merge request %d for project_id %d\n%v", iid, projectId, err)
	}

	plugin.Logger(ctx).Debug("getMergeRequest", "completed successfully")
	return mergeRequest, nil
}

func listMergeRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.EqualsQuals

	if q["project_id"] == nil &&
		q["assignee_id"] == nil &&
		q["author_id"] == nil &&
		q["reviewer_id"] == nil &&
		isPublicGitLab(d) {
		plugin.Logger(ctx).Error("listMergeRequests", "Public GitLab requires an '=' qualifier for at least one of the following columns 'reviewer_id', 'assignee_id', 'author_id', 'project_id' - none was provided")
		return nil, fmt.Errorf("when using the gitlab_merge_request table with GitLab Cloud, `List`" +
			"call requires an '=' qualifier for one or more of the following columns: 'project_id', 'author_id', 'assignee_id', 'reviewer_id'")
	}

	if q["project_id"] != nil {
		plugin.Logger(ctx).Debug("listMergeRequests", "project_id qualifier obtained, re-directing SDK call to ListProjectMergeRequests")
		return listProjectMergeRequests(ctx, d, h)
	}

	return listAllMergeRequests(ctx, d, h)
}

func listProjectMergeRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectMergeRequests", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	q := d.EqualsQuals

	projectId := int(q["project_id"].GetInt64Value())

	opt := &api.ListProjectMergeRequestsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	if q["assignee_id"] != nil {
		assigneeId := api.AssigneeID(q["assignee_id"].GetInt64Value())
		opt.AssigneeID = assigneeId
		plugin.Logger(ctx).Debug("listProjectMergeRequests", "filter[assignee_id]", assigneeId)
	}

	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opt.AuthorID = &authorId
		plugin.Logger(ctx).Debug("listProjectMergeRequests", "filter[author_id]", authorId)
	}

	if q["reviewer_id"] != nil {
		reviewerId := api.ReviewerID(q["reviewer_id"].GetInt64Value())
		opt.ReviewerID = reviewerId
		plugin.Logger(ctx).Debug("listProjectMergeRequests", "filter[reviewer_id]", reviewerId)
	}

	for {
		plugin.Logger(ctx).Debug("listProjectMergeRequests", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		mergeRequests, response, err := conn.MergeRequests.ListProjectMergeRequests(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectMergeRequests", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain merge requests for project_id %d\n%v", projectId, err)
		}

		for _, mergeRequest := range mergeRequests {
			d.StreamListItem(ctx, mergeRequest)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectMergeRequests", "completed successfully")
				return nil, nil
			}
		}

		if response.NextPage == 0 {
			break
		}

		opt.Page = response.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectMergeRequests", "completed successfully")
	return nil, nil
}

func listAllMergeRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listAllMergeRequests", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listAllMergeRequests", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	defaultScope := "all"

	opt := &api.ListMergeRequestsOptions{
		Scope: &defaultScope,
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	if q["assignee_id"] != nil {
		assigneeId := api.AssigneeID(q["assignee_id"].GetInt64Value())
		opt.AssigneeID = assigneeId
		plugin.Logger(ctx).Debug("listAllMergeRequests", "filter[assignee_id]", assigneeId)
	}

	if q["author_id"] != nil {
		authorId := int(q["author_id"].GetInt64Value())
		opt.AuthorID = &authorId
		plugin.Logger(ctx).Debug("listAllMergeRequests", "filter[author_id]", authorId)
	}

	if q["reviewer_id"] != nil {
		reviewerId := api.ReviewerID(q["reviewer_id"].GetInt64Value())
		opt.ReviewerID = reviewerId
		plugin.Logger(ctx).Debug("listAllMergeRequests", "filter[reviewer_id]", reviewerId)
	}

	for {
		plugin.Logger(ctx).Debug("listAllMergeRequests", "page", opt.Page, "perPage", opt.PerPage)
		mergeRequests, response, err := conn.MergeRequests.ListMergeRequests(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listAllMergeRequests", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain merge requests\n%v", err)
		}

		for _, mergeRequest := range mergeRequests {
			d.StreamListItem(ctx, mergeRequest)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listAllMergeRequests", "completed successfully")
				return nil, nil
			}
		}

		if response.NextPage == 0 {
			break
		}

		opt.Page = response.NextPage
	}

	plugin.Logger(ctx).Debug("listAllMergeRequests", "completed successfully")
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
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The global ID for the merge request.",
		},
		{
			Name:        "iid",
			Type:        proto.ColumnType_INT,
			Description: "The internal ID to the project for the merge request",
			Transform:   transform.FromField("IID").NullIfZero(),
		},
		{
			Name:        "target_branch",
			Type:        proto.ColumnType_STRING,
			Description: "The target branch for the merge request.",
		},
		{
			Name:        "source_branch",
			Type:        proto.ColumnType_STRING,
			Description: "The source branch of the merge request.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project containing the merge request - link to `gitlab_project.ID`",
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the merge request.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the merge request (open, closed, merged)",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of merge request creation.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of last update to the issue.",
		},
		{
			Name:        "upvotes",
			Type:        proto.ColumnType_INT,
			Description: "Count of up-votes received on the merge request.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "downvotes",
			Type:        proto.ColumnType_INT,
			Description: "Count of down-votes received on the issue.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "author_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the author - link to `gitlab_user.id`.",
			Transform:   transform.FromField("Author.ID"),
		},
		{
			Name:        "author_username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the author - link to `gitlab_user.username`.",
			Transform:   transform.FromField("Author.Username"),
		},
		{
			Name:        "author_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the author.",
			Transform:   transform.FromField("Author.Name"),
		},
		{
			Name:        "assignee_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the assignee - link to `gitlab_user.id`.",
			Transform:   transform.FromField("Assignee.ID"),
		},
		{
			Name:        "assignee_username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the assignee - link to `gitlab_user.username`.",
			Transform:   transform.FromField("Assignee.Username"),
		},
		{
			Name:        "assignee_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the assignee.",
			Transform:   transform.FromField("Assignee.Name"),
		},
		{
			Name:        "assignees",
			Type:        proto.ColumnType_JSON,
			Description: "An array of assigned usernames, for when more than one user is assigned.",
			Transform:   transform.FromField("Assignees").NullIfZero().Transform(parseBasicUserCollection),
		},
		{
			Name:        "reviewers",
			Type:        proto.ColumnType_JSON,
			Description: "An array of usernames who've been asked to review the merge request.",
			Transform:   transform.FromField("Reviewers").NullIfZero().Transform(parseBasicUserCollection),
		},
		{
			Name:        "source_project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the source project.",
		},
		{
			Name:        "target_project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the target project.",
		},
		{
			Name:        "labels",
			Type:        proto.ColumnType_JSON,
			Description: "An array of textual labels applied to the merge request.",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Description: "The description of the merge request.",
		},
		{
			Name:        "draft",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the merge request is a draft.",
		},
		{
			Name:        "work_in_progress",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the merge request is a work in progress.",
		},
		{
			Name:        "merge_when_pipeline_succeeds",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the merge request will be merged upon completion of CI/CD pipeline.",
		},
		{
			Name:        "merge_status",
			Type:        proto.ColumnType_STRING,
			Description: "Descriptive status about the ability of being able to merge the merge request.",
		},
		{
			Name:        "merge_error",
			Type:        proto.ColumnType_STRING,
			Description: "Error message if the merge request can not be merged.",
		},
		{
			Name:        "merged_by_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user who merged the merge request - link to `gitlab_user.id`.",
			Transform:   transform.FromField("MergedBy.ID"),
		},
		{
			Name:        "merged_by_username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user who merged the merge request - link to `gitlab_user.username`.",
			Transform:   transform.FromField("MergedBy.Username"),
		},
		{
			Name:        "merged_by_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the user whom merged the merge request.",
		},
		{
			Name:        "merged_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when merge request was merged.",
		},
		{
			Name:        "closed_by_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user who closed the merge request - link to `gitlab_user.id`.",
			Transform:   transform.FromField("ClosedBy.ID"),
		},
		{
			Name:        "closed_by_username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user who closed the merge request - link to `gitlab_user.username`.",
			Transform:   transform.FromField("ClosedBy.Username"),
		},
		{
			Name:        "closed_by_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the user who closed the merge request.",
			Transform:   transform.FromField("ClosedBy.Name"),
		},
		{
			Name:        "closed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when merge request was closed.",
		},
		{
			Name:        "subscribed",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the user associated to the token used to access the data is subscribed to the merge request.",
		},
		{
			Name:        "sha",
			Type:        proto.ColumnType_STRING,
			Description: "",
			Transform:   transform.FromField("SHA"),
		},
		{
			Name:        "merge_commit_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The hash of the merge commit.",
			Transform:   transform.FromField("MergeCommitSHA"),
		},
		{
			Name:        "squash_commit_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The hash of the squashed merge commit.",
			Transform:   transform.FromField("SquashCommitSHA"),
		},
		{
			Name:        "user_notes_count",
			Type:        proto.ColumnType_INT,
			Description: "A count of user notes on the merge request.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "changes_count",
			Type:        proto.ColumnType_STRING, // NOTE: This is string in SDK
			Description: "A count of changes contained within the merge request.",
		},
		{
			Name:        "should_remove_source_branch",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if source_branch should be deleted on merge.",
		},
		{
			Name:        "force_remove_source_branch",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if source_branch will be force deleted on merge.",
		},
		{
			Name:        "allow_collaboration",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if collaboration is allowed on the merge request.",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to access the merge request.",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "short_ref",
			Type:        proto.ColumnType_STRING,
			Description: "Short reference of the merge request.",
			Transform:   transform.FromField("References.Short"),
		},
		{
			Name:        "rel_ref",
			Type:        proto.ColumnType_STRING,
			Description: "Relative reference of the merge request.",
			Transform:   transform.FromField("References.Relative"),
		},
		{
			Name:        "full_ref",
			Type:        proto.ColumnType_STRING,
			Description: "Full reference of the merge request.",
			Transform:   transform.FromField("References.Full"),
		},
		{
			Name:        "discussion_locked",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the merge request has the discussions locked against new input.",
		},
		{
			Name:        "can_merge",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the current user can merge the merge request.",
			Transform:   transform.FromField("User.CanMerge"),
		},
		{
			Name:        "squash",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if a squash is requested.",
		},
		{
			Name:        "base_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The base sha for the diff.",
			Transform:   transform.FromField("DiffRefs.BaseSha"),
		},
		{
			Name:        "head_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The head sha for the diff.",
			Transform:   transform.FromField("DiffRefs.HeadSha"),
		},
		{
			Name:        "start_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The start sha for the diff.",
			Transform:   transform.FromField("DiffRefs.StartSha"),
		},
		{
			Name:        "diverged_commits_count",
			Type:        proto.ColumnType_INT,
			Description: "A count of commits diverged from target_branch.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "rebase_in_progress",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if a rebase is in progress.",
		},
		{
			Name:        "approvals_before_merge",
			Type:        proto.ColumnType_INT,
			Description: "The number of approvals required before merge can proceed.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "reference",
			Type:        proto.ColumnType_STRING,
			Description: "The reference code of the merge request (example: `!4`).",
		},
		{
			Name:        "first_contribution",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the merge request contains a first contribution to the project.",
		},
		{
			Name:        "has_conflicts",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the merge request has conflicts with the target_branch.",
		},
		{
			Name:        "blocking_discussions_resolved",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if blocking discussions have all been resolved.",
		},
		{
			Name:        "reviewer_id",
			Type:        proto.ColumnType_INT,
			Description: "Contains reviewer_id if passed as a qualifier for filtering, else null.",
			Transform:   transform.FromQual("reviewer_id").NullIfZero(),
		},
		// Pipeline
		{
			Name:        "pipeline_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the pipeline run against the merge request.",
			Transform:   transform.FromField("Pipeline.ID"),
		},
		{
			Name:        "pipeline_project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project the pipeline run against the merge request belongs to.",
			Transform:   transform.FromField("Pipeline.ProjectID"),
		},
		{
			Name:        "pipeline_status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the pipeline.",
			Transform:   transform.FromField("Pipeline.Status"),
		},
		{
			Name:        "pipeline_source",
			Type:        proto.ColumnType_STRING,
			Description: "The source of the pipeline.",
			Transform:   transform.FromField("Pipeline.Source"),
		},
		{
			Name:        "pipeline_ref",
			Type:        proto.ColumnType_STRING,
			Description: "The reference of the pipeline.",
			Transform:   transform.FromField("Pipeline.Ref"),
		},
		{
			Name:        "pipeline_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The commit sha that is run in the pipeline.",
			Transform:   transform.FromField("Pipeline.SHA"),
		},
		{
			Name:        "pipeline_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the pipeline.",
			Transform:   transform.FromField("Pipeline.WebURL"),
		},
		{
			Name:        "pipeline_created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the pipeline was created.",
			Transform:   transform.FromField("Pipeline.CreatedAt"),
		},
		{
			Name:        "pipeline_updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the pipeline was updated.",
			Transform:   transform.FromField("Pipeline.UpdatedAt"),
		},
		// Milestone
		{
			Name:        "milestone_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the milestone the merge request is placed into.",
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
	}
}
