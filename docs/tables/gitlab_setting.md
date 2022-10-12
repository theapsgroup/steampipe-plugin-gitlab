# Table: gitlab_setting

The `gitlab_setting` table is used to obtain the settings for the GitLab instance you're connected to.

> Note: When used with the [Public GitLab](https://gitlab.com/) may result in a 403 Forbidden error.

## Example

### Obtain settings

```sql
select
  *
from
  gitlab_setting;
```
