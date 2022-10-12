package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectRepositoryFile() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_repository_file",
		Description: "Details of a specific File from a Repository on a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
				{
					Name:    "file_path",
					Require: plugin.Required,
				},
				{
					Name:    "ref",
					Require: plugin.Optional,
				},
			},
			Hydrate: listRepoFile,
		},
		Columns: repoFileColumns(),
	}
}

func repoFileColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "file_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the file",
		},
		{
			Name:        "file_path",
			Type:        proto.ColumnType_STRING,
			Description: "The path of the file",
			Transform:   transform.FromQual("file_path"),
		},
		{
			Name:        "size",
			Type:        proto.ColumnType_INT,
			Description: "The size of the file",
		},
		{
			Name:        "encoding",
			Type:        proto.ColumnType_STRING,
			Description: "The encoding used on the Content field value",
		},
		{
			Name:        "content",
			Type:        proto.ColumnType_STRING,
			Description: "The content of the file, encoded as per the encoded field",
		},
		{
			Name:        "ref",
			Type:        proto.ColumnType_STRING,
			Description: "The repository ref (tag, branch, etc)",
		},
		{
			Name:        "blob_id",
			Type:        proto.ColumnType_STRING,
			Description: "The blob ID of the file, can be used to obtain the blob",
		},
		{
			Name:        "commit_id",
			Type:        proto.ColumnType_STRING,
			Description: "The ID of the last commit which affected the file, can be used to pull the commit",
		},
		{
			Name:        "content_sha256",
			Type:        proto.ColumnType_STRING,
			Description: "The SHA256 hash of the file content",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project this repository file belongs to - link `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}

// Hydrate Functions
func listRepoFile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	q := d.KeyColumnQuals

	projectId := int(q["project_id"].GetInt64Value())
	filePath := q["file_path"].GetStringValue()

	ref := "main"
	if q["ref"] != nil {
		ref = q["ref"].GetStringValue()
	}

	opt := api.GetFileOptions{
		Ref: &ref,
	}

	file, _, err := conn.RepositoryFiles.GetFile(projectId, filePath, &opt)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, file)

	return nil, nil
}
