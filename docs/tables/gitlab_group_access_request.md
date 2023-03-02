# Table: gitlab_group_access_request

The `gitlab_group_access_request` table can be used to query information about access requests for a specific group.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### List all access requests for a specific group

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
  gitlab_group_access_request
where
  group_id = 14597683;
```

### Get a specific access request by id

```sql
select
  username,
  requested_at
from
  gitlab_group_access_request
where
  group_id = 14597683
and
  id = 132;
```