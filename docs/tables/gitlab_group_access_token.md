# Table: gitlab_group_access_token

The `gitlab_group_access_token` table can be used to query group access tokens associated with a specific group.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### Obtain access tokens for a specific group

```sql
select
  id,
  name,
  scopes,
  active,
  revoked,
  expires_at
from
  gitlab_group_access_token
where
  group_id = 42;
```
