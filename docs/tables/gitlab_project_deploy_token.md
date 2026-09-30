# Table: gitlab_project_deploy_token

The `gitlab_project_deploy_token` table can be used to query deploy tokens associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain deploy tokens for a specific project

```sql
select
  id,
  name,
  username,
  scopes,
  revoked,
  expired,
  expires_at
from
  gitlab_project_deploy_token
where
  project_id = 1258;
```
