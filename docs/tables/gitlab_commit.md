# Table: gitlab_commit

Query commits from a specific project, **you must specify** a `project_id` in the where or join clause.

## Examples

### List commits

```sql
select
  *
from
  gitlab_commit
where
  project_id = 1;
```

### List commits (by newest first)

```sql
select
  *
from
  gitlab_commit
where
  project_id = 1
order by
  created_at desc;
```