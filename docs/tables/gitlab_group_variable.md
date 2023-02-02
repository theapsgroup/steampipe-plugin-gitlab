# Table: gitlab_group_variable

The `gitlab_group_variable` table can be used to view information about variables within GitLab at the Group level.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### List all variables held against a group

```sql
select
  key,
  value,
  variable_type,
  environment_scope,
  masked,
  protected,
  raw
from 
  gitlab_group_variable
where
  group_id = 42;
```

### Get a specific variable for a group by key

```sql
select
  key,
  value,
  variable_type,
  environment_scope,
  masked,
  protected,
  raw
from 
  gitlab_group_variable
where
  project_id = 42
and
  key = 'VARIABLE_NAME';
```
