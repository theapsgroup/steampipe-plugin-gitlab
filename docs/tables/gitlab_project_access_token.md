# Table: gitlab_project_access_token

The `gitlab_project_access_token` table can be used to query project access tokens associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain access tokens for a specific project

```sql
select
  id,
  name,
  scopes,
  active,
  revoked,
  expires_at
from
  gitlab_project_access_token
where
  project_id = 1258;
```

### Find active project access tokens that never expire

```sql
select
  project_id,
  name,
  scopes
from
  gitlab_project_access_token
where
  project_id = 1258
  and active = true
  and expires_at is null;
```
