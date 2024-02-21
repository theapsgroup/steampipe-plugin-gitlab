# Table: gitlab_project_pipeline_detail

The `gitlab_project_pipeline_detail` table can be used to query detailed information about a specific pipeline instance on a specific project.

However, **you must specify** a `project_id` and an `id` (for the pipeline) in the where or join clause.

## Examples

### Obtain information for a specific pipeline

```sql
select
  *
from
  gitlab_project_pipeline_detail
where
  project_id = 123
and
  id = 12345;
```
