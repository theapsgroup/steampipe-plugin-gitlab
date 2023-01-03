package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_repository",
		Description: "Repository for a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
			},
			Hydrate: listRepositoryTree,
		},
		Columns: []*plugin.Column{
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
				Name:        "project_id",
				Type:        proto.ColumnType_INT,
				Description: "The ID of the project this repository belongs to - link `gitlab_project.id`.",
				Transform:   transform.FromQual("project_id"),
			},
		},
	}
}

// Hydrate Functions

func listRepositoryTree(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())

	opt := &api.ListTreeOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 20,
		},
		Recursive: api.Bool(true),
	}

	for {
		nodes, resp, err := conn.Repositories.ListTree(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, node := range nodes {
			d.StreamListItem(ctx, node)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
