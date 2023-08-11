# Table: gitlab_my_events

The `gitlab_my_events` table can be used to query information about your activity on a gitlab system.

## Examples

### Get all activity you're associated with
```sql
select
  *
from
  gitlab_my_events;
```

### Get all activity since 2023-08-01
```sql
select
  *
from
  gitlab_my_events
where
  created_at > '2023-08-01';
```
