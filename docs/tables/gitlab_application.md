# Table: gitlab_application
The `gitlab_application` table can be used to query information about OAuth applications within the GitLab instance.

> THIS TABLE WILL ONLY RETURN DATA FOR ADMINISTRATORS

## Examples

### List all OAuth applications

```sql
select
  id,
  application_id,
  application_name,
  callback_url,
  confidential
from
  gitlab_application;
```