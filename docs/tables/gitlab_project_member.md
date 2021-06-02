# Table: gitlab_project_member

A project member is a user that is associated to a specific project.

The `gitlab_project_member` table can be used to query information members of a specific project.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all members for a specific project

```sql
select
  *
from
  gitlab_project_member
where project_id = 123;
```