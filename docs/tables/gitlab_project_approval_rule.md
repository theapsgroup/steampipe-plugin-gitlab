# Table: gitlab_project_approval_rule

The `gitlab_project_approval_rule` table can be used to query merge request approval rules associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain approval rules for a specific project

```sql
select
  id,
  name,
  rule_type,
  approvals_required,
  eligible_approvers,
  users,
  groups
from
  gitlab_project_approval_rule
where
  project_id = 1258;
```

### Find "any approver" rules that require fewer than 2 approvals

```sql
select
  project_id,
  name,
  approvals_required
from
  gitlab_project_approval_rule
where
  project_id = 1258
  and rule_type = 'any_approver'
  and approvals_required < 2;
```
