## v0.0.3 [WIP]

_Enhancements_

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
