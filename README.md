# GitLab plugin for Steampipe

## Query GitLab with SQL

Use SQL to query GitLab. Example:

```sql
select * from gitlab_project
```

## Get Started

Build & Installation from source:

```shell
go build -o steampipe-plugin-gitlab.plugin

mv steampipe-plugin-gitlab.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/gitlab@latest/steampipe-plugin-gitlab.plugin

cp config/gitlab.spc ~/.steampipe/config/gitlab.spc
```

Configuration is preferably done by ensuring you have the following Environment Variables set:

- `GITLAB_ADDR` for the address of your GitLab API endpoint (e.g `https://gitlab.mycompany.com/api/v4`)
- `GITLAB_TOKEN` for the API token used to access GitLab (private or personal access tokens accepted)

These can also be set in the configuration file:
`vi ~/.steampipe/config/gitlab.spc` 

## Documentation

Further documentation can he [found here](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/docs/index.md)

## Credits

GitLab API Wrapper [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) this is licensed separately using this [Apache License](https://github.com/xanzy/go-gitlab/blob/master/LICENSE)