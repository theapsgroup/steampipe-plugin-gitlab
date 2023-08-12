# Table: gitlab_user_event

The `gitlab_user_event` table can be used to query information about user activity on a gitlab system.

## Examples

### Get all activity
```sql
select
  *
from
  gitlab_user_event;
```

### Get all activity since 2023-08-01
```sql
select
  *
from
  gitlab_user_event
where
  created_at > '2023-08-01';
```
