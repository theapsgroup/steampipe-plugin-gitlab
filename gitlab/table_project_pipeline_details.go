package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"strings"
	"time"
)

type ProjectPipelineDetails struct {
	ID             int
	Status         string
	Ref            string
	SHA            string
	BeforeSHA      string
	Tag            bool
	YamlErrors     string
	UserID         int
	Username       string
	UpdatedAt      *time.Time
	CreatedAt      *time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
	CommittedAt    *time.Time
	Duration       int
	Coverage       string
	WebURL         string
	ProjectID      int
}

func tableProjectPipelineDetail() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_pipeline_detail",
		Description: "Pipeline details for a specific pipeline",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    listProjectPipelineDetails,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the pipeline."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the pipeline (success/failed/canceled)."},
			{Name: "ref", Type: proto.ColumnType_STRING, Description: "The reference associated with the pipeline (branch name or tag)."},
			{Name: "sha", Type: proto.ColumnType_STRING, Description: "The commit SHA at which the pipeline was run against.", Transform: transform.FromField("SHA")},
			{Name: "before_sha", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("BeforeSHA")},
			{Name: "tag", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "yaml_errors", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "The ID of the user which triggered the pipeline - link to `gitlab_user.ID`.", Transform: transform.FromField("UserID")},
			{Name: "username", Type: proto.ColumnType_STRING, Description: "The username of the user which triggered the pipeline - link to `gitlab_user.username`."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the pipeline was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the pipeline was last updated."},
			{Name: "started_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the pipeline started."},
			{Name: "finished_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the pipeline finished."},
			{Name: "committed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of when the commit used by the pipeline was created."},
			{Name: "duration", Type: proto.ColumnType_INT, Description: ""},
			{Name: "coverage", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to view the pipeline.", Transform: transform.FromField("WebURL")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project the pipeline was run against - link `gitlab_project.id`."},
		},
	}
}

func listProjectPipelineDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())
	pipelineId := int(d.KeyColumnQuals["id"].GetInt64Value())

	pipeline, _, err := conn.Pipelines.GetPipeline(projectId, pipelineId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, &ProjectPipelineDetails{
		ID: pipeline.ID,
		Status: pipeline.Status,
		Ref: pipeline.Ref,
		SHA: pipeline.SHA,
		BeforeSHA: pipeline.BeforeSHA,
		Tag: pipeline.Tag,
		YamlErrors: pipeline.YamlErrors,
		UserID: pipeline.User.ID,
		Username: pipeline.User.Username,
		CreatedAt: pipeline.CreatedAt,
		UpdatedAt: pipeline.UpdatedAt,
		StartedAt: pipeline.StartedAt,
		FinishedAt: pipeline.FinishedAt,
		CommittedAt: pipeline.CommittedAt,
		Duration: pipeline.Duration,
		Coverage: pipeline.Coverage,
		WebURL: pipeline.WebURL,
		ProjectID: projectId,
	})

	return nil, nil
}
