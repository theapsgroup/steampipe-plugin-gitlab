# Table: gitlab_project_repository

The `gitlab_project_repository` can be used to list out the files/folders within the repository.

However, **you must specify** a `project_id` in the where or join clause.

## Examples

### List all files/folders for a project repository

```sql
select
  *
from
  gitlab_project_repository
where
  project_id = 123;  
```
