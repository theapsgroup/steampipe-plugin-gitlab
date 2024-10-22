# Table: gitlab_user_event

The `gitlab_user_event` table can be used to query information about user activity on a gitlab system.

## Examples

### List all issues worked on in the last week by user alice

```sql
select
  project.full_path, 
  event.target_iid, 
  event.action_name
from
  gitlab_user_event as event, 
  gitlab_project as project, 
  gitlab_user as u
where
  event.project_id = project.id 
  and target_type = 'Issue' 
  and event.created_at > current_date - interval '7 days' 
  and event.author_id = u.id 
  and u.username = 'alice';
```

### Get activity counts by project over the last 30 days by user bob

```sql
select
  project.full_path, 
  count(project.full_path) as events
from
  gitlab_user_event as event, 
  gitlab_project as project,
  gitlab_user as u
where
  event.project_id = project.id 
  and event.created_at > current_date - interval '30 days' 
  and event.author_id = u.id 
  and u.username = 'bob'
group by
  project.full_path
order by
  events;
```
