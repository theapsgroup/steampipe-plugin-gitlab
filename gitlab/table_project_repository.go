package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_repository",
		Description: "Obtain information about a repository for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
				{
					Name:      "ref",
					Require:   plugin.Optional,
					Operators: []string{"="},
				},
			},
			Hydrate: listRepositoryTree,
		},
		Columns: repoColumns(),
	}
}

// Hydrate Functions
func listRepositoryTree(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listRepositoryTree", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listRepositoryTree", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListTreeOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
		Recursive: api.Bool(true),
	}

	if d.EqualsQualString("ref") != "" {
		ref := api.String(d.EqualsQualString("ref"))
		opt.Ref = ref
		plugin.Logger(ctx).Debug("listRepositoryTree", "filter[ref]", *ref)
	}

	for {
		plugin.Logger(ctx).Debug("listRepositoryTree", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		nodes, resp, err := conn.Repositories.ListTree(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listRepositoryTree", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain repository for project_id %d\n%v", projectId, err)
		}

		for _, node := range nodes {
			d.StreamListItem(ctx, node)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listRepositoryTree", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listRepositoryTree", "completed successfully")
	return nil, nil
}

// Column Function
func repoColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "The ID of the file or folder within the repository",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the file or folder within the repository",
		},
		{
			Name:        "type",
			Type:        proto.ColumnType_STRING,
			Description: "The type of the file or folder within the repository",
		},
		{
			Name:        "path",
			Type:        proto.ColumnType_STRING,
			Description: "The path of the file or folder within the repository",
		},
		{
			Name:        "mode",
			Type:        proto.ColumnType_STRING,
			Description: "The mode of the file or folder within the repository",
		},
		{
			Name:        "ref",
			Type:        proto.ColumnType_STRING,
			Description: "The name of a repository branch or tag or, if not given, the default branch",
			Transform:   transform.FromQual("ref"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project this repository belongs to - link `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
