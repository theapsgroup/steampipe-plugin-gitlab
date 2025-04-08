package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectPagesDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_pages_domain",
		Description: "Obtain information about pages domains for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectPagesDomains,
		},
		Columns: projectPageColumns(),
	}
}

// Hydrate Functions
func listProjectPagesDomains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectPagesDomains", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectPagesDomains", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListPagesDomainsOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectPagesDomains", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		domains, resp, err := conn.PagesDomains.ListPagesDomains(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectPagesDomains", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain page domains for project_id %d\n%v", projectId, err)
		}

		for _, domain := range domains {
			d.StreamListItem(ctx, domain)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectPagesDomains", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectPagesDomains", "completed successfully")
	return nil, nil
}

// Column Function
func projectPageColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "domain",
			Type:        proto.ColumnType_STRING,
			Description: "The custom domain configured for the pages.",
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "The url configured for the domain (on protocol).",
		},
		{
			Name:        "auto_ssl_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if SSL Certificates are auto-generated for the custom pages domain.",
		},
		{
			Name:        "certificate_expired",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the certificate has expired.",
			Transform:   transform.FromField("Certificate.Expired"),
		},
		{
			Name:        "certificate_expiration",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the certificate expires",
			Transform:   transform.FromField("Certificate.Expiration"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project this custom pages domain belongs to - link `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
		{
			Name:        "verified",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the domain is verified.",
		},
		{
			Name:        "verification_code",
			Type:        proto.ColumnType_STRING,
			Description: "The verification code associated with the pages domain.",
		},
		{
			Name:        "enabled_until",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp at which the pages domain is disabled.",
		},
	}
}
