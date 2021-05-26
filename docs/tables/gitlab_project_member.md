# Table: gitlab_project_member

Query project members for a specific project, **you must specify** a `project_id` in the where or join clause of the query.

## Examples

### List all members for a specific project

```sql
select
  *
from
  gitlab_project_member
where project_id = 123;
```