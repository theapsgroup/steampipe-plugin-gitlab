# Table: gitlab_epic

The `gitlab_epic` table can be used to query information about epics associated with a specific group.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### List all epics for a specific group

```sql
select
  id,
  title,
  description,
  state,
  author,
  due_date
from
  gitlab_epic
where
  group_id = 14597683;
```

### Get epics from a specific group and author

```sql
select
  id,
  title,
  description,
  state,
  due_date,
  upvotes,
  downvotes
from
  gitlab_epic
where
  group_id = 14597683
and
  author_id = 469853662;
```
