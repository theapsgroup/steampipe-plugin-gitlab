# Table: gitlab_merge_request

Query all Merge Requests in GitLab.

> Note: If you wish to obtain merge requests for a specific `project` 
> you should use the dedicated `gitlab_project_merge_request` table for better performance.

## Examples

### List all Merge Requests

```sql
select
  *
from
  gitlab_merge_request;
```

### List of Merge Requests which have been merged in the last week

```sql
select
  project_id,
  title,
  state,
  target_branch,
  source_branch,
  merged_at,
  merged_by_username
from
  gitlab_merge_request
where
  state = 'merged'
and
  merged_at >= (current_date - interval '7' day);
```
