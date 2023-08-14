# Table: gitlab_my_event

The `gitlab_my_event` table can be used to query information about your activity on a gitlab system.

## Examples

### List all event activity associated with the logged in user

```sql
select
  *
from
  gitlab_my_event;
```

### Get all activity since 2023-08-01
```sql
select
  *
from
  gitlab_my_event
where
  created_at > '2023-08-01';
```
