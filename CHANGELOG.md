## v0.2.0 [WIP]

_What's new?_

- New Table [gitlab_merge_request_change](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_merge_request_change) to see all changes in a merge request.
- Extended the `gitlab_branch` table with the following columns:
  - `commit_message`
- Extended the `gitlab_commit` table with the following columns:
  - `commit_additions`
  - `commit_deletions`
  - `commit_total_changes`
  - `pipeline_id`
  - `pipeline_status`
  - `pipeline_source`
  - `pipeline_ref`
  - `pipeline_sha`
  - `pipeline_url`
  - `pipeline_created`
  - `pipeline_updated`
- Extended the `gitlab_epic` table with the following columns:
  - `parent_id` 
  - `user_notes_count`
  - `author_name`
  - `author_url`
- Removed the following columns from `gitlab_epic` table as no longer in the SDK:
  - `reference`
- Extended the `gitlab_group` table with the following columns:
  - `default_branch_protection` 
  - `file_template_project_id`
  - `shared_runners_enabled`
  - `prevent_forking_outside_group`
  - `commit_count`
  - `storage_size`
  - `repository_size`
  - `wiki_size`
  - `lfs_objects_size`
  - `job_artifacts_size`
  - `pipeline_artifacts_size`
  - `packages_size`
  - `snippets_size`
  - `uploads_size`
- Extended the `gitlab_group_hook` table with the following columns:
  - `confidential_note_events`
  - `enable_ssl_verification` 
- Extended the `gitlab_issue` & `gitlab_my_issue` tables with the following columns:
  - `iid`
  - `author_name`
  - `weight`
  - `issue_type`
  - `subscribed`
  - `user_notes_count`
  - `merge_requests_count`
  - `milestone_id`
  - `milestone_iid`
  - `milestone_title`
  - `milestone_description`
  - `milestone_created_at`
  - `milestone_updated_at`
  - `milestone_start_date`
  - `milestone_due_date`
  - `milestone_state`
  - `milestone_expired`
  - `labels`
  - `short_ref`
  - `rel_ref`
  - `full_ref`
  - `time_estimate`
  - `total_time_spent`
  - `issue_link_id`
  - `epic_issue_id`
- Extended the `gitlab_merge_request` table with the following columns
  - `author_name`
  - `assignee_name`
  - `source_project_id`
  - `target_project_id`
  - `labels`
  - `draft`
  - `merged_by_name`
  - `closed_by_name`
  - `short_ref`
  - `rel_ref`
  - `full_ref`
  - `milestone_id`
  - `milestone_iid`
  - `milestone_title`
  - `milestone_description`
  - `milestone_created_at`
  - `milestone_updated_at`
  - `milestone_start_date`
  - `milestone_due_date`
  - `milestone_state`
  - `milestone_expired`
  - `can_merge`
  - `pipeline_id`
  - `pipeline_project_id`
  - `pipeline_status`
  - `pipeline_source`
  - `pipeline_ref`
  - `pipeline_sha`
  - `pipeline_url`
  - `pipeline_created_at`
  - `pipeline_updated_at`
  - `base_sha`
  - `head_sha`
  - `start_sha`
  - `first_contribution`
  - `blocking_discussions_resolved`
- Extended the `gitlab_project` & `gitlab_my_project` tables with the following columns:
  - `ssh_url`
  - `http_url`
  - `readme_url`
  - `owner_name`
  - `resolve_outdated_diff_discussions`
  - `container_registry_image_prefix`
  - `container_registry_access_level`
  - `container_expiration_policy`
  - `import_status`
  - `import_error`
  - `license_url`
  - `license`
  - `shared_runners_enabled`
  - `runners_token`
  - `public_jobs`
  - `allow_merge_on_skipped_pipeline`
  - `only_allow_merge_if_pipeline_succeeds`
  - `only_allow_merge_if_all_discussions_are_resolved`
  - `remove_source_branch_after_merge`
  - `repository_storage`
  - `merge_method`
  - `fork_parent_id`
  - `fork_parent_name`
  - `fork_parent_path`
  - `fork_parent_url`
  - `mirror`
  - `mirror_user_id`
  - `mirror_trigger_builds`
  - `only_mirror_protected_branches`
  - `mirror_overwrites_diverged_branches`
  - `autoclose_referenced_issues`
  - `ci_forward_deployment_enabled`
  - `ci_config_path`
  - `ci_separated_caches`
