# Table: gitlab_project_approval

The `gitlab_project_approval` table can be used to query the merge request approval configuration of a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain the approval configuration for a specific project

```sql
select
  project_id,
  approvals_before_merge,
  merge_requests_author_approval,
  merge_requests_disable_committers_approval,
  reset_approvals_on_push,
  disable_overriding_approvers_per_merge_request
from
  gitlab_project_approval
where
  project_id = 1258;
```

### Flag projects that allow authors to approve their own merge requests

```sql
select
  project_id
from
  gitlab_project_approval
where
  project_id = 1258
  and merge_requests_author_approval = true;
```
