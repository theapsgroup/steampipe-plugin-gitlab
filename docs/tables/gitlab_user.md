# Table: gitlab_user

Obtaining information about Users on the GitLab instance.

> Note: It is not advised to use this table without filtering for 
> specific user id or username if you're using the hosted GitLab.com instance since this will attempt to obtain **ALL** users.

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