---
organization: The APS Group
category: ["software development"]
icon_url: "/images/plugins/theapsgroup/gitlab.svg"
brand_color: "#003A75"
display_name: "GitLab"
short_name: "gitlab"
description: "Steampipe plugin for querying GitLab Repositories, Users and other resources."
og_description: Query GitLab with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/theapsgroup/gitlab-social-graphic.png"
---

# GitLab + Turbot Steampipe

[GitLab](https://about.gitlab.com/) is a provider of Internet hosting for software development and version control using Git. It offers the distributed version control and source code management (SCM) functionality of Git, plus its own features.

[Steampipe](https://steampipe.io/) is an open source CLI for querying cloud APIs using SQL from [Turbot](https://turbot.com/)

## Getting Started

### Build & Installation

Currently, you will need to build and install this plugin manually via the following:

```shell
go build -o steampipe-plugin-gitlab.plugin

mv steampipe-plugin-gitlab.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/gitlab@latest/steampipe-plugin-gitlab.plugin

cp config/gitlab.spc ~/.steampipe/config/gitlab.spc
```

### Prerequisites

- GitLab (either hosted or self-hosted)
- GitLab Token (either private or [personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html))

### Configuration

The preferred option is to use Environment Variables for configuration.

However, you can configure in the `~./steampipe/config/gitlab.spc` (this will take precedence).

Environment Variables:
- `GITLAB_ADDR` for the base url for the API endpoint (ex: `https://gitlab.mycompany.com/api/v4`)
- `GITLAB_TOKEN` for the API token (ex: `f7Ea3C3ojOY0GLzmhS5kE`)

Configuration File:

```hcl
connection "gitlab" {
  plugin  = "theapsgroup/gitlab"
  baseurl = "https://gitlab.mycompany.com/api/v4"
  token   = "f7Ea3C3ojOY0GLzmhS5kE"
}
```

### Testing

A quick test can be performed from your terminal with:

```shell
steampipe query "select * from gitlab_version"
```

## Tables

The following tables are currently available, follow links for more information.

- [gitlab_version](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/tables/gitlab_version.md)
- [gitlab_user](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/tables/gitlab_user.md)
- [gitlab_group](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/tables/gitlab_group.md)
- [gitlab_project](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/tables/gitlab_project.md)
- [gitlab_issue](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/tables/gitlab_issue.md)