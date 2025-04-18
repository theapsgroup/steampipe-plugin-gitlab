package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

// TODO: Figure out being able to use full_path as a key for the get function, currently seems to fail in gitlab api wrapper.

func tableProject() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project",
		Description: "Obtain information on projects in the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listProjects,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
				{Name: "owner_username", Require: plugin.Optional},
			},
		},
		Columns: projectColumns(),
	}
}

// Hydrate Functions
func listProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.EqualsQuals

	if q["owner_id"] == nil &&
		q["owner_username"] == nil &&
		q["id"] == nil &&
		isPublicGitLab(d) {
		plugin.Logger(ctx).Error("listProjects", "Public GitLab requires an '=' qualifier for at least one of the following columns 'id', 'owner_id', 'owner_username' - none was provided")
		return nil, fmt.Errorf("when using the gitlab_project table with GitLab Cloud, `List` call requires an '=' qualifier for one or more of the following columns: 'id', 'owner_id', 'owner_username'")
	}

	if q["id"] != nil {
		plugin.Logger(ctx).Debug("listProjects", "id qualifier obtained, re-directing SDK call to GetProject")
		return getProject(ctx, d, h)
	}

	if q["owner_id"] != nil || q["owner_username"] != nil {
		plugin.Logger(ctx).Debug("listProjects", "owner_id or owner_username qualifier obtained, re-directing SDK call to ListUserProjects")
		return listUserProjects(ctx, d, h)
	}

	return listAllProjects(ctx, d, h)
}

func listUserProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listUserProjects", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserProjects", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	q := d.EqualsQuals
	stats := true

	opt := &api.ListProjectsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	},
		Statistics: &stats,
	}

	var x interface{}
	if q["owner_id"] != nil {
		x = int(q["owner_id"].GetInt64Value())
		plugin.Logger(ctx).Debug("listUserProjects", "filter[owner_id]", x.(int))
	} else {
		x = q["owner_username"].GetStringValue()
		plugin.Logger(ctx).Debug("listUserProjects", "filter[owner_username]", x)
	}

	for {
		plugin.Logger(ctx).Debug("listUserProjects", "page", opt.Page, "perPage", opt.PerPage)
		projects, resp, err := conn.Projects.ListUserProjects(x, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listUserProjects", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain projects\n%v", err)
		}

		for _, project := range projects {
			d.StreamListItem(ctx, project)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listUserProjects", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listUserProjects", "completed successfully")
	return nil, nil
}

func listAllProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listAllProjects", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	stats := true
	opt := &api.ListProjectsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	},
		Statistics: &stats,
	}

	for {
		plugin.Logger(ctx).Debug("listAllProjects", "page", opt.Page, "perPage", opt.PerPage)
		projects, resp, err := conn.Projects.ListProjects(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listAllProjects", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain issues\n%v", err)
		}

		for _, project := range projects {
			d.StreamListItem(ctx, project)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listAllProjects", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listAllProjects", "completed successfully")
	return nil, nil
}

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getProject", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	q := d.EqualsQuals
	id := int(q["id"].GetInt64Value())
	stats := true
	opt := &api.GetProjectOptions{
		Statistics: &stats,
	}

	plugin.Logger(ctx).Debug("getProject", "id", id)
	project, _, err := conn.Projects.GetProject(id, opt)
	if err != nil {
		plugin.Logger(ctx).Error("getProject", "id", id, "error", err)
		return nil, fmt.Errorf("unable to obtain project with id %d\n%v", id, err)
	}

	d.StreamListItem(ctx, project)

	plugin.Logger(ctx).Debug("getProject", "completed successfully")
	return nil, nil
}

