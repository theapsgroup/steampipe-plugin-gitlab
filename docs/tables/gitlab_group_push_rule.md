# Table: gitlab_group_push_rule

The `gitlab_group_push_rule` table can be used to query information about the rules associated with pushing to projects/repos in a specific group.

However, **you must specify** a `group_id` in the where or join clause.

## Examples

### Obtain the push rules for a specific group

```sql
select
  id,
  created_at,
  commit_message_regex,
  commit_message_negative_regex,
  branch_name_regex,
  deny_delete_tag,
  member_check,
  prevent_secrets,
  author_email_regex,
  file_name_regex,
  max_file_size,
  commit_committer_check,
  reject_unsigned_commits
from
  gitlab_group_push_rule
where
  group_id = 14597683;
```
