# Table: gitlab_project_repository_file

The `gitlab_project_repository_file` can be used to obtain file information/contents for a single file from within a repository.

However, **you must specify** a `project_id` and a `file_path` for the file in the where or join clauses.

> NOTE: Optionally you may provide a `ref` in the where or join clauses to specify a specific branch, tag or commit - the default value for ref is `main`.

## Examples

### Obtain information about the README.md from a specific project

```sql
select
  *
from
  gitlab_project_repository_file
where
  project_id = 123
and
  file_path = 'README.md';
```

### Obtain information from a file in a sub-folder

```sql
select
  *
from
  gitlab_project_repository_file
where
  project_id = 123
and
  file_path = 'folder/other_folder/file.ext';
```

