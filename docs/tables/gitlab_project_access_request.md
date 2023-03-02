# Table: gitlab_project_access_request

The `gitlab_project_access_request` table can be used to query information about access requests for a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all access requests for a specific project

```sql
select
  id,
  username,
  name,
  state,
  access_level,
  created_at,
  requested_at
from
  gitlab_project_access_request
where
  project_id = 45453535;
```

### Get a specific access request by id

```sql
select
  username,
  requested_at
from
  gitlab_project_access_request
where
  project_id = 45453535
and
  id = 873;
```