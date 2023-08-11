# Table: gitlab_user_events

The `gitlab_user_events` table can be used to query information about user activity on a gitlab system.

## Examples

### Get all activity
```sql
select
  *
from
  gitlab_user_events;
```

### Get all activity since 2023-08-01
```sql
select
  *
from
  gitlab_user_events
where
  created_at > '2023-08-01';
```
