package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableMyProject() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_my_project",
		Description: "Obtain information about projects that the authenticated user is a member of within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listMyProjects,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getMyProject,
		},
		Columns: projectColumns(),
	}
}

func listMyProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listMyProjects", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listMyProjects", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	membership := true
	stats := true

	opt := &api.ListProjectsOptions{
		Membership: &membership,
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
		Statistics: &stats,
	}

	for {
		plugin.Logger(ctx).Debug("listMyProjects", "page", opt.Page, "perPage", opt.PerPage)
		projects, resp, err := conn.Projects.ListProjects(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listMyProjects", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain projects for current user\n%v", err)
		}

		for _, project := range projects {
			d.StreamListItem(ctx, project)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listMyProjects", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listMyProjects", "completed successfully")
	return nil, nil
}

func getMyProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getMyProject", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getMyProject", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}
	q := d.EqualsQuals
	id := int(q["id"].GetInt64Value())
	stats := true

	opt := &api.GetProjectOptions{Statistics: &stats}

	plugin.Logger(ctx).Debug("getMyProject", "id", id)
	project, _, err := conn.Projects.GetProject(id, opt)
	if err != nil {
		plugin.Logger(ctx).Error("getMyProject", "id", id, "error", err)
		return nil, fmt.Errorf("unable to obtain project with id %d\n%v", id, err)
	}

	plugin.Logger(ctx).Debug("getMyProject", "completed successfully")
	return project, nil
}
