package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"os"
	"strings"
	"time"
)

const publicGitLabBaseUrl = "https://gitlab.com/api/v4"

func connect(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	baseUrl := os.Getenv("GITLAB_ADDR")
	token := os.Getenv("GITLAB_TOKEN")

	gitlabConfig := GetConfig(d.Connection)
	if &gitlabConfig != nil {
		if gitlabConfig.BaseUrl != nil {
			baseUrl = *gitlabConfig.BaseUrl
		}
		if gitlabConfig.Token != nil {
			token = *gitlabConfig.Token
		}
	}

	if baseUrl == "" {
		plugin.Logger(ctx).Info(fmt.Sprintf("no baseUrl was passed in - using %s", publicGitLabBaseUrl))
		baseUrl = publicGitLabBaseUrl // Default to public GitLab if not set, rather than return an error.
	}
	if token == "" {
		plugin.Logger(ctx).Error("no token provided in configuration file nor GITLAB_TOKEN environment variable")
		return nil, fmt.Errorf("GitLab Private/Personal Access Token must be set either in GITLAB_TOKEN env var or in connection config file")
	}

	plugin.Logger(ctx).Debug("attempting to create new client", "baseUrl", baseUrl)
	client, err := api.NewClient(token, api.WithBaseURL(baseUrl))
	if err != nil {
		plugin.Logger(ctx).Error("unable to create client", "baseUrl", baseUrl, "error", err)
		return nil, err
	}

	return client, nil
}

// sanitizeUrl is a util func for stripping accidental double slashes in urls
func sanitizeUrl(url string) string {
	return strings.ReplaceAll(url, "//", "/")
}

// isoTimeTransform is a transformation func for *gitlab.ISOTime to *time.Time
func isoTimeTransform(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return nil, nil
	}
	x := input.Value.(*api.ISOTime).String()
	return time.Parse("2006-01-02", x)
}

// parseAccessLevel is a util func for returning a string description based on integer for access level
func parseAccessLevel(input int) string {
	switch input {
	case 0:
		return "No Permissions"
	case 5:
		return "Minimal Access"
	case 10:
		return "Guest"
	case 20:
		return "Reporter"
	case 30:
		return "Developer"
	case 40:
		return "Maintainer"
	case 50:
		return "Owner"
	default:
		return "No Permissions"
	}
}

// isPublicGitLab is a util function to determine if the API is the public GitLab
func isPublicGitLab(d *plugin.QueryData) bool {
	cfg := GetConfig(d.Connection)
	if &cfg != nil {
		if cfg.BaseUrl != nil && *cfg.BaseUrl == publicGitLabBaseUrl {
			return true
		}
	}

	return false
}
