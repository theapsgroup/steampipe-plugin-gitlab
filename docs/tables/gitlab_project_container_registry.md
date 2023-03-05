# Table: gitlab_project_container_registry

The `gitlab_project_container_registry` table can be used to query information about container registries for a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all container registries for a specific project

```sql
select
  id,
  name,
  path,
  location,
  created_at,
  cleanup_policy_started_at
from
  gitlab_project_container_registry
where
  project_id = 45453535;
```
