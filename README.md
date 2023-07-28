![image](https://github.com/theapsgroup/steampipe-plugin-gitlab/raw/main/docs/gitlab-plugin-social-graphic.png)

# GitLab plugin for Steampipe

* **[Get started →](https://hub.steampipe.io/plugins/theapsgroup/gitlbb)**
* Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables)
* Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
* Get involved: [Issues](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io/downloads):

```shell
steampipe plugin install theapsgroup/gitlab
```

[Configure the plugin](https://hub.steampipe.io/plugins/theapsgroup/gitlab#configuration) using the configuration file:

```shell
vi ~/.steampipe/gitlab.spc
```

Or environment variables:

```shell
export GITLAB_TOKEN=f7Ea3C3ojOY0GLzmhS5kE
```

Start Steampipe:

```shell
steampipe query
```

Run a query:

```sql
select
  full_path,
  visibility,
  forks_count,
  star_count
from
  gitlab_my_project;
```

## Developing

Prerequisites:

* [Steampipe](https://steampipe.io/downloads)
* [Golang](https://golang.org/doc/install)
* GitLab (either hosted or self-hosted)
* GitLab Token (either private or [personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html))

Clone:

```sh
git clone https://github.com/theapsgroup/steampipe-plugin-gitlab.git
cd steampipe-plugin-gitlab
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/gitlab.spc
```

Try it!

```shell
steampipe query
> .inspect gitlab
```

Further reading:

* [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
* [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

All contributions are subject to the [Apache 2.0 open source license](https://github.com/theapsgroup/steampipe-plugin-gitlab/blob/main/LICENSE).

`help wanted` issues:

* [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
* [GitLab Plugin](https://github.com/theapsgroup/steampipe-plugin-gitlab/labels/help%20wanted)

## Credits

GitLab API Wrapper [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) (licensed separately using this [Apache License](https://github.com/xanzy/go-gitlab/blob/master/LICENSE))