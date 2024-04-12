## v0.56.0 [2024-04-12]

_What's new?_

- New tables added
  - [azure_backup_policy](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_backup_policy) ([#739](https://github.com/turbot/steampipe-plugin-azure/pull/739))

_Bug fixes_

- Fixed the plugin's Postgres FDW Extension crash [issue](https://github.com/turbot/steampipe-postgres-fdw/issues/434).

## v0.55.0 [2024-03-22]

_What's new?_

- New tables added
  - [azure_compute_virtual_machine_metric_available_memory](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_available_memory) ([#729](https://github.com/turbot/steampipe-plugin-azure/pull/729))
  - [azure_compute_virtual_machine_metric_available_memory_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_available_memory_daily) ([#729](https://github.com/turbot/steampipe-plugin-azure/pull/729))
  - [azure_compute_virtual_machine_metric_available_memory_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_available_memory_hourly) ([#729](https://github.com/turbot/steampipe-plugin-azure/pull/729))

_Enhancements_

- Added the `access_keys` column to `azure_storage_account` table. ([#730](https://github.com/turbot/steampipe-plugin-azure/pull/730))
- Added the `disable_local_auth` column to `azure_cosmosdb_account` table. ([#736](https://github.com/turbot/steampipe-plugin-azure/pull/736))

## v0.54.0 [2024-02-02]

_What's new?_

- New tables added
  - [azure_api_management_backend](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_api_management_backend) ([#689](https://github.com/turbot/steampipe-plugin-azure/pull/689))
  - [azure_consumption_usage](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_consumption_usage) ([#721](https://github.com/turbot/steampipe-plugin-azure/pull/721))
  - [azure_monitor_log_profile](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_monitor_log_profile) ([#717](https://github.com/turbot/steampipe-plugin-azure/pull/717))

_Enhancements_

- Added the `authorization_rules` column to `azure_servicebus_namespace` table. ([#719](https://github.com/turbot/steampipe-plugin-azure/pull/719))

## v0.53.0 [2024-01-22]

_Enhancements_

- Added the `audit_policy` column to `azure_sql_database` and `azure_sql_server` tables. ([#711](https://github.com/turbot/steampipe-plugin-azure/pull/711))
- Added the `webhooks` column to `azure_container_registry` table. ([#710](https://github.com/turbot/steampipe-plugin-azure/pull/710))
- Added the `disable_local_auth` and `status` columns to `azure_servicebus_namespace` table. ([#715](https://github.com/turbot/steampipe-plugin-azure/pull/715))

_Bug fixes_

- Fixed the `azure_key_vault_secret` table to correctly return data when keyvault name is in camel-case. ([#638](https://github.com/turbot/steampipe-plugin-azure/pull/638))

## v0.52.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/install/steampipe.sh), as a [Postgres FDW](https://steampipe.io/install/postgres.sh), as a [SQLite extension](https://steampipe.io/install/sqlite.sh) and as a standalone [exporter](https://steampipe.io/install/export.sh).
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension.
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-azure/blob/main/docs/LICENSE).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server enacapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#699](https://github.com/turbot/steampipe-plugin-azure/pull/699))

## v0.51.0 [2023-11-02]

_What's new?_

- New tables added
  - [azure_alert_management](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_alert_management) ([#685](https://github.com/turbot/steampipe-plugin-azure/pull/685))
  - [azure_databricks_workspace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_databricks_workspace) ([#692](https://github.com/turbot/steampipe-plugin-azure/pull/692))
  - [azure_monitor_activity_log_event](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_monitor_activity_log_event) ([#684](https://github.com/turbot/steampipe-plugin-azure/pull/684))
  - [azure_recovery_services_backup_job](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_recovery_services_backup_job) ([#681](https://github.com/turbot/steampipe-plugin-azure/pull/681))

## v0.50.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#679](https://github.com/turbot/steampipe-plugin-azure/pull/679))

## v0.50.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#676](https://github.com/turbot/steampipe-plugin-azure/pull/676))
- Recompiled plugin with Go version `1.21`. ([#676](https://github.com/turbot/steampipe-plugin-azure/pull/676))

## v0.49.0 [2023-08-31]

_What's new?_

- New tables added
  - [azure_postgresql_flexible_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_postgresql_flexible_server) ([#659](https://github.com/turbot/steampipe-plugin-azure/pull/659))

## v0.48.0 [2023-08-24]

_Enhancements_

- Added the `server_security_alert_policy` column to `azure_mysql_server` table. ([#656](https://github.com/turbot/steampipe-plugin-azure/pull/656))
- Added the `storage_info_value` column to `azure_app_service_web_app` table. ([#657](https://github.com/turbot/steampipe-plugin-azure/pull/657))

_Bug fixes_

- Fixed the `disable_local_auth` column in `azure_eventgrid_domain` table to correctly return data instead of `null`. ([#658](https://github.com/turbot/steampipe-plugin-azure/pull/658))

## v0.47.0 [2023-08-17]

_Enhancements_

- Added the `server_security_alert_policy` column to `azure_postgresql_server` table. ([#651](https://github.com/turbot/steampipe-plugin-azure/pull/651))
- Added the `identity` column to `azure_kusto_cluster` table. ([#652](https://github.com/turbot/steampipe-plugin-azure/pull/652))
- Added the `site_config_resource` column to `azure_app_service_web_app_slot` table. ([#653](https://github.com/turbot/steampipe-plugin-azure/pull/653))

_Bug fixes_

- Fixed the `GetConfig` of `azure_app_service_web_app_slot` table to correctly return data instead of an empty row. ([#654](https://github.com/turbot/steampipe-plugin-azure/pull/654))

## v0.46.0 [2023-07-19]

_What's new?_

- New tables added
  - [azure_container_group](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_container_group) ([#634](https://github.com/turbot/steampipe-plugin-azure/pull/634))

## v0.45.1 [2023-07-14]

_Bug fixes_

- Fixed the `private_endpoint_connections` column of the `azure_mariadb_server` table to correctly return data instead of null. ([#631](https://github.com/turbot/steampipe-plugin-azure/pull/631))

## v0.45.0 [2023-07-06]

_What's new?_

- New tables added
  - [azure_kubernetes_service_version](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_kubernetes_service_version) ([#623](https://github.com/turbot/steampipe-plugin-azure/pull/623))

## v0.44.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#619](https://github.com/turbot/steampipe-plugin-azure/pull/619))

## v0.43.0 [2023-05-30]

_Enhancements_

- Added columns `table_logging_delete`, `table_logging_read`, `table_logging_retention_policy`, `table_logging_version`, `table_logging_write`, and `table_properties` to `azure_storage_account` table. ([#614](https://github.com/turbot/steampipe-plugin-azure/pull/614))

## v0.42.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#604](https://github.com/turbot/steampipe-plugin-azure/pull/604))

## v0.41.1 [2023-05-04]

_Bug fixes_

- Fixed the `intrusion_detection_mode` column of `azure_firewall_policy table` to correctly return data instead of `null`. ([#609](https://github.com/turbot/steampipe-plugin-azure/pull/609))

## v0.41.0 [2023-04-21]

_Enhancements_

- Added columns `status_of_primary` and `status_of_secondary` to `azure_storage_account` table. ([#605](https://github.com/turbot/steampipe-plugin-azure/pull/605)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v0.40.1 [2023-04-05]

_Bug fixes_

- Fixed the `throughput_settings` column of `azure_cosmosdb_mongo_collection` and `azure_cosmosdb_mongo_database` tables to correctly return data instead of an error, when default throughput setting is used for databases or collections. ([#602](https://github.com/turbot/steampipe-plugin-azure/pull/602))

## v0.40.0 [2023-03-31]

_What's new?_

- New tables added
  - [azure_app_service_web_app_slot](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_app_service_web_app_slot) ([#592](https://github.com/turbot/steampipe-plugin-azure/pull/592))
  - [azure_cosmosdb_restorable_database_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_cosmosdb_restorable_database_account) ([#596](https://github.com/turbot/steampipe-plugin-azure/pull/596))
  - [azure_firewall_policy](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_firewall_policy) ([#598](https://github.com/turbot/steampipe-plugin-azure/pull/598))

_Enhancements_

- Added column `restore_parameters` to `azure_cosmosdb_account` table. ([#594](https://github.com/turbot/steampipe-plugin-azure/pull/594))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#599](https://github.com/turbot/steampipe-plugin-azure/pull/599))

## v0.39.0 [2023-03-24]

_What's new?_

- New tables added
  - [azure_cosmosdb_mongo_collection](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_cosmosdb_mongo_collection) ([#589](https://github.com/turbot/steampipe-plugin-azure/pull/589))
  - [azure_private_dns_zone](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_private_dns_zone) ([#583](https://github.com/turbot/steampipe-plugin-azure/pull/583)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Enhancements_

- Added column `throughput_settings` to `azure_cosmosdb_mongo_database` table. ([#587](https://github.com/turbot/steampipe-plugin-azure/pull/587))
- Added column `backup_policy` to `azure_cosmosdb_account` table. ([#585](https://github.com/turbot/steampipe-plugin-azure/pull/585))

## v0.38.0 [2023-03-15]

_What's new?_

- New tables added
  - [azure_application_insight](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_application_insight) ([#579](https://github.com/turbot/steampipe-plugin-azure/pull/579))
  - [azure_bastion_host](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_bastion_host) ([#580](https://github.com/turbot/steampipe-plugin-azure/pull/580))
  - [azure_compute_ssh_key](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_ssh_key) ([#560](https://github.com/turbot/steampipe-plugin-azure/pull/560)) (Thanks [@srgg](https://github.com/srgg) for the contribution!)
  - [azure_dns_zone](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_dns_zone) ([#575](https://github.com/turbot/steampipe-plugin-azure/pull/575)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Bug fixes_

- Fixed column name typo in `azure_compute_virtual_machine_scale_set_network_interface` table. ([#573](https://github.com/turbot/steampipe-plugin-azure/pull/573)) (Thanks [@jackdelab](https://github.com/jackdelab) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v520-2023-03-02) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#574](https://github.com/turbot/steampipe-plugin-azure/pull/574))

## v0.37.0 [2023-02-22]

_What's new?_

- New tables added
  - [azure_automation_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_automation_account) ([#569](https://github.com/turbot/steampipe-plugin-azure/pull/569))
  - [azure_automation_variable](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_automation_variable) ([#571](https://github.com/turbot/steampipe-plugin-azure/pull/571))

## v0.36.0 [2023-02-10]

_Enhancements_

- Added column `vnet_rules` to `azure_mysql_server` table. ([#558](https://github.com/turbot/steampipe-plugin-azure/pull/558))

_Bug fixes_

- Fixed the `ip_configurations` column in `azure_firewall` table to return `null` instead of a panic error when no IP configuration is set on the Azure firewall. ([#561](https://github.com/turbot/steampipe-plugin-azure/pull/561)) (Thanks [@mdaguete](https://github.com/mdaguete) for the contribution!!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.12](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4112-2023-02-09) which fixes the query caching functionality. ([#565](https://github.com/turbot/steampipe-plugin-azure/pull/565))

## v0.35.1 [2023-01-10]

_Bug fixes_

- Fixed the `vulnerability_assessment_scan_records` column of the `azure_sql_database` table to return `nil` instead of an error when vulnerability assessment settings are unavailable for a SQL database. ([#552](https://github.com/turbot/steampipe-plugin-azure/pull/552))
- Fixed the `ip_configurations` column of the `azure_subnet` table to return `nil` instead of an error when IP configuration details are unavailable for a subnet. ([#556](https://github.com/turbot/steampipe-plugin-azure/pull/556))

## v0.35.0 [2022-11-25]

_What's new?_

- New tables added
  - [azure_compute_virtual_machine_scale_set_network_interface](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_scale_set_network_interface) ([#537](https://github.com/turbot/steampipe-plugin-azure/pull/537))
  - [azure_key_vault_key_version](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_key_version) ([#549](https://github.com/turbot/steampipe-plugin-azure/pull/549))

## v0.34.2 [2022-11-11]

_Bug fixes_

- Fixed the authentication flow to correctly refresh expired Azure CLI token credentials. ([#544](https://github.com/turbot/steampipe-plugin-azure/pull/544))

## v0.34.1 [2022-11-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-aws/pull/1384) which increases the default open file limit. ([#543](https://github.com/turbot/steampipe-plugin-azure/pull/543))

## v0.34.0 [2022-11-07]

_What's new?_

- New tables added
  - [azure_nat_gateway](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_nat_gateway) ([#533](https://github.com/turbot/steampipe-plugin-azure/pull/533))

_Enhancements_

- Updated the `docs/index.md` file to include multi-subscription configuration examples. ([#538](https://github.com/turbot/steampipe-plugin-azure/pull/538))
- Added column `apps` to the `azure_app_service_plan` table. ([#536](https://github.com/turbot/steampipe-plugin-azure/pull/536))

_Bug fixes_

- Fixed the column `virtual_network_rules` of the `azure_cosmosdb_account` table to correctly return data instead of `null`. ([#532](https://github.com/turbot/steampipe-plugin-azure/pull/532))
- Fixed the `azure_app_configuration table` table to correctly return results instead of returning an error. ([#531](https://github.com/turbot/steampipe-plugin-azure/pull/531))
- Fixed `capabilities`, `costs`, `location_info` and `restrictions` columns in `azure_compute_resource_sku` table to correctly return data instead of an empty array. ([#528](https://github.com/turbot/steampipe-plugin-azure/pull/528))
- Fixed invalid references to `GCP Monitor` in various `azure_compute_*` table documents. ([#520](https://github.com/turbot/steampipe-plugin-azure/pull/520)) (Thanks [@JoshRosen](https://github.com/JoshRosen) for the fix!)

## v0.33.0 [2022-09-29]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#514](https://github.com/turbot/steampipe-plugin-azure/pull/514))
- Recompiled plugin with Go version `1.19`. ([#514](https://github.com/turbot/steampipe-plugin-azure/pull/514))

## v0.32.0 [2022-09-15]

_Enhancements_

- Added column `diagnostic_logs_configuration` to `azure_app_service_web_app` table. ([#517](https://github.com/turbot/steampipe-plugin-azure/pull/517))

_Bug fixes_

- Fixed `location` -> `region` column names in `azure_virtual_network` table doc examples. ([#511](https://github.com/turbot/steampipe-plugin-azure/pull/511))

_Deprecated_

- Updated `azure_ad_group`, `azure_ad_service_principal`, and `azure_ad_user` tables (deprecated in v0.20.0) to no longer return results and instead return an error and suggest the respective replacement table. These tables will be removed entirely from this plugin in a future version.

## v0.31.0 [2022-07-22]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#508](https://github.com/turbot/steampipe-plugin-azure/pull/508))

## v0.30.0 [2022-07-01]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v331--2022-06-30). ([#500](https://github.com/turbot/steampipe-plugin-azure/pull/500))

## v0.29.0 [2022-06-27]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v330--2022-6-22). ([#498](https://github.com/turbot/steampipe-plugin-azure/pull/498))

## v0.28.0 [2022-06-09]

_What's new?_

- Added `ignore_error_codes` config arg to provide users the ability to set a list of additional Azure error codes to ignore while running queries. For instance, to ignore some common access denied errors, which is helpful when running with limited permissions, set the argument `ignore_error_codes = ["UnauthorizedOperation", "InsufficientAccountPermissions"]`. For more information, please see [Azure plugin configuration](https://hub.steampipe.io/plugins/turbot/azure#configuration). ([#495](https://github.com/turbot/steampipe-plugin-azure/pull/495))

_Bug fixes_

- Fixed `azure_mssql_elasticpool`, `azure_mysql_flexible_server`, `azure_policy_assignment`, `azure_public_ip` and `azure_storage_share_file` tables to correctly return data instead of an empty row. ([#495](https://github.com/turbot/steampipe-plugin-azure/pull/495))

## v0.27.2 [2022-06-01]

_Bug fixes_

- Fixed the `access_control` column in `azure_logic_app_workflow` table to consistently return `null` instead of intermittently returning `{}` when no data is available. ([#486](https://github.com/turbot/steampipe-plugin-azure/pull/486))

## v0.27.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#489](https://github.com/turbot/steampipe-plugin-azure/pull/489))

## v0.27.0 [2022-05-05]

_What's new?_

- New tables added
  - [azure_security_center_automation](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_automation) ([#481](https://github.com/turbot/steampipe-plugin-azure/pull/481))
  - [azure_security_center_sub_assessment](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_sub_assessment) ([#482](https://github.com/turbot/steampipe-plugin-azure/pull/482))

_Enhancements_

- Added `ip_configurations` column to `azure_subnet` table. ([#483](https://github.com/turbot/steampipe-plugin-azure/pull/483))

## v0.26.0 [2022-04-27]

_What's new?_

- New tables added
  - [azure_mysql_flexible_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mysql_flexible_server) ([#476](https://github.com/turbot/steampipe-plugin-azure/pull/476))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#471](https://github.com/turbot/steampipe-plugin-azure/pull/471))
- Added support for native Linux ARM and Mac M1 builds. ([#479](https://github.com/turbot/steampipe-plugin-azure/pull/479))

_Bug fixes_

- `azure_storage_share_file` table has been updated to handle `FeatureNotSupportedForAccount` error when the storage type is `BlockBlobStorage` ([#478](https://github.com/turbot/steampipe-plugin-azure/pull/478))

## v0.25.0 [2022-04-05]

_What's new?_

- New tables added
  - [azure_management_group](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_management_group) ([#460](https://github.com/turbot/steampipe-plugin-azure/pull/460))
  - [azure_storage_share_file](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_share_file) ([#464](https://github.com/turbot/steampipe-plugin-azure/pull/464))

_Enhancements_

- Fixed the `network_access_policy` column of `azure_compute_disk` table to correctly return the network access policy instead of null ([#447](https://github.com/turbot/steampipe-plugin-azure/pull/447))
- Updated the data type of `create_mode` and `read_scale` columns to `string` in `azure_sql_database` table ([#459](https://github.com/turbot/steampipe-plugin-azure/pull/459))

_Bug fixes_

- `azure_storage_table` and `azure_storage_queue` tables have been updated to handle the `FeatureNotSupportedForAccount` error when the storage type is `BlockBlobStorage` ([#467](https://github.com/turbot/steampipe-plugin-azure/pull/467)) ([#465](https://github.com/turbot/steampipe-plugin-azure/pull/465))

## v0.24.0 [2022-03-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#449](https://github.com/turbot/steampipe-plugin-azure/pull/449))

## v0.23.2 [2022-01-19]

_Bug fixes_

- Fixed: Authentication session credentials are now cached correctly ([#442](https://github.com/turbot/steampipe-plugin-azure/pull/442))

## v0.23.1 [2022-01-07]

_Bug fixes_

- Renamed column `environment_name` to `cloud_environment` across all the tables to remain consistent with Azure documentation ([#438](https://github.com/turbot/steampipe-plugin-azure/pull/438))

## v0.23.0 [2022-01-05]

_What's new?_

- New tables added
  - [azure_compute_virtual_machine_scale_set_vm](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_scale_set_vm) ([#425](https://github.com/turbot/steampipe-plugin-azure/pull/425))

_Enhancements_

- Added column `environment_name` across all the tables ([#424](https://github.com/turbot/steampipe-plugin-azure/pull/424))
- Added column `diagnostic_settings` to `azure_storage_account` table ([#432](https://github.com/turbot/steampipe-plugin-azure/pull/432))
- Added column `server_configurations` to `azure_mysql_server` table ([#429](https://github.com/turbot/steampipe-plugin-azure/pull/429))
- `azure_storage_account`, `azure_storage_blob_service`, `azure_storage_container`, `azure_storage_queue`, `azure_storage_table` and `azure_storage_table_service` tables have been updated to handle the `FeatureNotSupportedForAccount` error when the storage type is `FileStorage` ([#418](https://github.com/turbot/steampipe-plugin-azure/pull/418))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#431](https://github.com/turbot/steampipe-plugin-azure/pull/431))

_Bug fixes_

- `immutability_policy` column in `azure_storage_container` table will now correctly return the data instead of returning null when immutability policy is set to a storage container ([#434](https://github.com/turbot/steampipe-plugin-azure/pull/434))

## v0.22.0 [2021-12-15]

_Enhancements_

- Added `os_name` and `os_version` columns to the `azure_compute_virtual_machine` table ([#420](https://github.com/turbot/steampipe-plugin-azure/pull/420))

## v0.21.1 [2021-12-02]

_Enhancements_

- Updated the descriptions for `azure_ad_group`, `azure_ad_service_principal` and `azure_ad_user` tables to reflect the deprecation status ([#414](https://github.com/turbot/steampipe-plugin-azure/pull/414))

_Bug fixes_

- Fixed `azure_ad_group`, `azure_ad_service_principal` and `azure_ad_user` tables to use `GraphEndpoint` instead of `ResourceManagerEndpoint` to make the API calls ([#413](https://github.com/turbot/steampipe-plugin-azure/pull/413))

## v0.21.0 [2021-11-23]

_Enhancements_

- Recompiled plugin Go version 1.17 ([#408](https://github.com/turbot/steampipe-plugin-azure/pull/408))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#406](https://github.com/turbot/steampipe-plugin-azure/pull/406))

## v0.20.0 [2021-10-26]

_Enhancements_

- Updated: Add context cancellation handling to all the tables ([#343](https://github.com/turbot/steampipe-plugin-azure/pull/343))
- Recompiled plugin with [steampipe-plugin-sdk v1.7.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v170--2021-10-18) ([#400](https://github.com/turbot/steampipe-plugin-azure/pull/400))
- The configuration section in the docs/index.md file now includes additional information on different methods of setting up credentials in the `azure.spc` file ([268](https://github.com/turbot/steampipe-plugin-azure/pull/268))


_Bug fixes_

- Fixed: Authentication now works properly for non-public cloud environments ([268](https://github.com/turbot/steampipe-plugin-azure/pull/268))

_Deprecated_

- The following tables have been deprecated since they are now maintained in the [azuread plugin](https://hub.steampipe.io/plugins/turbot/azuread/tables). These tables will be removed in the next major version. We recommend updating any scripts or workflows that use these tables to use the equivalent tables in the Azure AD plugin instead.
  - [azure_ad_group](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_ad_group) (replaced by [azuread_group](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_group))
  - [azure_ad_service_principal](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_ad_service_principal) (replaced by [azuread_service_principal](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_service_principal))
  - [azure_ad_user](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_ad_user) (replaced by [azuread_user](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_user))

## v0.19.0 [2021-10-07]

_What's new?_

- New tables added
  - [azure_app_configuration](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_app_configuration) ([#344](https://github.com/turbot/steampipe-plugin-azure/pull/344))
  - [azure_application_gateway](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_application_gateway) ([#316](https://github.com/turbot/steampipe-plugin-azure/pull/316))
  - [azure_cognitive_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_cognitive_account) ([#301](https://github.com/turbot/steampipe-plugin-azure/pull/301))
  - [azure_compute_disk_access](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_access) ([#288](https://github.com/turbot/steampipe-plugin-azure/pull/288))
  - [azure_databox_edge_device](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_databox_edge_device) ([#377](https://github.com/turbot/steampipe-plugin-azure/pull/377))
  - [azure_eventgrid_domain](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_eventgrid_domain) ([#314](https://github.com/turbot/steampipe-plugin-azure/pull/314))
  - [azure_eventgrid_topic](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_eventgrid_topic) ([#352](https://github.com/turbot/steampipe-plugin-azure/pull/352))
  - [azure_frontdoor](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_frontdoor) ([#362](https://github.com/turbot/steampipe-plugin-azure/pull/362))
  - [azure_hdinsight_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hdinsight_cluster) ([#395](https://github.com/turbot/steampipe-plugin-azure/pull/395))
  - [azure_healthcare_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_healthcare_service) ([#345](https://github.com/turbot/steampipe-plugin-azure/pull/345))
  - [azure_hpc_cache](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hpc_cache) ([#374](https://github.com/turbot/steampipe-plugin-azure/pull/374))
  - [azure_hybrid_compute_machine](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hybrid_compute_machine) ([#378](https://github.com/turbot/steampipe-plugin-azure/pull/378))
  - [azure_hybrid_kubernetes_connected_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_hybrid_kubernetes_connected_cluster) ([#376](https://github.com/turbot/steampipe-plugin-azure/pull/376))
  - [azure_iothub_dps](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_iothub_dps) ([#363](https://github.com/turbot/steampipe-plugin-azure/pull/363))
  - [azure_kusto_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_kusto_cluster) ([#369](https://github.com/turbot/steampipe-plugin-azure/pull/369))
  - [azure_machine_learning_workspace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_machine_learning_workspace) ([#315](https://github.com/turbot/steampipe-plugin-azure/pull/315))
  - [azure_mssql_virtual_machine](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mssql_virtual_machine) ([#313](https://github.com/turbot/steampipe-plugin-azure/pull/313))
  - [azure_service_fabric_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_service_fabric_cluster) ([#310](https://github.com/turbot/steampipe-plugin-azure/pull/310))
  - [azure_signalr_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_signalr_service) ([#328](https://github.com/turbot/steampipe-plugin-azure/pull/328))
  - [azure_spring_cloud_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_spring_cloud_service) ([#347](https://github.com/turbot/steampipe-plugin-azure/pull/347))
  - [azure_storage_sync](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_sync) ([#326](https://github.com/turbot/steampipe-plugin-azure/pull/326))
  - [azure_synapse_workspace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_synapse_workspace) ([#346](https://github.com/turbot/steampipe-plugin-azure/pull/346))

_Enhancements_

- Added `encryption_scope` column to `azure_storage_account` table ([#392](https://github.com/turbot/steampipe-plugin-azure/pull/392))
- Added `security_profile` column to `azure_compute_virtual_machine` table ([#387](https://github.com/turbot/steampipe-plugin-azure/pull/387))
- Added `public_network_access` and `minimal_tls_version` column to `azure_sql_server` table ([#371](https://github.com/turbot/steampipe-plugin-azure/pull/371))
- Added `guest_configuration_assignments` column to `azure_compute_virtual_machine` table ([#353](https://github.com/turbot/steampipe-plugin-azure/pull/353)) ([#380](https://github.com/turbot/steampipe-plugin-azure/pull/380))
- Added `cluster_settings` column to `azure_app_service_environment` table ([#360](https://github.com/turbot/steampipe-plugin-azure/pull/360))
- Added `developer_portal_url`, `disable_gateway`, `enable_client_certificate`, `api_version_constraint`, `certificates`, `custom_properties`, `identity_user_assigned_identities`, `virtual_network_type`, `restore`, `scm_url`, `zones` and `diagnostic_settings` columns to `azure_api_management` table ([#336](https://github.com/turbot/steampipe-plugin-azure/pull/336))
- Added `security_alert_policies` column to `azure_mssql_managed_instance` table ([#333](https://github.com/turbot/steampipe-plugin-azure/pull/333))
- Added `private_endpoint_connections` column to `azure_eventhub_namespace` table ([#331](https://github.com/turbot/steampipe-plugin-azure/pull/331))
- Added `private_endpoint_connections` column to `azure_servicebus_namespace` table ([#334](https://github.com/turbot/steampipe-plugin-azure/pull/334))
- Added `server_keys` column to `azure_mysql_server` table ([#337](https://github.com/turbot/steampipe-plugin-azure/pull/337))
- Added `identity` column to `azure_compute_virtual_machine` table ([#341](https://github.com/turbot/steampipe-plugin-azure/pull/341))
- Added `server_keys` column to `azure_postgresql_server` table ([#299](https://github.com/turbot/steampipe-plugin-azure/pull/299))
- Added `private_endpoint_connections` column to `azure_key_vault` table ([#306](https://github.com/turbot/steampipe-plugin-azure/pull/306)) ([#342](https://github.com/turbot/steampipe-plugin-azure/pull/342))
- Added `private_endpoint_connections` column to `azure_sql_server` table ([#300](https://github.com/turbot/steampipe-plugin-azure/pull/300))
- Added `private_endpoint_connections` column to `azure_data_factory` table ([#298](https://github.com/turbot/steampipe-plugin-azure/pull/298))

_Bug fixes_

- Querying column `encryption` in table `azure_servicebus_namespace` now render all available properties ([#366](https://github.com/turbot/steampipe-plugin-azure/pull/366))
- Querying column `cluster_arm_id` in table `azure_eventhub_namespace` no longer return `nil` if available ([#351](https://github.com/turbot/steampipe-plugin-azure/pull/351))
- Querying column `private_endpoint_connections` in table `azure_mysql_server` now render all available properties ([#338](https://github.com/turbot/steampipe-plugin-azure/pull/338))
- Querying column `private_endpoint_connections` in table `azure_postgresql_server` now render all available properties ([#339](https://github.com/turbot/steampipe-plugin-azure/pull/339))
- Querying column `encryption_protector` in table `azure_sql_server` now render all available properties ([#361](https://github.com/turbot/steampipe-plugin-azure/pull/361))


## v0.18.0 [2021-08-25]

_What's new?_

- New tables added
  - [azure_key_vault_deleted_vault](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_deleted_vault) ([#263](https://github.com/turbot/steampipe-plugin-azure/pull/263))
  - [azure_lb_backend_address_pool](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_backend_address_pool) ([#266](https://github.com/turbot/steampipe-plugin-azure/pull/266))
  - [azure_lb_nat_rule](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_nat_rule) ([#267](https://github.com/turbot/steampipe-plugin-azure/pull/267))
  - [azure_lb_outbound_rule](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_outbound_rule) ([#264](https://github.com/turbot/steampipe-plugin-azure/pull/264))
  - [azure_mssql_elasticpool](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mssql_elasticpool) ([#276](https://github.com/turbot/steampipe-plugin-azure/pull/276))
  - [azure_mssql_managed_instance](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mssql_managed_instance) ([#277](https://github.com/turbot/steampipe-plugin-azure/pull/277))

_Enhancements_

- Updated: Add `vulnerability_assessments` and `vulnerability_assessment_scan_records` columns in `azure_sql_database` table ([#279](https://github.com/turbot/steampipe-plugin-azure/pull/279))


## v0.17.0 [2021-08-13]

_What's new?_

- New tables added
  - [azure_compute_virtual_machine_scale_set](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_scale_set) ([#249](https://github.com/turbot/steampipe-plugin-azure/pull/249))
  - [azure_data_lake_analytics_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_lake_analytics_account) ([#253](https://github.com/turbot/steampipe-plugin-azure/pull/253))
  - [azure_lb](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb) ([#231](https://github.com/turbot/steampipe-plugin-azure/pull/231))
  - [azure_lb_probe](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_probe) ([#238](https://github.com/turbot/steampipe-plugin-azure/pull/238))
  - [azure_lb_rule](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_lb_rule) ([#235](https://github.com/turbot/steampipe-plugin-azure/pull/235))
  - [azure_resource_link](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_resource_link) ([#252](https://github.com/turbot/steampipe-plugin-azure/pull/252))
  - [azure_search_service](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_search_service) ([#257](https://github.com/turbot/steampipe-plugin-azure/pull/257))

_Enhancements_

- Updated: Add `retention_policy_id`, `retention_policy_name`, `retention_policy_type` and `retention_policy_property` columns in `azure_sql_database` table ([#255](https://github.com/turbot/steampipe-plugin-azure/pull/255))

_Bug fixes_

- Fixed: Integration test issues for several tables ([#259](https://github.com/turbot/steampipe-plugin-azure/pull/259))
- Fixed: Expired CLI authentication tokens will now automatically be refreshed ([#234](https://github.com/turbot/steampipe-plugin-azure/pull/234))

## v0.16.0 [2021-08-11]

_What's new?_

- New tables added
  - [azure_batch_account](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_batch_account) ([#242](https://github.com/turbot/steampipe-plugin-azure/pull/242))
  - [azure_data_lake_store](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_lake_store) ([#217](https://github.com/turbot/steampipe-plugin-azure/pull/217))
  - [azure_iothub](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_iothub) ([#232](https://github.com/turbot/steampipe-plugin-azure/pull/232))
  - [azure_key_vault_managed_hardware_security_module](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_managed_hardware_security_module) ([#236](https://github.com/turbot/steampipe-plugin-azure/pull/236))
  - [azure_logic_app_workflow](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_logic_app_workflow) ([#230](https://github.com/turbot/steampipe-plugin-azure/pull/230))
  - [azure_recovery_services_vault](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_recovery_services_vault) ([#243](https://github.com/turbot/steampipe-plugin-azure/pull/243))
  - [azure_stream_analytics_job](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_stream_analytics_job) ([#246](https://github.com/turbot/steampipe-plugin-azure/pull/246))

_Enhancements_

- Updated: Add `diagnostic_settings` column in `azure_network_security_group` table ([#247](https://github.com/turbot/steampipe-plugin-azure/pull/247))
- Updated: Add `ExtensionType` property in `extensions` column for `azure_compute_virtual_machine` table ([#229](https://github.com/turbot/steampipe-plugin-azure/pull/229))
- Updated: Add `enable_automatic_updates`, `provision_vm_agent_windows`, `time_zone`, `additional_unattend_content`, `patch_settings` and `win_rm` columns in `azure_compute_virtual_machine` table ([#223](https://github.com/turbot/steampipe-plugin-azure/pull/223))
- Updated: Add `diagnostic_settings` column in `azure_servicebus_namespace` table ([#225](https://github.com/turbot/steampipe-plugin-azure/pull/225))
- Updated: Add `diagnostic_settings` column in `azure_eventhub_namespace` table ([#226](https://github.com/turbot/steampipe-plugin-azure/pull/226))
- Updated: Add `network_acls` column in `azure_key_vault` table ([#220](https://github.com/turbot/steampipe-plugin-azure/pull/220))
- Updated: Add `virtual_network_rules` column in `azure_sql_server` table ([#227](https://github.com/turbot/steampipe-plugin-azure/pull/227))
- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v141--2021-07-20) ([#207](https://github.com/turbot/steampipe-plugin-azure/pull/207))

_Bug fixes_

- Fixed: Pagination for listing resources for all tables ([#254](https://github.com/turbot/steampipe-plugin-azure/pull/254))<br/>
  _This bug impacted all the tables in plugin. Now tables will not go into infinite loop for large number of resources._
- Fixed: Improve properties of `firewall_rules` column in `azure_sql_server` table ([#237](https://github.com/turbot/steampipe-plugin-azure/pull/237))

## v0.15.0 [2021-07-31]

_What's new?_

- New tables added
  - [azure_container_registry](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_container_registry) ([#196](https://github.com/turbot/steampipe-plugin-azure/pull/196))
  - [azure_eventhub_namespace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_eventhub_namespace) ([#201](https://github.com/turbot/steampipe-plugin-azure/pull/201))
  - [azure_mariadb_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mariadb_server) ([#206](https://github.com/turbot/steampipe-plugin-azure/pull/206))
  - [azure_redis_cache](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_redis_cache) ([#203](https://github.com/turbot/steampipe-plugin-azure/pull/203))
  - [azure_security_center_jit_network_access_policy](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_jit_network_access_policy) ([#192](https://github.com/turbot/steampipe-plugin-azure/pull/192))
  - [azure_servicebus_namespace](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_servicebus_namespace) ([#200](https://github.com/turbot/steampipe-plugin-azure/pull/200))

_Enhancements_

- Updated: Add column `vnet_connection` to `azure_app_service_web_app` table ([#204](https://github.com/turbot/steampipe-plugin-azure/pull/204))

## v0.14.0 [2021-07-22]

_What's new?_

- New tables added
  - [azure_data_factory_pipeline](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_factory_pipeline) ([#169](https://github.com/turbot/steampipe-plugin-azure/pull/169))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v131--2021-07-15)

_Bug fixes_

- Fixed: `azure_virtual_network_gateway` table's parent hydrate now lists resource groups instead of virtual networks to prevent duplicate rows ([#181](https://github.com/turbot/steampipe-plugin-azure/pull/181))

## v0.13.0 [2021-07-13]

_What's new?_

- New tables added
  - [azure_compute_disk_metric_read_ops](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_read_ops) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_read_ops_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_read_ops_daily) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_read_ops_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_read_ops_hourly) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_write_ops](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_write_ops) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_write_ops_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_write_ops_daily) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_disk_metric_write_ops_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_disk_metric_write_ops_hourly) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_virtual_machine_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_cpu_utilization) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_virtual_machine_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_cpu_utilization_daily) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_compute_virtual_machine_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_compute_virtual_machine_metric_cpu_utilization_hourly) ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
  - [azure_data_factory](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_factory) ([#158](https://github.com/turbot/steampipe-plugin-azure/pull/158))
  - [azure_data_factory_dataset](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_data_factory_dataset) ([#168](https://github.com/turbot/steampipe-plugin-azure/pull/168))
  - [azure_express_route_circuit](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_express_route_circuit) ([#170](https://github.com/turbot/steampipe-plugin-azure/pull/170))
  - [azure_virtual_network_gateway](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_virtual_network_gateway) ([#157](https://github.com/turbot/steampipe-plugin-azure/pull/157))

_Enhancements_

- Updated: `azure-sdk-for-go` to `v55.4.0+incompatible` ([#166](https://github.com/turbot/steampipe-plugin-azure/pull/166))
- Updated: Change several metric table function names to be consistent with naming standards ([#173](https://github.com/turbot/steampipe-plugin-azure/pull/173))

_Bug fixes_

- Fixed: Integration tests for several tables and remove unused integration tests ([#175](https://github.com/turbot/steampipe-plugin-azure/pull/175))

## v0.12.0 [2021-07-01]

_Enhancements_

- Updated: Add `lifecycle_management_policy` column to `azure_storage_account` table ([#155](https://github.com/turbot/steampipe-plugin-azure/pull/155))

## v0.11.0 [2021-06-03]

_What's new?_

- New tables added
  - [azure_policy_definition](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_policy_definition) ([#141](https://github.com/turbot/steampipe-plugin-azure/pull/141))
  - [azure_tenant](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_tenant) ([#142](https://github.com/turbot/steampipe-plugin-azure/pull/142))

## v0.10.0 [2021-05-27]

_What's new?_

- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)
- New tables added
  - [azure_storage_blob](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_blob) ([#133](https://github.com/turbot/steampipe-plugin-azure/pull/133))

_Bug fixes_

- Fixed: Improved error messages when we fail to get credentials from the Azure CLI ([#137](https://github.com/turbot/steampipe-plugin-azure/pull/137))

## v0.9.0 [2021-05-20]

_What's new?_

- New tables added
  - [azure_policy_assignment](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_policy_assignment) ([#123](https://github.com/turbot/steampipe-plugin-azure/pull/123))
  - [azure_security_center_setting](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_setting) ([#115](https://github.com/turbot/steampipe-plugin-azure/pull/115))
  - [azure_security_center_subscription_pricing](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_subscription_pricing) ([#135](https://github.com/turbot/steampipe-plugin-azure/pull/135))
  - [azure_subscription](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_subscription) ([#132](https://github.com/turbot/steampipe-plugin-azure/pull/132))

## v0.8.0 [2021-05-13]

_What's new?_

- New tables added
  - [azure_security_center_auto_provisioning](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_auto_provisioning) ([#117](https://github.com/turbot/steampipe-plugin-azure/pull/117))
  - [azure_security_center_contact](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_security_center_contact) ([#121](https://github.com/turbot/steampipe-plugin-azure/pull/121))

_Enhancements_

- Updated: README.md and docs/index.md now contain links to our Slack community ([#129](https://github.com/turbot/steampipe-plugin-azure/pull/129))
- Updated: Bump lodash from 4.17.20 to 4.17.21 in /azure-test ([#126](https://github.com/turbot/steampipe-plugin-azure/pull/126))

## v0.7.0 [2021-05-06]

_What's new?_

- New tables added
  - [azure_kubernetes_cluster](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_kubernetes_cluster) ([#83](https://github.com/turbot/steampipe-plugin-azure/pull/83))

_Enhancements_

- Updated: Add `identity` column to `azure_app_service_web_app` table ([#90](https://github.com/turbot/steampipe-plugin-azure/pull/90))
- Updated: Add `os_disk_vhd_uri` column to `azure_compute_virtual_machine` table ([#88](https://github.com/turbot/steampipe-plugin-azure/pull/88))

_Bug fixes_

- Fixed: The get calls in the `azure_key_vault_secret` table should not fail for disabled secrets ([#111](https://github.com/turbot/steampipe-plugin-azure/pull/111))

## v0.6.0 [2021-04-29]

_What's new?_

- New tables added
  - [azure_key_vault_key](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_key) ([#73](https://github.com/turbot/steampipe-plugin-azure/pull/73))
  - [azure_network_watcher_flow_log](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_network_watcher_flow_log) ([#61](https://github.com/turbot/steampipe-plugin-azure/pull/61))

_Enhancements_

- Updated: Add `extensions` column to `azure_compute_virtual_machine` table ([#85](https://github.com/turbot/steampipe-plugin-azure/pull/85))

_Bug fixes_

- Fixed: `azure_key_vault` table queries no longer crash when getting vault diagnostic settings ([#107](https://github.com/turbot/steampipe-plugin-azure/pull/107))
- Fixed: `deleted_time` and `last_modified_time` columns now show the correct dates in `azure_storage_container` table ([#106](https://github.com/turbot/steampipe-plugin-azure/pull/106))
- Fixed: `encryption_key_vault_properties_last_rotation_time` column now shows the correct date in `azure_storage_account` table ([#101](https://github.com/turbot/steampipe-plugin-azure/pull/101))
- Fixed: `subscription_id` column now displays the correct subscription ID in `azure_diagnostic_setting` table ([#99](https://github.com/turbot/steampipe-plugin-azure/pull/99))
- Fixed: Column name in example query in `azure_key_vault_secret` table doc ([#108](https://github.com/turbot/steampipe-plugin-azure/pull/108))

## v0.5.0 [2021-04-22]

_What's new?_

- New tables added
  - [azure_diagnostic_setting](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_diagnostic_setting) ([#70](https://github.com/turbot/steampipe-plugin-azure/pull/70))
  - [azure_key_vault_secret](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_key_vault_secret) ([#76](https://github.com/turbot/steampipe-plugin-azure/pull/76))
  - [azure_log_alert](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_log_alert) ([#70](https://github.com/turbot/steampipe-plugin-azure/pull/70))
  - [azure_log_profile](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_log_profile) ([#70](https://github.com/turbot/steampipe-plugin-azure/pull/70))
  - [azure_sql_database](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_sql_database) ([#93](https://github.com/turbot/steampipe-plugin-azure/pull/93))

_Enhancements_

- Updated: Add `diagnostic_settings` column to `azure_key_vault` table ([#96](https://github.com/turbot/steampipe-plugin-azure/pull/96))
- Updated: Add `queue` prefix to various queue logging columns in `azure_storage_account` table ([#94](https://github.com/turbot/steampipe-plugin-azure/pull/94))

_Bug fixes_

- Fixed: List calls should not infinitely loop in `azure_key_vault` table ([#82](https://github.com/turbot/steampipe-plugin-azure/pull/82))

## v0.4.0 [2021-04-08]

_What's new?_

- New tables added
  - [azure_mysql_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_mysql_server) ([#66](https://github.com/turbot/steampipe-plugin-azure/pull/66))
  - [azure_postgresql_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_postgresql_server) ([#64](https://github.com/turbot/steampipe-plugin-azure/pull/64))
  - [azure_storage_container](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_storage_container) ([#71](https://github.com/turbot/steampipe-plugin-azure/pull/71))

_Enhancements_

- Updated: Add `auth_settings` and `configuration` columns to `azure_app_service_function_app` table ([#77](https://github.com/turbot/steampipe-plugin-azure/pull/77))
- Updated: Add `auth_settings` and `configuration` columns to `azure_app_service_web_app` table ([#77](https://github.com/turbot/steampipe-plugin-azure/pull/77))
- Updated: Add `blob_service_logging` column to `azure_storage_account` table ([#80](https://github.com/turbot/steampipe-plugin-azure/pull/80))

_Bug fixes_

- Fixed: The table `azure_sql_server` should return `null` instead of an empty object for columns with missing data ([#68](https://github.com/turbot/steampipe-plugin-azure/pull/68))

## v0.3.1 [2021-03-25]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v0.2.6](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v026-2021-03-18)

_Documentation_

- Fixed various example queries in `azure_app_service_web_app` table document ([#59](https://github.com/turbot/steampipe-plugin-azure/pull/59))

## v0.3.0 [2021-03-11]

_What's new?_

- New tables added
  - [azure_sql_server](https://hub.steampipe.io/plugins/turbot/azure/tables/azure_sql_server)

_Bug fixes_

- Removed use of deprecated `ItemFromKey` function from all tables

## v0.2.2 [2021-03-02]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve issue:
  - Fix tables failing with error similar to `Error: pq: rpc error: code = Internal desc = get hydrate function getAdGroup failed with panic interface conversion: interface {} is nil, not *graphrbac.ADGroup`([#29](https://github.com/turbot/steampipe-plugin-azure/issues/29)).

## v0.2.1 [2021-02-25]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve SDK issues:
  - Fix error for missing required quals [#40](https://github.com/turbot/steampipe-plugin-sdk/issues/42).
  - Queries fail with error socket: too many open files [#190](https://github.com/turbot/steampipe/issues/190)

## v0.2.0 [2021-02-18]

_What's new?_

- Added support for [connection configuration](https://github.com/turbot/steampipe-plugin-azure/blob/main/docs/index.md#connection-configuration). You may specify azure `tenant_id`, `subscription_id`, `client_id`, `client_secret`, `client_certificate_path`, `client_certificate_password`, `username` and `password` for each connection in a configuration file. You can have multiple azure connections, each configured for a different azure subscription.

_Enhancements_

- Updates
  - Added columns power_state, private_ips and public_ips to azure_compute_virtual_machine table ([#17](https://github.com/turbot/steampipe-plugin-azure/pull/17))

_Bug fixes_

- Breaking changes

  - Renamed earlier `azure_storage_blob` table to `azure_storage_blob_service` ([#7](https://github.com/turbot/steampipe-plugin-azure/pull/7))
  - Renamed earlier `azure_storage_table` table to `azure_storage_table_service` ([#10](https://github.com/turbot/steampipe-plugin-azure/pull/10))
  - Removed columns managed_disk_storage_account_type and os_disk_size_gb from `azure_compute_virtual_machine` table ([#17](https://github.com/turbot/steampipe-plugin-azure/pull/17))
