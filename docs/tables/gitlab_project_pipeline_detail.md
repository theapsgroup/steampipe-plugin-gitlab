# Table: gitlab_project_pipeline_detail

Obtain detailed pipeline information for a single pipeline from a single project, **you must specify** a `project_id` and an `id` (for the pipeline) in the where or join clause.

## Examples

### Obtain information for a specific pipeline

```sql
select
  *
from
  gitlab_project_pipeline
where
  project_id = 123
and
  id = 12345;
```