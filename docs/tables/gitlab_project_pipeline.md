# Table: gitlab_project_pipeline

Obtain basic pipeline information for all pipelines against a specific project, **you must specify** a `project_id` in the where or join clause.

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
