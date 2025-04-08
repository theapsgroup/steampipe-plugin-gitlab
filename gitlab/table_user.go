package gitlab

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableUser() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_user",
		Description: "Obtain information about users within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "username"}),
			Hydrate:    getUser,
		},
		Columns: userColumns(),
	}
}

// Hydrate Functions
func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listUsers", "started")
	if isPublicGitLab(d) {
		plugin.Logger(ctx).Error("listUsers", "Public GitLab requires an '=' qualifier for at least one of the following columns 'id', 'username' - none was provided")
		return nil, fmt.Errorf("when using the gitlab_user table with GitLab Cloud, `List` call requires an '=' qualifier for " +
			"one or more of the following columns: 'id', 'username'")
	}

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUsers", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	opt := &api.ListUsersOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 100,
	}}

	for {
		plugin.Logger(ctx).Debug("listUsers", "page", opt.Page, "perPage", opt.PerPage)
		users, resp, err := conn.Users.ListUsers(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listUsers", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain users\n%v", err)
		}

		for _, user := range users {
			d.StreamListItem(ctx, user)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listUsers", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listUsers", "completed successfully")
	return nil, nil
}

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getUser", "started")
	userId := int(d.EqualsQuals["id"].GetInt64Value())
	userName := d.EqualsQuals["username"].GetStringValue()

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getUser", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	if userId != 0 {
		var b = false
		opt := api.GetUsersOptions{WithCustomAttributes: &b}
		plugin.Logger(ctx).Debug("getUser", "filter[id]", userId)

		user, _, err := conn.Users.GetUser(userId, opt)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				plugin.Logger(ctx).Warn("getUser", "id", userId, "no user was found, returning empty result set")
				return nil, nil
			}
			plugin.Logger(ctx).Error("getUser", "id", userId, "error", err)
			return nil, fmt.Errorf("unable to obtain user with id %d\n%v", userId, err)
		}

		plugin.Logger(ctx).Debug("getUser", "completed successfully")
		return user, nil
	} else {
		opt := &api.ListUsersOptions{Username: &userName, ListOptions: api.ListOptions{Page: 1, PerPage: 50}}
		plugin.Logger(ctx).Debug("getUser", "filter[username]", userName)

		user, _, err := conn.Users.ListUsers(opt)
		if err != nil {
			plugin.Logger(ctx).Error("getUser", "userName", userName, "error", err)
			return nil, fmt.Errorf("unable to obtain user %s\n%v", userName, err)
		}

		if len(user) > 0 {
			plugin.Logger(ctx).Warn("getUser", fmt.Sprintf("found %d users - returning only first user", len(user)))
			plugin.Logger(ctx).Debug("getUser", "completed successfully")
			return user[0], nil
		}

		plugin.Logger(ctx).Debug("getUser", "completed successfully")
		return nil, nil
	}
}

// Column Functions
func userColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the user.",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The login/username of the user.",
		},
		{
			Name:        "email",
			Type:        proto.ColumnType_STRING,
			Description: "The primary email address of the user.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the user.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the user active, blocked, etc",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url for GitLab profile of user",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The timestamp when the user was created.",
		},
		{
			Name:        "bio",
			Type:        proto.ColumnType_STRING,
			Description: "The biography of the user.",
		},
		{
			Name:        "location",
			Type:        proto.ColumnType_STRING,
			Description: "The geographic location of the user.",
		},
		{
			Name:        "public_email",
			Type:        proto.ColumnType_STRING,
			Description: "The public email address of the user.",
		},
		{
			Name:        "skype",
			Type:        proto.ColumnType_STRING,
			Description: "The Skype address of the user.",
		},
		{
			Name:        "linkedin",
			Type:        proto.ColumnType_STRING,
			Description: "The LinkedIn account of the user.",
		},
		{
			Name:        "twitter",
			Type:        proto.ColumnType_STRING,
			Description: "The Twitter handle of the user.",
		},
		{
			Name:        "website_url",
			Type:        proto.ColumnType_STRING,
			Description: "The personal website of the user.",
			Transform:   transform.FromField("WebsiteURL"),
		},
		{
			Name:        "organization",
			Type:        proto.ColumnType_STRING,
			Description: "The organization of the user.",
		},
		{
			Name:        "ext_id",
			Type:        proto.ColumnType_STRING,
			Description: "The external ID of the user.",
			Transform:   transform.FromField("ExternUID"),
		},
		{
			Name:        "provider",
			Type:        proto.ColumnType_STRING,
			Description: "The external provider of the user.",
		},
		{
			Name:        "theme_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the users chosen theme.",
			Transform:   transform.FromField("ThemeID"),
		},
		{
			Name:        "last_activity_on",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The date user was last active.",
			Transform:   transform.FromField("LastActivityOn").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "color_scheme_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the users chosen color scheme.",
			Transform:   transform.FromField("ColorSchemeID"),
		},
		{
			Name:        "is_admin",
			Type:        proto.ColumnType_BOOL,
			Description: "Is the user an Administrator",
		},
		{
			Name:        "avatar_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the users avatar.",
			Transform:   transform.FromField("AvatarURL"),
		},
		{
			Name:        "can_create_group",
			Type:        proto.ColumnType_BOOL,
			Description: "The user has permissions to create groups.",
		},
		{
			Name:        "can_create_project",
			Type:        proto.ColumnType_BOOL,
			Description: "The user has permissions to create projects",
		},
		{
			Name:        "projects_limit",
			Type:        proto.ColumnType_INT,
			Description: "The limit of personal projects the user can create.",
		},
		{
			Name:        "current_sign_in_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The timestamp of users current signed in session.",
		},
		{
			Name:        "last_sign_in_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The timestamp of users last sign in.",
		},
		{
			Name:        "confirmed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The timestamp of user confirmation.",
		},
		{
			Name:        "two_factor_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Has the user enabled 2FA/MFA",
		},
		{
			Name:        "note",
			Type:        proto.ColumnType_STRING,
			Description: "The notes against the user.",
		},
		{
			Name:        "identities",
			Type:        proto.ColumnType_JSON,
			Description: "JSON Array of identity information for federated/IdP accounts",
		},
		{
			Name:        "external",
			Type:        proto.ColumnType_BOOL,
			Description: "Is the user an external entity",
		},
		{
			Name:        "private_profile",
			Type:        proto.ColumnType_BOOL,
			Description: "Is the users profile set to private.",
		},
		{
			Name:        "shared_runners_minutes_limit",
			Type:        proto.ColumnType_INT,
			Description: "Limit in minutes of time the user can utilise shared runner resources.",
		},
		{
			Name:        "extra_shared_runners_minutes_limit",
			Type:        proto.ColumnType_INT,
			Description: "Limit in minutes of extra time the user can utilise shared runner resources.",
		},
		{
			Name:        "using_license_seat",
			Type:        proto.ColumnType_BOOL,
			Description: "Is the user utilising a seat/slot on the license.",
		},
		{
			Name:        "custom_attributes",
			Type:        proto.ColumnType_JSON,
			Description: "JSON Array of custom attributes held against the user.",
		},
		{
			Name:        "bot",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if user account is a bot.",
		},
		{
			Name:        "job_title",
			Type:        proto.ColumnType_STRING,
			Description: "The users job title.",
		},
	}
}
