# Table: gitlab_my_event

The `gitlab_my_event` table can be used to query information about your activity on a gitlab system.

## Examples

### List all issues worked on in the last week

```sql
select
  project.full_path, 
  event.target_iid, 
  event.action_name
from
  gitlab_my_event as event, 
  gitlab_project as project
where
  event.project_id = project.id 
  and target_type = 'Issue' 
  and event.created_at > current_date - interval '7 days'
```

### Get activity counts by project over the last 30 days

```sql
select
  project.full_path, 
  count(project.full_path) as events
from
  gitlab_my_event as event, 
  gitlab_project as project
where
  event.project_id = project.id 
  and event.created_at > current_date - interval '30 days'
group by
  project.full_path
order by
  events;
```
