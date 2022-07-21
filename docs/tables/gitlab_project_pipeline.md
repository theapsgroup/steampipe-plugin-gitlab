# Table: gitlab_project_pipeline

The `gitlab_project_pipeline` table can be used to query information about pipelines on a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List pipeline information for a project

```sql
select
  *
from
  gitlab_project_pipeline
where
  project_id = 123;
```
