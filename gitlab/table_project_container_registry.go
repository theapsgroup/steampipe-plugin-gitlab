package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectContainerRegistry() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_container_registry",
		Description: "Obtain information on the container registry associated to a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listProjectContainerRegistries,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: projectContainerRegistryColumns(),
	}
}

// Hydrate Functions
func listProjectContainerRegistries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectContainerRegistries", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectContainerRegistries", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListRegistryRepositoriesOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		plugin.Logger(ctx).Debug("listProjectContainerRegistries", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		crs, resp, err := conn.ContainerRegistry.ListProjectRegistryRepositories(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectContainerRegistries", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain container registries for project_id %d\n%v", projectId, err)
		}

		for _, cr := range crs {
			d.StreamListItem(ctx, cr)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectContainerRegistries", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectContainerRegistries", "completed successfully")
	return nil, nil
}

// Column Functions
func projectContainerRegistryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the container registry.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the container registry.",
		},
		{
			Name:        "path",
			Type:        proto.ColumnType_STRING,
			Description: "The path of the container registry.",
		},
		{
			Name:        "location",
			Type:        proto.ColumnType_STRING,
			Description: "The location (full path) of the container registry.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the container registry was created.",
		},
		{
			Name:        "cleanup_policy_started_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the cleanup policy was started.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project - link to `gitlab_project.id",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