- Expanded the `gitlab_project_job` table with the following columns:
  - `queued_duration`
  - `user_name`
  - `pipeline_project_id`
  - `pipeline_ref`
  - `pipeline_sha`
  - `pipeline_status`
  - `artifacts`
  - `runner_id`
  - `runner_name`
  - `runner_description`
  - `runner_active`
  - `runner_is_shared`
  - `commit_id`
  - `commit_short_id`
  - `allow_failure`
  - `failure_reason`
  - `tag`
- Expanded the `gitlab_project_member` table with the following columns:
  - `created_at`
- Expanded the `gitlab_project_pages_domain` table with the following columns:
  - `verified`
  - `verification_code`
  - `enabled_until`
- Expanded the `gitlab_project_pipeline` table with the following columns:
  - `source`
- Expanded the `gitlab_project_pipeline_detail` table with the following columns:
  - `iid`
  - `source`
- Expanded the `gitlab_project_repository_file` table with the following columns:
  - `execute_filemode`
- Expanded the `gitlab_settings` table with the following columns:
  - `abuse_notification_email`
  - `after_sign_up_text`
  - `allow_group_owners_to_manage_ldap`
  - `automatic_purchased_storage_allocation`
  - `can_create_group`
  - `container_registry_cleanup_tags_service_max_list_size`
  - `container_registry_delete_tags_service_timeout`
  - `container_registry_expiration_policies_caching`
  - `container_registry_expiration_policies_worker_capacity`
  - `container_registry_import_created_before`
  - `container_registry_import_max_retries`
  - `container_registry_import_max_step_duration`
  - `container_registry_import_max_tags_count`
  - `custom_http_clone_url_root`
  - `deactivate_dormant_users`
  - `default_ci_config_path`
  - `default_project_deletion_protection`
  - `delayed_group_deletion`
  - `delayed_project_deletion`
  - `delete_inactive_projects`
  - `deletion_adjourned_period`
  - `diff_max_files`
  - `diff_max_lines`
  - `diff_max_patch_bytes`
  - `disable_feed_token`
  - `disable_overriding_approvers_per_merge_request`
  - `domain_allowlist`
  - `domain_denylist`
  - `domain_denylist_enabled`
  - `eks_integration_enabled`
  - `eks_account_id`
  - `eks_access_key_id`
  - `eks_secret_access_key`
  - `email_additional_text`
  - `email_restrictions`
  - `email_restrictions_enabled`
  - `enforce_namespace_storage_limit`
  - `enforce_pat_expiration`
  - `enforce_ssh_key_expiration`
  - `external_pipeline_validation_service_timeout`
  - `external_pipeline_validation_service_token`
  - `external_pipeline_validation_service_url`
  - `floc_enabled`
  - `geo_node_allowed_ips`
  - `geo_status_timeout`
  - `gitpod_enabled`
  - `gitpod_url`
  - `git_rate_limit_users_allowlist`
  - `group_owners_can_manage_default_branch_protection`
  - `group_runner_token_expiration_interval`
  - `inactive_projects_delete_after_months`
  - `inactive_projects_min_size_mb`
  - `inactive_projects_send_warning_email_after_months`
  - `in_product_marketing_emails_enabled`
  - `invisible_captcha_enabled`
  - `issues_create_limit`
  - `keep_latest_artifact`
  - `lock_memberships_to_ldap`
  - `login_recaptcha_protection_enabled`
  - `maintenance_mode`
  - `maintenance_mode_message`
  - `max_export_size`
  - `max_import_size`
  - `max_personal_access_token_lifetime`
  - `max_ssh_key_lifetime`
  - `max_yaml_depth`
  - `max_yaml_size_bytes`
  - `minimum_password_length`
  - `notes_create_limit`
  - `password_number_required`
  - `password_symbol_required`
  - `password_uppercase_required`
  - `password_lowercase_required`
  - `performance_bar_enabled`
  - `personal_access_token_prefix`
  - `prevent_merge_request_author_approval`
  - `prevent_merge_request_committers_approval`
  - `project_runner_token_expiration_interval`
  - `pseudonymizer_enabled`
  - `rate_limiting_response_text`
  - `require_admin_approval_after_user_signup`
  - `runner_token_expiration_interval`
  - `search_rate_limit`
  - `search_rate_limit_unauthenticated`
  - `updating_name_disabled_for_users`
  - `usage_ping_features_enabled`
  - `user_deactivation_emails_enabled`
