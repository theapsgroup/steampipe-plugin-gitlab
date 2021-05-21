package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"os"
	"strings"
	"time"
)

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
		return nil, fmt.Errorf("GitLab Base Address must be set either in GITLAB_ADDR env var or in connection config file")
	}
	if token == "" {
		return nil, fmt.Errorf("GitLab Private/Personal Access Token must be set either in GITLAB_TOKEN env var or in connection config file")
	}

	client, err := api.NewClient(token, api.WithBaseURL(baseUrl))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// sanitizeUrl is a util func for stripping accidental double slashes in urls
func sanitizeUrl(url string) string {
	return strings.ReplaceAll(url, "//","/")
}

// isoTimeTransform is a transformation func for *gitlab.ISOTime to *time.Time
func isoTimeTransform(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return nil, nil
	}
	x := input.Value.(*api.ISOTime).String()
	return time.Parse("2006-01-02", x)
}