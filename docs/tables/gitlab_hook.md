# Table: gitlab_group_hook

The `gitlab_group_hook` table can be used to query information about the webhooks in a specific group.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### List hooks for a specific group

```sql
select
  id,
  url,
  push_events,
  push_events_branch_filter,
  issues_events,
  confidential_issues_events,
  merge_requests_events,
  tag_push_events,
  note_events,
  confidential_note_events,
  job_events,
  pipeline_events,
  wiki_page_events,
  deployment_events,
  releases_events,
  subgroup_events,
  created_at
from
  gitlab_group_hook
where
  group_id = 14597683;
```
