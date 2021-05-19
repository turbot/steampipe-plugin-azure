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
