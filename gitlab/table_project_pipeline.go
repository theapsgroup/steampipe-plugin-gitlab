package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectPipeline() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_pipeline",
		Description: "Obtain information about pipelines for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
				{
					Name:      "updated_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">=", "=", "<", "<="},
				},
				{
					Name:      "status",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
			},
			Hydrate: listProjectPipelines,
		},
		Columns: projectPipelineColumns(),
	}
}

func listProjectPipelines(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectPipelines", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectPipelines", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectPipelinesOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	if d.Quals["updated_at"] != nil {
		for _, q := range d.Quals["updated_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			beforeTime := givenTime.Add(time.Duration(-1) * time.Second)
			afterTime := givenTime.Add(time.Second * 1)

			switch q.Operator {
			case ">":
				opt.UpdatedAfter = &afterTime
			case ">=":
				opt.UpdatedAfter = &givenTime
			case "=":
				opt.UpdatedAfter = &beforeTime
				opt.UpdatedBefore = &afterTime
			case "<=":
				opt.UpdatedBefore = &givenTime
			case "<":
				opt.UpdatedBefore = &beforeTime
			}
		}
	}

	if d.EqualsQuals["status"] != nil {
		s := d.EqualsQuals["status"].GetStringValue()

		switch strings.ToLower(s) {
		case "pending":
			opt.Status = api.BuildState(api.Pending)
		case "created":
			opt.Status = api.BuildState(api.Created)
		case "running":
			opt.Status = api.BuildState(api.Canceled)
		case "success":
			opt.Status = api.BuildState(api.Success)
		case "failed":
			opt.Status = api.BuildState(api.Failed)
		case "canceled":
			opt.Status = api.BuildState(api.Canceled)
		case "skipped":
			opt.Status = api.BuildState(api.Skipped)
		case "manual":
			opt.Status = api.BuildState(api.Manual)
		default:
			return nil, nil
		}
	}

	for {
		plugin.Logger(ctx).Debug("listProjectPipelines", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		pipelines, resp, err := conn.Pipelines.ListProjectPipelines(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectPipelines", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain pipelines for project_id %d\n%v", projectId, err)
		}

		for _, pipeline := range pipelines {
			d.StreamListItem(ctx, pipeline)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectPipelines", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectPipelines", "completed successfully")
	return nil, nil
}

// Column Function
func projectPipelineColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the pipeline.",
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the pipeline (success/failed/canceled).",
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
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to view the pipeline.",
			Transform:   transform.FromField("WebURL"),
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
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project the pipeline was run against - link `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
		{
			Name:        "source",
			Type:        proto.ColumnType_STRING,
			Description: "The source associated with the pipeline.",
		},
	}
}
