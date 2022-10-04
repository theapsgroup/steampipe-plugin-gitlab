package gitlab

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

type ProjectJob struct {
	ID         int
	Name       string
	Status     string
	Ref        string
	Stage      string
	CreatedAt  *time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Duration   float64
	WebURL     string
	UserID     int
	Username   string
	PipelineID int
	ProjectID  int
}

func tableProjectJob() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_job",
		Description: "Jobs for a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectJobs,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the job."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The Name of the job."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the job (success/failed/canceled)."},
			{Name: "ref", Type: proto.ColumnType_STRING, Description: "The reference associated with the pipeline (branch name or tag)."},
			{Name: "stage", Type: proto.ColumnType_STRING, Description: "The stage of the job (build/test/staging/production)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the job was created."},
			{Name: "started_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the job was started."},
			{Name: "finished_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the job finished."},
			{Name: "duration", Type: proto.ColumnType_DOUBLE, Description: "Running duration of the job (seconds)."},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to view the job.", Transform: transform.FromField("WebURL")},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The ID of the user wcho triggered the job - link to `gitlab_user.ID`.", Transform: transform.FromField("UserID")},
			{Name: "username", Type: proto.ColumnType_STRING, Description: "The NAME of the user wcho triggered the job - link to `gitlab_user.username`.", Transform: transform.FromField("Username")},
			{Name: "pipeline_id", Type: proto.ColumnType_INT, Description: "The ID of the pipeline which the jobs belongs to - link `gitlab_pipeline.id`."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project the job was run against - link `gitlab_project.id`."},
		},
	}
}

func listProjectJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())

	opt := &api.ListJobsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 20,
	}}

	for {
		jobs, resp, err := conn.Jobs.ListProjectJobs(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, job := range jobs {
			d.StreamListItem(ctx, &ProjectJob{
				ID:         job.ID,
				Name:       job.Name,
				Status:     job.Status,
				Ref:        job.Ref,
				Stage:      job.Stage,
				CreatedAt:  job.CreatedAt,
				StartedAt:  job.StartedAt,
				FinishedAt: job.FinishedAt,
				Duration:   job.Duration,
				WebURL:     job.WebURL,
				UserID:     job.User.ID,
				Username:   job.User.Username,
				PipelineID: job.Pipeline.ID,
				ProjectID:  projectId,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
