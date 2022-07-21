![image](https://github.com/theapsgroup/steampipe-plugin-gitlab/raw/main/docs/gitlab-plugin-social-graphic.png)

# GitLab plugin for Steampipe

## Query GitLab with SQL

Use SQL to query GitLab. Example:

```sql
select * from gitlab_project
```

## Get Started

### Installation

```shell
steampipe plugin install theapsgroup/gitlab
```

Or if you prefer, you can clone this repository and build/install from source directly.

```shell
go build -o steampipe-plugin-gitlab.plugin

mv steampipe-plugin-gitlab.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/gitlab@latest/steampipe-plugin-gitlab.plugin

cp config/gitlab.spc ~/.steampipe/config/gitlab.spc
```

Configuration is preferably done by ensuring you have the following Environment Variables set:

- `GITLAB_ADDR` for the address of your GitLab API endpoint (e.g `https://gitlab.mycompany.com/api/v4`)
- `GITLAB_TOKEN` for the API token used to access GitLab (private or [personal access](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) tokens accepted)

These can also be set in the configuration file:
`vi ~/.steampipe/config/gitlab.spc` 

> Note: If `GITLAB_ADDR` (`baseurl` in the config file) is not set it will default to the public cloud-hosted GitLab instance -> https://gitlab.com/api/v4

## Documentation

Further documentation can he [found here](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/index.md)

## Credits

GitLab API Wrapper [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) (licensed separately using this [Apache License](https://github.com/xanzy/go-gitlab/blob/master/LICENSE))