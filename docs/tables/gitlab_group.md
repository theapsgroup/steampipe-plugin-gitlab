# Table: gitlab_group

Query the Groups (& SubGroups) of GitLab.

## Examples

### Get all Groups

```sql
select
  *
from
  gitlab_group;
```

### Get all top level groups

```sql
select
  *
from
  gitlab_group
where
  parent_id is null;
```

### Get all private groups

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