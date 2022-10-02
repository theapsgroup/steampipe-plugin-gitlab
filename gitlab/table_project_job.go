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
	ID        int
	Name      string
	Status    string
	CreatedAt *time.Time
	WebURL    string
	ProjectID int
	// Stage      string
	// FinishedAt *time.Time
	// User
	// Duration
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
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to view the job.", Transform: transform.FromField("WebURL")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the job was created."},
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
				ID:        job.ID,
				Name:      job.Name,
				Status:    job.Status,
				WebURL:    job.WebURL,
				CreatedAt: job.CreatedAt,
				ProjectID: projectId,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
