# Table: gitlab_user

Obtaining information about Users on the GitLab instance.

## Examples

### List all users

```sql
select
  *
from
  gitlab_user;
```

### Obtain a list of usernames that are currently blocked

```sql
select
  username,
  state
from
  gitlab_user
where
  state = 'blocked';
```

### Obtain a list of GitLab admins

```sql
select
  username,
  name,
  state,
  email
from
  gitlab_user
where
  is_admin = true;
```