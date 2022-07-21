# Table: gitlab_user

Obtaining information about Users on the GitLab instance.

> Note: When used with the [Public GitLab](https://gitlab.com) you must specify an `=` qualifier for at least one of the following fields.
> - `id`
> - `username`
>
> This is to prevent attempting to return **ALL** users which would result in an error.

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