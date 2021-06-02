# Table: gitlab_issue

Issues are used to track bugs, feature requests, tasks, etc on GitLab.

The `gitlab_issue` table can be used to query information against all issues in the GitLab instance.

> Note: It's recommended that you use the `gitlab_my_issue` table for performance.
> 
>It is not advised to use this table if you're using the hosted GitLab.com instance since this will attempt to obtain **ALL** public issues.

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