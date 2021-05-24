# Table: gitlab_project

Query projects within the GitLab instance.

## Examples

### Get all projects
```sql
select
  *
from
  gitlab_project;
```

### Get top 10 projects based on stars
```sql
select
  *
from
  gitlab_project
order by
  star_count desc
limit 1;  
```

### Get project creation information
```sql
select
  u.username as creator,
  p.full_path as project,
  p.created_at as created
from
  gitlab_project p
inner join
  gitlab_user u
on 
  p.creator_id = u.id;
```