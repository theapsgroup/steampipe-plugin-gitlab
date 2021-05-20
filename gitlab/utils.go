package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	api "github.com/xanzy/go-gitlab"
	"os"
	"strings"
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

func sanitizeUrl(url string) string {
	return strings.ReplaceAll(url, "//","/")
}