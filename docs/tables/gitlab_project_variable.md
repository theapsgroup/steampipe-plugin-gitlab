# Table: gitlab_project_variable

The `gitlab_project_variable` table can be used to view information about variables within GitLab at the Project level.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all variables for a project

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
  gitlab_project_variable
where
  project_id = 173;
```

### Get a specific variable for a project by key

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
  gitlab_project_variable
where
  project_id = 173
and
  key = 'VARIABLE_NAME';
```
