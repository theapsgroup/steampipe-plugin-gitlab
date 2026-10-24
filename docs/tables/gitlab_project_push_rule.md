# Table: gitlab_project_push_rule

The `gitlab_project_push_rule` table can be used to query the push rules configured for a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain push rules for a specific project

```sql
select
  id,
  reject_unsigned_commits,
  commit_committer_check,
  member_check,
  prevent_secrets
from
  gitlab_project_push_rule
where
  project_id = 1258;
```
