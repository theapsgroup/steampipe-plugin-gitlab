package gitlab

import (
    "context"
    "github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
    "github.com/turbot/steampipe-plugin-sdk/v5/plugin"
    api "github.com/xanzy/go-gitlab"
)

func tableInstanceVariable() *plugin.Table {
    return &plugin.Table{
        Name:        "gitlab_instance_variable",
        Description: "Variables held against the instance of GitLab (only available in self-hosted model).",
        List: &plugin.ListConfig{
            Hydrate: listInstanceVars,
        },
        Columns: instanceVarColumns(),
    }
}

func listInstanceVars(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
    // Not available on public, only self-hosted.
    if isPublicGitLab(d) {
        return nil, nil
    }

    conn, err := connect(ctx, d)
    if err != nil {
        return nil, err
    }

    opt := &api.ListInstanceVariablesOptions{
        Page:    1,
        PerPage: 20,
    }

    for {
        vars, resp, err := conn.InstanceVariables.ListVariables(opt)
        if err != nil {
            return nil, err
        }

        for _, v := range vars {
            d.StreamListItem(ctx, v)
        }

        if resp.NextPage == 0 {
            break
        }

        opt.Page = resp.NextPage
    }

    return nil, nil
}

func instanceVarColumns() []*plugin.Column {
    return []*plugin.Column{
        {
            Name:        "key",
            Type:        proto.ColumnType_STRING,
            Description: "The key of the variable.",
        },
        {
            Name:        "value",
            Type:        proto.ColumnType_STRING,
            Description: "The value of the variable.",
        },
        {
            Name:        "variable_type",
            Type:        proto.ColumnType_STRING,
            Description: "The type of the variable (env var, etc).",
        },
        {
            Name:        "environment_scope",
            Type:        proto.ColumnType_STRING,
            Description: "The environment(s) that this variable is in scope for.",
        },
        {
            Name:        "protected",
            Type:        proto.ColumnType_BOOL,
            Description: "Indicates if the variable is only applied to protected branches.",
        },
        {
            Name:        "masked",
            Type:        proto.ColumnType_BOOL,
            Description: "Indicates if the variable is masked (hidden) in job logs.",
        },
        {
            Name:        "raw",
            Type:        proto.ColumnType_BOOL,
            Description: "Indicates if the variable is is a raw format.",
        },
    }
}
