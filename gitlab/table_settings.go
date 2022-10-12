package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableSetting() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_setting",
		Description: "GitLab Settings",
		List: &plugin.ListConfig{
			Hydrate: listSettings,
		},
		Columns: settingsColumns(),
	}
}

// Hydrate Functions
func listSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	settings, _, err := conn.Settings.GetSettings()
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, settings)

	return nil, nil
}

func settingsColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the settings.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the settings were created.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the settings were last updated.",
		},
		{
			Name:        "admin_mode",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if admins must re-authenticate to perform administrative functions.",
		},
		{
			Name:        "admin_notification_email",
			Type:        proto.ColumnType_STRING,
			Description: "The email address for administrative notifications.",
		},
		{
			Name:        "after_sign_out_path",
			Type:        proto.ColumnType_STRING,
			Description: "The location users are redirect to upon logging out.",
		},
		{
			Name:        "akismet_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if Akismet spam protection is enabled.",
		},
		{
			Name:        "akismet_api_key",
			Type:        proto.ColumnType_STRING,
			Description: "The API key used for Akisment if enabled.",
		},
		{
			Name:        "allow_local_requests_from_web_hooks_and_services",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if requests can be sent to the local network from web hooks and services.",
		},
		{
			Name:        "allow_local_requests_from_system_hooks",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if requests can be sent to the local network from system hooks.",
		},
		{
			Name:        "archive_builds_in_human_readable",
			Type:        proto.ColumnType_STRING,
			Description: "The human readable representation of when jobs are regarded as expired, if null they never expire.",
		},
		{
			Name:        "asset_proxy_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if asset proxying is enabled.",
		},
		{
			Name:        "asset_proxy_url",
			Type:        proto.ColumnType_STRING,
			Description: "The URL of the asset proxy server.",
			Transform:   transform.FromField("AssetProxyURL"),
		},
		{
			Name:        "asset_proxy_secret_key",
			Type:        proto.ColumnType_STRING,
			Description: "The secret key used to provide access to the asset proxy server.",
		},
		{
			Name:        "asset_proxy_allowlist",
			Type:        proto.ColumnType_JSON,
			Description: "An array of domains which are not proxied.",
		},
		{
			Name:        "authorized_keys_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the authorized_keys file is supported for SSH within GitLab instance.",
		},
		{
			Name:        "auto_devops_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if auto devops is enabled for projects by default.",
			Transform:   transform.FromField("AutoDevOpsEnabled"),
		},
		{
			Name:        "auto_devops_domain",
			Type:        proto.ColumnType_STRING,
			Description: "The domain used by default for all projects Auto Review Apps and Auto Deploy stages.",
			Transform:   transform.FromField("AutoDevOpsDomain"),
		},
		{
			Name:        "container_expiration_policies_enable_historic_entries",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if cleanup policies are enabled for all projects.",
		},
		{
			Name:        "container_registry_token_expire_delay",
			Type:        proto.ColumnType_INT,
			Description: "Container registry token expiration in minutes.",
		},
		{
			Name:        "default_artifacts_expire_in",
			Type:        proto.ColumnType_STRING,
			Description: "The human-readable default expiration time of each jobs artifacts.",
		},
		{
			Name:        "default_branch_name",
			Type:        proto.ColumnType_STRING,
			Description: "Instance level initial branch name.",
		},
		{
			Name:        "default_branch_protection",
			Type:        proto.ColumnType_INT,
			Description: "Level of default branch protection.",
		},
		{
			Name:        "default_group_visibility",
			Type:        proto.ColumnType_STRING,
			Description: "The default visibility of groups, can be private, internal or public.",
		},
		{
			Name:        "default_project_creation",
			Type:        proto.ColumnType_INT,
			Description: "Default project creation protection level (0 no one, 1 maintainers, 2 developers & maintainers).",
		},
		{
			Name:        "default_project_visibility",
			Type:        proto.ColumnType_STRING,
			Description: "The default visibility of projects, can be private, internal or public.",
		},
		{
			Name:        "default_projects_limit",
			Type:        proto.ColumnType_INT,
			Description: "Limit of personal projects each user can create in the instance.",
		},
		{
			Name:        "default_snippet_visibility",
			Type:        proto.ColumnType_STRING,
			Description: "The default visibility of snippets, can be private, internal or public.",
		},
		{
			Name:        "disabled_oauth_sign_in_sources",
			Type:        proto.ColumnType_JSON,
			Description: "Disabled OAuth sign-in sources.",
			Transform:   transform.FromField("DisabledOauthSignInSources"),
		},
		{
			Name:        "dns_rebinding_protection_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if DNS rebinding attack protection is enforced.",
			Transform:   transform.FromField("DNSRebindingProtectionEnabled"),
		},
		{
			Name:        "dsa_key_restriction",
			Type:        proto.ColumnType_INT,
			Description: "The minimum allowed bit length of an uploaded DSA key. Default is 0 (no restriction). -1 disables DSA keys.",
			Transform:   transform.FromField("DSAKeyRestriction"),
		},
		{
			Name:        "ecdsa_key_restriction",
			Type:        proto.ColumnType_INT,
			Description: "The minimum allowed curve size (in bits) of an uploaded ECDSA key. Default is 0 (no restriction). -1 disables ECDSA keys.",
			Transform:   transform.FromField("ECDSAKeyRestriction"),
		},
		{
			Name:        "ed25519_key_restriction",
			Type:        proto.ColumnType_INT,
			Description: "The minimum allowed curve size (in bits) of an uploaded ED25519 key. Default is 0 (no restriction). -1 disables ED25519 keys.",
			Transform:   transform.FromField("Ed25519KeyRestriction"),
		},
		{
			Name:        "email_author_in_body",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the email of the author of issue, MR, comment, etc is included in the email body instead of replacing the email of sender.",
		},
		{
			Name:        "enabled_git_access_protocol",
			Type:        proto.ColumnType_STRING,
			Description: "The enabled protocols for git access, values are ssh, http or nil.",
		},
		{
			Name:        "enforce_terms",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if terms are enforced as application terms of service.",
		},
		{
			Name:        "terms",
			Type:        proto.ColumnType_STRING,
			Description: "The terms of service.",
		},
		{
			Name:        "external_auth_client_cert",
			Type:        proto.ColumnType_STRING,
			Description: "The certificate used to authenticate with external authorization service.",
		},
		{
			Name:        "external_auth_client_key_pass",
			Type:        proto.ColumnType_STRING,
			Description: "The passphrase to use for private key - is encrypted when stored.",
		},
		{
			Name:        "external_auth_client_key",
			Type:        proto.ColumnType_STRING,
			Description: "The private key for the certificate used for authentication - this is encrypted when stored.",
		},
		{
			Name:        "external_authorization_service_default_label",
			Type:        proto.ColumnType_STRING,
			Description: "The default classification label to use when requesting authorization and no classification label has been specified on the project.",
		},
		{
			Name:        "external_authorization_service_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if using an external authorization service for accessing projects.",
		},
		{
			Name:        "external_authorization_service_timeout",
			Type:        proto.ColumnType_DOUBLE,
			Description: "The timeout after which an authorization request is aborted, in seconds.",
		},
		{
			Name:        "external_authorization_service_url",
			Type:        proto.ColumnType_STRING,
			Description: "The URL to which authorization requests are directed.",
			Transform:   transform.FromField("ExternalAuthorizationServiceURL"),
		},
		/* Note: not available in returned object from v0.55.0
		   {
		   	Name:        "external_pipeline_validation_service_timeout",
		   	Type:        proto.ColumnType_INT,
		   	Description: "How long to wait for a response from the pipeline validation service.",
		   },
		   {
		   	Name:        "external_pipeline_validation_service_token",
		   	Type:        proto.ColumnType_STRING,
		   	Description: "The token to include as the X-Gitlab-Token header in requests to the URL in external_pipeline_validation_service_url.",
		   },
		   {
		   	Name:        "external_pipeline_validation_service_url",
		   	Type:        proto.ColumnType_STRING,
		   	Description: "The URL to use for pipeline validation requests.",
		   },
		*/
		{
			Name:        "first_day_of_week",
			Type:        proto.ColumnType_INT,
			Description: "Start day of the week for calendar views and date pickers. Valid values are 0 (default) for Sunday, 1 for Monday, and 6 for Saturday.",
		},
		{
			Name:        "gitaly_timeout_default",
			Type:        proto.ColumnType_INT,
			Description: "Default Gitaly timeout, in seconds.",
		},
		{
			Name:        "gitaly_timeout_medium",
			Type:        proto.ColumnType_INT,
			Description: "Medium Gitaly timeout, in seconds.",
		},
		{
			Name:        "gitaly_timeout_fast",
			Type:        proto.ColumnType_INT,
			Description: "Gitaly fast operation timeout, in seconds.",
		},
		{
			Name:        "grafana_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if grafana is enabled.",
		},
		{
			Name:        "grafana_url",
			Type:        proto.ColumnType_STRING,
			Description: "The URL of the grafana instance.",
			Transform:   transform.FromField("GrafanaURL"),
		},
		{
			Name:        "gravatar_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if gravatar is enabled.",
		},
		{
			Name:        "hashed_storage_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if hashed storage is enabled.",
		},
		{
			Name:        "help_page_hide_commercial_content",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if marketing-related entries are hidden from help.",
		},
		{
			Name:        "help_page_support_url",
			Type:        proto.ColumnType_STRING,
			Description: "Alternate support URL for help page and help dropdown.",
			Transform:   transform.FromField("HelpPageSupportURL"),
		},
		{
			Name:        "help_page_text",
			Type:        proto.ColumnType_STRING,
			Description: "Custom text displayed on the help page.",
		},
		{
			Name:        "hide_third_party_offers",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if third-party offers are hidden in GitLab.",
		},
		{
			Name:        "home_page_url",
			Type:        proto.ColumnType_STRING,
			Description: "The location users are sent to if not logged in.",
			Transform:   transform.FromField("HomePageURL"),
		},
		{
			Name:        "housekeeping_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if housekeeping is enabled.",
		},
		{
			Name:        "housekeeping_bitmaps_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Git pack file bitmap creation is always enabled and cannot be changed via API and UI. This API field is deprecated and always returns true.",
		},
		{
			Name:        "housekeeping_full_repack_period",
			Type:        proto.ColumnType_INT,
			Description: "Number of Git pushes after which a full git repack is run.",
		},
		{
			Name:        "housekeeping_gc_period",
			Type:        proto.ColumnType_INT,
			Description: "Number of Git pushes after which git gc is run.",
		},
		{
			Name:        "housekeeping_incremental_repack_period",
			Type:        proto.ColumnType_INT,
			Description: "Number of Git pushes after which an incremental git repack is run.",
		},
		{
			Name:        "html_emails_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if HTML emails are enabled.",
			Transform:   transform.FromField("HTMLEmailsEnabled"),
		},
		{
			Name:        "import_sources",
			Type:        proto.ColumnType_JSON,
			Description: "An array is strings used to define sources from which projects can be imported.",
		},
		{
			Name:        "local_markdown_version",
			Type:        proto.ColumnType_INT,
			Description: "Increase this value when any cached Markdown should be invalidated.",
		},
		{
			Name:        "max_artifacts_size",
			Type:        proto.ColumnType_INT,
			Description: "Maximum artifacts size in MB.",
		},
		{
			Name:        "max_attachment_size",
			Type:        proto.ColumnType_INT,
			Description: "Limit of attachment size in MB.",
		},
		{
			Name:        "max_pages_size",
			Type:        proto.ColumnType_INT,
			Description: "Maximum size of pages repositories in MB.",
		},
		{
			Name:        "metrics_method_call_threshold",
			Type:        proto.ColumnType_INT,
			Description: "A method call is only tracked when it takes longer than the given amount of milliseconds.",
		},
		{
			Name:        "mirror_available",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if Maintainers can mirror repositories.",
		},
		{
			Name:        "mirror_capacity_threshold",
			Type:        proto.ColumnType_INT,
			Description: "Minimum capacity to be available before scheduling more mirrors preemptively.",
		},
		{
			Name:        "mirror_max_capacity",
			Type:        proto.ColumnType_INT,
			Description: "Maximum number of mirrors that can be synchronizing at the same time.",
		},
		{
			Name:        "mirror_max_delay",
			Type:        proto.ColumnType_INT,
			Description: "Maximum time (in minutes) between updates that a mirror can have when scheduled to synchronize.",
		},
		{
			Name:        "outbound_local_requests_whitelist",
			Type:        proto.ColumnType_JSON,
			Description: "An array of trusted domains or IP addresses to which local requests are allowed when local requests for hooks and services are disabled.",
		},
		{
			Name:        "pages_domain_verification_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if users need to prove ownership of custom domain.",
		},
		{
			Name:        "password_authentication_enabled_for_git",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if authentication for Git over HTTP(S) via a GitLab account password is enabled.",
		},
		{
			Name:        "password_authentication_enabled_for_web",
			Type:        proto.ColumnType_BOOL,
			Description: "indicates if authentication for the web interface via a GitLab account password is enabled.",
		},
		{
			Name:        "performance_bar_allowed_group_path",
			Type:        proto.ColumnType_STRING,
			Description: "Path of the group that is allowed to toggle the performance bar.",
		},
		{
			Name:        "plantuml_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if Plant UML integration is enabled.",
		},
		{
			Name:        "plantuml_url",
			Type:        proto.ColumnType_STRING,
			Description: "The PlantUML instance URL for integration.",
			Transform:   transform.FromField("PlantumlURL"),
		},
		{
			Name:        "polling_interval_multiplier",
			Type:        proto.ColumnType_DOUBLE,
			Description: "Interval multiplier used by endpoints that perform polling.",
		},
		{
			Name:        "project_export_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if project exporting is enabled.",
		},
		{
			Name:        "prometheus_metrics_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if prometheus metrics are enabled.",
		},
		{
			Name:        "protected_ci_variables",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if CI/CD variables are protected by default.",
			Transform:   transform.FromField("ProtectedCIVariables"),
		},
		{
			Name:        "push_event_hooks_limit",
			Type:        proto.ColumnType_INT,
			Description: "Number of changes (branches or tags) in a single push to determine whether webhooks and services fire or not.",
		},
		{
			Name:        "push_event_activities_limit",
			Type:        proto.ColumnType_INT,
			Description: "Number of changes (branches or tags) in a single push to determine whether individual push events or bulk push events are created.",
		},
		{
			Name:        "recaptcha_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if reCAPTCHA is enabled.",
		},
		{
			Name:        "recaptcha_private_key",
			Type:        proto.ColumnType_STRING,
			Description: "Private key for reCAPTCHA.",
		},
		{
			Name:        "recaptcha_site_key",
			Type:        proto.ColumnType_STRING,
			Description: "Site key for reCAPTCHA.",
		},
		{
			Name:        "receive_max_input_size",
			Type:        proto.ColumnType_INT,
			Description: "Maximum push size (MB).",
		},
		{
			Name:        "repository_checks_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if GitLab periodically runs git fsck in all project and wiki repositories to look for silent disk corruption issues.",
		},
		{
			Name:        "repository_size_limit",
			Type:        proto.ColumnType_INT,
			Description: "Size limit per repository (MB).",
		},
		{
			Name:        "repository_storages",
			Type:        proto.ColumnType_JSON,
			Description: "An array of names of enabled storage paths, taken from gitlab.yml.",
		},
		{
			Name:        "require_two_factor_authentication",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if 2FA is required for all users.",
		},
		{
			Name:        "restricted_visibility_levels",
			Type:        proto.ColumnType_JSON,
			Description: "An array of levels that cannot be used by non-Administrator users for groups, projects or snippets.",
		},
		{
			Name:        "rsa_key_restriction",
			Type:        proto.ColumnType_INT,
			Description: "The minimum allowed bit length of an uploaded RSA key. Default is 0 (no restriction). -1 disables RSA keys.",
		},
		{
			Name:        "send_user_confirmation_email",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if users receive a confirmation email on sign-up.",
		},
		{
			Name:        "session_expire_delay",
			Type:        proto.ColumnType_INT,
			Description: "Session duration in minutes.",
		},
		{
			Name:        "shared_runners_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if shared runners are enabled for new projects.",
		},
		{
			Name:        "shared_runners_minutes",
			Type:        proto.ColumnType_INT,
			Description: "The maximum number of CI/CD minutes that a group can use on shared runners per month.",
		},
		{
			Name:        "shared_runners_text",
			Type:        proto.ColumnType_STRING,
			Description: "Shared runners text.",
		},
		{
			Name:        "sign_in_text",
			Type:        proto.ColumnType_STRING,
			Description: "Text on the login page.",
		},
		{
			Name:        "signup_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if user registration is allowed.",
		},
		{
			Name:        "snowplow_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if snowplow tracking is enabled.",
		},
		{
			Name:        "snowplow_collector_hostname",
			Type:        proto.ColumnType_STRING,
			Description: "The Snowplow collector hostname.",
		},
		{
			Name:        "snowplow_cookie_domain",
			Type:        proto.ColumnType_STRING,
			Description: "The Snowplow cookie domain.",
		},
		{
			Name:        "snowplow_site_id",
			Type:        proto.ColumnType_STRING,
			Description: "The Snowplow site name / application ID.",
		},
		{
			Name:        "terminal_max_session_time",
			Type:        proto.ColumnType_INT,
			Description: "Maximum time for web terminal websocket connection (in seconds). 0 for unlimited time.",
		},
		{
			Name:        "throttle_authenticated_api_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if authenticated API request rate limit is enabled.",
			Transform:   transform.FromField("ThrottleAuthenticatedAPIEnabled"),
		},
		{
			Name:        "throttle_authenticated_api_period_in_seconds",
			Type:        proto.ColumnType_INT,
			Description: "Rate limit period (in seconds).",
			Transform:   transform.FromField("ThrottleAuthenticatedAPIPeriodInSeconds"),
		},
		{
			Name:        "throttle_authenticated_api_requests_per_period",
			Type:        proto.ColumnType_INT,
			Description: "Maximum requests per period per user.",
			Transform:   transform.FromField("ThrottleAuthenticatedAPIRequestsPerPeriod"),
		},
		{
			Name:        "throttle_authenticated_web_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if authenticated web request rate limit is enabled.",
			Transform:   transform.FromField("ThrottleAuthenticatedWebEnabled"),
		},
		{
			Name:        "throttle_authenticated_web_period_in_seconds",
			Type:        proto.ColumnType_INT,
			Description: "Rate limit period (in seconds).",
			Transform:   transform.FromField("ThrottleAuthenticatedWebPeriodInSeconds"),
		},
		{
			Name:        "throttle_authenticated_web_requests_per_period",
			Type:        proto.ColumnType_INT,
			Description: "Maximum requests per period per user.",
			Transform:   transform.FromField("ThrottleAuthenticatedWebRequestsPerPeriod"),
		},
		{
			Name:        "throttle_unauthenticated_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if unauthenticated web request rate limit is enabled.",
			Transform:   transform.FromField("ThrottleUnauthenticatedEnabled"),
		},
		{
			Name:        "throttle_unauthenticated_period_in_seconds",
			Type:        proto.ColumnType_INT,
			Description: "Rate limit period (in seconds).",
			Transform:   transform.FromField("ThrottleUnauthenticatedPeriodInSeconds"),
		},
		{
			Name:        "throttle_unauthenticated_requests_per_period",
			Type:        proto.ColumnType_INT,
			Description: "Maximum requests per period per IP.",
			Transform:   transform.FromField("ThrottleUnauthenticatedRequestsPerPeriod"),
		},
		{
			Name:        "time_tracking_limit_to_hours",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if time tracking units is limited to hours only.",
		},
		{
			Name:        "two_factor_grace_period",
			Type:        proto.ColumnType_INT,
			Description: "Amount of time (in hours) that users are allowed to skip forced configuration of two-factor authentication.",
		},
		{
			Name:        "unique_ips_limit_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if users are limited to sign in from different IPs.",
			Transform:   transform.FromField("UniqueIPsLimitEnabled"),
		},
		{
			Name:        "unique_ips_limit_per_user",
			Type:        proto.ColumnType_INT,
			Description: "Maximum number of IPs per user.",
			Transform:   transform.FromField("UniqueIPsLimitPerUser"),
		},
		{
			Name:        "unique_ips_limit_time_window",
			Type:        proto.ColumnType_INT,
			Description: "How many seconds an IP is counted towards the limit.",
			Transform:   transform.FromField("UniqueIPsLimitTimeWindow"),
		},
		{
			Name:        "usage_ping_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if every week GitLab reports license usage back to GitLab, Inc.",
		},
		{
			Name:        "user_default_external",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if newly registered users are external by default.",
		},
		{
			Name:        "user_default_internal_regex",
			Type:        proto.ColumnType_STRING,
			Description: "Specify an email address regex pattern to identify default internal users.",
		},
		{
			Name:        "user_oauth_applications",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if users can register any application to use GitLab as an OAuth provider.",
		},
		{
			Name:        "user_show_add_ssh_key_message",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the warning is shown to users when they haven't configured an SSH key.",
			Transform:   transform.FromField("UserShowAddSSHKeyMessage"),
		},
		{
			Name:        "version_check_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if GitLab checks for new versions and informs about available updates.",
		},
		{
			Name:        "web_ide_clientside_preview_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if Live Preview (allow live previews of JavaScript projects in the Web IDE using CodeSandbox Live Preview) is enabled.",
			Transform:   transform.FromField("WebIDEClientsidePreviewEnabled"),
		},
	}
}
