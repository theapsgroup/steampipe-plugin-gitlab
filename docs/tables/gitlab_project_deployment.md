# Table: gitlab_project_deployment

The `gitlab_project_deployment` table can be used to obtain information about deployments associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all deployments for a specific project

```sql
select
  id,
  iid,
  ref,
  sha,
  status,
  created_at,
  updated_at,
  user_id,
  user_username,
  environment_id,
  environment_name,
  deployable_id,
  deployable_status,
  deployable_stage,
  deployable_name,
  deployable_ref,
  deployable_commit_id,
  deployable_pipeline_id
from
  gitlab_project_deployment
where
  project_id = 14597683;
```

### Get a specific deployment for a project

```sql
select
  id,
  iid,
  ref,
  sha,
  status,
  created_at,
  updated_at,
  user_id,
  user_username,
  environment_id,
  environment_name,
  deployable_id,
  deployable_status,
  deployable_stage,
  deployable_name,
  deployable_ref,
  deployable_commit_id,
  deployable_pipeline_id
from
  gitlab_project_deployment
where
  project_id = 14597683
and
  id = 1486132;
```

### Get information about the commit associated with a specific deployment

```sql
select
  d.id,
  d.environment_name,
  c.author_email,
  c.short_id,
  c.title,
  c.message
from
  gitlab_project_deployment d
left outer join 
  gitlab_commit c
on
  d.project_id = c.project_id
and
  d.deployable_commit_id = c.id
where
  d.project_id = 14597683;
```
