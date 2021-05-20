# Table: gitlab_user

Obtaining information about Users on the GitLab instance.

## Columns

| Column | Description |
| - | - |
| id | GitLabs internal user ID. |
| username | The login/username of the user. |
| email | The primary email address of the user. |
| name | The name of the user. |
| state | The state of the user active, blocked, etc) |
| web_url | The url for GitLab profile of user |
| created_at | The timestamp when the user was created. |
| bio | The biography of the user. |
| location | The geographic location of the user. |
| public_email | The public email address of the user. |
| skype | The Skype address of the user. |
| linkedin | The LinkedIn account of the user. |
| twitter | The Twitter handle of the user. |
| website_url | The personal website of the user. |
| organization | The organization of the user. |
| ext_id | The external ID of the user. |
| provider | The external provider of the user. |
| theme_id | The ID of the users chosen theme. |
| last_activity_on | The date user was last active. |
| color_scheme_id | The ID of the users chosen color scheme. |
| is_admin | Is the user an Administrator |
| avatar_url | The url of the users avatar. |
| can_create_group | The user has permissions to create groups. |
| can_create_project | The user has permissions to create projects |
| projects_limit | The limit of personal projects the user can create. |
| current_sign_in_at | The timestamp of users current signed in session. |
| last_sign_in_at | The timestamp of users last sign in. |
| confirmed_at | The timestamp of user confirmation. |
| two_factor_enabled | Has the user enabled 2FA/MFA |
| note | The notes against the user. |
| identities | JSON Array of identity information for federated/IdP accounts |
| external | Is the user an external entity |
| private_profile | Is the users profile set to private. |
| shared_runners_minutes_limit | Limit in minutes of time the user can utilise shared runner resources. |
| extra_shared_runners_minutes_limit | Limit in minutes of extra time the user can utilise shared runner resources. |
| using_license_seat | Is the user utilising a seat/slot on the license. |
| custom_attributes | JSON Array of custom attributes held against the user. |

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