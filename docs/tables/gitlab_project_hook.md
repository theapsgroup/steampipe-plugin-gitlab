# Table: gitlab_project_hook

The `gitlab_project_hook` table can be used to query webhooks configured for a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### Obtain hooks for a specific project

```sql
select
  id,
  url,
  push_events,
  merge_requests_events,
  enable_ssl_verification
from
  gitlab_project_hook
where
  project_id = 1258;
```

### Find hooks that do not verify SSL

```sql
select
  project_id,
  id,
  url
from
  gitlab_project_hook
where
  project_id = 1258
  and enable_ssl_verification = false;
```