- Expanded the `gitlab_snippet` table with the following columns:
  - `author_name`
  - `files`
- Expanded the `gitlab_user` table with the following columns:
  - `bot`
  - `job_title`

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v5.0.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v500-2022-11-16)
- Updated: Recompiled with [xanzy/go-gitlab v0.77.0](https://github.com/xanzy/go-gitlab/releases/tag/v0.77.0)

## v0.1.3 [2022-12-12]

_What's new?_

- Extended the `gitlab_project` & `gitlab_my_project` tables with namespace fields as below:
  - `namespace_id` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_name` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_kind` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_path` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)
  - `namespace_full_path` [#32](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/32)

## v0.1.2 [2022-10-13]

_What's new?_

- New table [gitlab_project_job](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_job) - Thanks [@hiepph](https://github.com/hiepph)

## v0.1.1 [2022-10-12]

_Enhancements_

- Updated: Recompiled with [golang version 1.19](https://tip.golang.org/doc/go1.19)
- Updated: Recompiled with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08)

_What's new?_

- New tables added
  - [gitlab_epic](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_epic) *Premium License Required* [#13](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/13)
  - [gitlab_group_iteration](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_group_iteration) *Premium License Required* [#13](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/13)
  - [gitlab_project_iteration](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_iteration) *Premium License Required* [#13](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/13)
  - [gitlab_group_push_rule](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_group_push_rule) *Premium License Required* [#19](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/19)
  - [gitlab_hook](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_hook) *Premium License Required* [#21](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/21)
  - [gitlab_project_protected_branch](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_protected_branch) [#18](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/18)
  - [gitlab_project_pages_domain](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_pages_domain) [#23](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/23)
  - [gitlab_setting](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_setting) [#17](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/17)
- New columns added to [gitlab_project](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project) & [gitlab_my_project](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_my_project) tables 
  - `issues_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `repository_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `merge_requests_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `forking_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `wiki_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `builds_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `snippets_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `pages_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `operations_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `analytics_access_level` [#20](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/20)
  - `topics` [#14](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/14)


## v0.1.0 [2022-05-05]

_Enhancements_

- Updated: Recompiled with [golang version 1.18](https://tip.golang.org/doc/go1.18)
- Updated: Recompiled with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30)

## v0.0.5 [2022-03-25]

_What's new?_

- New tables added
  - [gitlab_project_repository](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_repository)
  - [gitlab_project_repository_file](https://hub.steampipe.io/plugins/theapsgroup/gitlab/tables/gitlab_project_repository_file)

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23)
- Updated: Recompiled with [go-gitlab v0.55.0](https://github.com/xanzy/go-gitlab/releases/tag/v0.55.0)

## v0.0.4 [2021-11-29]

_Enhancements_

- Updated: Recompiled with [golang version 1.17](https://tip.golang.org/doc/go1.17)
- Updated: Recompiled with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22)
- Updated: Recompiled with [go-gitlab v0.52.2](https://github.com/xanzy/go-gitlab/releases/tag/v0.52.2)

## v0.0.3 [2021-09-16]

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v1.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v151--2021-09-13)
- Updated: Recompiled with [go-gitlab v0.50.4](https://github.com/xanzy/go-gitlab/releases/tag/v0.50.4)
- Updated: Added `commit_count`, `storage_size`, `repository_size`, `lfs_objects_size` & `job_artifacts_size` columns to `gitlab_project` & `gitlab_my_project` tables ([#5](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/5))

## v0.0.2 [2021-07-23]

_What's new?_

- Set default API Url to the hosted GitLab to prevent needing to manually define this as per Issue [#3](https://github.com/theapsgroup/steampipe-plugin-gitlab/issues/3)
- Utilising a new feature `optional qualifiers` to allow for faster targeted queries on certain tables `table_issue`, `table_merge_request` & `table_project`
- Enforced specific qualifiers **ONLY** when using the public hosted GitLab instance (to prevent errors / excessively long-running queries)
  - `table_issue` requires at least one of the following `assignee`, `assignee_id`, `author_id` or `project_id`
  - `table_merge_request` requires at least one of the following `assignee_id`, `author_id`, `reviewer_id` or `project_id`
  - `table_project` requires at least one of the following `owner_id` or `owner_username`
  - `table_user` requires at least one of the following `id` or `username`
