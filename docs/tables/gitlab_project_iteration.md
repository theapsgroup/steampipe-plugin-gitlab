# Table: gitlab_project_iteration

The `gitlab_project_iteration` table can be used to query information about iterations for specific projects.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all iterations for a specific project

```sql
select
  id,
  project_id,
  title,
  description,
  due_date
from
  gitlab_project_iteration
where
  project_id = 14597683;
```

### Get a specific iteration by iid

```sql
select
  id,
  iid,
  project_id,
  title,
  description,
  due_date
from
  gitlab_project_iteration
where
  project_id = 14597683
and
  iid = 15;
```
