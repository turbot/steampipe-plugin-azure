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
