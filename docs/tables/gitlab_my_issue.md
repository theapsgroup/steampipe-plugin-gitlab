# Table: gitlab_my_issue

Issues are used to track bugs, feature requests, tasks, etc on GitLab.

The `gitlab_my_issue` table can be used to query information against issues created by or assigned to the authenticated user in the GitLab instance.

## Examples

### Obtain all your issues

```sql
select
  *
from
  gitlab_my_issue;
```
