# Table: gitlab_snippet

Query against the snippets of the currently authenticated user only.

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
