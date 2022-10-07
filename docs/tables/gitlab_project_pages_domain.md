# Table: gitlab_project_pages_domain

The `gitlab_project_pages_domain` table can be used to query information on custom domains used for pages associated with a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all custom pages domains for a specific project

```sql
select
  domain,
  url,
  certificate_expiration,
  certificate_expired,
  auto_ssl_enabled
from
  gitlab_project_pages_domain
where
  project_id = 1337;
```
