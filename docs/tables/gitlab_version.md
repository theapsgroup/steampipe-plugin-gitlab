# Table: gitlab_version

Version information about the GitLab instance.

## Columns

| Column | Description |
| - | - |
| version | The GitLab version number (example: `13.0.14`)|
| revision | The revision hash for the version (example: `ad4adc9d0e1`) |

## Examples

### Get version information for the Gitlab instance

```sql
select
  *
from
  gitlab_version;
```