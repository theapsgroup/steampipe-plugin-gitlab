# Table: gitlab_merge_request_change

The `gitlab_merge_request_change` table can be used to view all changes associated with a single merge request.

However, **you must specify** both an `iid` of a merge request as well as it's `project_id` in the where or join clause.

## Examples

### Obtain all changes associated to a specific merge request

```sql
select
  old_path,
  new_path,
  a_mode,
  b_mode,
  diff,
  new_file,
  renamed_file,
  deleted_file
from
  gitlab_merge_request_change
where
  iid = 123
and
  project_id = 42;
```
