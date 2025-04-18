package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectDeployment() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_deployment",
		Description: "Obtain information about deployments associated with a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listProjectDeployments,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getProjectDeployment,
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
		},
		Columns: projectDeploymentColumns(),
	}
}

// Hydrate Functions
func listProjectDeployments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectDeployments", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectDeployments", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	projectId := int(q["project_id"].GetInt64Value())
	opt := &api.ListProjectDeploymentsOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		plugin.Logger(ctx).Debug("listProjectDeployments", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		deps, resp, err := conn.Deployments.ListProjectDeployments(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectDeployments", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain deployments for project_id %d\n%v", projectId, err)
		}

		for _, dep := range deps {
			d.StreamListItem(ctx, dep)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectDeployments", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectDeployments", "completed successfully")
	return nil, nil
}

func getProjectDeployment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getProjectDeployment", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getProjectDeployment", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	id := int(d.EqualsQuals["id"].GetInt64Value())
	plugin.Logger(ctx).Debug("getProjectDeployment", "projectId", projectId, "id", id)

	dep, _, err := conn.Deployments.GetProjectDeployment(projectId, id)
	if err != nil {
		plugin.Logger(ctx).Error("getProjectDeployment", "projectId", projectId, "id", id, "error", err)
		return nil, fmt.Errorf("unable to obtain deployment %d for project_id %d\n%v", id, projectId, err)
	}

	plugin.Logger(ctx).Debug("getProjectDeployment", "completed successfully")
	return dep, nil
}

// Column Function
func projectDeploymentColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the deployment.",
		},
		{
			Name:        "iid",
			Type:        proto.ColumnType_INT,
			Description: "The internal ID of the deployment.",
			Transform:   transform.FromField("IID"),
		},
		{
			Name:        "ref",
			Type:        proto.ColumnType_STRING,
			Description: "The reference associated with the deployment (branch name or tag).",
		},
		{
			Name:        "sha",
			Type:        proto.ColumnType_STRING,
			Description: "The commit SHA at which the deployment was run against.",
			Transform:   transform.FromField("SHA"),
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the deployment (running/success/failed/canceled).",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the deployment was created.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the deployment was last updated.",
		},
		{
			Name:        "user_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user whom triggered the deployment.",
			Transform:   transform.FromField("User.ID"),
		},
		{
			Name:        "user_username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user whom triggered the deployment.",
			Transform:   transform.FromField("User.Username"),
		},
		{
			Name:        "environment_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the environment the deployment is deployed to.",
			Transform:   transform.FromField("Environment.ID"),
		},
		{
			Name:        "environment_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the environment the deployment is deployed to.",
			Transform:   transform.FromField("Environment.Name"),
		},
		{
			Name:        "deployable_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the deployable.",
			Transform:   transform.FromField("Deployable.ID"),
		},
		{
			Name:        "deployable_status",
			Type:        proto.ColumnType_STRING,
			Description: "The status of the deployable.",
			Transform:   transform.FromField("Deployable.Status"),
		},
		{
			Name:        "deployable_stage",
			Type:        proto.ColumnType_STRING,
			Description: "The stage of the deployable.",
			Transform:   transform.FromField("Deployable.Stage"),
		},
		{
			Name:        "deployable_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the deployable.",
			Transform:   transform.FromField("Deployable.Name"),
		},
		{
			Name:        "deployable_ref",
			Type:        proto.ColumnType_STRING,
			Description: "The ref of the deployable.",
			Transform:   transform.FromField("Deployable.Ref"),
		},
		{
			Name:        "deployable_commit_id",
			Type:        proto.ColumnType_STRING,
			Description: "The ID of the commit for the deployable.",
			Transform:   transform.FromField("Deployable.Commit.ID"),
		},
		{
			Name:        "deployable_pipeline_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the pipeline for the deployable.",
			Transform:   transform.FromField("Deployable.Pipeline.ID"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project - link to `gitlab_project.id",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
