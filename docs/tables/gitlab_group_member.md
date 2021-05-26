# Table: gitlab_group_member

Query group members for a specific group, **you must specify** a `group_id` in the where or join clause of the query.

## Examples

### List all members for a specific group

```sql
select
  *
from
  gitlab_group_member
where group_id = 123;
```