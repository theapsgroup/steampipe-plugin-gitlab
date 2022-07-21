# Table: gitlab_branch

A branch is essentially a unique set of code changes with a unique name.

The `gitlab_branch` table can be used to query information about any branch.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List branches

```sql
select
  *
from
  gitlab_branch
where
 project_id = 1;
```

### Get branch information for a specific set of projects

```sql
select
  p.name as project_name,
  p.full_path as project_path,
  b.name as branch_name,
  b.default as is_default_branch,
  b.commit_short_id as commit_hash
from
  gitlab_branch b
inner join
  gitlab_my_project p
on
  b.project_id = p.id
where b.project_id in (
  select
    id
  from
    gitlab_my_project
  where
    full_path like '%service%'
);
```