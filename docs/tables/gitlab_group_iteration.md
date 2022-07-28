# Table: gitlab_group_iteration

The `gitlab_group_iteration` table can be used to query information about iterations for specific groups.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### List all iterations for a specific group

```sql
select
  id,
  group_id,
  title,
  description,
  due_date
from
  gitlab_group_iteration
where
  group_id = 14597683;
```

### Get a specific iteration by iid

```sql
select
  id,
  iid,
  group_id,
  title,
  description,
  due_date
from
  gitlab_group_iteration
where
  group_id = 14597683
and
  iid = 15;
```
