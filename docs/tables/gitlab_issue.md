# Table: gitlab_issue

Query issues within the GitLab instance.

## Examples

### Obtain all issues

```sql
select
  *
from
  gitlab_issue;
```

### Obtain a list of Confidential Issues

```sql
select
  *
from
  gitlab_issue
where 
  confidential = true;
```

### Obtain counts of issues by state

```sql
select
  state
  count(*) as quantity
from
  gitlab_issue
group by
  state
```