# Table: gitlab_project_merge_request

The `gitlab_project_merge_request` table can be used to query all merge requests against a specific project.

However, **you must specify** a `project_id` in the where or join clause.

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
