# Table: gitlab_project_protected_branch

The `gitlab_project_protected_branch` table can be used to query information on protected branches associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain protected branches for a specific project

```sql
select
  id,
  name,
  allow_force_push,
  code_owner_approval_required,
  push_access_levels,
  merge_access_levels,
  unprotect_access_levels
from
  gitlab_project_protected_branch
where
  project_id = 1258;
```
