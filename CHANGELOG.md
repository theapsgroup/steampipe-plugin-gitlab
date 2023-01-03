## v0.2.0 [WIP]

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v5.0.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v500-2022-11-16)
- 

## v0.1.3 [2022-12-12]

_What's new?_

- Extended the `gitlab_project` & `gitlab_my_project` tables with namespace fields as below:
  - `namespace_id` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_name` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_kind` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_path` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_full_path` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)

## v0.1.2 [2022-10-13]

_What's new?_

- New table [gitlab_project_job](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_job) - Thanks [@hiepph](https://github.com/hiepph)

## v0.1.1 [2022-10-12]

_Enhancements_

- Updated: Recompiled with [golang version 1.19](https://tip.golang.org/doc/go1.19)
- Updated: Recompiled with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08)

_What's new?_

- New tables added
  - [gitlab_epic](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_epic) *Premium License Required* [#13](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/13)
  - [gitlab_group_iteration](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_group_iteration) *Premium License Required* [#13](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/13)
  - [gitlab_project_iteration](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_iteration) *Premium License Required* [#13](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/13)
  - [gitlab_group_push_rule](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_group_push_rule) *Premium License Required* [#19](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/19)
  - [gitlab_hook](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_hook) *Premium License Required* [#21](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/21)
  - [gitlab_project_protected_branch](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_protected_branch) [#18](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/18)
  - [gitlab_project_pages_domain](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_pages_domain) [#23](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/23)
  - [gitlab_setting](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_setting) [#17](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/17)
- New columns added to [gitlab_project](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project) & [gitlab_my_project](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_my_project) tables 
  - `issues_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `repository_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `merge_requests_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `forking_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `wiki_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `builds_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `snippets_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `pages_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `operations_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `analytics_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `topics` [#14](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/14)


## v0.1.0 [2022-05-05]

_Enhancements_

- Updated: Recompiled with [golang version 1.18](https://tip.golang.org/doc/go1.18)
- Updated: Recompiled with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30)

## v0.0.5 [2022-03-25]

_What's new?_

- New tables added
  - [gitlab_project_repository](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_repository)
  - [gitlab_project_repository_file](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_repository_file)

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23)
- Updated: Recompiled with [go-gitlab v0.55.0](https://github.com/xanzy/go-gitlab/releases/tag/v0.55.0)

## v0.0.4 [2021-11-29]

_Enhancements_

- Updated: Recompiled with [golang version 1.17](https://tip.golang.org/doc/go1.17)
- Updated: Recompiled with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22)
- Updated: Recompiled with [go-gitlab v0.52.2](https://github.com/xanzy/go-gitlab/releases/tag/v0.52.2)

## v0.0.3 [2021-09-16]

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v1.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v151--2021-09-13)
- Updated: Recompiled with [go-gitlab v0.50.4](https://github.com/xanzy/go-gitlab/releases/tag/v0.50.4)
- Updated: Added `commit_count`, `storage_size`, `repository_size`, `lfs_objects_size` & `job_artifacts_size` columns to `gitlab_project` & `gitlab_my_project` tables ([#5](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/5))

## v0.0.2 [2021-07-23]

_What's new?_

- Set default API Url to the hosted GitLab to prevent needing to manually define this as per Issue [#3](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/3)
- Utilising a new feature `optional qualifiers` to allow for faster targeted queries on certain tables `table_issue`, `table_merge_request` & `table_project`
- Enforced specific qualifiers **ONLY** when using the public hosted GitLab instance (to prevent errors / excessively long-running queries)
  - `table_issue` requires at least one of the following `assignee`, `assignee_id`, `author_id` or `project_id`
  - `table_merge_request` requires at least one of the following `assignee_id`, `author_id`, `reviewer_id` or `project_id`
  - `table_project` requires at least one of the following `owner_id` or `owner_username`
  - `table_user` requires at least one of the following `id` or `username`
