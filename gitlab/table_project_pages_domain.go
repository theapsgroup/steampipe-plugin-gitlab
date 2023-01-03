package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectPagesDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_pages_domain",
		Description: "Pages Domains associated with a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectPagesDomains,
		},
		Columns: []*plugin.Column{
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
		},
	}
}

func listProjectPagesDomains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListPagesDomainsOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		domains, resp, err := conn.PagesDomains.ListPagesDomains(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, domain := range domains {
			d.StreamListItem(ctx, domain)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
