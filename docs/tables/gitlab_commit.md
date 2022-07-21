# Table: gitlab_commit

A commit is a change-set to the code.

The `gitlab_commit` table can be used to query information about any commit.

However, **you must specify** a `project_id` in the where or join clause.

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

### Obtain an individual commit

```sql
select
  *
from
  gitlab_commit
where
  project_id = 1
and
  id = '73012177d1c8eb765bfd952ccfc50c679f147d12';
```

### Contributions by author

```sql
select
  author_email,
  count(*)
from
  gitlab_commit
where
  project_id = 1
group by
  author_email
order by
  count desc;
```