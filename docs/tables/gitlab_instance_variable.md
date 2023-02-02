# Table: gitlab_instance_variable

The `gitlab_instance_variable` table can be used to view variables that apply to the Self-Hosted GitLab instance, this feature isn't available on the public hosted GitLab & will return empty if queried there.

## Examples

### List variables associated with the entire instance.

```sql
select
  key,
  value,
  variable_type,
  masked,
  protected,
  raw
from 
  gitlab_instance_variable
```
