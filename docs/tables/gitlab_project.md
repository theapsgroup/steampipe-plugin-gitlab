# Table: gitlab_project

Query projects within GitLab instance.

## Columns

| Column | Description |
| - | - |
| id | GitLabs internal project ID. |
| name | The projects name. |
| path | The projects path. |
| description | The projects description. |
| default_branch | The projects default branch name. |
| full_name | The projects name including namespace. |
| full_path | The projects path including namespace. |
| public | Is the project public |
| visibility | The projects visibility level (private/public/internal) |
| web_url | The projects url. |
| tag_list | An array of tags associated to the project. |
| issues_enabled | Indicates if project has issues enabled. |
| open_issues_count | A count of open issues on the project. |
| merge_requests_enabled | Indicates if merge requests are enabled on the project |
| approvals_before_merge | The project setting for number of approvals required before a merge request can be merged. |
| jobs_enabled | Indicates if the project has jobs enabled. |
| wiki_enabled | Indicates if the project has the wiki enabled. |
| snippets_enabled | Indicates if the project has snippets enabled. |
| container_registry_enabled | Indicates if the project has the container registry enabled. |
| creator_id | User ID of the projects creator. |
| created_at | Timestamp of when project was created. |
| last_activity_at | Timestamp of when last activity happened on the project. |
| marked_for_deletion_at | Timestamp of when project was marked for deletion. |
| empty_repo | Indicates if the repository of the project is empty. |
| archived | Indicates if the project is archived. |
| avatar_url | The url for the projects avatar. |
| forks_count | The number of forks of the project. |
| star_count | The number of stars given to the project. |
| lfs_enabled | Indicates if the project has large file system enabled. |
| request_access_enabled | Indicates if the project has request access enabled. |
| packages_enabled | Indicates if the project has packages enabled. |
| owner_id | The projects owner ID. (null if owned by a group) |
| owner_username | The projects owner username. (null if owned by a group) |

## Examples

### Get all projects
```sql
select
  *
from
  gitlab_project;
```

### Get top 10 projects based on stars
```sql
select
  *
from
  gitlab_project
order by
  star_count desc
limit 1;  
```

### Get project creation information
```sql
select
  u.username as creator,
  p.full_path as project,
  p.created_at as created
from
  gitlab_project p
inner join
  gitlab_user u
on 
  p.creator_id = u.id;
```