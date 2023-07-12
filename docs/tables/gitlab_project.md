# Table: gitlab_project

The `gitlab_project` table will obtain information from all projects the user would be able to see, public/associated.

> **Note**: When used with the [Public GitLab](https://gitlab.com) you must specify an `=` qualifier for at least one of the following fields.
> - `id`
> - `owner_id`
> - `owner_username`
>
> This is to prevent attempting to return ALL public projects which would result in an error.

## Examples

### Get all projects

Note: This query will not work on the hosted SaaS (See note at top of page).

```sql
select
  *
from
  gitlab_project;
```

### Get all projects for a specific owner

```sql
select
  id,
  namespace_full_path as project
from
  gitlab_project
where
  owner_username = 'test';
```

### Get top 10 projects based on stars

Note: This query will not work on the hosted SaaS (See note at top of page).

```sql
select
  namespace_full_path as project,
  star_count
from
  gitlab_project
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
  gitlab_project p
inner join
  gitlab_user u
on 
  p.creator_id = u.id
and
  u.username = 'test';
```