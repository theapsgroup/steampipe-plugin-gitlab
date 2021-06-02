# Table: gitlab_group

A group is a collection of projects & members.

The `gitlab_group` table can be used to query groups (where you are a member of for authenticated user, unless authenticated user is an admin in which case, all groups will be available).

## Examples

### Get all Groups

```sql
select
  *
from
  gitlab_group;
```

### Get top level groups

```sql
select
  *
from
  gitlab_group
where
  parent_id is null;
```

### Get private groups

```sql
select
  *
from
  gitlab_group
where
  visibility = 'private';
```

### Obtain a count of different visibility levels 

```sql
select
  visibility,
  count(id) as group_count
from
  gitlab_group
group by
  visibility
```