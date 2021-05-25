# Table: gitlab_project_merge_request

Query a specific projects merge requests, **you must specify** a `project_id` in the where or join clause.

## Examples

### List of merge requests for a project

```sql
select
  *
from
  gitlab_project_merge_request
where
  project_id = 1;
```
