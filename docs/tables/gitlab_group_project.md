# Table: gitlab_group_project

The `gitlab_group_project`  table will obtain information from all projects associated to the group (& it's sub-groups).

## Examples

### List all projects for a group and its subgroups

```sql
select
  id,
  name,
  full_path,
  description,
  default_branch,
  public,
  visibility,
  archived.
  web_url
from
  gitlab_group_project
where
  group_id = 1234;
```

### List all projects for a specific group only

```sql
select
  id,
  name,
  full_path,
  description,
  default_branch,
  public,
  visibility,
  archived.
  web_url
from
  gitlab_group_project
where
  group_id = 1234
and
  namespace_id = 1234;
```
