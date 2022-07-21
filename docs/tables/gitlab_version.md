# Table: gitlab_version

The `gitlab_version` table can be used to query version information about the GitLab instance.

> Note: Should only return a single row of data.

## Examples

### Get version information for the Gitlab instance

```sql
select
  *
from
  gitlab_version;
```