// Column Functions
func projectColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The projects name.",
		},
		{
			Name:        "path",
			Type:        proto.ColumnType_STRING,
			Description: "The projects path.",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Description: "The projects description.",
		},
		{
			Name:        "default_branch",
			Type:        proto.ColumnType_STRING,
			Description: "The projects default branch name.",
		},
		{
			Name:        "full_name",
			Type:        proto.ColumnType_STRING,
			Description: "The projects name including namespace.",
			Transform:   transform.FromField("NameWithNamespace"),
		},
		{
			Name:        "full_path",
			Type:        proto.ColumnType_STRING,
			Description: "The projects path including namespace.",
			Transform:   transform.FromField("PathWithNamespace"),
		},
		{
			Name:        "public",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project is public",
		},
		{
			Name:        "visibility",
			Type:        proto.ColumnType_STRING,
			Description: "The projects visibility level (private/public/internal)",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The projects url.",
		},
		{
			Name:        "ssh_url",
			Type:        proto.ColumnType_STRING,
			Description: "The ssh url for the project repository.",
			Transform:   transform.FromField("SSHURLToRepo"),
		},
		{
			Name:        "http_url",
			Type:        proto.ColumnType_STRING,
			Description: "The http url to the project repository.",
			Transform:   transform.FromField("HTTPURLToRepo"),
		},
		{
			Name:        "readme_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for the projects readme file",
			Transform:   transform.FromField("ReadmeURL"),
		},
		{
			Name:        "tag_list",
			Type:        proto.ColumnType_JSON,
			Description: "An array of tags associated to the project.",
		},
		{
			Name:        "topics",
			Type:        proto.ColumnType_JSON,
			Description: "An array of topics associated to the project.",
		},
		{
			Name:        "issues_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if project has issues enabled.",
		},
		{
			Name:        "open_issues_count",
			Type:        proto.ColumnType_INT,
			Description: "A count of open issues on the project.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "merge_requests_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if merge requests are enabled on the project",
		},
		{
			Name:        "approvals_before_merge",
			Type:        proto.ColumnType_INT,
			Description: "The project setting for number of approvals required before a merge request can be merged.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "jobs_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has jobs enabled.",
		},
		{
			Name:        "wiki_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has the wiki enabled.",
		},
		{
			Name:        "snippets_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has snippets enabled.",
		},
		{
			Name:        "container_registry_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has the container registry enabled.",
		},
		{
			Name:        "container_registry_image_prefix",
			Type:        proto.ColumnType_STRING,
			Description: "The image prefix for the container registry.",
		},
		{
			Name:        "container_registry_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "Access level of current user for the container registry.",
		},
		{
			Name:        "container_expiration_policy",
			Type:        proto.ColumnType_JSON,
			Description: "JSON Object outlining the expiration policy attached to the container registry.",
		},
		{
			Name:        "creator_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the projects creator. - link to `gitlab_user.id`",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when project was created.",
		},
		{
			Name:        "last_activity_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when last activity happened on the project.",
		},
		{
			Name:        "marked_for_deletion_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when project was marked for deletion.",
			Transform:   transform.FromField("MarkedForDeletionAt").NullIfZero().Transform(isoTimeTransform)},
		{
			Name:        "empty_repo",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the repository of the project is empty.",
		},
		{
			Name:        "archived",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project is archived.",
		},
		{
			Name:        "avatar_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for the projects avatar.",
		},
		{
			Name:        "forks_count",
			Type:        proto.ColumnType_INT,
			Description: "The number of forks of the project.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "star_count",
			Type:        proto.ColumnType_INT,
			Description: "The number of stars given to the project.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "lfs_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has large file system enabled.",
			Transform:   transform.FromField("LFSEnabled"),
		},
		{
			Name:        "request_access_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has request access enabled.",
		},
		{
			Name:        "packages_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has packages enabled.",
		},
		{
			Name:        "owner_id",
			Type:        proto.ColumnType_INT,
			Description: "The projects owner ID. (null if owned by a group) - link to `gitlab_user.id`",
			Transform:   transform.FromField("Owner.ID"),
		},
		{
			Name:        "owner_username",
			Type:        proto.ColumnType_STRING,
			Description: "The projects owner username. (null if owned by a group) - link to `gitlab_user.username`",
			Transform:   transform.FromField("Owner.Username"),
		},
		{
			Name:        "owner_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name for the projects owner.",
			Transform:   transform.FromField("Owner.Name"),
		},
		{
			Name:        "commit_count",
			Type:        proto.ColumnType_INT,
			Description: "The number of commits on the project.",
			Transform:   transform.FromField("Statistics.CommitCount"),
		},
		{
			Name:        "storage_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of the project on disk.",
			Transform:   transform.FromField("Statistics.StorageSize"),
		},
		{
			Name:        "repository_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of the projects repository on disk.",
			Transform:   transform.FromField("Statistics.RepositorySize"),
		},
		{
			Name:        "lfs_objects_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of the projects LFS objects on disk.",
			Transform:   transform.FromField("Statistics.LFSObjectsSize"),
		},
		{
			Name:        "job_artifacts_size",
			Type:        proto.ColumnType_INT,
			Description: "The size of projects job artifacts on disk.",
			Transform:   transform.FromField("Statistics.JobArtifactsSize"),
		},
		{
			Name:        "issues_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for issues.",
		},
		{
			Name:        "repository_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for the repository.",
		},
		{
			Name:        "merge_requests_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for merge requests.",
		},
		{
			Name:        "forking_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for forks/forking.",
		},
		{
			Name:        "wiki_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for the wiki.",
		},
		{
			Name:        "builds_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for builds.",
		},
		{
			Name:        "snippets_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for snippets.",
		},
		{
			Name:        "pages_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for pages.",
		},
		{
			Name:        "operations_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for operations.",
		},
		{
			Name:        "analytics_access_level",
			Type:        proto.ColumnType_STRING,
			Description: "The access level for analytics.",
		},
		{
			Name:        "namespace_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the namespace to which the project belongs.",
			Transform:   transform.FromField("Namespace.ID"),
		},
		{
			Name:        "namespace_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the namespace to which the project belongs.",
			Transform:   transform.FromField("Namespace.Name"),
		},
		{
			Name:        "namespace_path",
			Type:        proto.ColumnType_STRING,
			Description: "The path of the namespace to which the project belongs.",
			Transform:   transform.FromField("Namespace.Path"),
		},
		{
			Name:        "namespace_full_path",
			Type:        proto.ColumnType_STRING,
			Description: "The full path of the namespace to which the project belongs.",
			Transform:   transform.FromField("Namespace.FullPath"),
		},
		{
			Name:        "namespace_kind",
			Type:        proto.ColumnType_STRING,
			Description: "The kind of the namespace to which the project belongs.",
			Transform:   transform.FromField("Namespace.Kind"),
		},
		{
			Name:        "resolve_outdated_diff_discussions",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if outdated diff discussions should be resolved.",
		},
		{
			Name:        "import_status",
			Type:        proto.ColumnType_STRING,
			Description: "Status of project import.",
		},
		{
			Name:        "import_error",
			Type:        proto.ColumnType_STRING,
			Description: "Error of importing project (if any).",
		},
		{
			Name:        "license_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for the license of the project.",
			Transform:   transform.FromField("LicenseURL"),
		},
		{
			Name:        "license",
			Type:        proto.ColumnType_STRING,
			Description: "The projects license type.",
			Transform:   transform.FromField("License.Name"),
		},
		{
			Name:        "license_key",
			Type:        proto.ColumnType_STRING,
			Description: "The projects license spdx id/key.",
			Transform:   transform.FromField("License.Key"),
		},
		{
			Name:        "shared_runners_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has shared runners enabled.",
		},
		{
			Name:        "runners_token",
			Type:        proto.ColumnType_STRING,
			Description: "The token used for runners by the project.",
		},
		{
			Name:        "public_jobs",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project has/allows public jobs.",
		},
		{
			Name:        "allow_merge_on_skipped_pipeline",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if merges are allowed if the pipeline is skipped.",
		},
		{
			Name:        "only_allow_merge_if_pipeline_succeeds",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if merges are only allowed when the pipeline succeeds.",
		},
		{
			Name:        "only_allow_merge_if_all_discussions_are_resolved",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if merges are only allowed when all discussions are resolved.",
		},
		{
			Name:        "remove_source_branch_after_merge",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if source branches are removed after merge by default on the project.",
		},
		{
			Name:        "repository_storage",
			Type:        proto.ColumnType_STRING,
			Description: "The type of storage used by the repository.",
		},
		{
			Name:        "merge_method",
			Type:        proto.ColumnType_STRING,
			Description: "The projects default merge method (merge, squash, rebase, etc).",
		},
		{
			Name:        "fork_parent_id",
			Type:        proto.ColumnType_INT,
			Description: "ID of the fork parent.",
			Transform:   transform.FromField("ForkedFromProject.ID"),
		},
		{
			Name:        "fork_parent_name",
			Type:        proto.ColumnType_STRING,
			Description: "Full name of the fork parent.",
			Transform:   transform.FromField("ForkedFromProject.NameWithNamespace"),
		},
		{
			Name:        "fork_parent_path",
			Type:        proto.ColumnType_STRING,
			Description: "Full path of the fork parent.",
			Transform:   transform.FromField("ForkedFromProject.PathWithNamespace"),
		},
		{
			Name:        "fork_parent_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the fork parent.",
			Transform:   transform.FromField("ForkedFromProject.WebURL"),
		},
		{
			Name:        "mirror",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the project is a mirror",
		},
		{
			Name:        "mirror_user_id",
			Type:        proto.ColumnType_INT,
			Description: "ID of the user whom configured the mirror",
		},
		{
			Name:        "mirror_trigger_builds",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the mirror can trigger builds.",
		},
		{
			Name:        "only_mirror_protected_branches",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if only protected branches are mirrored.",
		},
		{
			Name:        "mirror_overwrites_diverged_branches",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the mirror can overwrite diverged branches.",
		},
		{
			Name:        "autoclose_referenced_issues",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if referenced issues will be automatically closed by merges in the project.",
			Transform:   transform.FromField("AutocloseReferencedIssues"),
		},
		{
			Name:        "ci_forward_deployment_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if ci forward deployments are enabled.",
			Transform:   transform.FromField("CIForwardDeploymentEnabled"),
		},
		{
			Name:        "ci_config_path",
			Type:        proto.ColumnType_STRING,
			Description: "The path of the CI configuration.",
			Transform:   transform.FromField("CIConfigPath"),
		},
		{
			Name:        "ci_separated_caches",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the CI uses separate caches.",
			Transform:   transform.FromField("CISeperateCache"),
		},
	}
}
