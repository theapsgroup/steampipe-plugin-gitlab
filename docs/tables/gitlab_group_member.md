# Table: gitlab_group_member

A group member is a user that is associated to a specific group.

The `gitlab_group_member` table can be used to query information members of a specific group.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### List all members for a specific group

```sql
select
  *
from
  gitlab_group_member
where group_id = 123;
```