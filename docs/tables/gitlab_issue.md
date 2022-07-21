# Table: gitlab_issue

Issues are used to track bugs, feature requests, tasks, etc on GitLab.

The `gitlab_issue` table can be used to query information against all issues in the GitLab instance.

> Note: When used with the [Public GitLab](https://gitlab.com) you must specify an `=` qualifier for at least one of the following fields.
> - `assignee`
> - `assignee_id`
> - `author_id`
> - `project_id`
>
> This is to prevent attempting to return **ALL** public issues which would result in an error.

## Examples

### Obtain all issues

```sql
select
  *
from
  gitlab_issue;
```

### Obtain a list of Confidential Issues

```sql
select
  *
from
  gitlab_issue
where 
  confidential = true;
```

### Obtain counts of issues by state

```sql
select
  state
  count(*) as quantity
from
  gitlab_issue
group by
  state
```

### Get all issues for your projects
```sql
select 
  p.name as project_name,
  p.id as project_id,
  i.title as issue_title,
  i.id as issue_id,
  i.description as description,
  i.state as state
from 
  gitlab_my_project as p,
  gitlab_issue as i
where
  p.id = i.project_id
```
