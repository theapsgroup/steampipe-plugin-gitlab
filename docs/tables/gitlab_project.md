# Table: gitlab_project

The `gitlab_project` table will obtain information from all projects the user would be able to see, public/associated.

> Note: It's recommended that you use the `gitlab_my_project` table for performance.
>
>It is not advised to use this table if you're using the hosted GitLab.com instance since this will attempt to obtain **ALL** public projects.

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