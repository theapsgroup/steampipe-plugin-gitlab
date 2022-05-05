package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	api "github.com/xanzy/go-gitlab"
)

func tableMyProject() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_my_project",
		Description: "Projects in the GitLab Instance where authenticated user is a member.",
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
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
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
		projects, resp, err := conn.Projects.ListProjects(opt)
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			d.StreamListItem(ctx, project)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

func getMyProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	q := d.KeyColumnQuals
	id := int(q["id"].GetInt64Value())
	stats := true

	opt := &api.GetProjectOptions{Statistics: &stats}

	project, _, err := conn.Projects.GetProject(id, opt)
	if err != nil {
		return nil, err
	}
	return project, nil
}
