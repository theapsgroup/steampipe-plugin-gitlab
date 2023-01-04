package gitlab

import (
    "context"
    "github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
    "github.com/turbot/steampipe-plugin-sdk/v5/plugin"
    "github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
    api "github.com/xanzy/go-gitlab"
)

func tableMergeRequestChange() *plugin.Table {
    return &plugin.Table{
        Name:        "gitlab_merge_request_change",
        Description: "Get all changes associated with a merge request.",
        List: &plugin.ListConfig{
            Hydrate:    listChanges,
            KeyColumns: plugin.AllColumns([]string{"iid", "project_id"}),
        },
        Columns: mergeRequestChangeColumns(),
    }
}

// Hydrate Functions
func listChanges(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    conn, err := connect(ctx, d)
    if err != nil {
        return nil, err
    }

    q := d.EqualsQuals
    iid := int(q["iid"].GetInt64Value())
    projectId := int(q["project_id"].GetInt64Value())

    mergeRequest, _, err := conn.MergeRequests.GetMergeRequest(projectId, iid, &api.GetMergeRequestsOptions{})
    if err != nil {
        return nil, err
    }

    for _, change := range mergeRequest.Changes {
        d.StreamListItem(ctx, change)
    }

    return nil, nil
}

func mergeRequestChangeColumns() []*plugin.Column {
    return []*plugin.Column{
        {
            Name:        "iid",
            Type:        proto.ColumnType_INT,
            Description: "Internal ID of the merge request to which the change belongs.",
            Transform:   transform.FromQual("iid"),
        },
        {
            Name:        "project_id",
            Type:        proto.ColumnType_INT,
            Description: "ID of the project to which the merge request belongs.",
            Transform:   transform.FromQual("project_id"),
        },
        {
            Name:        "old_path",
            Type:        proto.ColumnType_STRING,
            Description: "Old path of the file.",
        },
        {
            Name:        "new_path",
            Type:        proto.ColumnType_STRING,
            Description: "New path of the file.",
        },
        {
            Name:        "a_mode",
            Type:        proto.ColumnType_STRING,
            Description: "The a mode associated with the change.",
        },
        {
            Name:        "b_mode",
            Type:        proto.ColumnType_STRING,
            Description: "The b mode associated with the change.",
        },
        {
            Name:        "diff",
            Type:        proto.ColumnType_STRING,
            Description: "The change diff.",
        },
        {
            Name:        "new_file",
            Type:        proto.ColumnType_BOOL,
            Description: "Indicates if it is a new file added.",
        },
        {
            Name:        "renamed_file",
            Type:        proto.ColumnType_BOOL,
            Description: "Indicates if the file has been renamed.",
        },
        {
            Name:        "deleted_file",
            Type:        proto.ColumnType_BOOL,
            Description: "Indicates if the file has been deleted.",
        },
    }
}
