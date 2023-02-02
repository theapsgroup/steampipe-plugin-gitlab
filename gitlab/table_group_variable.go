package gitlab

import (
    "context"
    "github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
    "github.com/turbot/steampipe-plugin-sdk/v5/plugin"
    "github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
    api "github.com/xanzy/go-gitlab"
)

func tableGroupVariable() *plugin.Table {
    return &plugin.Table{
        Name:        "gitlab_group_variable",
        Description: "Variables for a GitLab Group",
        List: &plugin.ListConfig{
            KeyColumns: []*plugin.KeyColumn{
                {
                    Name:    "group_id",
                    Require: plugin.Required,
                },
            },
            Hydrate: listGroupVars,
        },
        Get: &plugin.GetConfig{
            KeyColumns: plugin.AllColumns([]string{"group_id", "key"}),
            Hydrate:    getGroupVar,
        },
        Columns: groupVarColumns(),
    }
}

func listGroupVars(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
    conn, err := connect(ctx, d)
    if err != nil {
        return nil, err
    }

    groupId := int(d.EqualsQuals["group_id"].GetInt64Value())

    opt := &api.ListGroupVariablesOptions{
        Page:    1,
        PerPage: 20,
    }

    for {
        vars, resp, err := conn.GroupVariables.ListVariables(groupId, opt)
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

func getGroupVar(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
    conn, err := connect(ctx, d)
    if err != nil {
        return nil, err
    }

    projectId := int(d.EqualsQuals["group_id"].GetInt64Value())
    key := d.EqualsQuals["key"].GetStringValue()

    v, _, err := conn.GroupVariables.GetVariable(projectId, key)
    if err != nil {
        return nil, err
    }

    return v, nil
}

func groupVarColumns() []*plugin.Column {
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
        {
            Name:        "group_id",
            Type:        proto.ColumnType_INT,
            Description: "The ID of the group this repository belongs to - link `gitlab_group.id`.",
            Transform:   transform.FromQual("group_id"),
        },
    }
}