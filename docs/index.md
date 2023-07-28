---
organization: The APS Group
category: ["software development"]
icon_url: "/images/plugins/theapsgroup/gitlab.svg"
brand_color: "#FCA121"
display_name: "GitLab"
short_name: "gitlab"
description: "Steampipe plugin for querying GitLab Repositories, Users and other resources."
og_description: Query GitLab with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/theapsgroup/gitlab-social-graphic.png"
---

# GitLab + Turbot Steampipe

[GitLab](https://about.gitlab.com/) is a provider of Internet hosting for software development and version control using Git. It offers the distributed version control and source code management (SCM) functionality of Git, plus its own features.

[Steampipe](https://steampipe.io/) is an open source CLI for querying cloud APIs using SQL from [Turbot](https://turbot.com/)

## Documentation

- [Table definitions / examples](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables)

## Get started

### Install

Download and install the latest GitLab plugin:

```shell
steampipe plugin install theapsgroup/gitlab
```

### Configuration

Installing the latest GitLab plugin will create a config file (`~/.steampipe/config/gitlab.spc`) with a single connection named `gitlab`:

```hcl
connection "gitlab" {
  plugin = "theapsgroup/gitlab"

  # The baseUrl of your GitLab Instance API (ignore if set in GITLAB_ADDR env var)
  # baseurl = "https://gitlab.company.com/api/v4"

  # Access Token for which to use for the API (ignore if set in GITLAB_TOKEN env var)
  # token = "x11x1xXxXx1xX1Xx11"
}
```

- `token` - [Personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for your GitLab account. This can also be set via the `GITLAB_TOKEN` environment variable.
- `baseurl` - GitLab URL (e.g. `https://gitlab.company.com/api/v4`). Not required for GitLab cloud. This can also be via the `GITLAB_ADDR` environment variable.

#### Configuration file example

```hcl
connection "gitlab" {
  plugin  = "theapsgroup/gitlab"
  baseurl = "https://gitlab.mycompany.com/api/v4"
  token   = "f7Ea3C3ojOY0GLzmhS5kE"
}
```

## Get involved

- Open source: https://github.com/theapsgroup/steampipe-plugin-gitlab
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
