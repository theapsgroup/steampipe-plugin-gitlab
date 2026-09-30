# Table: gitlab_project_deploy_key

The `gitlab_project_deploy_key` table can be used to query deploy keys associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain deploy keys for a specific project

```sql
select
  id,
  title,
  can_push,
  created_at
from
  gitlab_project_deploy_key
where
  project_id = 1258;
```

### Find deploy keys with write (push) access

```sql
select
  project_id,
  title
from
  gitlab_project_deploy_key
where
  project_id = 1258
  and can_push = true;
```
