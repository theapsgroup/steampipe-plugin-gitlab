# Table: gitlab_project

The `gitlab_project` table will obtain information from all projects the user would be able to see, public/associated.

> Note: When used with the [Public GitLab](https://gitlab.com) you must specify an `=` qualifier for at least one of the following fields.
> - `id`
> - `owner_id`
> - `owner_username`
>
> This is to prevent attempting to return ALL public projects which would result in an error.

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