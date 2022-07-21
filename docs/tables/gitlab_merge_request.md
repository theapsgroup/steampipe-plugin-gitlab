# Table: gitlab_merge_request

The `gitlab_merge_request` table can be used to query all merge requests in the GitLab instance.

> Note: When used with the [Public GitLab](https://gitlab.com) you must specify an `=` qualifier for at least one of the following fields.
> - `reviewer_id`
> - `assignee_id`
> - `author_id`
> - `project_id`
>
> This is to prevent attempting to return **ALL** public merge requests which would result in an error.

## Examples

### List all merge requests

```sql
select
  *
from
  gitlab_merge_request;
```

### List of merge requests for a specific project

```sql
select
  *
from
  gitlab_merge_request
where
  project_id = 12345;
```

### List of merge requests which have been merged in the last week

```sql
select
  project_id,
  title,
  state,
  target_branch,
  source_branch,
  merged_at,
  merged_by_username
from
  gitlab_merge_request
where
  state = 'merged'
and
  merged_at >= (current_date - interval '7' day);
```

### List all merge requests for your projects

```sql
select
  p.name as project_name,
  p.id as project_id,
  m.id as mr_id,
  m.title as mr_title,
  m.state as mr_state,
  m.assignee_username as mr_assignee
from
  gitlab_my_project as p,
  gitlab_merge_request as m
```
