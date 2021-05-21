# Table: gitlab_group

Query the Groups (& SubGroups) of GitLab.

## Columns

| Column | Description |
| - | - |
| id | GitLabs internal project ID. |
| name | The group name. |
| path | The group path. |
| description | The groups description. |
| membership_lock | Determines if membership of the group is locked. |
| visibility | The groups visibility (private/internal/public) |
| lfs_enabled | Does the group have Large File System enabled. |
| avatar_url | The url for the groups avatar. |
| web_url | The url for the group. |
| request_access_enabled | Does the group allow access requests. |
| full_name | The full name of the group. |
| full_path | The full path of the group |
| parent_id | The ID of the groups parent group (for sub-groups) |
| custom_attributes | An array of custom attributes. |
| share_with_group_lock | Can this group be shared with other groups |
| require_two_factor_authentication | Does this group require 2fa. |
| two_factor_grace_period | Grace Period (in hours) for 2fa. |
| project_creation_level | The level at which project creation is permitted developer/maintainer/owner |
| auto_devops_enabled | Does the group have auto devops. |
| subgroup_creation_level | The level at which sub-group creation is permitted developer/maintainer/owner |
| emails_disabled | Does this group have email notifications disabled. |
| mentions_disabled | Does this group have mention notifications disabled. |
| runners_token | The groups runner token. |
| ldap_cn | LDAP CN associated with group. |
| ldap_access | LDAP Access associated with group. |
| ldap_group_links | LDAP groups linked to the group. |
| shared_runners_minutes_limit | Limit in minutes of time the group can utilise shared runner resources. |
| extra_shared_runners_minutes_limit | Limit in minutes of extra time the group can utilise shared runner resources. |
| marked_for_deletion_on | Timestamp for when the group was marked to be deleted. |
| created_at | Timestamp for when the group was created. |

## Examples

### Get all Groups

```sql
select
  *
from
  gitlab_group;
```

## Get all top level groups

```sql
select
  *
from
  gitlab_group
where
  parent_id is null;
```

## Get all private groups

```sql
select
  *
from
  gitlab_group
where
  visibility = 'private';
```

### Obtain a count of different visibility levels 

```sql
select
  visibility,
  count(id) as group_count
from
  gitlab_group
group by
  visibility
```