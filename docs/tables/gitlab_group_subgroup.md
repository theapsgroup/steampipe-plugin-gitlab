# Table: gitlab_group_subgroup

The `gitlab_group_subgroup` table will obtain information about subgroups for a specific group.

However, **you must specify** a `parent_id` in the where or join clause.

## Examples

### List all subgroups of a group

```sql
select
  id,
  name,
  full_path,
  description,
  visibility,
  parent_id,
  created_at
from
  gitlab_group_subgroup
where
  parent_id = 34234;
```
