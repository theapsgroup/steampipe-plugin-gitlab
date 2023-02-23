package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

func tableProjectPipelineDetail() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_pipeline_detail",
		Description: "Obtain details for a specific pipeline within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    listProjectPipelineDetails,
		},
		Columns: projectPipelineDetailColumns(),
	}
}

// Hydration Functions
func listProjectPipelineDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectPipelineDetails", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectPipelineDetails", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	pipelineId := int(d.EqualsQuals["id"].GetInt64Value())

	plugin.Logger(ctx).Debug("listProjectPipelineDetails", "projectId", projectId, "pipelineId", pipelineId)
	pipeline, _, err := conn.Pipelines.GetPipeline(projectId, pipelineId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			plugin.Logger(ctx).Warn("listProjectPipelineDetails", "projectId", projectId, "pipelineId", pipelineId, "no project/pipeline was found, returning empty result set")
			return nil, nil
		}
		plugin.Logger(ctx).Error("listProjectPipelineDetails", "projectId", projectId, "pipelineId", pipelineId, "error", err)
		return nil, fmt.Errorf("unable to obtain pipeline details for project_id %d - id %d\n%v", projectId, pipelineId, err)
	}

	d.StreamListItem(ctx, pipeline)

	plugin.Logger(ctx).Debug("listProjectPipelineDetails", "completed successfully")
	return nil, nil
}

// Column Function
func projectPipelineDetailColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the pipeline.",
		},
		{
			Name:        "iid",
			Type:        proto.ColumnType_INT,
			Description: "The internal ID of the pipeline.",
			Transform:   transform.FromField("IID"),
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the pipeline (success/failed/canceled).",
		},
		{
			Name:        "source",
			Type:        proto.ColumnType_STRING,
			Description: "The source of the pipeline.",
		},
		{
			Name:        "ref",
			Type:        proto.ColumnType_STRING,
			Description: "The reference associated with the pipeline (branch name or tag).",
		},
		{
			Name:        "sha",
			Type:        proto.ColumnType_STRING,
			Description: "The commit SHA at which the pipeline was run against.",
			Transform:   transform.FromField("SHA"),
		},
		{
			Name:        "before_sha",
			Type:        proto.ColumnType_STRING,
			Description: "",
			Transform:   transform.FromField("BeforeSHA"),
		},
		{
			Name:        "tag",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the pipeline was triggered by a tag.",
		},
		{
			Name:        "yaml_errors",
			Type:        proto.ColumnType_STRING,
			Description: "",
		},
		{
			Name:        "user_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user which triggered the pipeline - link to `gitlab_user.ID`.",
			Transform:   transform.FromField("User.ID"),
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user which triggered the pipeline - link to `gitlab_user.username`.",
			Transform:   transform.FromField("User.Username"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the pipeline was created.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the pipeline was last updated.",
		},
		{
			Name:        "started_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the pipeline started.",
		},
		{
			Name:        "finished_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the pipeline finished.",
		},
		{
			Name:        "committed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the commit used by the pipeline was created.",
		},
		{
			Name:        "duration",
			Type:        proto.ColumnType_INT,
			Description: "Time in seconds the pipeline took to run.",
		},
		{
			Name:        "queued_duration",
			Type:        proto.ColumnType_INT,
			Description: "Time in seconds the pipeline was queued awaiting running.",
		},
		{
			Name:        "coverage",
			Type:        proto.ColumnType_STRING,
			Description: "Coverage",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to view the pipeline.",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project the pipeline was run against - link `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
