package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectJob() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_job",
		Description: "Obtain information about jobs for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectJobs,
		},
		Columns: projectJobColumns(),
	}
}

// Hydrate Functions
func listProjectJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectJobs", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectJobs", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListJobsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 20,
	}}

	for {
		plugin.Logger(ctx).Debug("listProjectJobs", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		jobs, resp, err := conn.Jobs.ListProjectJobs(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectJobs", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain jobs for project_id %d\n%v", projectId, err)
		}

		for _, job := range jobs {
			d.StreamListItem(ctx, job)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectJobs", "completed successfully")
	return nil, nil
}

// Column Function
func projectJobColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the job.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The Name of the job.",
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the job (success/failed/canceled).",
		},
		{
			Name:        "ref",
			Type:        proto.ColumnType_STRING,
			Description: "The reference associated with the pipeline (branch name or tag).",
		},
		{
			Name:        "stage",
			Type:        proto.ColumnType_STRING,
			Description: "The stage of the job (build/test/staging/production).",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the job was created.",
		},
		{
			Name:        "started_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the job was started.",
		},
		{
			Name:        "finished_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the job finished.",
		},
		{
			Name:        "duration",
			Type:        proto.ColumnType_DOUBLE,
			Description: "Running duration of the job (seconds).",
		},
		{
			Name:        "queued_duration",
			Type:        proto.ColumnType_DOUBLE,
			Description: "Duration in seconds the job was queued before running.",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to view the job.",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "user_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user who triggered the job - link to `gitlab_user.ID`.",
			Transform:   transform.FromField("User.ID")},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user who triggered the job - link to `gitlab_user.username`.",
			Transform:   transform.FromField("User.Username"),
		},
		{
			Name:        "user_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the user who triggered the job.",
			Transform:   transform.FromField("User.Name"),
		},
		{
			Name:        "pipeline_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the pipeline which the jobs belongs to - link `gitlab_pipeline.id`.",
			Transform:   transform.FromField("Pipeline.ID"),
		},
		{
			Name:        "pipeline_project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project to which the pipeline belongs - link `gitlab_project.id`.",
			Transform:   transform.FromField("Pipeline.ProjectID"),
		},
		{
			Name:        "pipeline_ref",
			Type:        proto.ColumnType_STRING,
			Description: "The ref of the pipeline.",
			Transform:   transform.FromField("Pipeline.Ref"),
		},
		{
			Name:        "pipeline_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The sha of the commit the pipeline ran against.",
			Transform:   transform.FromField("Pipeline.Sha"),
		},
		{
			Name:        "pipeline_status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the pipeline.",
			Transform:   transform.FromField("Pipeline.Status"),
		},
		{
			Name:        "artifacts",
			Type:        proto.ColumnType_JSON,
			Description: "An array of artifact information",
		},
		{
			Name:        "runner_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the runner assigned to the job.",
			Transform:   transform.FromField("Runner.ID"),
		},
		{
			Name:        "runner_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the runner assigned to the job.",
			Transform:   transform.FromField("Runner.Name"),
		},
		{
			Name:        "runner_description",
			Type:        proto.ColumnType_STRING,
			Description: "The description of the runner assigned to the job.",
			Transform:   transform.FromField("Runner.Description"),
		},
		{
			Name:        "runner_active",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the runner is active.",
			Transform:   transform.FromField("Runner.Active"),
		},
		{
			Name:        "runner_is_shared",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the runner is shared.",
			Transform:   transform.FromField("Runner.IsShared"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project the job was run against - link `gitlab_project.id`.",
		},
		{
			Name:        "commit_id",
			Type:        proto.ColumnType_STRING,
			Description: "The ID of the commit.",
			Transform:   transform.FromField("Commit.ID"),
		},
		{
			Name:        "commit_short_id",
			Type:        proto.ColumnType_STRING,
			Description: "The short ID of the commit.",
			Transform:   transform.FromField("Commit.ShortID"),
		},
		{
			Name:        "allow_failure",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the job is allowed to fail and allow the pipeline to proceed.",
		},
		{
			Name:        "failure_reason",
			Type:        proto.ColumnType_STRING,
			Description: "The reason for the jobs failure (if failed).",
		},
		{
			Name:        "tag",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the job was started by a tag.",
		},
	}
}
