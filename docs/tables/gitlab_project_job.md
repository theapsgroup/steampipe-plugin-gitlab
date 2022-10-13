# Table: gitlab_project_job

The `gitlab_project_job` table can be used to query information about jobs on a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

+ List jobs for a project (adding `limit` is highly recommended):

```sql
select
  *
from
  gitlab_project_job
where
  project_id = '123'
limit 10;
```

+ List jobs where name contains `deploy` string:

```sql
select
  *
from
  gitlab_project_job
where
  project_id = '123'
  and name like '%deploy%'
limit 10;
```

## Reference

+ https://docs.gitlab.com/ee/api/jobs.html
