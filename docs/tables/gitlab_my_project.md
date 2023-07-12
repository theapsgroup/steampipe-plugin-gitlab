# Table: gitlab_my_project

The `gitlab_my_project` table can be used to query information against projects that the authenticated user is a member of.

## Examples

### Get all projects you're associated with
```sql
select
  *
from
  gitlab_my_project;
```

### Get top 10 projects you're associated with based on stars
```sql
select
  *
from
  gitlab_my_project
order by
  star_count desc
limit 10;  
```

### Get project creation information
```sql
select
  u.username as creator,
  p.full_path as project,
  p.created_at as created
from
  gitlab_my_project p
inner join
  gitlab_user u
on 
  p.creator_id = u.id;
```