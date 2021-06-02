# Table: gitlab_snippet

The `gitlab_snippet` table can be used to query information about snippets owned by the currently authenticated user.

## Examples

### List all your snippets

```sql
select
  *
from
  gitlab_snippet;
```

### Obtain a count of your snippets

```sql
select
  count(*) as snippet_count
from
  gitlab_snippet;
```
