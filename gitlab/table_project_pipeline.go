package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"time"
)

type ProjectPipeline struct {
	ID        int
	Status    string
	Ref       string
	SHA       string
	UpdatedAt *time.Time
	CreatedAt *time.Time
	WebURL    string
	ProjectID int
	Source    string
}

func tableProjectPipeline() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_pipeline",
		Description: "Pipelines for a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectPipelines,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the pipeline."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the pipeline (success/failed/canceled)."},
			{Name: "ref", Type: proto.ColumnType_STRING, Description: "The reference associated with the pipeline (branch name or tag)."},
			{Name: "sha", Type: proto.ColumnType_STRING, Description: "The commit SHA at which the pipeline was run against.", Transform: transform.FromField("SHA")},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to view the pipeline.", Transform: transform.FromField("WebURL")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the pipeline was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the pipeline was last updated."},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project the pipeline was run against - link `gitlab_project.id`."},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "The source associated with the pipeline."},
		},
	}
}

func listProjectPipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())

	opt := &api.ListProjectPipelinesOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 20,
	}}

	for {
		pipelines, resp, err := conn.Pipelines.ListProjectPipelines(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, pipeline := range pipelines {
			d.StreamListItem(ctx, &ProjectPipeline{
				ID:        pipeline.ID,
				Status:    pipeline.Status,
				Ref:       pipeline.Ref,
				SHA:       pipeline.SHA,
				WebURL:    pipeline.WebURL,
				CreatedAt: pipeline.CreatedAt,
				UpdatedAt: pipeline.UpdatedAt,
				ProjectID: projectId,
				Source:    pipeline.Source,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
